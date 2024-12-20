package main

import (
	"embed"
	"image_proce/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 创建应用管理器
	manager := app.NewAppManager()

	err := wails.Run(&options.App{
		Title:  "图片处理工具",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        manager.Startup,
		Bind: []interface{}{
			manager.ImageApp,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
