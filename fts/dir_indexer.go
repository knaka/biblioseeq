package fts

import (
	"bufio"
	"context"
	"database/sql"
	"github.com/knaka/biblioseeq/db/sqlcgen"
	ftslog "github.com/knaka/biblioseeq/log"
	. "github.com/knaka/go-utils"
	"github.com/rjeczalik/notify"
	"github.com/samber/lo"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type DirIndexer struct {
	absPath        string
	evalPath       string
	fileExtensions []string
	indexMutex     sync.Mutex
	dbConn         *sql.DB
	chStop         chan any
	chStarted      chan any
	chScanned      chan any
}

func NewDirIndexer(dbConn *sql.DB, absPath string, evalPath string, fileExtensions []string) *DirIndexer {
	return &DirIndexer{
		dbConn:         dbConn,
		chStop:         make(chan any),
		chStarted:      make(chan any, 2),
		chScanned:      make(chan any, 2),
		absPath:        absPath,
		evalPath:       evalPath,
		fileExtensions: fileExtensions,
	}
}

func extractTitleAndTags(path string, body *string) (title string, tags []string) {
	switch filepath.Ext(path) {
	case ".md", ".MD":
		reader := strings.NewReader(*body)
		scanner := bufio.NewScanner(reader)
		firstLine := true
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if firstLine {
				if line == "---" {
					firstLine = false
					continue
				} else {
					break
				}
			}
			if line == "---" {
				break
			}
			divs := strings.SplitN(line, ":", 2)
			if len(divs) != 2 {
				continue
			}
			key := strings.TrimSpace(strings.ToLower(divs[0]))
			value := strings.TrimSpace(divs[1])
			switch key {
			case "title":
				title = value
			case "tags":
				tags = strings.Split(strings.ToLower(value), ",")
			}
		}
	}
	return
}

func (dirIndexer *DirIndexer) IndexExistingFile(path string) {
	ctx := context.Background()
	store := sqlcgen.New(dirIndexer.dbConn)
	if lo.NoneBy(dirIndexer.fileExtensions, func(ext string) bool {
		return filepath.Ext(path) == ext
	}) {
		return
	}
	stat, err := os.Stat(path)
	// broken symlink?
	if err != nil {
		return
	}
	if stat.IsDir() {
		return
	}
	if file := R(store.GetFile(ctx, &sqlcgen.GetFileParams{Path: path})).NilIf(sql.ErrNoRows); file == nil {
		addCache(path, stat, store)
	} else {
		updateCacheIfModified(file, stat, store)
	}
}

func addCache(path string, stat os.FileInfo, store *sqlcgen.Queries) {
	body := string(V(os.ReadFile(path)))
	title, tags := extractTitleAndTags(path, &body)
	bodyTokenized := tokenizeJapanese(body)
	ftslog.Println("Adding index for file:", path)
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
}

func updateCacheIfModified(file *sqlcgen.File, stat os.FileInfo, store *sqlcgen.Queries) {
	path := file.Path
	if file.ModifiedAt.Equal(stat.ModTime()) && file.Size == stat.Size() {
		ftslog.Println("Skipped file:", path)
		return
	}
	ftslog.Println("Updating index for file:", path)
	body := string(V(os.ReadFile(path)))
	title, tags := extractTitleAndTags(path, &body)
	bodyTokenized := tokenizeJapanese(body)
	V0(store.UpdateFtsFile(context.Background(), &sqlcgen.UpdateFtsFileParams{
		Path: path,
		Body: bodyTokenized,
	}))
	V0(store.UpdateFile(context.Background(), &sqlcgen.UpdateFileParams{
		Path:       path,
		Title:      title,
		Tags:       strings.Join(tags, ","),
		ModifiedAt: stat.ModTime(),
		Size:       stat.Size(),
	}))
}

func (dirIndexer *DirIndexer) DeleteIndexRecursively(path string) {
	dirIndexer.indexMutex.Lock()
	defer dirIndexer.indexMutex.Unlock()
	store := sqlcgen.New(dirIndexer.dbConn)
	ctx := context.Background()
	V0(store.DeleteFtsFiles(ctx, &sqlcgen.DeleteFtsFilesParams{OptPathPrefix: &path}))
	V0(store.DeleteFiles(ctx, &sqlcgen.DeleteFilesParams{OptPathPrefix: &path}))
}

func (dirIndexer *DirIndexer) onEvent(eventInfo notify.EventInfo) {
	dirIndexer.indexMutex.Lock()
	defer dirIndexer.indexMutex.Unlock()
	event := eventInfo.Event()
	path := eventInfo.Path()
	if last, found := strings.CutPrefix(path, dirIndexer.evalPath); found {
		path = dirIndexer.absPath + last
	}
	//sys := eventInfo.Sys()
	if event&notify.Write != 0 {
		// ファイルへの書き込み。単体更新
		dirIndexer.IndexExistingFile(path)
	}
	if event&notify.Remove != 0 {
		// log.Println("ファイルもしくはディレクトリの削除。パス前方一致削除。再帰で通知が来ていたら、配下のファイルはすでに消されていると思われるので、何もしなくて良いこともある")
		// 消えたのがディレクトリかは分からない。もう無いので
		dirIndexer.DeleteIndexRecursively(path)
	}
	if event&notify.Rename != 0 {
		stat := R(os.Stat(path)).NilIf(os.ErrNotExist)
		if stat == nil {
			dirIndexer.DeleteIndexRecursively(path)
		}
		if stat.IsDir() {
			dirIndexer.indexDirectory(path)
		} else {
			dirIndexer.IndexExistingFile(path)
		}
	}
}

func (dirIndexer *DirIndexer) indexExistingFileWalkFunc(path string, info fs.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	dirIndexer.IndexExistingFile(path)
	return nil
}

func (dirIndexer *DirIndexer) indexDirectory(path string) {
	V0(filepath.Walk(path, dirIndexer.indexExistingFileWalkFunc))
}

func (dirIndexer *DirIndexer) IndexAll() {
	dirIndexer.indexMutex.Lock()
	defer dirIndexer.indexMutex.Unlock()

	ctx := context.Background()
	store := sqlcgen.New(dirIndexer.dbConn)
	ch := make(chan *sqlcgen.File)
	go (func() {
		// todo: pagination
		fileList, err := store.GetFiles(ctx)
		if err != nil {
			return
		}
		for _, file := range fileList {
			ch <- file
		}
		close(ch)
	})()
	remove := func(path string) {
		ftslog.Println("Removing index for file:", path)
		V0(store.DeleteFtsFiles(ctx, &sqlcgen.DeleteFtsFilesParams{OptPathPrefix: &path}))
		V0(store.DeleteFiles(ctx, &sqlcgen.DeleteFilesParams{OptPathPrefix: &path}))
	}
	cachedPaths := make(map[string]any)
	for file := range ch {
		if lo.NoneBy(dirIndexer.fileExtensions, func(ext string) bool {
			return filepath.Ext(file.Path) == ext
		}) {
			remove(file.Path)
			continue
		}
		if stat, err := os.Stat(file.Path); err == nil {
			if stat.IsDir() {
				remove(file.Path)
				continue
			}
			updateCacheIfModified(file, stat, store)
			cachedPaths[file.Path] = nil
		} else {
			remove(file.Path)
			continue
		}
	}

	V0(filepath.Walk(dirIndexer.absPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if _, found := cachedPaths[path]; found {
			return nil
		}
		if lo.NoneBy(dirIndexer.fileExtensions, func(ext string) bool {
			return filepath.Ext(path) == ext
		}) {
			return nil
		}
		addCache(path, info, store)
		return nil
	}))

	dirIndexer.chScanned <- "scanned"
}

func (dirIndexer *DirIndexer) WatchContinuously(ctx context.Context) {
	ch := make(chan notify.EventInfo, 10)
	V0(notify.Watch(filepath.Join(dirIndexer.absPath, "..."), ch,
		notify.Write, notify.Remove, notify.Rename))
	defer notify.Stop(ch)
	defer close(ch)
	dirIndexer.chStarted <- "started"
outer:
	for {
		select {
		case eventInfo := <-ch:
			dirIndexer.onEvent(eventInfo)
		case <-dirIndexer.chStop:
			break outer
		case <-ctx.Done():
			break outer
		}
	}
}

func (dirIndexer *DirIndexer) WaitForWatchingStarted(ctx context.Context) {
	select {
	case started := <-dirIndexer.chStarted:
		dirIndexer.chStarted <- started
	case <-ctx.Done():
	}
}

func (dirIndexer *DirIndexer) StopWatching() {
	dirIndexer.chStop <- "stop"
}

func (dirIndexer *DirIndexer) Watching() bool {
	return len(dirIndexer.chStop) == 0 && len(dirIndexer.chStarted) > 0
}

func (dirIndexer *DirIndexer) WaitForInitialScanFinished(ctx context.Context) {
	select {
	case scanned := <-dirIndexer.chScanned:
		dirIndexer.chScanned <- scanned
	case <-ctx.Done():
	}
}

func (dirIndexer *DirIndexer) InitialScanFinished() bool {
	return len(dirIndexer.chScanned) > 0
}
