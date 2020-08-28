package memlibrary

// #cgo LDFLAGS: -ldl -lrt
//#include "memlibrary.h"
//#include <stdlib.h>
import "C"
import (
	"errors"
	"fmt"
	"github.com/djhohnstein/librarian/pkg/datagram"
	"github.com/djhohnstein/librarian/pkg/memfile"
	"sync"
	"unsafe"
)

type ModuleFunctions struct {
	mainFunction unsafe.Pointer
	mainCallback unsafe.Pointer
}

type InMemoryModule struct {
	Name       string
	MemoryFile *memfile.InMemoryFile
	Exports    *ModuleFunctions
}

type moduleManager struct {
	moduleMap map[string]*InMemoryModule
	mtx       *sync.Mutex
}

var modManager = &moduleManager{
	moduleMap: make(map[string]*InMemoryModule),
	mtx:       &sync.Mutex{},
}

func (modMan *moduleManager) addModule(mod *InMemoryModule) {
	modMan.mtx.Lock()
	defer modMan.mtx.Unlock()
	modMan.moduleMap[mod.Name] = mod
}

func (modMan *moduleManager) removeModule(name string) {
	modMan.mtx.Lock()
	defer modMan.mtx.Unlock()
	modMan.moduleMap[name].Dispose()
	delete(modMan.moduleMap, name)
}

func InvokeLoadedModule(name string, data []byte) ([]byte, error) {
	modManager.mtx.Lock()
	defer modManager.mtx.Unlock()
	if _, ok := modManager.moduleMap[name]; ok {
		dg := datagram.New(data, 0)
		results := (*datagram.Datagram)(C.call_module_function(modManager.moduleMap[name].Exports.mainFunction, unsafe.Pointer(dg)))
		data := C.GoBytes(results.Results, C.int((*results).Length))
		return data, nil
	}
	return nil, errors.New(fmt.Sprintf("Module '%s' is not loaded.", name))
}

func RemoveLoadedModule(name string) {
	modManager.removeModule(name)
}

//export RouteDataFromModule
func RouteDataFromModule(cData unsafe.Pointer) {
	var dg *datagram.Datagram
	if cData != nil {
		dg = (*datagram.Datagram)(cData)
		data := C.GoBytes((*dg).Results, C.int((*dg).Length))
		fmt.Printf("App Core got data from module: %s\n", data)
	}
}

func (memMod *InMemoryModule) Dispose() {
	if !memMod.MemoryFile.Closed {
		memMod.MemoryFile.Close()
	}
	C.free(memMod.Exports.mainFunction)
	C.free(memMod.Exports.mainCallback)
}

func (memMod *InMemoryModule) SendDataToModule(data []byte) {
	dg := datagram.New(data, 0)
	defer dg.Dispose()
	ptr := (*memMod).Exports.mainCallback
	C.call_module_callback(ptr, unsafe.Pointer(dg))
}

func (memMod *InMemoryModule) Invoke(arguments []byte) []byte {
	dg := datagram.New(arguments, 0)
	defer dg.Dispose()
	results := (*datagram.Datagram)(C.call_module_function(
		memMod.Exports.mainFunction, unsafe.Pointer(dg)))
	data := C.GoBytes(results.Results, C.int((*results).Length))
	return data
}

func New(mFile *memfile.InMemoryFile,
	moduleName, functionName, callbackName string) (*InMemoryModule, error) {
	cPath := C.CString(mFile.Path)
	cFuncName := C.CString(functionName)
	cCbName := C.CString(callbackName)
	defer func() {
		C.free(unsafe.Pointer(cPath))
		C.free(unsafe.Pointer(cFuncName))
		C.free(unsafe.Pointer(cCbName))
	}()
	res := C.load_module(cPath, cFuncName, cCbName)
	if res == nil {
		return nil, errors.New("Failed to acquire function exports.")
	}
	result := &InMemoryModule{
		Name:       moduleName,
		MemoryFile: mFile,
		Exports:    (*ModuleFunctions)(res),
	}
	modManager.addModule(result)
	return result, nil
}
