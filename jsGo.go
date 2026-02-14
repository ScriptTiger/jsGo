package jsGo

import "syscall/js"

var (
	// Global aliases
	Global = js.Global()
	Set = Global.Set
	Get = Global.Get
	Call = Global.Call

	// Global property aliases
	DevicePixelRatio = Get("devicePixelRatio")
	InnerHeight = Get("innerHeight")
	InnerWidth = Get("innerWidth")
	OuterHeight = Get("outerHeight")
	OuterWidth = Get("outerWidth")

	// Constructor aliases
	Array = Get("Array")
	AudioContext = Get("AudioContext")
	Blob = Get("Blob")
	Date = Get("Date")
	Error = Get("Error")
	FileReader = Get("FileReader")
	Number = Get("Number")
	Object = Get("Object")
	Response = Get("Response")
	String = Get("String")
	TextDecoder = Get("TextDecoder")
	TextEncoder = Get("TextEncoder")
	Uint8Array = Get("Uint8Array")
	URL = Get("URL")
	URLSearchParams = Get("URLSearchParams")

	// BOM namespace object aliases
	Crypto = Get("crypto")
	History = Get("history")
	IndexedDB = Get("indexedDB")
	Intl = Get("Intl")
	IDBKeyRange = Get("IDBKeyRange")
	JSON = Get("JSON")
	Location = Get("location")
	Math = Get("Math")
	Navigator = Get("navigator")
	Performance = Get("performance")
	Screen = Get("screen")
	Subtle = Crypto.Get("subtle")
	WebAssembly = Get("WebAssembly")

	// DOM aliases
	Document = Get("document")
	Head = Document.Get("head")
	Body = Document.Get("body")

	// Method aliases
	Alert = Get("alert").Invoke
	Atob = Get("atob").Invoke
	Btoa = Get("btoa").Invoke
	ClearInterval = Get("clearInterval").Invoke
	ClearTimeout = Get("clearTimeout").Invoke
	Fetch = Get("fetch").Invoke
	GetRandomValues = Crypto.Get("getRandomValues").Invoke
	Log = Get("console").Get("log").Invoke
	MatchMedia = Get("matchMedia").Invoke
	ParseInt = Get("parseInt").Invoke
	SetInterval = Get("setInterval").Invoke
	SetTimeout = Get("setTimeout").Invoke
	ShowSaveFilePicker = Get("showSaveFilePicker").Invoke

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

// Create a span element containing both a hidden input element of type file with an onchange callback and also a button to trigger it
func CreateLoadFileButton(text, accept string, multiple bool, onchange func(loadFiles js.Value)) (span js.Value) {
	span = CreateElement("span")
	input := CreateElement("input")
	input.Set("type", "file")
	input.Set("accept", accept)
	input.Set("multiple", multiple)
	input.Set("style", "display: none;")
	input.Set("onchange", ProcOf(func(event []js.Value) {
		onchange(event[0].Get("target").Get("files"))
		input.Set("value", nil)
	}))
	span.Call("appendChild", input)
	button := CreateButton(text, func() {input.Call("click")})
	span.Call("appendChild", button)
	return
}

// Create a button which calls showSaveFilePicker with a callback
func CreateSaveFileButton(text string, options map[string]any, saveFileCallback func(saveFile js.Value)) (button js.Value) {
	button = CreateButton(text, func() {
		ThenableChain(
			ShowSaveFilePicker(options),
			func(saveFile js.Value) (any) {
				saveFileCallback(saveFile)
				return nil
			},
		)
	})
	return
}

// 32-bit variation of the FNV-1a non-cryptographic hashing algorithm which immediately/synchonously returns the hash as a JS ArrayBuffer
func FNV1a32(data []byte) (js.Value) {
	var hash uint32 = 0x811c9dc5
	for _, char := range data {
		hash ^= uint32(char)
		hash *= 0x01000193
	}
	bytes := make([]byte, 4)
	for i, _ := range bytes {
		bytes[i] = byte(hash >> (24 - (i * 8)))
	}
	jsBytes := Uint8Array.New(4)
	js.CopyBytesToJS(jsBytes, bytes)
	return jsBytes.Get("buffer")
}

// 64-bit variation of the FNV-1a non-cryptographic hashing algorithm which immediately/synchonously returns the hash as a JS ArrayBuffer
func FNV1a64(data []byte) (js.Value) {
	var hash uint64 = 0xcbf29ce484222325
	for _, char := range data {
		hash ^= uint64(char)
		hash *= 0x100000001b3
	}
	bytes := make([]byte, 8)
	for i, _ := range bytes {
		bytes[i] = byte(hash >> (56 - (i * 8)))
	}
	jsBytes := Uint8Array.New(8)
	js.CopyBytesToJS(jsBytes, bytes)
	return jsBytes.Get("buffer")
}

// Base function for SHA-2 family functions which takes the bit length, data byte slice, and callback to asynchonrously handle the hash as a JS ArrayBuffer
func sha2(length int, data []byte, shaCallback func(hash js.Value)) {
	jsData := Uint8Array.New(len(data))
	js.CopyBytesToJS(jsData, data)
	ThenableChain(
		Subtle.Call("digest", "SHA-"+String.Invoke(length).String(), jsData),
		func(hash js.Value) (any) {
			shaCallback(hash)
			return nil
		},
	)
	return
}

// 256-bit variation of SHA-2 which takes the data byte slice and callback to asynchonrously handle the hash as a JS ArrayBuffer
func SHA256(data []byte, shaCallback func(hash js.Value)) {sha2(256, data, shaCallback)}

// 384-bit variation of SHA-2 which takes the data byte slice and callback to asynchonrously handle the hash as a JS ArrayBuffer
func SHA384(data []byte, shaCallback func(hash js.Value)) {sha2(384, data, shaCallback)}

// 512-bit variation of SHA-2 which takes the data byte slice and callback to asynchonrously handle the hash as a JS ArrayBuffer
func SHA512(data []byte, shaCallback func(hash js.Value)) {sha2(512, data, shaCallback)}

// Takes a permission descriptor, such as "camera", "microphone", "geolocation", etc., and a callback to handle the returned PermissionStatus object
func Permissions(str string, permissionsCallback func(permissionStatus js.Value)) {
	ThenableChain(
		Navigator.Get("permissions").Call("query", map[string]any{"name": str}),
		func(permissionStatus js.Value) (any) {
			permissionsCallback(permissionStatus)
			return nil
		},
	)
}
