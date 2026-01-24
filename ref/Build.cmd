@echo off

rem BUILD REQUIREMENTS
rem Go, TinyGo, GopherJS, and terser (via Node.js with npx) must all be installed and in your path
rem Make sure you set your GOPHERJS_GOROOT environmental variable as needed
rem set GOPHERJS_GOROOT=C:\path-to-gopherjs-goroot

set app=jsGo

if not exist Release md Release

:ecmascript
echo Building JS...
set GOOS=js
set GOARCH=ecmascript
call gopherjs build -o Release/%app%.js
echo Minifying JS...
call npx terser Release/%app%.js -c -m -o Release/%app%.m.js

:wasm
set GOARCH=wasm
call :build_wasm moduleA
call :build_wasm moduleB

:html
if not exist Release\index.html (
	copy index.html Release
)

:server
set GOOS=
set GOARCH=
if not exist Release\server.exe (
	echo Building server.exe...
	call go build -ldflags="-s -w" -o Release/server.exe
)

:done
del Release\%app%.js Release\%app%.js.map
echo Done
pause
exit /b

:build_wasm
echo Building %1 WASM...
call tinygo build -no-debug -gc=leaking -panic=trap -tags wasm,%1 -o Release/%1.wasm
exit /b
