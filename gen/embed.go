package gen

import (
	_ "embed"
)

//go:embed template/single.tpl
var singleApisTemplate string

//go:embed template/group_readme.tpl
var groupReadmeTemplate string

//go:embed template/group_apis.tpl
var groupApisTemplate string
