package fts

import (
	"github.com/rjeczalik/notify"
	"log"

	. "github.com/knaka/go-utils"
)

type Engine struct {
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Watch() {
	ch := make(chan notify.EventInfo, 10)
	V0(notify.Watch("/path/to/watch", ch, notify.Write, notify.Remove, notify.Rename))
	defer notify.Stop(ch)
	for {
		eventInfo := <-ch
		path := eventInfo.Path()
		event := eventInfo.Event()
		sys := eventInfo.Sys()
		log.Println("Gotten event info:", event, path, sys)
	}
}
