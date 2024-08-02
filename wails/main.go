package main

// 2238e24

import (
	"context"
	"embed"
	"github.com/knaka/biblioseeq/conf"
	"github.com/knaka/biblioseeq/fts"
	"github.com/knaka/biblioseeq/web"
	. "github.com/knaka/go-utils"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"log"
)

//go:embed frontend/src
var assets embed.FS

//go:embed build/appicon.png
var macIcon []byte

// e02a8bd

func openWindowAndWait(host string, port int) (err error) {
	defer Catch(&err)
	// Create an instance of the app structure
	app := NewApp(host, port)
	// Create application with options
	V0(wails.Run(&options.App{
		Title:             "BiblioSeeQ",
		Width:             1024,
		Height:            768,
		MinWidth:          1024,
		MinHeight:         768,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Menu:             nil,
		Logger:           nil,
		LogLevel:         logger.DEBUG,
		OnStartup:        app.startup,
		OnDomReady:       app.domReady,
		OnBeforeClose:    app.beforeClose,
		OnShutdown:       app.shutdown,
		WindowStartState: options.Normal,
		Bind: []any{
			app,
		},
		// Windows platform specific options.
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			// DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "",
		},
		// Mac platform specific options.
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "BiblioSeeQ",
				Message: "",
				Icon:    macIcon,
			},
		},
		// Linux platform specific options.
		Linux: nil,
	}))
	return
}

func main() {
	//ftslog.SetOutput(os.Stderr)
	host := ""
	port := V(web.GetFreePort())
	config := V(conf.ReadConfig())
	ftsOpts := []fts.Option{
		fts.WithDefaultDBFilePath(),
		fts.MigratesDB(),
	}
	for _, confDir := range config.Directories {
		ftsOpts = append(ftsOpts, fts.WithTargetDirectory(
			confDir.AbsPath,
			confDir.FileExtensions,
		))
	}
	ftsIndexer := fts.NewIndexer(ftsOpts...)
	go (func() {
		ftsIndexer.WatchContinuously(context.Background())
	})()
	ftsIndexer.WaitForInitialScanFinished(context.Background())
	log.Println("c743caa Finished")
	server := V(web.NewServer(host, port, ftsIndexer))
	go func() { V0(server.ListenAndServe()) }()
	V0(openWindowAndWait(V2(web.ParseServerAddr(server.Addr))))
}
