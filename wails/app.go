package main

import (
	"context"
	"fmt"
	"github.com/knaka/biblioseeq/web"
	neturl "net/url"
	"strings"

	. "github.com/knaka/go-utils"
)

// App struct
type App struct {
	ctx  context.Context
	host string
	port int
}

// NewApp creates a new App application struct
func NewApp(host string, port int) *App {
	return &App{
		host: Elvis(host, "localhost"),
		port: port,
	}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

//goland:noinspection GoUnusedExportedFunction
func (a *App) GetAccessInfo() string {
	url := neturl.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", a.host, a.port),
		Path:   "/login",
		RawQuery: strings.Join([]string{
			fmt.Sprintf("password=%s", web.LocalPassword),
			"path=/ap/",
		}, "&"),
	}
	return url.String()
}
