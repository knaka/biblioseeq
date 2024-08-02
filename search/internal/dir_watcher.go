package internal

import (
	"bufio"
	"context"
	. "github.com/knaka/go-utils"
	"github.com/rjeczalik/notify"
	"path/filepath"
	"strings"
)

//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
type DirWatcher struct {
	dirPath   string
	chStarted chan any
	cancel    context.CancelFunc
	onEvent   func(eventInfo notify.EventInfo)
}

func (dw *DirWatcher) DirPath() string {
	return dw.dirPath
}

func NewDirWatcher(dirPath string, onEvent func(info notify.EventInfo)) *DirWatcher {
	return &DirWatcher{
		dirPath:   dirPath,
		chStarted: make(chan any, 2),
		cancel:    nil,
		onEvent:   onEvent,
	}
}

func (dw *DirWatcher) WatchContinuously(ctx context.Context) {
	ctx, dw.cancel = context.WithCancel(ctx)
	defer dw.cancel()
	chEventInfo := make(chan notify.EventInfo, 10)
	defer close(chEventInfo)
	V0(notify.Watch(filepath.Join(dw.dirPath, "..."), chEventInfo,
		notify.Write, notify.Remove, notify.Rename))
	defer notify.Stop(chEventInfo)
	dw.chStarted <- "started"
outer:
	for {
		select {
		case eventInfo := <-chEventInfo:
			if dw.onEvent != nil {
				dw.onEvent(eventInfo)
			}
		case <-ctx.Done():
			break outer
		}
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

func (dw *DirWatcher) WaitForWatchingStarted(ctx context.Context) {
	select {
	case started := <-dw.chStarted:
		dw.chStarted <- started
	case <-ctx.Done():
	}
}

func (dw *DirWatcher) watching() bool {
	return len(dw.chStarted) > 0
}

func (dw *DirWatcher) stop() {
	if dw.cancel != nil {
		dw.cancel()
	}
}
