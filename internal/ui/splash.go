// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package ui

import (
	"fmt"
	"strings"

	"github.com/liupzmin/tview"
	"github.com/liupzmin/weewoe/internal/config"
)

// LogoSmall WeeWoe small log.
var LogoSmall = []string{

	`__  _  __ ____   ______  _  ______   ____`,
	`\ \/ \/ // __ \_/ __ \ \/ \/ /  _ \_/ __ \ `,
	` \     /\  ___/\  ___/\     (  <_> )  ___/`,
	`  \/\_/  \___  >\___  >\/\_/ \____/ \___  >`,
}

// LogoBig WeeWoe big logo for splash page.
var LogoBig = []string{
	` __  __  __     __     __   __  __  __    ___      __`,
	`/\ \/\ \/\ \  /'__'\ /'__'\/\ \/\ \/\ \  / __'\  /'__'\ `,
	`\ \ \_/ \_/ \/\  __//\  __/\ \ \_/ \_/ \/\ \L\ \/\  __/ `,
	` \ \___x___/'\ \____\ \____\\ \___x___/'\ \____/\ \____\`,
	`  \/__//__/   \/____/\/____/ \/__//__/   \/___/  \/____/`,
}

// Splash represents a splash screen.
type Splash struct {
	*tview.Flex
}

// NewSplash instantiates a new splash screen with product and company info.
func NewSplash(styles *config.Styles, version string) *Splash {
	s := Splash{Flex: tview.NewFlex()}
	s.SetBackgroundColor(styles.BgColor())

	logo := tview.NewTextView()
	logo.SetDynamicColors(true)
	logo.SetTextAlign(tview.AlignCenter)
	s.layoutLogo(logo, styles)

	vers := tview.NewTextView()
	vers.SetDynamicColors(true)
	vers.SetTextAlign(tview.AlignCenter)
	s.layoutRev(vers, version, styles)

	s.SetDirection(tview.FlexRow)
	s.AddItem(logo, 10, 1, false)
	s.AddItem(vers, 1, 1, false)

	return &s
}

func (s *Splash) layoutLogo(t *tview.TextView, styles *config.Styles) {
	logo := strings.Join(LogoBig, fmt.Sprintf("\n[%s::b]", styles.Body().LogoColor))
	fmt.Fprintf(t, "%s[%s::b]%s\n",
		strings.Repeat("\n", 2),
		styles.Body().LogoColor,
		logo)
}

func (s *Splash) layoutRev(t *tview.TextView, rev string, styles *config.Styles) {
	fmt.Fprintf(t, "[%s::b]Revision [red::b]%s", styles.Body().FgColor, rev)
}
