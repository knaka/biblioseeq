package rpc

import (
	"context"
	"github.com/bufbuild/connect-go"
	v1 "github.com/knaka/biblioseeq/pbgen/v1"
	"github.com/knaka/biblioseeq/pbgen/v1/v1connect"
	weblib "github.com/knaka/biblioseeq/web/lib"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	. "github.com/knaka/go-utils"
)

// MainServiceHandlerImpl is a concrete implementation of v1connect.MainServiceHandler.
type MainServiceHandlerImpl struct{}

// V MainServiceHandlerImpl implements v1connect.MainServiceHandler
var _ v1connect.MainServiceHandler = &MainServiceHandlerImpl{}

func newResponseWithMsg[T any](_ *connect.Response[T]) *connect.Response[T] {
	var msg T
	return connect.NewResponse(&msg)
}

func (s *MainServiceHandlerImpl) VersionInfo(ctx context.Context, req *connect.Request[v1.VersionInfoRequest]) (resp *connect.Response[v1.VersionInfoResponse], err error) {
	resp = newResponseWithMsg(resp)

	resp.Msg.VersionInfo = &v1.VersionInfo{
		Version: "0.1.2",
	}

	return
}

func (s *MainServiceHandlerImpl) Status(ctx context.Context, req *connect.Request[v1.StatusRequest]) (resp *connect.Response[v1.StatusResponse], err error) {
	resp = newResponseWithMsg(resp)
	ctxValue := V(weblib.GetCtxValue(ctx))

	resp.Msg.InitialScanFinished = ctxValue.FtsIndexer.InitialScanFinished()

	return
}

type User struct {
}

var defaultUser User

func DefaultUser() *User {
	return &defaultUser
}

// CurrentTime returns the current time.
func (s *MainServiceHandlerImpl) CurrentTime(ctx context.Context, req *connect.Request[v1.CurrentTimeRequest]) (resp *connect.Response[v1.CurrentTimeResponse], err error) {
	resp = newResponseWithMsg(resp)

	resp.Msg.Timestamp = timestamppb.New(time.Now())

	return
}
