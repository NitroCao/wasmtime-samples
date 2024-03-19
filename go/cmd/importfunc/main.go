package main

import (
	"flag"
	"log"
	"os"

	"github.com/bytecodealliance/wasmtime-go/v18"
)

const (
	moduleName = "foo"
)

var (
	wasmFile string
)

func init() {
	flag.StringVar(&wasmFile, "wasm", "", "wasm file")
}

func main() {
	flag.Parse()

	store := wasmtime.NewStore(wasmtime.NewEngine())
	defer store.Close()
	wasm, err := os.ReadFile(wasmFile)
	if err != nil {
		log.Fatalf("failed to read wasm %s: %s", wasmFile, err)
	}

	linker := wasmtime.NewLinker(store.Engine)
	module, err := wasmtime.NewModule(store.Engine, wasm)
	if err != nil {
		log.Fatalf("failed to create module: %s", err)
	}
	defer module.Close()

	if err = linker.FuncWrap("", "", func(n int32) int32 {
		return n * 2
	}); err != nil {
		log.Fatalf("failed to wrap an anonymous function: %s", err)
	}
	if err = linker.FuncWrap(moduleName, "bar", func(n int32) int32 {
		return n + 1
	}); err != nil {
		log.Fatalf("failed to wrap a named function: %s", err)
	}
	instance, err := linker.Instantiate(store, module)
	if err != nil {
		log.Fatalf("failed to create instance: %s", err)
	}

	run := instance.GetFunc(store, "add")
	if run == nil {
		log.Printf("not a function")
	}

	if result, err := run.Call(store, 2, 3); err != nil {
		log.Fatalf("failed to run: %s", err)
	} else {
		log.Printf("result: %d", result)
	}
}
