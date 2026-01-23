//go:build wasm && moduleB

package main

import "github.com/ScriptTiger/jsGo"

// Test function to return a string
func funcB() (any) {
	return "<p>Test body string from moduleB</p>"
}

func main() {

	// Expose funcB to JS globally
	jsGo.SetSimpleFunc("funcB", funcB)

	// Keep this module open so that it can be used as a library and its functions can be called externally as needed
	select {}
}