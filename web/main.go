package web

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/knaka/biblioseeq"
	"github.com/knaka/biblioseeq/pbgen/v1/v1connect"
	"github.com/knaka/biblioseeq/web/rpc"
	. "github.com/knaka/go-utils"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"io/fs"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func GetFreePort() (port int, err error) {
	var addr *net.TCPAddr
	if addr, err = net.ResolveTCPAddr("tcp", "127.0.0.1:0"); err != nil {
		return
	}
	var listener *net.TCPListener
	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		return
	}
	defer (func() { _ = listener.Close() })()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

type ctxKey struct{}

type CtxValue struct{}

func GetCtxValue(ctx context.Context) (ctxValue *CtxValue, err error) {
	defer Catch(&err)
	var ok bool
	if ctxValue, ok = ctx.Value(ctxKey{}).(*CtxValue); !ok {
		return nil, fmt.Errorf("invalid context")
	}
	return
}

func ParseServerAddr(addr string) (host string, port int, err error) {
	defer Catch(&err)
	divs := strings.SplitN(addr, ":", 2)
	host = divs[0]
	port = V(strconv.Atoi(divs[1]))
	return
}

func NewServer(
	ctx context.Context,
	host string,
	port int,
) (
	server *http.Server,
	err error,
) {
	defer Catch(&err)
	if port == 0 {
		port = V(GetFreePort())
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	ctxNew := context.WithValue(ctx, ctxKey{}, &CtxValue{})
	server = &http.Server{
		Addr:        addr,
		Handler:     GetWrappedRouter(),
		BaseContext: func(_ net.Listener) context.Context { return ctxNew },
	}
	return
}

const TokenName = "BiblioSeeQToken"

var LocalPassword string = uuid.New().String()
var tokenValue string = uuid.New().String()

func AuthRequired(c *gin.Context) {
	token, _ := c.Cookie(TokenName)
	if token != tokenValue {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized 5a07fdc"})
		return
	}
	c.Next()
}

func getApiRouter() http.Handler {
	ginRoute := gin.Default()
	clientStaticFs := V(fs.Sub(biblioseeq.Client, "client/build"))
	httpFsClient := http.FS(clientStaticFs)
	authAp := ginRoute.Group("ap")
	authAp.Use(AuthRequired)
	authAp.StaticFS("", httpFsClient)
	ginRoute.Any("login", func(c *gin.Context) {
		password := c.Query("password")
		if password != LocalPassword {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized 3f72df4"})
			return
		}
		path := c.Query("path")
		c.SetCookie(TokenName, tokenValue, 0, "/", "localhost", false, false)
		c.Redirect(http.StatusFound, path)
	})
	return ginRoute
}

func GetWrappedRouter() http.Handler {
	mux := http.NewServeMux()
	mux.Handle(v1connect.NewMainServiceHandler(&rpc.MainServiceHandlerImpl{}))
	ginRoute := getApiRouter()
	mux.Handle("/ap/", ginRoute)
	mux.Handle("/login", ginRoute)
	mux.Handle("/api/", ginRoute)
	corsHandler := cors.New(cors.Options{
		Debug: false,
		// “Credentials are cookies, authorization headers, or TLS client certificates.” とのことなので、credential を gRPC の request body にしか乗せないのであれば false でも良いか // Access-Control-Allow-Credentials - HTTP | MDN https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials
		AllowCredentials: true,
		AllowedOrigins: []string{
			"*",
			"http://localhost:3000",
			"http://localhost:8000",
		},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{"*"},
	})
	return h2c.NewHandler(corsHandler.Handler(mux), &http2.Server{})
}
