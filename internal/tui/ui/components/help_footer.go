package components

import (
	"strings"

	"charm.land/bubbles/v2/key"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

// RenderHelpFooter renders a list of key bindings as a help footer
func RenderHelpFooter(keys []key.Binding) string {
	var parts []string

	for _, k := range keys {
		if !k.Enabled() {
			continue
		}
		help := k.Help()
		if help.Key == "" {
			continue
		}
		keyStr := styles.ForegroundStyle().Render(help.Key)
		descStr := styles.MutedStyle().Render(help.Desc)
		parts = append(parts, keyStr+" "+descStr)
	}

	separator := styles.DimStyle().Render(" · ")
	return strings.Join(parts, separator)
}
