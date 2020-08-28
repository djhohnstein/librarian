package memfile

// #cgo LDFLAGS: -lrt -ldl
//#include <unistd.h>
//#include <stdlib.h>
//#include "memfile.h"
import "C"
import (
	"errors"
	"fmt"
	"github.com/djhohnstein/librarian/pkg/osinfo"
	"os"
	"sync"
	"unsafe"
)

//export IsModernKernel
func IsModernKernel() bool {
	os, err := osinfo.New()
	if err != nil {
		panic(err)
	}
	if os.Release.Major < 3 {
		return false
	} else if os.Release.Minor > 3 {
		return false
	}

	if os.Release.Minor < 17 {
		return false
	} else {
		return true
	}
}

type InMemoryFile struct {
	FileDescriptor C.int
	Path           string
	Closed         bool
	mtx            *sync.Mutex
}

func New(shmName string) (*InMemoryFile, error) {
	cStr := C.CString(shmName)
	defer C.free(unsafe.Pointer(cStr))
	shmFD := C.open_ramfs(cStr)

	if int(shmFD) < 0 {
		return nil, errors.New(fmt.Sprintf("Failed to acquire file handle with return code %d", shmFD))
	}

	if !IsModernKernel() {
		return &InMemoryFile{
			FileDescriptor: shmFD,
			Path:           fmt.Sprintf("/proc/%d/fd/%d", os.Getpid(), int(shmFD)),
			Closed:         false,
			mtx:            &sync.Mutex{},
		}, nil
	} else {
		C.close(shmFD)
		return &InMemoryFile{
			FileDescriptor: shmFD,
			Path:           fmt.Sprintf("/dev/shm/%s", shmName),
			Closed:         true,
			mtx:            &sync.Mutex{},
		}, nil
	}
}

func (memFile *InMemoryFile) Write(data []byte) int {
	if memFile.Closed {
		panic(errors.New("File is closed can't write!"))
	}
	memFile.mtx.Lock()
	dataPtr := C.CBytes(data)
	dataLen := C.ulong(len(data))
	defer func() {
		C.free(dataPtr)
		memFile.mtx.Unlock()
	}()
	return int(C.write(memFile.FileDescriptor, dataPtr, dataLen))
}

func (memFile *InMemoryFile) Close() {
	if memFile.Closed {
		return
	}
	memFile.mtx.Lock()
	defer memFile.mtx.Unlock()
	C.close(memFile.FileDescriptor)
	memFile.Closed = true
}
