package main

import (
	"fmt"
	"log"

	"github.com/dop251/goja"
)

type GOJA struct {
	vm *goja.Runtime
}

func NewGOJA() *GOJA {
	vm := goja.New()
	_, err := vm.RunString("window = {};")
	if err != nil {
		log.Panic(err)
	}
	console := vm.NewObject()
	consoleLog := func(call goja.FunctionCall) goja.Value {
		var args []interface{}
		for _, arg := range call.Arguments {
			args = append(args, arg.String())
		}
		log.Println(args...)
		return goja.Undefined()
	}
	llama := vm.NewObject()
	generate := func(call goja.FunctionCall) goja.Value {
		return vm.ToValue(Query(call.Arguments[0].String()))
	}
	err = console.Set("log", consoleLog)
	if err != nil {
		panic(err)
	}
	err = vm.Set("console", console)
	if err != nil {
		panic(err)
	}
	err = llama.Set("generate", generate)
	if err != nil {
		panic(err)
	}
	err = vm.Set("llama", llama)
	if err != nil {
		panic(err)
	}
	return &GOJA{
		vm: vm,
	}
}

func (g *GOJA) Run(i int, code string) error {
	program, err := goja.Compile(fmt.Sprintf("code%d", i), code, true)
	if err != nil {
		return err
	}
	_, err = g.vm.RunProgram(program)
	if err != nil {
		return err
	}
	return nil
}
