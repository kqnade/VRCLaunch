package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func (a *App) layoutProfileEditor(gtx layout.Context) layout.Dimensions {
	return panelInsets().Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := "New Profile"
				if a.state.ProfileForm.ID != "" {
					title = "Edit Profile"
				}
				h := material.H6(a.th, title)
				h.Color = a.th.Palette.Fg
				return layout.Inset{Bottom: unit.Dp(12)}.Layout(gtx, h.Layout)
			}),
			layout.Rigid(a.editorRow("Name", &a.editNameEd)),
			layout.Rigid(a.editorRow("Profile Index (--profile=N)", &a.editIndexEd)),
			layout.Rigid(a.editorRow("FPS (0=unset)", &a.editFPSEd)),
			layout.Rigid(a.editorRow("Screen Width", &a.editWidthEd)),
			layout.Rigid(a.editorRow("Screen Height", &a.editHeightEd)),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					sw := material.CheckBox(a.th, &a.editFullscreenBl, "Fullscreen (Desktop mode)")
					return sw.Layout(gtx)
				})
			}),
			layout.Rigid(a.editorRow("Custom Args (space-separated)", &a.editCustomEd)),
			layout.Rigid(a.layoutStatus),
			layout.Rigid(a.layoutEditorButtons),
		)
	})
}

func (a *App) editorRow(label string, ed *widget.Editor) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					l := material.Caption(a.th, label)
					l.Color = mutedColor(a.th)
					return l.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: unit.Dp(2)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						e := material.Editor(a.th, ed, "")
						e.Color = a.th.Palette.Fg
						return e.Layout(gtx)
					})
				}),
			)
		})
	}
}

func (a *App) layoutEditorButtons(gtx layout.Context) layout.Dimensions {
	return layout.Inset{Top: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.th, &a.editCancelBtn, "Cancel")
				return btn.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if a.state.ProfileForm.ID == "" {
					return layout.Dimensions{}
				}
				btn := material.Button(a.th, &a.editDeleteBtn, "Delete")
				return btn.Layout(gtx)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: gtx.Constraints.Min}
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.th, &a.editSaveBtn, "Save")
				return btn.Layout(gtx)
			}),
		)
	})
}
