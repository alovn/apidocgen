package gen

import (
	_ "embed"
)

//go:embed template/default.tpl
// var defaultTemplate embed.FS
var defaultTemplate string
