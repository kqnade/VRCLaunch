package ui

import (
	"image/color"
	"log"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/kqnade/VRCLaunch/internal/config"
	"github.com/kqnade/VRCLaunch/internal/launcher"
	"github.com/kqnade/VRCLaunch/internal/uistate"
)

// Persister saves the config when the user-visible state changes.
type Persister func(*config.Config) error

// Launcher is the function that spawns VRChat. Injected for testability.
type Launcher func(launchPath string, p config.Profile, mode launcher.Mode) error

type App struct {
	state   *uistate.State
	th      *material.Theme
	persist Persister
	launch  Launcher

	// Top-bar / main view
	settingsBtn      widget.Clickable
	addProfileBtn    widget.Clickable
	launchVRBtn      widget.Clickable
	launchDesktopBtn widget.Clickable
	profileList      layout.List
	profileWidgets   map[string]*profileRowWidgets

	// Profile editor
	editNameEd       widget.Editor
	editIndexEd      widget.Editor
	editFPSEd        widget.Editor
	editWidthEd      widget.Editor
	editHeightEd     widget.Editor
	editCustomEd     widget.Editor
	editFullscreenBl widget.Bool
	editSaveBtn      widget.Clickable
	editCancelBtn    widget.Clickable
	editDeleteBtn    widget.Clickable

	// Settings
	settingsLaunchPathEd widget.Editor
	settingsSaveBtn      widget.Clickable
	settingsCancelBtn    widget.Clickable

	currentView uistate.View
	viewLoaded  bool
}

type profileRowWidgets struct {
	selectBtn widget.Clickable
	editBtn   widget.Clickable
}

func NewApp(state *uistate.State, th *material.Theme, persist Persister, launch Launcher) *App {
	a := &App{
		state:          state,
		th:             th,
		persist:        persist,
		launch:         launch,
		profileList:    layout.List{Axis: layout.Vertical},
		profileWidgets: make(map[string]*profileRowWidgets),
	}
	a.editNameEd.SingleLine = true
	a.editIndexEd.SingleLine = true
	a.editFPSEd.SingleLine = true
	a.editWidthEd.SingleLine = true
	a.editHeightEd.SingleLine = true
	a.editCustomEd.SingleLine = true
	a.settingsLaunchPathEd.SingleLine = true
	return a
}

func (a *App) Run(w *app.Window) error {
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			a.handleEvents(gtx)
			a.layoutBackground(gtx)
			a.layoutRoot(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (a *App) layoutBackground(gtx layout.Context) {
	bg := a.th.Palette.Bg
	defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: bg}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func (a *App) layoutRoot(gtx layout.Context) layout.Dimensions {
	if a.state.View != a.currentView {
		a.currentView = a.state.View
		a.viewLoaded = false
	}

	switch a.state.View {
	case uistate.ViewMain:
		return a.layoutMain(gtx)
	case uistate.ViewEditProfile:
		if !a.viewLoaded {
			a.bindProfileForm()
			a.viewLoaded = true
		}
		return a.layoutProfileEditor(gtx)
	case uistate.ViewSettings:
		if !a.viewLoaded {
			a.bindSettingsForm()
			a.viewLoaded = true
		}
		return a.layoutSettings(gtx)
	}
	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func (a *App) bindProfileForm() {
	a.editNameEd.SetText(a.state.ProfileForm.Name)
	a.editIndexEd.SetText(a.state.ProfileForm.Index)
	a.editFPSEd.SetText(a.state.ProfileForm.FPS)
	a.editWidthEd.SetText(a.state.ProfileForm.ScreenWidth)
	a.editHeightEd.SetText(a.state.ProfileForm.ScreenHeight)
	a.editCustomEd.SetText(a.state.ProfileForm.CustomArgs)
	a.editFullscreenBl.Value = a.state.ProfileForm.ScreenFullscreen
}

func (a *App) bindSettingsForm() {
	a.settingsLaunchPathEd.SetText(a.state.SettingsForm.LaunchPath)
}

func (a *App) handleEvents(gtx layout.Context) {
	switch a.state.View {
	case uistate.ViewMain:
		a.handleMainEvents(gtx)
	case uistate.ViewEditProfile:
		a.handleProfileEditorEvents(gtx)
	case uistate.ViewSettings:
		a.handleSettingsEvents(gtx)
	}
}

func clicked(c *widget.Clickable, gtx layout.Context) bool {
	hit := false
	for c.Clicked(gtx) {
		hit = true
	}
	return hit
}

func (a *App) handleMainEvents(gtx layout.Context) {
	if clicked(&a.settingsBtn, gtx) {
		a.state.GotoSettings()
	}
	if clicked(&a.addProfileBtn, gtx) {
		a.state.GotoNewProfile()
	}
	for id, w := range a.profileWidgets {
		if clicked(&w.selectBtn, gtx) {
			a.state.SelectProfile(id)
			a.savePersist()
		}
		if clicked(&w.editBtn, gtx) {
			if err := a.state.GotoEditProfile(id); err != nil {
				a.state.Status = err.Error()
			}
		}
	}
	if clicked(&a.launchVRBtn, gtx) {
		a.launchSelected(launcher.ModeVR)
	}
	if clicked(&a.launchDesktopBtn, gtx) {
		a.launchSelected(launcher.ModeDesktop)
	}
}

func (a *App) handleProfileEditorEvents(gtx layout.Context) {
	a.state.ProfileForm.Name = a.editNameEd.Text()
	a.state.ProfileForm.Index = a.editIndexEd.Text()
	a.state.ProfileForm.FPS = a.editFPSEd.Text()
	a.state.ProfileForm.ScreenWidth = a.editWidthEd.Text()
	a.state.ProfileForm.ScreenHeight = a.editHeightEd.Text()
	a.state.ProfileForm.CustomArgs = a.editCustomEd.Text()
	a.state.ProfileForm.ScreenFullscreen = a.editFullscreenBl.Value

	if clicked(&a.editCancelBtn, gtx) {
		a.state.GotoMain()
	}
	if clicked(&a.editSaveBtn, gtx) {
		if err := a.state.SaveProfileFromForm(); err != nil {
			a.state.Status = err.Error()
		} else {
			a.state.Status = ""
			a.savePersist()
		}
	}
	if clicked(&a.editDeleteBtn, gtx) && a.state.ProfileForm.ID != "" {
		if err := a.state.DeleteProfile(a.state.ProfileForm.ID); err != nil {
			a.state.Status = err.Error()
		} else {
			a.state.GotoMain()
			a.savePersist()
		}
	}
}

func (a *App) handleSettingsEvents(gtx layout.Context) {
	a.state.SettingsForm.LaunchPath = a.settingsLaunchPathEd.Text()

	if clicked(&a.settingsCancelBtn, gtx) {
		a.state.GotoMain()
	}
	if clicked(&a.settingsSaveBtn, gtx) {
		a.state.SaveSettingsFromForm()
		a.savePersist()
	}
}

func (a *App) savePersist() {
	if a.persist == nil {
		return
	}
	if err := a.persist(a.state.Config); err != nil {
		log.Printf("persist config: %v", err)
		a.state.Status = "保存に失敗しました: " + err.Error()
	}
}

func (a *App) launchSelected(mode launcher.Mode) {
	if a.launch == nil {
		a.state.Status = "launcher not configured"
		return
	}
	if a.state.SelectedID == "" {
		a.state.Status = "プロファイルを選択してください"
		return
	}
	p := a.state.Config.FindByID(a.state.SelectedID)
	if p == nil {
		a.state.Status = "プロファイルが見つかりません"
		return
	}
	if err := a.launch(a.state.Config.LaunchPath, *p, mode); err != nil {
		a.state.Status = "起動失敗: " + err.Error()
		return
	}
	modeName := "Desktop"
	if mode == launcher.ModeVR {
		modeName = "VR"
	}
	a.state.Status = "起動: " + p.Name + " (" + modeName + ")"
}

// --- helper styling ---

func dividerColor(th *material.Theme) color.NRGBA {
	c := th.Palette.Fg
	c.A = 0x20
	return c
}

func panelInsets() layout.Inset {
	return layout.UniformInset(unit.Dp(16))
}
