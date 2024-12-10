package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xackery/shinsparxx/config"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
	"github.com/xackery/wlk/win"
)

var (
	cfg                    *config.CritSprinklerConfiguration
	settingsWnd            *walk.MainWindow
	wndPlayerDefaultOption *walk.RadioButton
	wndPlayerExpOption     *walk.RadioButton
	wndPlayerAAOption      *walk.RadioButton
)

const (
	wndPlayerExpPath     = "player/exp/EQUI_PlayerWindow.xml"
	wndPlayerAAExpPath   = "player/aaexp/EQUI_PlayerWindow.xml"
	wndPlayerDefaultPath = "player/EQUI_PlayerWindow.xml"
)

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var err error
	exePath := os.Args[0]

	wd := filepath.Dir(exePath)

	cfg, err = config.LoadCritSprinklerConfig(wd + "/critsprinkler.ini")
	if err != nil {
		fmt.Printf("load crit sprinkler config: %v\n", err)
	}

	cmw := cpl.MainWindow{
		Title:    "Shin Sparxx Config",
		Name:     "settings",
		AssignTo: &settingsWnd,
		Size:     cpl.Size{Width: cfg.SettingsW, Height: cfg.SettingsH},
		Layout:   cpl.VBox{},

		Children: []cpl.Widget{
			cpl.GroupBox{
				Title:     "Player Window",
				Layout:    cpl.VBox{},
				Alignment: cpl.AlignHVDefault,
				Children: []cpl.Widget{
					cpl.RadioButton{Text: "Default", AssignTo: &wndPlayerDefaultOption, Enabled: true},
					cpl.RadioButton{Text: "Show Exp", AssignTo: &wndPlayerExpOption},
					cpl.RadioButton{Text: "Show AA Exp", AssignTo: &wndPlayerAAOption},
				},
			},
			cpl.PushButton{
				Text:    "Save",
				MaxSize: cpl.Size{Width: 45},
				OnClicked: func() {
					err := updateSave()
					if err != nil {
						fmt.Printf("Error: %v\n", err)
					}
					os.Exit(0)

				},
			},
		},
		OnSizeChanged: func() {

		},
		OnMouseMove: func(x, y int, button walk.MouseButton) {
		},

		OnBoundsChanged: func() {
		},

		Visible: cfg.LogPath == "",
	}

	err = cmw.Create()
	if err != nil {
		return fmt.Errorf("create main window: %w", err)
	}

	settingsWnd.Closing().Attach(func(isCancel *bool, reason byte) {
		*isCancel = true
		err := updateSave()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		settingsWnd.SetVisible(false)
	})

	settingsWnd.SetWidth(cfg.SettingsW)
	settingsWnd.SetHeight(cfg.SettingsH)
	settingsWnd.SetX(cfg.SettingsX)
	settingsWnd.SetY(cfg.SettingsY)

	code := settingsWnd.Run()
	if code != 0 {
		fmt.Printf("mainWalk Error: %v\n", code)
	}
	if cfg.LogPath == "" {

		win.SetForegroundWindow(settingsWnd.Handle())
		win.SetActiveWindow(settingsWnd.Handle())

	}

	return nil
}

func updateSave() error {
	if cfg == nil {
		return fmt.Errorf("config not loaded")
	}
	//cfg.SettingsY = game.settingsWindowY

	err := cfg.Save()
	if err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	return nil
}
