package fts

import (
	"context"
	"database/sql"
	"github.com/knaka/biblioseeq/db/sqlcgen"
	"github.com/samber/lo"
	"path/filepath"
	"sync"
)

type Indexer struct {
	dbConn      *sql.DB
	dirIndexers []*DirIndexer
}

type Option func(*Indexer)

func WithDirectory(path string, fileExtensions []string) Option {
	return func(e *Indexer) {
		e.dirIndexers = append(e.dirIndexers, NewDirIndexer(e.dbConn, path, fileExtensions))
	}
}

func NewIndexer(dbConn *sql.DB, opts ...Option) (ret *Indexer) {
	ret = &Indexer{
		dbConn: dbConn,
	}
	for _, opt := range opts {
		opt(ret)
	}
	return
}

func (indexer *Indexer) WatchContinuously() {
	wg := sync.WaitGroup{}
	for _, dirIndexer := range indexer.dirIndexers {
		wg.Add(1)
		go (func() {
			dirIndexer.WatchContinuously()
			wg.Done()
		})()
		dirIndexer.WaitForWatchingStarted()
		go (func() {
			dirIndexer.IndexAll()
		})()
	}
	wg.Wait()
}

func (indexer *Indexer) WaitForInitialScanFinished() {
	for _, dirIndexer := range indexer.dirIndexers {
		dirIndexer.WaitForInitialScanFinished()
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
	if err != nil {
		return
	}
	return lo.Map(rows, func(item *sqlcgen.QueryRow, _ int) *QueryResult {
		return &QueryResult{
			Path:    item.File.Path,
			Snippet: item.Snippet,
		}
	}), nil
}
