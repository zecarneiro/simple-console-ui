package simpleconsoleui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zecarneiro/golangutils"
)

const (
	modalId = "modal"
	ModalTextColor = tcell.ColorWhite
)

var (
	isModalOpen bool
)

func showModal(modal *tview.Modal) {
	SaveMainWindowState()
	app.SetFocus(modal)
	if golangutils.InArray(appPages.GetPageNames(false), modalId) {
		appPages.RemovePage(modalId)
	}
	appPages.AddPage(modalId, modal, true, true)
	isModalOpen = true
}

func hideModal() {
	id, _ := appPages.GetFrontPage()
	if id == modalId {
		RestoreSavedWindowState()
		appPages.SwitchToPage(appMainId)
	}
	isModalOpen = false
}

func buildModal(message string, color tcell.Color) *tview.Modal {
	modal := tview.NewModal().SetFocus(0).SetBackgroundColor(color).SetTextColor(ModalTextColor)
	modal.SetText(message)
	return modal
}

func buildCloseModalOnly(title string, message string, closeLabel string, color tcell.Color, callback func()) *tview.Modal {
	if len(closeLabel) == 0 {
		closeLabel = "Close"
	}
	modal := buildModal(message, color)
	modal.SetTitle(title).SetTitleAlign(tview.AlignCenter)
	modal.AddButtons([]string{closeLabel}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if callback != nil {
			callback()
		}
		hideModal()
	})
	return modal
}

/* -------------------------------------------------------------------------- */
/*                                 PUBLIC AREA                                */
/* -------------------------------------------------------------------------- */
func Confirm(message string, confirmLabel string, cancelLabel string, callback func(canContinue bool)) {
	if len(confirmLabel) == 0 {
		confirmLabel = "Ok"
	}
	if len(cancelLabel) == 0 {
		cancelLabel = "Cancel"
	}
	modal := buildModal(message, tcell.ColorGray)
	modal.AddButtons([]string{confirmLabel, cancelLabel}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if callback != nil {
			callback(buttonLabel == confirmLabel)
		}
		hideModal()
	})
	showModal(modal)
}

func Error(message string, closeLabel string, callback func()) {
	showModal(buildCloseModalOnly("Error", message, closeLabel, tcell.ColorRed, callback))
}

func Information(message string, closeLabel string, callback func()) {
	showModal(buildCloseModalOnly("Information", message, closeLabel, tcell.ColorDarkBlue, callback))
}

func Ok(message string, closeLabel string, callback func()) {
	showModal(buildCloseModalOnly("Success", message, closeLabel, tcell.ColorGreen, callback))
}

func Warn(message string, closeLabel string, callback func()) {
	modal := buildCloseModalOnly("Warnning", message, closeLabel, tcell.ColorYellow, callback)
	modal.SetTextColor(tcell.ColorDarkCyan)
	showModal(modal)
}
