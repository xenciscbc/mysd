package cmd

import _ "embed"

//go:embed hooks/mysd-statusline.js
var statuslineHookBytes []byte
