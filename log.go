package simpleconsoleui

import (
	"fmt"

	"github.com/rivo/tview"
)

var (
	labelMessage *tview.TextView
)

func loadLabelMessage() {
	if labelMessage == nil {
		labelMessage = tview.NewTextView()
	}
	labelMessage.Clear().SetTitle("Logs").SetBorder(true)
	labelMessage.SetDynamicColors(true)
}

func setLabelMessageData(data string) {
	oldData := labelMessage.GetText(false)
	if len(oldData) > 0 {
		data = oldData + "\n" + data
	}
	labelMessage.SetText(fmt.Sprintf("[white]%s", data))
}

/* -------------------------------------------------------------------------- */
/*                                 PUBLIC AREA                                */
/* -------------------------------------------------------------------------- */
func ClearLog() {
	labelMessage.SetText("[white]")
	loadLabelMessage()
}
func LogLog(data string) {
	setLabelMessageData(data)
}

func DebugLog(data string) {
	setLabelMessageData(fmt.Sprintf("[DEBUG] %s", data))
}

func WarnLog(data string) {
	setLabelMessageData(fmt.Sprintf("[[yellow]WARN[white]] %s", data))
}

func ErrorLog(data string) {
	setLabelMessageData(fmt.Sprintf("[[red]ERROR[white]] %s", data))
}

func InfoLog(data string) {
	setLabelMessageData(fmt.Sprintf("[[blue]INFO[white]] %s", data))
}

func OkLog(data string) {
	setLabelMessageData(fmt.Sprintf("[[green]OK[white]] %s", data))
}

func PromptLog(data string) {
	setLabelMessageData(fmt.Sprintf("[gray]>>> %s[white]", data))
}
