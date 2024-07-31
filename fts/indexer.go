package fts

import (
	"context"
	"database/sql"
	"github.com/knaka/biblioseeq/db"
	"github.com/knaka/biblioseeq/db/sqlcgen"
	_ "github.com/mattn/go-sqlite3"
	"github.com/samber/lo"
	"path/filepath"
	"sync"

	. "github.com/knaka/go-utils"
)

type Indexer struct {
	dbConn          *sql.DB
	dirIndexers     []*DirIndexer
	shouldMigrateDB bool
}

type Option func(*Indexer)

func WithTargetDirectory(path string, fileExtensions []string) Option {
	return func(e *Indexer) {
		e.dirIndexers = append(e.dirIndexers, NewDirIndexer(e.dbConn, path, fileExtensions))
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
		wg.Add(1)
		go (func() {
			dirIndexer.IndexAll()
			wg.Done()
		})()
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
		if _, errSub := filepath.Rel(dir.path, path); errSub == nil {
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

type QueryResult struct {
	Path    string
	Snippet string
}

func (indexer *Indexer) Query(query string) (results []*QueryResult, err error) {
	ctx := context.Background()
	store := sqlcgen.New(indexer.dbConn)
	rows, err := store.Query(ctx, &sqlcgen.QueryParams{
		Query: query,
	})
	// ToDo: Should return error for syntax error?
	if err != nil {
		return nil, nil
	}
	return lo.Map(rows, func(item *sqlcgen.QueryRow, _ int) *QueryResult {
		return &QueryResult{
			Path:    item.File.Path,
			Snippet: item.Snippet,
		}
	}), nil
}

func (indexer *Indexer) InitialScanFinished() bool {
	return lo.EveryBy(indexer.dirIndexers, func(dirIndexer *DirIndexer) bool {
		return dirIndexer.InitialScanFinished()
	})
}
