# wasm demo for ToolGUI


* xnix

```shell
GOOS=js GOARCH=wasm go build -o main.wasm
```

* win

```powershell
$env:GOOS = 'js'; $env:GOARCH = 'wasm'; go build -o main.wasm
```
