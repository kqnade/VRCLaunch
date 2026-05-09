package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func (a *App) layoutSettings(gtx layout.Context) layout.Dimensions {
	return panelInsets().Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				h := material.H6(a.th, "Settings")
				h.Color = a.th.Palette.Fg
				return layout.Inset{Bottom: unit.Dp(12)}.Layout(gtx, h.Layout)
			}),
			layout.Rigid(a.editorRow("Steam launch.exe path", &a.settingsLaunchPathEd)),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				hint := material.Caption(a.th, "例: F:\\SteamLibrary\\steamapps\\common\\VRChat\\launch.exe")
				hint.Color = mutedColor(a.th)
				return hint.Layout(gtx)
			}),
			layout.Rigid(a.layoutStatus),
			layout.Rigid(a.layoutSettingsButtons),
		)
	})
}

func (a *App) layoutSettingsButtons(gtx layout.Context) layout.Dimensions {
	return layout.Inset{Top: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.th, &a.settingsCancelBtn, "Cancel")
				return btn.Layout(gtx)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: gtx.Constraints.Min}
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.th, &a.settingsSaveBtn, "Save")
				return btn.Layout(gtx)
			}),
		)
	})
}
