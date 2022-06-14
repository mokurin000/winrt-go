// Code generated by winrt-go-gen. DO NOT EDIT.

//go:build windows

//nolint
package streams

import (
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
)

type Buffer struct {
	IBuffer
}

const GUIDiBufferFactory string = "71af914d-c10f-484b-bc50-14bc623b3a27"

type iBufferFactory struct {
	ole.IInspectable
}

type iBufferFactoryVtbl struct {
	ole.IInspectableVtbl

	Create uintptr
}

func (v *iBufferFactory) VTable() *iBufferFactoryVtbl {
	return (*iBufferFactoryVtbl)(unsafe.Pointer(v.RawVTable))
}

func Create(capacity uint32) (*Buffer, error) {
	inspectable, err := ole.RoGetActivationFactory("Windows.Storage.Streams.Buffer", ole.NewGUID(GUIDiBufferFactory))
	if err != nil {
		return nil, err
	}
	v := (*iBufferFactory)(unsafe.Pointer(inspectable))

	var out *Buffer
	hr, _, _ := syscall.SyscallN(
		v.VTable().Create,
		0,                             // this is a static func, so there's no this
		uintptr(capacity),             // in capacity
		uintptr(unsafe.Pointer(&out)), // out *Buffer
	)

	if hr != 0 {
		return nil, ole.NewError(hr)
	}

	return out, nil
}
