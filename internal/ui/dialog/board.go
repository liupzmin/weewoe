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

// ShowBoard pops a confirmation dialog.
func ShowBoard(styles config.Dialog, pages *ui.Pages, msg string) {
	f := tview.NewForm()
	f.SetItemPadding(0)
	f.SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(styles.ButtonBgColor.Color()).
		SetButtonTextColor(styles.ButtonFgColor.Color()).
		SetLabelColor(styles.LabelFgColor.Color()).
		SetFieldTextColor(tcell.ColorIndianRed)
	f.AddButton("OK", func() {
		dismissBoard(pages)
	})
	if b := f.GetButton(0); b != nil {
		b.SetBackgroundColorActivated(styles.ButtonFocusBgColor.Color())
		b.SetLabelColorActivated(styles.ButtonFocusFgColor.Color())
	}
	f.SetFocus(0)
	modal := tview.NewModalForm("<board>", f)
	modal.SetText(cowTalk2(msg))
	modal.SetTextColor(tcell.ColorOrangeRed)
	modal.SetDoneFunc(func(int, string) {
		dismissBoard(pages)
	})
	pages.AddPage(confirmKey, modal, false, false)
	pages.ShowPage(confirmKey)
}

func dismissBoard(pages *ui.Pages) {
	pages.RemovePage(confirmKey)
}

func cowTalk2(says string) string {
	msg := fmt.Sprintf("< Woo? %s >", says)
	buff := make([]string, 0, len(cow)+3)
	buff = append(buff, msg)
	buff = append(buff, cow...)

	return strings.Join(buff, "\n")
}
