package scripts

import "embed"

//go:embed server.sh tui.sh
var Files embed.FS
