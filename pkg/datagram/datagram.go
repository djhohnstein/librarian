package datagram

import "C"
import "unsafe"

type Datagram struct {
	Length      C.int
	Results     unsafe.Pointer
	Name        *C.char
	MessageType C.int
}

func New(results []byte, messageType int) *Datagram {
	return &Datagram{
		Length:      C.int(len(results)),
		Results:     C.CBytes(results),
		Name:        C.CString(""),
		MessageType: C.int(messageType),
	}
}

func (dg *Datagram) Dispose() {
	C.free(dg.Results)
	C.free(dg.Name)
}
