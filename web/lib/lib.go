package lib

import (
	"context"
	"fmt"
	"github.com/knaka/biblioseeq/fts"
	"github.com/knaka/go-utils"
)

type CtxKey struct{}

type CtxValue struct {
	SearchEngine *fts.SearchEngine
}

func GetCtxValue(ctx context.Context) (ctxValue *CtxValue, err error) {
	defer utils.Catch(&err)
	var ok bool
	if ctxValue, ok = ctx.Value(CtxKey{}).(*CtxValue); !ok {
		return nil, fmt.Errorf("invalid context")
	}
	return
}
