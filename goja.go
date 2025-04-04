package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/dop251/goja"
)

type CAS interface {
	Compile(i int, algebrite []byte) error
	Load() error
	Run(line string) (string, error)
}

type GOJA struct {
	vm      *goja.Runtime
	program *goja.Program
}

func NewGOJA() CAS {
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
	err = console.Set("log", consoleLog)
	if err != nil {
		panic(err)
	}
	err = vm.Set("console", console)
	if err != nil {
		panic(err)
	}
	return &GOJA{
		vm: vm,
	}
}

func (g *GOJA) Compile(i int, algebrite []byte) error {
	program, err := goja.Compile(fmt.Sprintf("code%d", i), string(algebrite), true)
	if err != nil {
		return err
	}
	g.program = program
	return nil
}

func (g *GOJA) Load() error {
	_, err := g.vm.RunProgram(g.program)
	if err != nil {
		return err
	}
	return nil
}

func (g *GOJA) Run(line string) (string, error) {
	vm := g.vm
	window := vm.Get("window").ToObject(vm)
	alg := window.Get("Algebrite").ToObject(vm)
	run, valid := goja.AssertFunction(alg.Get("run"))
	if !valid {
		return "", errors.New("window.Algebrite.run is not a function")
	}
	result, err := run(goja.Null(), vm.ToValue(line))
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
