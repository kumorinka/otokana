package main

import (
	"embed"
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"otokana/internal/audio"
	"otokana/internal/store"
	"otokana/view"
)

//go:embed assets/fonts/*.ttf
var fontFiles embed.FS

// fonts are the embedded files of each family (SIL OFL, licenses beside
// them in assets/fonts).
var fonts = []struct{ file, family string }{
	{"KleeOne-SemiBold.ttf", view.FontKana},
	{"ZenKakuGothicNew-Regular.ttf", view.FontUI},
	{"ZenKakuGothicNew-Medium.ttf", view.FontUI},
	{"ZenKakuGothicNew-Bold.ttf", view.FontUI},
}

func init() {
	for _, f := range fonts {
		data, err := fontFiles.ReadFile("assets/fonts/" + f.file)
		if err != nil {
			log.Fatal(err)
		}
		if err := ui.RegisterFont(data, f.family); err != nil {
			log.Fatalf("%s: %v", f.file, err)
		}
	}
}

func main() {
	app := view.NewApp()
	var win *mygo.Window
	app.Pick = func(dir bool) (string, error) {
		opts := mygo.OpenDialogOptions{Parent: win, Title: "Choose a voicebank folder", Directory: true}
		if !dir {
			opts = mygo.OpenDialogOptions{Parent: win, Title: "Choose a voicebank zip file",
				Filters: []mygo.FileFilter{{Name: "Zip", Extensions: []string{"zip"}}}}
		}
		paths, err := mygo.Dialog.Open(opts)
		if err != nil || len(paths) == 0 {
			return "", err
		}
		return paths[0], nil
	}
	app.Update = func(fn func()) { win.Update(fn) }
	app.Go = func(fn func()) { go fn() }
	app.Out = audio.Open()
	mygo.App.WhenReady(func() {
		dir, err := mygo.App.Path(mygo.PathUserData)
		if err != nil {
			log.Fatal(err)
		}
		app.Store = &store.Store{Dir: dir}
		// Without its settings or voicebank the app still runs, with the
		// defaults, without sound.
		if err := app.Load(); err != nil {
			log.Printf("data: %v", err)
		}
		win = mygo.NewWindow(mygo.WindowOptions{
			Title:     "OtoKana",
			Width:     1280,
			Height:    800,
			MinWidth:  1024,
			MinHeight: 720,
			// The desk, not white, until the first frame.
			BackgroundColor: "#22201C",
			// Opens where the user left it last time.
			StateKey: "main",
			// The window shows the interface MyGo draws, not a web page.
			Content: ui.View(app.View),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
