package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/knaka/biblioseeq/web"
	"github.com/webui-dev/go-webui/v2"
	neturl "net/url"
	"strings"

	. "github.com/knaka/go-utils"
)

//go:embed index.html
var webUiTop string

func openWindowAndWait(host string, port int) (err error) {
	defer Catch(&err)
	// Create a new window.
	win := webui.NewWindow()
	defer win.Destroy()
	webui.Bind(win, "getAccessInfo", func(_ webui.Event) (ret struct {
		URL string `json:"url"`
	}) {
		url := neturl.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%d", Ternary(host != "", host, "localhost"), port),
			Path:   "/login",
			RawQuery: strings.Join([]string{
				fmt.Sprintf("password=%s", web.LocalPassword),
				"path=/ap/",
			}, "&"),
		}
		ret.URL = url.String()
		return
	})
	preferredBrowserStr := "AnyBrowser"
	//preferredBrowserStr := "Chrome"
	//preferredBrowserStr := "Chromium"
	preferredBrowser := V(StrToBrowser(preferredBrowserStr))
	// An empty `name` and `path` means the default user profile.
	// Needs to be called before `webui_show()`.
	win.SetProfile("", "")
	V0(win.ShowBrowser(webUiTop, preferredBrowser))
	webui.Wait()
	return
}

func main() {
	server := V(web.NewServer(context.Background(), "", 0))
	go (func() { V0(server.ListenAndServe()) })()
	V0(openWindowAndWait(V2(web.ParseServerAddr(server.Addr))))
}
