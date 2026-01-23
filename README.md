[![Say Thanks!](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg)](https://docs.google.com/forms/d/e/1FAIpQLSfBEe5B_zo69OBk19l3hzvBmz3cOV6ol1ufjh0ER1q3-xd2Rg/viewform)

**DISCLAIMER!!!: THIS MODULE IS STILL IN ITS EARLY DEVELOPMENT, SO USE AT YOUR OWN RISK!**

# jsGo
The jsGo package is a collection of convenience functions and variables to make working with the standard Go `syscall/js` package in both GopherJS and Go/TinyGo WASM a bit more natural. For the most part, most code should be easily transferable from GopherJS to Go/TinyGo WASM, and vice versa, with minimal changes needed, if any.

**For Go/TinyGo WASM, jsGo requires the "glue code" to be loaded so that it can use the defined `Go` class when instantiating the WASM module, and also so that the WASM module can make call-outs to it should it need to access the JS APIs outside of the WASM sandbox.**

To import jsGo into your project:  
`go get github.com/ScriptTiger/jsGo`  
Then just `import "github.com/ScriptTiger/jsGo"` and get started with using its functions and variables.

Please refer to the dev package docs and reference implementation for more details and ideas on how to integrate jsGo into your project.

Dev package docs:  
https://pkg.go.dev/github.com/ScriptTiger/jsGo

Reference implementation:  
https://github.com/ScriptTiger/jsGo/blob/main/ref

# More About ScriptTiger

For more ScriptTiger scripts and goodies, check out ScriptTiger's GitHub Pages website:  
https://scripttiger.github.io/
