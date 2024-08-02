package search

import (
	"context"
	"fmt"
	"github.com/knaka/biblioseeq/config"
	"github.com/knaka/biblioseeq/db"
	"github.com/knaka/biblioseeq/log"
	"github.com/knaka/biblioseeq/search/internal"
	"os"
	"path/filepath"
	"sync"

	. "github.com/knaka/go-utils"
)

type Engine struct {
	confFilePath    string
	idxMgr          *internal.IdxMgr
	dirWatchers     []*internal.DirWatcher
	shouldMigrateDB bool
	chSynchronized  chan any
}

type Option func(*Engine)

func ShouldMigratesDB(f bool) Option {
	return func(se *Engine) {
		se.shouldMigrateDB = f
	}
}

func NewEngine(opts ...Option) (se *Engine) {
	se = &Engine{
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
	se.idxMgr = internal.NewIdxMgr(dbFilePath)
	conf_ := V(config.Read(se.confFilePath))
	for _, confDir := range conf_.Directories {
		se.dirWatchers = append(se.dirWatchers, internal.NewDirWatcher(
			confDir.Path,
			se.idxMgr.OnNotifyEvent,
		))
		se.idxMgr.AddTarget(confDir.Path, confDir.FileExtensions)
	}
	return
}

func (se *Engine) Serve(ctx context.Context) {
	wg := sync.WaitGroup{}
	for _, dirWatcher := range se.dirWatchers {
		wg.Add(1)
		go (func() {
			dirWatcher.WatchContinuously(ctx)
			wg.Done()
		})()
		dirWatcher.WaitForWatchingStarted(ctx)
		log.Println("Started watching for directory", dirWatcher.DirPath())
	}
	se.idxMgr.SynchronizeIndexesToFiles()
	se.chSynchronized <- "yes"
	log.Println("Finished initial file index update")
	wg.Wait()
}

func (se *Engine) Query(query string) ([]*internal.QueryResult, error) {
	return se.idxMgr.Query(query)
}

// InitialScanFinished returns true if all directories have been scanned at least once.
func (se *Engine) InitialScanFinished() bool {
	return len(se.chSynchronized) > 0
}

// WaitForInitialScanFinished blocks until all directories have been scanned at least once.
func (se *Engine) WaitForInitialScanFinished(ctx context.Context) {
	select {
	case <-ctx.Done():
	case <-se.chSynchronized:
		se.chSynchronized <- "yeah"
	}
}
