package main

import (
	"flag"
	"log"
	"os"

	"github.com/bytecodealliance/wasmtime-go/v18"
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
	wasm, err := os.ReadFile(wasmFile)
	if err != nil {
		log.Fatalf("failed to read wasm %s: %s", wasmFile, err)
	}

	module, err := wasmtime.NewModule(store.Engine, wasm)
	if err != nil {
		log.Fatalf("failed to create module: %s", err)
	}

	doDouble := wasmtime.WrapFunc(store, func(n int32) int32 {
		return n * 2
	})

	instance, err := wasmtime.NewInstance(store, module, []wasmtime.AsExtern{doDouble})
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
