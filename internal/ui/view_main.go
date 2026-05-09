package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/kqnade/VRCLaunch/internal/config"
)

func (a *App) layoutMain(gtx layout.Context) layout.Dimensions {
	return panelInsets().Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(a.layoutHeader),
			layout.Flexed(1, a.layoutProfileList),
			layout.Rigid(a.layoutLaunchBar),
			layout.Rigid(a.layoutStatus),
		)
	})
}

func (a *App) layoutHeader(gtx layout.Context) layout.Dimensions {
	return layout.Inset{Bottom: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				h := material.H5(a.th, "VRCLaunch")
				h.Color = a.th.Palette.Fg
				return h.Layout(gtx)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, 0)}
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.th, &a.addProfileBtn, "+ Profile")
				return btn.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(a.th, &a.settingsBtn, "Settings")
					return btn.Layout(gtx)
				})
			}),
		)
	})
}

func (a *App) layoutProfileList(gtx layout.Context) layout.Dimensions {
	if len(a.state.Config.Profiles) == 0 {
		lbl := material.Body1(a.th, "プロファイルがありません。+ Profile から追加してください。")
		lbl.Color = a.th.Palette.Fg
		lbl.Alignment = text.Middle
		return layout.Center.Layout(gtx, lbl.Layout)
	}

	a.syncProfileWidgets()
	profiles := a.state.Config.Profiles

	return a.profileList.Layout(gtx, len(profiles), func(gtx layout.Context, i int) layout.Dimensions {
		return a.layoutProfileRow(gtx, profiles[i])
	})
}

func (a *App) syncProfileWidgets() {
	seen := make(map[string]bool, len(a.state.Config.Profiles))
	for _, p := range a.state.Config.Profiles {
		seen[p.ID] = true
		if _, ok := a.profileWidgets[p.ID]; !ok {
			a.profileWidgets[p.ID] = &profileRowWidgets{}
		}
	}
	for id := range a.profileWidgets {
		if !seen[id] {
			delete(a.profileWidgets, id)
		}
	}
}

func (a *App) layoutProfileRow(gtx layout.Context, p config.Profile) layout.Dimensions {
	w := a.profileWidgets[p.ID]
	selected := a.state.SelectedID == p.ID

	return layout.Inset{Bottom: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				if !selected {
					return layout.Dimensions{Size: gtx.Constraints.Min}
				}
				return drawSelectionBg(gtx, a.th)
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return material.Clickable(gtx, &w.selectBtn, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										title := material.Body1(a.th, p.Name)
										title.Color = a.th.Palette.Fg
										return title.Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										sub := material.Caption(a.th, sprintProfileSubtitle(p))
										sub.Color = mutedColor(a.th)
										return sub.Layout(gtx)
									}),
								)
							})
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(a.th, &w.editBtn, "Edit")
							btn.TextSize = unit.Sp(12)
							btn.Inset = layout.UniformInset(unit.Dp(6))
							return btn.Layout(gtx)
						}),
					)
				})
			}),
		)
	})
}

func sprintProfileSubtitle(p config.Profile) string {
	out := "index " + itoa(p.Index)
	if p.Options.FPS > 0 {
		out += "  " + itoa(p.Options.FPS) + " fps"
	}
	if p.Options.ScreenWidth > 0 && p.Options.ScreenHeight > 0 {
		out += "  " + itoa(p.Options.ScreenWidth) + "x" + itoa(p.Options.ScreenHeight)
	}
	return out
}

func (a *App) layoutLaunchBar(gtx layout.Context) layout.Dimensions {
	return layout.Inset{Top: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEnd}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.th, &a.launchDesktopBtn, "Launch Desktop")
				return btn.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.th, &a.launchVRBtn, "Launch VR")
				return btn.Layout(gtx)
			}),
		)
	})
}

func (a *App) layoutStatus(gtx layout.Context) layout.Dimensions {
	if a.state.Status == "" {
		return layout.Dimensions{}
	}
	return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		lbl := material.Caption(a.th, a.state.Status)
		lbl.Color = mutedColor(a.th)
		return lbl.Layout(gtx)
	})
}

func drawSelectionBg(gtx layout.Context, th *material.Theme) layout.Dimensions {
	c := th.Palette.ContrastBg
	c.A = 0x40
	rect := image.Rectangle{Max: gtx.Constraints.Min}
	defer clip.UniformRRect(rect, gtx.Dp(4)).Push(gtx.Ops).Pop()
	paint.ColorOp{Color: c}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: rect.Max}
}

func mutedColor(th *material.Theme) color.NRGBA {
	c := th.Palette.Fg
	c.A = 0xa0
	return c
}
