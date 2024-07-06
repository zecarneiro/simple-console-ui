package simpleconsoleui

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	appMainId = "main"
)

/* -------------------------------------------------------------------------- */
/*                                 MODEL AREA                                 */
/* -------------------------------------------------------------------------- */
type Window struct {
	MenuName string
	Callback func()
	HasLog   bool
	MenuPage func() tview.Primitive
	Shurtcut rune // When use integer, set rune(integer+'0'), empty ser 0
	page     tview.Primitive
}
/* ----------------------------- END MODEL AREA ----------------------------- */

var (
	app         *tview.Application
	appPages *tview.Pages

	// Main Window Area
	windows     []Window
	mainFlex        *tview.Flex
	mainPages       *tview.Pages
	mainMenu    *tview.List
	title       string
	description string
	mainMenuWidth int

	// Selected Window menu id state
	selectedMenuId string
)

/* -------------------------------------------------------------------------- */
/*                                 WINDOW AREA                                */
/* -------------------------------------------------------------------------- */
func (w *Window) addPage(pageWidget tview.Primitive) {
	w.page = pageWidget
}

func (w *Window) getScreenPage() tview.Primitive {
	return w.page
}

/* -------------------------------------------------------------------------- */
/*                                   UI AREA                                  */
/* -------------------------------------------------------------------------- */
func appStart() {
	done := make(chan struct{})
	go func() {
		defer func() {
			done <- struct{}{}
		}()
		if err := app.Run(); err != nil {
			panic(err)
		}
	}()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(
			signals,
			syscall.SIGILL, syscall.SIGINT, syscall.SIGTERM,
		)
		<-signals
		app.Stop()
	}()
	<-done
}

func createMenu(pages *tview.Pages) *tview.List {
	mainMenu = tview.NewList()
	mainMenuWidth = 0
	for _, window := range windows {
		if len(window.MenuName) > 0 {
			if mainMenuWidth < len(window.MenuName) {
				mainMenuWidth = len(window.MenuName)
			}
			if window.MenuPage != nil {
				mainMenu.AddItem(window.MenuName, "", window.Shurtcut, nil)
			} else if window.Callback != nil {
				mainMenu.AddItem(window.MenuName, "", window.Shurtcut, window.Callback)
			}
		}
	}
	mainMenu.ShowSecondaryText(false)
	mainMenu.SetChangedFunc(func(_ int, page string, _ string, _ rune) {
		canSwitch := false
		if !isModalOpen {
			for _, window := range windows {
				if window.Callback == nil {
					canSwitch = true
					break
				}
			}
		}
		if canSwitch {
			pages.SwitchToPage(page)
		}
	})
	mainMenu.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if !isModalOpen {
			return action, event
		}
		return tview.MouseConsumed, nil
	})
	mainMenu.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if !isModalOpen {
			if event.Key() != tcell.KeyRune {
				return event
			}
			if event.Rune() == 'j' {
				return tcell.NewEventKey(tcell.KeyDown, event.Rune(), event.Modifiers())
			}
			if event.Rune() == 'k' {
				return tcell.NewEventKey(tcell.KeyUp, event.Rune(), event.Modifiers())
			}
			return event
		}
		return nil
	})

	return mainMenu
}

func createSidebar(pages *tview.Pages) tview.Primitive {
	menu := createMenu(pages)
	frame := tview.NewFrame(menu)
	frame.SetBorder(true)
	frame.SetBorders(0, 0, 1, 1, 1, 1)
	frame.AddText(title, true, tview.AlignLeft, tcell.ColorWhite)
	frame.AddText(description, true, tview.AlignLeft, tcell.ColorBlue)
	return frame
}

func createPages() {
	mainPages = tview.NewPages()
	for i, window := range windows {
		if window.Callback == nil { 
			var view tview.Primitive
			if window.MenuPage != nil {
				view = window.getScreenPage()
			} else {
				view = tview.NewTextView().SetText(window.MenuName)
			}
			mainPages.AddPage(window.MenuName, view, true, i == 0)
		}
	}

	mainPages.SetBorder(true)
}

func createUi() {
	appPages = tview.NewPages()
	createPages()
	menu := createSidebar(mainPages)

	mainFlex = tview.NewFlex()
	mainFlex.AddItem(menu, mainMenuWidth+4, 1, false)
	mainFlex.AddItem(mainPages, 0, 1, true)
	appPages.AddPage(appMainId, mainFlex, true, true)

	app.SetRoot(appPages, true)
	app.EnableMouse(true)

	focusingMenu := false
	app.SetFocus(mainPages)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.Stop()
		}
		key := event.Key()
		if key != tcell.KeyTAB && key != tcell.KeyBacktab {
			return event
		}

		if focusingMenu {
			app.SetFocus(mainPages)
		} else {
			app.SetFocus(menu)
		}
		focusingMenu = !focusingMenu
		return nil
	})
}

func buildPage() {
	for index, window := range windows {
		if window.MenuPage != nil {
			if window.HasLog {
				flexPage := tview.NewFlex()
				flexPage.SetDirection(tview.FlexRow)
				flexPage.AddItem(window.MenuPage(), 0, 10, true)
				flexPage.AddItem(labelMessage, 0, 3, false)
				windows[index].addPage(flexPage)
			} else {
				windows[index].addPage(window.MenuPage())
			}
		}
	}
}

func validateAndSetVars() {
	if app == nil {
		panic("Invalid app")
	}
	for _, window := range windows {
		count := 0
		for _, wwindowToValidate := range windows {
			if window.MenuName == wwindowToValidate.MenuName {
				count++
			}
		}
		if count > 1 {
			panic("Duplicate Menu name entry")
		}
	}
}

/* -------------------------------------------------------------------------- */
/*                                 PUBLIC AREA                                */
/* -------------------------------------------------------------------------- */
func InitUi(theme tview.Theme) {
	if (tview.Theme{}) == theme {
		theme = tview.Theme{
			PrimitiveBackgroundColor:    tcell.ColorBlack,
			ContrastBackgroundColor:     tcell.ColorDarkBlue,
			MoreContrastBackgroundColor: tcell.ColorGreen,
			BorderColor:                 tcell.ColorWhite,
			TitleColor:                  tcell.ColorWhite,
			GraphicsColor:               tcell.ColorWhite,
			PrimaryTextColor:            tcell.ColorGhostWhite,
			SecondaryTextColor:          tcell.ColorYellow,
			TertiaryTextColor:           tcell.ColorGreen,
			InverseTextColor:            tcell.ColorDeepSkyBlue,
			ContrastSecondaryTextColor:  tcell.ColorDarkCyan,
		}
	}
	tview.Styles = theme
}

func Start(application *tview.Application, givenWindows []Window, appName string, appLitleDesc string) {
	app = application
	windows = givenWindows
	title = appName
	description = appLitleDesc
	if windows == nil {
		windows = []Window{}
	}
	validateAndSetVars()
	loadLabelMessage()
	buildPage()
	createUi()
	appStart()
}

func Refresh() {
	mainFlex.Clear()
	ClearLog()
	buildPage()
	createUi()
}

func RefreshAndKeepOnPage() {
	app.SetFocus(mainPages)
	pageName, _ := mainPages.GetFrontPage()
	Refresh()
	mainMenu.Blur()
	GoToPage(pageName)
}

func GoToPage(pageName string) {
	mainPages.SwitchToPage(pageName)
	for index, window := range windows {
		if window.MenuName == pageName {
			mainMenu.SetCurrentItem(index)
			mainPages.SwitchToPage(pageName)
		}
	}
}

func ResetFocus() {
	app.SetRoot(mainFlex, true).SetFocus(mainPages)
}

func SaveMainWindowState() {
	index := mainMenu.GetCurrentItem()
	selectedMenuId = windows[index].MenuName
}

func RestoreSavedWindowState() {
	GoToPage(selectedMenuId)
}
