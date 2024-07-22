package main

import (
	"context"
	"github.com/knaka/biblioseeq/web"
	"log"
)

func main() {
	ctx := context.Background()
	err := web.ListenAndServe(
		ctx,
		":3000",
	)
	if err != nil {
		log.Panicf("panic 9b807b2 (%v)", err)
	}
}
