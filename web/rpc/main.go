package rpc

import (
	"context"
	"github.com/bufbuild/connect-go"
	v1 "github.com/knaka/biblioseeq/pbgen/v1"
	"github.com/knaka/biblioseeq/pbgen/v1/v1connect"
)

// MainServiceHandlerImpl is a concrete implementation of v1connect.MainServiceHandler.
type MainServiceHandlerImpl struct{}

// V MainServiceHandlerImpl implements v1connect.MainServiceHandler
var _ v1connect.MainServiceHandler = &MainServiceHandlerImpl{}

func newResponseWithMsg[T any](_ *connect.Response[T]) *connect.Response[T] {
	var msg T
	return connect.NewResponse(&msg)
}

func (s *MainServiceHandlerImpl) GetVersionInfo(ctx context.Context, req *connect.Request[v1.GetVersionInfoRequest]) (resp *connect.Response[v1.GetVersionInfoResponse], err error) {
	resp = newResponseWithMsg(resp)

	resp.Msg.VersionInfo = &v1.VersionInfo{
		Version: "0.1.2",
	}

	return
}
