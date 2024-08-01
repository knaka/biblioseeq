package fts

import (
	"context"
	"database/sql"
	"github.com/knaka/biblioseeq/db"
	"github.com/knaka/biblioseeq/db/sqlcgen"
	ftslog "github.com/knaka/biblioseeq/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/samber/lo"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	. "github.com/knaka/go-utils"
)

type Indexer struct {
	dbConn          *sql.DB
	dirIndexers     []*DirIndexer
	shouldMigrateDB bool
}

type Option func(*Indexer)

func WithTargetDirectory(absPath string, evalPath string, fileExtensions []string) Option {
	return func(e *Indexer) {
		e.dirIndexers = append(e.dirIndexers, NewDirIndexer(e.dbConn, absPath, evalPath, fileExtensions))
	}
}

func WithDBConn(dbConn *sql.DB) Option {
	return func(e *Indexer) {
		e.dbConn = dbConn
	}
}

func MigratesDB() Option {
	return func(e *Indexer) {
		e.shouldMigrateDB = true
	}
}

func WithDefaultDBFilePath() Option {
	return func(e *Indexer) {
		e.dbConn = V(sql.Open("sqlite3", V(GetDefaultDBFilePath())))
	}
}

func WithDBFile(path string) Option {
	return WithDBConn(V(sql.Open("sqlite3", path)))
}

func NewIndexer(opts ...Option) (ret *Indexer) {
	ret = &Indexer{}
	for _, opt := range opts {
		opt(ret)
	}
	if ret.shouldMigrateDB {
		dbPath := V(GetDefaultDBFilePath())
		db.Migrate(dbPath)
	}
	if ret.dbConn == nil {
		panic("dbConn is not set")
	}
	return
}

func (indexer *Indexer) WatchContinuously(ctx context.Context) {
	wg := sync.WaitGroup{}
	for _, dirIndexer := range indexer.dirIndexers {
		wg.Add(1)
		go (func() {
			dirIndexer.WatchContinuously(ctx)
			wg.Done()
		})()
		dirIndexer.WaitForWatchingStarted(ctx)
		if true {
			wg.Add(1)
			go (func() {
				dirIndexer.IndexAll()
				wg.Done()
			})()
		} else {
			dirIndexer.chScanned <- "scanned"
		}
	}
	wg.Wait()
}

func (indexer *Indexer) WaitForInitialScanFinished(ctx context.Context) {
	for _, dirIndexer := range indexer.dirIndexers {
		dirIndexer.WaitForInitialScanFinished(ctx)
	}
}

func (indexer *Indexer) getDirectoryIndexer(path string) (directory *DirIndexer) {
	for _, dir := range indexer.dirIndexers {
		if _, errSub := filepath.Rel(dir.absPath, path); errSub == nil {
			return dir
		}
	}
	return
}

func (indexer *Indexer) indexFile(filePath string) {
	dirIndexer := indexer.getDirectoryIndexer(filePath)
	if lo.NoneBy(dirIndexer.fileExtensions, func(ext string) bool {
		return filepath.Ext(filePath) == ext
	}) {
		return
	}
}

var reSpaces = sync.OnceValue(func() *regexp.Regexp { return regexp.MustCompile(`\s+`) })

type QueryResult struct {
	Path       string
	Title      string
	Tags       []string
	DirPath    string
	Snippet    string
	ModifiedAt time.Time
}

func (indexer *Indexer) Query(query string) (results []*QueryResult, err error) {
	defer Catch(&err)
	ctx := context.Background()
	store := sqlcgen.New(indexer.dbConn)
	if query == "" {
		rows, err := store.LatestEntries(ctx, &sqlcgen.LatestEntriesParams{})
		if err != nil {
			return nil, nil
		}
		dirPath := ""
		for _, dirIndexer := range indexer.dirIndexers {
			if _, errSub := filepath.Rel(dirIndexer.absPath, rows[0].Path); errSub == nil {
				dirPath = dirIndexer.absPath
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
			ftslog.Println("Invalid UTF-8 snippet string", item.File.Path)
			return nil, false
		}
		if !utf8.ValidString(item.File.Path) {
			ftslog.Println("Invalid UTF-8 path string", item.File.Path)
			return nil, false
		}
		snippet = RemoveZwsp(snippet)
		snippet = reSpaces().ReplaceAllString(snippet, " ")
		dirPath := ""
		for _, dirIndexer := range indexer.dirIndexers {
			if _, errSub := filepath.Rel(dirIndexer.absPath, item.File.Path); errSub == nil {
				dirPath = dirIndexer.absPath
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

func (indexer *Indexer) InitialScanFinished() bool {
	return lo.EveryBy(indexer.dirIndexers, func(dirIndexer *DirIndexer) bool {
		return dirIndexer.InitialScanFinished()
	})
}
