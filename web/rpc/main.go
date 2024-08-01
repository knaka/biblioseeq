package rpc

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/knaka/biblioseeq/fts"
	v1 "github.com/knaka/biblioseeq/pbgen/v1"
	"github.com/knaka/biblioseeq/pbgen/v1/v1connect"
	weblib "github.com/knaka/biblioseeq/web/lib"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"os/exec"
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

// Query
func (s *MainServiceHandlerImpl) Query(ctx context.Context, req *connect.Request[v1.QueryRequest]) (resp *connect.Response[v1.QueryResponse], err error) {
	resp = newResponseWithMsg(resp)
	ctxValue := V(weblib.GetCtxValue(ctx))

	queryRaw := req.Msg.Query
	query := fts.SeparateJapanese(queryRaw)
	results := V(ctxValue.FtsIndexer.Query(query))
	for _, result := range results {
		resp.Msg.Results = append(resp.Msg.Results, &v1.QueryResult{
			Path:       result.Path,
			Title:      result.Title,
			Tags:       result.Tags,
			DirPath:    result.DirPath,
			Snippet:    result.Snippet,
			ModifiedAt: timestamppb.New(result.ModifiedAt),
		})
	}

	return
}

func (s *MainServiceHandlerImpl) Content(ctx context.Context, req *connect.Request[v1.ContentRequest]) (resp *connect.Response[v1.ContentResponse], err error) {
	resp = newResponseWithMsg(resp)

	content := V(os.ReadFile(req.Msg.Path))
	resp.Msg.Content = string(content)

	return
}

// LaunchPath
func (s *MainServiceHandlerImpl) LaunchPath(ctx context.Context, req *connect.Request[v1.LaunchPathRequest]) (resp *connect.Response[v1.LaunchPathResponse], err error) {
	resp = newResponseWithMsg(resp)

	cmd := exec.Command("open", req.Msg.Path)
	V0(cmd.Run())

	return
}
