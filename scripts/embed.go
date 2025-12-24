package scripts

import "embed"

//go:embed server.sh tui.sh uninstall.sh
var Files embed.FS
