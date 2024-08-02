package fts

import (
	"context"
	"database/sql"
	"github.com/knaka/biblioseeq/db/sqlcgen"
	"github.com/knaka/biblioseeq/log"
	"github.com/rjeczalik/notify"
	"github.com/samber/lo"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	. "github.com/knaka/go-utils"
)

type idxTargetDir struct {
	path           string
	fileExtensions []string
}

type IdxMgr struct {
	dbConn          *sql.DB
	targetDirs      []*idxTargetDir
	shouldMigrateDB bool
	mu              sync.Mutex
}

func (im *IdxMgr) addTarget(absPath string, fileExtensions []string) {
	im.targetDirs = append(im.targetDirs,
		&idxTargetDir{path: absPath, fileExtensions: fileExtensions},
	)
}

func newIdxMgr(dbFilePath string) (im *IdxMgr) {
	return &IdxMgr{
		dbConn: V(sql.Open("sqlite3", dbFilePath)),
	}
}

func (im *IdxMgr) indexExistingFileWhileWalking(path string, info fs.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	im.indexExistingFile(path) // walk
	return nil
}

func (im *IdxMgr) indexDirectory(dirPath string) {
	V0(filepath.Walk(dirPath, im.indexExistingFileWhileWalking))
}

func (im *IdxMgr) indexExistingFile(pathFile string) {
	if !im.matchesToTarget(pathFile) {
		return
	}
	ctx := context.Background()
	store := im.newStore()
	stat, err := os.Stat(pathFile)
	if err != nil {
		return
	}
	if stat.IsDir() {
		return
	}
	if file := R(store.GetFile(ctx, &sqlcgen.GetFileParams{Path: pathFile})).NilIf(sql.ErrNoRows); file == nil {
		im.addCache(pathFile, stat) // notify
	} else {
		im.updateCacheIfModified(file, stat) // notify
	}
}

func (im *IdxMgr) onNotifyEvent(eventInfo notify.EventInfo) {
	im.mu.Lock()
	defer im.mu.Unlock()
	event := eventInfo.Event()
	path := eventInfo.Path()
	if event&notify.Write != 0 {
		im.indexExistingFile(path) // notify write
	}
	if event&notify.Remove != 0 {
		// The path can be for file or directory because the removed entry does not exist.
		im.deleteFileOrDirIndex(path)
	}
	if event&notify.Rename != 0 {
		stat := R(os.Stat(path)).NilIf(os.ErrNotExist)
		if stat == nil {
			// The path can be for file or directory because the moved source does not exist.
			im.deleteFileOrDirIndex(path)
		}
		if stat.IsDir() {
			im.indexDirectory(path)
		} else {
			im.indexExistingFile(path) // notify rename
		}
	}
}

func (im *IdxMgr) matchesToTarget(filePath string) bool {
	return lo.SomeBy(im.targetDirs, func(dir *idxTargetDir) bool {
		return strings.HasPrefix(filePath, dir.path) && lo.SomeBy(dir.fileExtensions, func(ext string) bool {
			return filepath.Ext(filePath) == ext
		})
	})
}

func (im *IdxMgr) newStore() *sqlcgen.Queries {
	return sqlcgen.New(im.dbConn)
}

type txStore struct {
	*sqlcgen.Queries
	tx *sql.Tx
}

func (ts *txStore) commit() {
	V0(ts.tx.Commit())
}

func (ts *txStore) rollback() {
	Expect(ts.tx.Rollback(), nil, sql.ErrTxDone)
}

func (im *IdxMgr) newTxStore() (store txStore) {
	store.tx = V(im.dbConn.Begin())
	store.Queries = sqlcgen.New(store.tx)
	return
}

func (im *IdxMgr) addCache(path string, stat os.FileInfo) {
	store := im.newTxStore()
	defer store.rollback()
	bodyBytes, err := os.ReadFile(path)
	if err != nil {
		log.Println("Failed to read file:", path)
		return
	}
	body := string(bodyBytes)
	title, tags := extractTitleAndTags(path, &body)
	bodyTokenized := tokenizeJapanese(body)
	log.Println("Adding index for file:", path)
	idFtsFile := V(store.AddFtsFile(context.Background(), &sqlcgen.AddFtsFileParams{
		Body: bodyTokenized,
	}))
	V0(store.AddFile(context.Background(), &sqlcgen.AddFileParams{
		Path:       path,
		Title:      title,
		Tags:       strings.Join(tags, ","),
		FtsFileID:  idFtsFile,
		ModifiedAt: stat.ModTime(),
		Size:       stat.Size(),
	}))
	store.commit()
}

func (im *IdxMgr) updateCacheIfModified(file *sqlcgen.File, stat os.FileInfo) {
	if file.ModifiedAt.Equal(stat.ModTime()) && file.Size == stat.Size() {
		log.Println("Skipped file:", file.Path)
		return
	}
	ctx := context.Background()
	store := im.newTxStore()
	defer store.rollback()
	log.Println("Updating index for file:", file.Path)
	body := string(V(os.ReadFile(file.Path)))
	title, tags := extractTitleAndTags(file.Path, &body)
	bodyTokenized := tokenizeJapanese(body)
	V0(store.UpdateFtsFile(ctx, &sqlcgen.UpdateFtsFileParams{
		Path: file.Path,
		Body: bodyTokenized,
	}))
	V0(store.UpdateFile(ctx, &sqlcgen.UpdateFileParams{
		Path:       file.Path,
		Title:      title,
		Tags:       strings.Join(tags, ","),
		ModifiedAt: stat.ModTime(),
		Size:       stat.Size(),
	}))
	store.commit()
}

// synchronizeIndexesToFiles removes indexes for files that do not exist and updates indexes for files that exist. This make consistency between the file system and the database.
func (im *IdxMgr) synchronizeIndexesToFiles() {
	ctx := context.Background()
	store := sqlcgen.New(im.dbConn)
	chFileIndex := make(chan *sqlcgen.File)
	go (func() {
		// todo: paginate
		fileIndexList := V(store.GetFiles(ctx))
		for _, file := range fileIndexList {
			chFileIndex <- file
		}
		close(chFileIndex)
	})()
	indexedPaths := make(map[string]any)

	// Remove indexes for files that do not exist.

	removeFileIndex := func(filePath string) {
		storeSub := im.newTxStore()
		defer storeSub.rollback()
		V0(storeSub.DeleteFtsFiles(ctx, &sqlcgen.DeleteFtsFilesParams{OptFilePath: &filePath}))
		V0(storeSub.DeleteFiles(ctx, &sqlcgen.DeleteFilesParams{OptFilePath: &filePath}))
		log.Println("Removed index for file:", filePath)
		storeSub.commit()
	}
	for file := range chFileIndex {
		if !im.matchesToTarget(file.Path) {
			removeFileIndex(file.Path)
			continue
		}
		if stat, err := os.Stat(file.Path); err == nil {
			if stat.IsDir() {
				removeFileIndex(file.Path)
				continue
			}
			im.updateCacheIfModified(file, stat) // initial
			indexedPaths[file.Path] = nil
		} else {
			removeFileIndex(file.Path)
			continue
		}
	}

	// Add or update indexes for files that exist.

	for _, dir := range im.targetDirs {
		V0(filepath.Walk(dir.path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if _, ok := indexedPaths[path]; !ok {
				im.indexExistingFile(path) // synchronize
			}
			return nil
		}))
	}
}

func (im *IdxMgr) deleteFileOrDirIndex(path string) {
	store := im.newTxStore()
	defer store.rollback()
	ctx := context.Background()
	V0(store.DeleteFtsFiles(ctx, &sqlcgen.DeleteFtsFilesParams{
		OptFilePath: &path,
		OptDirPath:  &path,
	}))
	V0(store.DeleteFiles(ctx, &sqlcgen.DeleteFilesParams{
		OptFilePath: &path,
		OptDirPath:  &path,
	}))
	store.commit()
}

//func (im *IdxMgr) WaitForInitialScanFinished(ctx context.Context) {
//	for _, dirIndexer := range im.dirIndexers {
//		dirIndexer.WaitForInitialScanFinished(ctx)
//	}
//}

//func (im *IdxMgr) getDirectoryIndexer(path string) (directory *DirWatcher) {
//	for _, dir := range im.dirIndexers {
//		if _, errSub := filepath.Rel(dir.dirPath, path); errSub == nil {
//			return dir
//		}
//	}
//	return
//}

//func (im *IdxMgr) indexFile(filePath string) {
//	dirIndexer := im.getDirectoryIndexer(filePath)
//	if lo.NoneBy(dirIndexer.fileExtensions, func(ext string) bool {
//		return filepath.Ext(filePath) == ext
//	}) {
//		return
//	}
//}

var reSpaces = sync.OnceValue(func() *regexp.Regexp { return regexp.MustCompile(`\s+`) })

type QueryResult struct {
	Path       string
	Title      string
	Tags       []string
	DirPath    string
	Snippet    string
	ModifiedAt time.Time
}

func (im *IdxMgr) Query(query string) (results []*QueryResult, err error) {
	//defer Catch(&err)
	ctx := context.Background()
	store := im.newStore()
	if query == "" {
		rows, errSub := store.LatestEntries(ctx, &sqlcgen.LatestEntriesParams{})
		if errSub != nil {
			return nil, nil
		}
		if len(rows) == 0 {
			return
		}
		dirPath := ""
		for _, dir := range im.targetDirs {
			if strings.HasPrefix(rows[0].Path, dir.path) {
				dirPath = dir.path
				break
			}
		}
		return lo.FilterMap(rows, func(item *sqlcgen.File, _ int) (result *QueryResult, ok bool) {
			return &QueryResult{
				Path:       item.Path,
				Title:      item.Title,
				Tags:       strings.Split(item.Tags, ","),
				DirPath:    dirPath,
				Snippet:    "",
				ModifiedAt: item.ModifiedAt,
			}, true
		}), nil
	}
	rows, err := store.Query(ctx, &sqlcgen.QueryParams{
		Query: query,
	})
	// ToDo: Should return error for syntax error?
	if err != nil {
		return nil, nil
	}
	return lo.FilterMap(rows, func(item *sqlcgen.QueryRow, _ int) (*QueryResult, bool) {
		snippet := item.Snippet
		if !utf8.ValidString(snippet) {
			log.Println("Invalid UTF-8 snippet string", item.File.Path)
			return nil, false
		}
		if !utf8.ValidString(item.File.Path) {
			log.Println("Invalid UTF-8 path string", item.File.Path)
			return nil, false
		}
		snippet = RemoveZwsp(snippet)
		snippet = reSpaces().ReplaceAllString(snippet, " ")
		dirPath := ""
		for _, dir := range im.targetDirs {
			if strings.HasPrefix(item.File.Path, dir.path) {
				dirPath = dir.path
				break
			}
		}
		return &QueryResult{
			Path:       item.File.Path,
			Title:      item.File.Title,
			Tags:       strings.Split(item.File.Tags, ","),
			DirPath:    dirPath,
			Snippet:    snippet,
			ModifiedAt: item.File.ModifiedAt,
		}, true
	}), nil
}

func (im *IdxMgr) InitialScanFinished() bool {
	//return lo.EveryBy(im.dirIndexers, func(dirIndexer *DirWatcher) bool {
	//	return dirIndexer.InitialScanFinished()
	//})
	return false
}

//func NewIndexer(opts ...Option) (ret *IdxMgr) {
//	ret = &IdxMgr{}
//	for _, opt := range opts {
//		opt(ret)
//	}
//	if ret.shouldMigrateDB {
//		dbPath := V(getDefaultDBFilePath())
//		db.Migrate(dbPath)
//	}
//	if ret.dbConn == nil {
//		panic("dbConn is not set")
//	}
//	return
//}

//return lo.SomeBy(indexer.dirIndexers, func(dirIndexer *DirWatcher) bool {
//	return strings.HasPrefix(path, dirIndexer.dirPath) && lo.SomeBy(dirIndexer.fileExtensions, func(ext string) bool {
//		return filepath.Ext(path) == ext
//	})
//})
