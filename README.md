syscall/js: wrapped functions are preempted within critical section

1. Clone the repo
2. Run `go run main.go` in ./server directory: it listens on port :8088
3. Open your web browser at 127.0.0.1:8088
4. Click where indicated and look at the console messages

If you have modified  the source code and want to recompile,
in the wasmtestissue folder, you can run: 
`GOOS=js GOARCH=wasm go build -o  server/assets/app.wasm`

