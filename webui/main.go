package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/knaka/biblioseeq/web"
	ui "github.com/webui-dev/go-webui/v2"
	neturl "net/url"
	"strings"

	. "github.com/knaka/go-utils"
)

//go:embed index.html
var webUiTop string

func openWindowAndWait(host string, port int) {
	// Create a new window.
	w := ui.NewWindow()
	defer w.Destroy()

	ui.Bind(w, "getAccessInfo", func(_ ui.Event) (ret struct {
		URL string `json:"url"`
	}) {
		url := neturl.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%d", host, port),
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
	preferredBrowser := V(StrToBrowser(preferredBrowserStr))
	// An empty `name` and `path` means the default user profile.
	// Needs to be called before `webui_show()`.
	w.SetProfile("", "")
	V0(w.ShowBrowser(webUiTop, preferredBrowser))

	ui.Wait()
}

func main() {
	ctx := context.Background()
	server, host, port := V3(web.NewServer(ctx, "", 0))
	go func() {
		_ = server.ListenAndServe()
	}()
	openWindowAndWait(host, port)
}
