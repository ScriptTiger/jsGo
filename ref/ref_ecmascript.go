//go:build ecmascript

package main

import (
	"syscall/js"

	"github.com/ScriptTiger/jsGo"
)

// Then() function for jsGo.LoadWASM()
func testFunc() {
	array := jsGo.Call("funcA")
	jsGo.Document.Set("title", array.Index(0).String())
	jsGo.Body.Set("innerHTML", array.Index(1).String()+jsGo.Call("funcB").String())
}

// Catch(err js.Value) function for jsGo.LoadWASM()
func errorFunc(err js.Value) {
	jsGo.Body.Set("innerHTML", "There was an error loading the WASM module!: "+err.Get("message").String())
}

func main() {

	// Load TinyGo "glue code"
	jsGo.LoadJS("https://cdn.jsdelivr.net/gh/tinygo-org/tinygo@0.40.1/targets/wasm_exec.js", func() {

		// Load WASM modules within thenable chains with callbacks
		jsGo.LoadWASM("moduleA.wasm", func() {jsGo.LoadWASM("moduleB.wasm", testFunc, errorFunc)}, errorFunc)

	})

}