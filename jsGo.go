package jsGo

import "syscall/js"

var (
	// Global aliases
	Global = js.Global()
	Set = Global.Set
	Get = Global.Get
	Call = Global.Call

	// Constructor aliases
	Array = Get("Array")
	Error = Get("Error")
	Object = Get("Object")
	String = Get("String")
	TextDecoder = Get("TextDecoder")
	TextEncoder = Get("TextEncoder")
	Uint8Array = Get("Uint8Array")
	URLSearchParams = Get("URLSearchParams")

	// DOM aliases
	Document = Get("document")
	Head = Document.Get("head")
	Body = Document.Get("body")

	// BOM aliases
	Location = Get("location")

	// Method aliases
	Atob = Get("atob").Invoke
	Btoa = Get("btoa").Invoke
	ClearInterval = Get("clearInterval").Invoke
	Fetch = Get("fetch").Invoke
	GetRandomValues = Crypto.Get("getRandomValues").Invoke
	Log = Get("console").Get("log").Invoke
	Now = Get("Date").Get("now").Invoke
	ParseInt = Get("parseInt").Invoke
	SetInterval = Get("setInterval").Invoke
	SetTimeout = Get("setTimeout").Invoke

	// Namespace object aliases
	Crypto = Get("crypto")
	Math = Get("Math")
	Subtle = Crypto.Get("subtle")
	WebAssembly = Get("WebAssembly")

	// Convenience variables
	Params = URLSearchParams.New(Location.Get("search").String())
)

// DOM method handlers

func CreateElement(tag string) (js.Value) {return Document.Call("createElement", tag)}
func CreateElementNS(ns, tag string) (js.Value) {return Document.Call("createElementNS", ns, tag)}
func GetElementById(id string) (js.Value) {return Document.Call("getElementById", id)}
func AppendChild(child js.Value) {Body.Call("appendChild", child)}
func Append(child js.Value) {Body.Call("append", child)}
func Prepend(child js.Value) {Body.Call("prepend", child)}

// JS type-specific method handlers

func IsError(err js.Value) (bool) {return Error.Call("isError", err).Bool()}
func HasOwn(object js.Value, str string) (bool) {return Object.Call("hasOwn", object, str).Bool()}

// Wrap a Go function so it can be used by JS
func FuncOf(fn func(args []js.Value) (any)) (js.Func) {
	return js.FuncOf(func(this js.Value, args []js.Value) (any) {return fn(args)})
}

// Wrap a Go function which has no arguments so it can be used by JS
func SimpleFuncOf(fn func() (any)) (js.Func) {
	return FuncOf(func(args []js.Value) (any) {return fn()})
}

// Wrap a Go procedure so it can be used by JS
func ProcOf(fn func(args []js.Value)) (js.Func) {
	return FuncOf(func(args []js.Value) (any) {
		fn(args)
		return nil
	})
}

// Wrap a Go procedure which has no arguments so it can be used by JS
func SimpleProcOf(fn func()) (js.Func) {
	return ProcOf(func(args []js.Value) {fn()})
}

// Expose a Go function to JS globally
func SetFunc(str string, fn func(args []js.Value) (any)) (js.Func) {
	jsFunc := FuncOf(fn)
	Set(str, jsFunc)
	return jsFunc
}

// Expose a Go function which has no arguments to JS globally
func SetSimpleFunc(str string, fn func() (any)) (js.Func) {
	jsFunc := SimpleFuncOf(func() (any) {return fn()})
	Set(str, jsFunc)
	return jsFunc
}

// Expose a Go procedure to JS globally
func SetProc(str string, fn func(args []js.Value)) (js.Func) {
	jsProc := ProcOf(fn)
	Set(str, jsProc)
	return jsProc
}

// Expose a Go procedure which has no arguments to JS globally
func SetSimpleProc(str string, fn func()) (js.Func) {
	jsProc := SimpleProcOf(fn)
	Set(str, jsProc)
	return jsProc
}

// Function to handle the then, catch, and finally methods for a given thenable's thenable chain
func ThenableChain(thenable js.Value, thenFunc func(arg js.Value) (any), funcs... func(arg js.Value)) (any) {
	thenReturn := thenable.Call(
		"then",
		FuncOf(func(args []js.Value) (any) {return thenFunc(args[0])}),
	)
	var catchReturn js.Value
	if len(funcs) > 0 {
		catchReturn = thenReturn.Call(
			"catch",
			ProcOf(func(args []js.Value) {if IsError(args[0]) {funcs[0](args[0])}}),
		)
	} else {return thenReturn}
	if len(funcs) > 1 {
		finallyReturn := catchReturn.Call(
			"finally",
			SimpleProcOf(func() {funcs[1](js.Value{})}),
		)
		return finallyReturn
	} else {return catchReturn}
}

// Load a classic JS script with an onload callback
func LoadJS(src string, onload func()) (classic js.Value) {
	classic = CreateElement("script")
	classic.Set("src", src)
	classic.Set("onload", SimpleProcOf(onload))
	classic.Set("async", true)
	Head.Call("appendChild", classic)
	return
}

// Load a WASM module with a thenable chain
func LoadWASM(str string, then func(), methods... func(err js.Value)) {
		goInstance := Get("Go").New()
		ThenableChain(
			WebAssembly.Call(
				"instantiateStreaming",
				Fetch(str),
				goInstance.Get("importObject"),
			),
			func(module js.Value) (any) {
				goInstance.Call("run", module.Get("instance"))
				then()
				return nil
			},
			methods...,
		)
}

// Create a button element with an onclick callback
func CreateButton(str string, onclick func()) (button js.Value) {
	button = CreateElement("button")
	button.Set("textContent", str)
	button.Set("onclick", ProcOf(func(event []js.Value) {
		event[0].Call("preventDefault")
		onclick()
	}))
	return
}

// 32-bit variation of the FNV-1a non-cryptographic hashing algorithm
func FNV1a32(data []byte) (hash uint32) {
	hash = 0x811c9dc5
	for _, char := range data {
		hash ^= uint32(char)
		hash *= 0x01000193
	}
	return
}

// 64-bit variation of the FNV-1a non-cryptographic hashing algorithm
func FNV1a64(data []byte) (hash uint64) {
	hash = 0xcbf29ce484222325
	for _, char := range data {
		hash ^= uint64(char)
		hash *= 0x100000001b3
	}
	return
}
