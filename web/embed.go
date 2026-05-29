package web

import "embed"

//go:embed app.js app.css favicon.svg resistor.wasm wasm_exec.js
var FS embed.FS
