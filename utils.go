package simpleconsoleui

import "github.com/rivo/tview"

/* -------------------------------------------------------------------------- */
/*                                 PUBLIC AREA                                */
/* -------------------------------------------------------------------------- */
func CenterScreen(widget tview.Primitive, width, height int) tview.Primitive {
	flexCenter := tview.NewFlex()
	flexCenter.SetDirection(tview.FlexRow)
	flexCenter.AddItem(nil, 0, 1, false)
	flexCenter.AddItem(tview.NewFlex().AddItem(nil, 0, 1, false).AddItem(widget, width, 0, true).AddItem(nil, 0, 1, false), height, 0, true)
	flexCenter.AddItem(nil, 0, 1, false)
	return flexCenter
}
