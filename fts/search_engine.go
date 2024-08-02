package fts

import (
	"context"
	"fmt"
	"github.com/knaka/biblioseeq/conf"
	"github.com/knaka/biblioseeq/db"
	"github.com/knaka/biblioseeq/log"
	"os"
	"path/filepath"
	"sync"

	. "github.com/knaka/go-utils"
)

type SearchEngine struct {
	confFilePath    string
	idxMgr          *IdxMgr
	dirWatchers     []*DirWatcher
	shouldMigrateDB bool
	chSynchronized  chan any
}

type Option func(*SearchEngine)

//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func ShouldMigratesDB(f bool) Option {
	return func(se *SearchEngine) {
		se.shouldMigrateDB = f
	}
}

func NewSearchEngine(opts ...Option) (se *SearchEngine) {
	se = &SearchEngine{
		chSynchronized: make(chan any, 2),
	}
	for _, opt := range opts {
		opt(se)
	}
	dbFilePath := filepath.Join(
		V(os.UserCacheDir()),
		fmt.Sprintf("%s.sqlite3", appName),
	)
	if se.shouldMigrateDB {
		db.Migrate(dbFilePath)
	}
	se.idxMgr = newIdxMgr(dbFilePath)
	confFile := V(conf.ReadConfig(se.confFilePath))
	for _, confDir := range confFile.Directories {
		se.dirWatchers = append(se.dirWatchers, newDirWatcher(
			confDir.Path,
			se.idxMgr.onNotifyEvent,
		))
		se.idxMgr.addTarget(confDir.Path, confDir.FileExtensions)
	}
	return
}

func (se *SearchEngine) Serve(ctx context.Context) {
	wg := sync.WaitGroup{}
	for _, dirWatcher := range se.dirWatchers {
		wg.Add(1)
		go (func() {
			dirWatcher.WatchContinuously(ctx)
			wg.Done()
		})()
		dirWatcher.waitForWatchingStarted(ctx)
		log.Println("Started watching for directory", dirWatcher.dirPath)
	}
	se.idxMgr.synchronizeIndexesToFiles()
	se.chSynchronized <- "yes"
	log.Println("Finished initial file index update")
	wg.Wait()
}

func (se *SearchEngine) Query(query string) ([]*QueryResult, error) {
	return se.idxMgr.Query(query)
}

// InitialScanFinished returns true if all directories have been scanned at least once.
func (se *SearchEngine) InitialScanFinished() bool {
	return len(se.chSynchronized) > 0
}

// WaitForInitialScanFinished blocks until all directories have been scanned at least once.
func (se *SearchEngine) WaitForInitialScanFinished(ctx context.Context) {
	select {
	case <-ctx.Done():
	case <-se.chSynchronized:
		se.chSynchronized <- "yeah"
	}
}
