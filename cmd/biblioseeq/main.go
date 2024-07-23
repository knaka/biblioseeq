package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"github.com/knaka/biblioseeq"
	"github.com/knaka/biblioseeq/db"
	"github.com/knaka/biblioseeq/db/sqlcgen"
	"github.com/knaka/biblioseeq/web"
	_ "github.com/mattn/go-sqlite3"
	ui "github.com/webui-dev/go-webui/v2"
	"log"
	"net"
	neturl "net/url"
	"os"
	"path/filepath"
	"strings"

	. "github.com/knaka/go-utils"
)

//go:embed index.html
var webUiTop string

func openWindowAndWait(port int) {
	// Create a new window.
	w := ui.NewWindow()
	defer w.Destroy()

	ui.Bind(w, "getAccessInfo", func(_ ui.Event) (ret struct {
		URL string `json:"url"`
	}) {
		url := neturl.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("localhost:%d", port),
			Path:   "/login",
			RawQuery: strings.Join([]string{
				fmt.Sprintf("password=%s", web.LocalPassword),
				"path=/ap/",
			}, "&"),
		}
		url.Path = "/login"
		ret.URL = url.String()
		return
	})

	preferredBrowserStr := "AnyBrowser"
	//preferredBrowserStr := "Chrome"
	//preferredBrowserStr := "Chromium"
	preferredBrowser := V(biblioseeq.StrToBrowser(preferredBrowserStr))
	// An empty `name` and `path` means the default user profile.
	// Needs to be called before `webui_show()`.
	w.SetProfile("", "")
	V0(w.ShowBrowser(webUiTop, preferredBrowser))

	ui.Wait()
}

func GetFreePort() (port int, err error) {
	var addr *net.TCPAddr
	if addr, err = net.ResolveTCPAddr("tcp", "localhost:0"); err != nil {
		return
	}
	var listener *net.TCPListener
	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		return
	}
	defer (func() { _ = listener.Close() })()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func main() {
	if false {
		dbFile := V(filepath.Abs(os.Args[1]))
		db.Migrate(dbFile)
		ctx := context.Background()
		dbConn := V(sql.Open("sqlite3", dbFile))
		store := sqlcgen.New(dbConn)
		version := V(store.GetVersion(ctx))
		log.Printf("sqlite version: %s", version)
	}
	port := V(GetFreePort())
	go func() {
		ctx := context.Background()
		_ = web.ListenAndServe(
			ctx,
			fmt.Sprintf(":%d", port),
		)
	}()
	openWindowAndWait(port)
}
