package cli

import (
	_ "embed"
)

//go:embed config.tpl.yaml
var configTemplate []byte
