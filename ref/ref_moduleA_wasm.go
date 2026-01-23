//go:build wasm && moduleA

package main

import "github.com/ScriptTiger/jsGo"

// Test function to return a string array
func funcA() (any) {
	array := jsGo.Array.New(2)
	array.SetIndex(0, "Test title string from moduleA")
	array.SetIndex(1, "<p>Test body string from moduleA</p>")
	return array
}

func main() {

	// Expose funcA to JS globally
	jsGo.SetSimpleFunc("funcA", funcA)

	// Keep this module open so that it can be used as a library and its functions can be called externally as needed
	select {}
}