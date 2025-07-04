// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package dialog

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/liupzmin/tview"
	"github.com/liupzmin/weewoe/internal/config"
	"github.com/liupzmin/weewoe/internal/ui"
)

// ShowConfirm pops a confirmation dialog.
func ShowError(styles config.Dialog, pages *ui.Pages, msg string) {
	f := tview.NewForm()
	f.SetItemPadding(0)
	f.SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(styles.ButtonBgColor.Color()).
		SetButtonTextColor(styles.ButtonFgColor.Color()).
		SetLabelColor(styles.LabelFgColor.Color()).
		SetFieldTextColor(tcell.ColorIndianRed)
	f.AddButton("OK", func() {
		dismissError(pages)
	})
	if b := f.GetButton(0); b != nil {
		b.SetBackgroundColorActivated(styles.ButtonFocusBgColor.Color())
		b.SetLabelColorActivated(styles.ButtonFocusFgColor.Color())
	}
	f.SetFocus(0)
	modal := tview.NewModalForm("<error>", f)
	modal.SetText(cowTalk(msg))
	modal.SetTextColor(tcell.ColorOrangeRed)
	modal.SetDoneFunc(func(int, string) {
		dismissError(pages)
	})
	pages.AddPage(confirmKey, modal, false, false)
	pages.ShowPage(confirmKey)
}

func dismissError(pages *ui.Pages) {
	pages.RemovePage(confirmKey)
}

func cowTalk(says string) string {
	msg := fmt.Sprintf("< Ruroh? %s >", says)
	buff := make([]string, 0, len(cow)+3)
	buff = append(buff, msg)
	buff = append(buff, cow...)

	return strings.Join(buff, "\n")
}

var cow = []string{
	`\   ^__^            `,
	` \  (oo)\_______    `,
	`    (__)\       )\/\`,
	`        ||----w |   `,
	`        ||     ||   `,
}
