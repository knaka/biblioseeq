package fts

import (
	. "github.com/knaka/go-utils"
	"github.com/rjeczalik/notify"
	"github.com/samber/lo"
	"log"
	"path/filepath"
)

type Directory struct {
	Path           string
	FileExtensions []string
}

type Index struct {
	Directories []*Directory
}

type Option func(*Index)

func WithDirectory(path string, fileExtensions []string) Option {
	return func(e *Index) {
		e.Directories = append(e.Directories, &Directory{
			Path:           path,
			FileExtensions: fileExtensions,
		})
	}
}

func NewIndex(opts ...Option) (ret *Index) {
	ret = &Index{}
	for _, opt := range opts {
		opt(ret)
	}
	return
}

type data struct {
	index     int
	eventInfo notify.EventInfo
}

func (index *Index) Watch() {
	ch := make(chan data, 100)
	var subChannels []chan notify.EventInfo
	for i, directory := range index.Directories {
		notifyCh := make(chan notify.EventInfo, 10)
		V0(notify.Watch(filepath.Join(directory.Path, "..."), notifyCh,
			notify.Write, notify.Remove, notify.Rename))
		go func() {
			for {
				ch <- data{i, <-notifyCh}
			}
		}()
		subChannels = append(subChannels, notifyCh)
	}
	defer func() {
		for _, subCh := range subChannels {
			notify.Stop(subCh)
		}
	}()
	for {
		d := <-ch
		eventInfo := d.eventInfo
		path := eventInfo.Path()
		event := eventInfo.Event()
		sys := eventInfo.Sys()
		log.Println("Got event:", event, path, sys)
		if event&notify.Write != 0 {
			// ファイルへの書き込み。単体更新
			index.UpdateFile(path)
		}
	}
}

func (index *Index) getDirConf(path string) (directory *Directory) {
	for _, dir := range index.Directories {
		if _, errSub := filepath.Rel(dir.Path, path); errSub == nil {
			return dir
		}
	}
	return
}

func (index *Index) UpdateFile(filePath string) {
	dirConf := index.getDirConf(filePath)
	if lo.NoneBy(dirConf.FileExtensions, func(ext string) bool {
		return filepath.Ext(filePath) == ext
	}) {
		return
	}

}
