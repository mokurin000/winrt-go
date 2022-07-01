// Code generated by winrt-go-gen. DO NOT EDIT.

//go:build windows

//nolint
package advertisement

import (
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/saltosystems/winrt-go/windows/foundation"
)

type BluetoothLEAdvertisementWatcher struct {
	iBluetoothLEAdvertisementWatcher

	iBluetoothLEAdvertisementWatcher2
}

func NewBluetoothLEAdvertisementWatcher() (*BluetoothLEAdvertisementWatcher, error) {
	inspectable, err := ole.RoActivateInstance("Windows.Devices.Bluetooth.Advertisement.BluetoothLEAdvertisementWatcher")
	if err != nil {
		return nil, err
	}
	return (*BluetoothLEAdvertisementWatcher)(unsafe.Pointer(inspectable)), nil
}

const GUIDiBluetoothLEAdvertisementWatcher string = "a6ac336f-f3d3-4297-8d6c-c81ea6623f40"

type iBluetoothLEAdvertisementWatcher struct {
	ole.IInspectable
}

type iBluetoothLEAdvertisementWatcherVtbl struct {
	ole.IInspectableVtbl

	GetMinSamplingInterval  uintptr
	GetMaxSamplingInterval  uintptr
	GetMinOutOfRangeTimeout uintptr
	GetMaxOutOfRangeTimeout uintptr
	GetStatus               uintptr
	GetScanningMode         uintptr
	SetScanningMode         uintptr
	GetSignalStrengthFilter uintptr
	SetSignalStrengthFilter uintptr
	GetAdvertisementFilter  uintptr
	SetAdvertisementFilter  uintptr
	Start                   uintptr
	Stop                    uintptr
	AddReceived             uintptr
	RemoveReceived          uintptr
	AddStopped              uintptr
	RemoveStopped           uintptr
}

func (v *iBluetoothLEAdvertisementWatcher) VTable() *iBluetoothLEAdvertisementWatcherVtbl {
	return (*iBluetoothLEAdvertisementWatcherVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *iBluetoothLEAdvertisementWatcher) GetStatus() (BluetoothLEAdvertisementWatcherStatus, error) {
	var out BluetoothLEAdvertisementWatcherStatus
	hr, _, _ := syscall.SyscallN(
		v.VTable().GetStatus,
		uintptr(unsafe.Pointer(v)),    // this
		uintptr(unsafe.Pointer(&out)), // out BluetoothLEAdvertisementWatcherStatus
	)

	if hr != 0 {
		return BluetoothLEAdvertisementWatcherStatusCreated, ole.NewError(hr)
	}

	return out, nil
}

func (v *iBluetoothLEAdvertisementWatcher) Start() error {
	hr, _, _ := syscall.SyscallN(
		v.VTable().Start,
		uintptr(unsafe.Pointer(v)), // this
	)

	if hr != 0 {
		return ole.NewError(hr)
	}

	return nil
}

func (v *iBluetoothLEAdvertisementWatcher) Stop() error {
	hr, _, _ := syscall.SyscallN(
		v.VTable().Stop,
		uintptr(unsafe.Pointer(v)), // this
	)

	if hr != 0 {
		return ole.NewError(hr)
	}

	return nil
}

func (v *iBluetoothLEAdvertisementWatcher) AddReceived(handler *foundation.TypedEventHandler) (foundation.EventRegistrationToken, error) {
	var out foundation.EventRegistrationToken
	hr, _, _ := syscall.SyscallN(
		v.VTable().AddReceived,
		uintptr(unsafe.Pointer(v)),       // this
		uintptr(unsafe.Pointer(handler)), // in foundation.TypedEventHandler
		uintptr(unsafe.Pointer(&out)),    // out foundation.EventRegistrationToken
	)

	if hr != 0 {
		return foundation.EventRegistrationToken{}, ole.NewError(hr)
	}

	return out, nil
}

func (v *iBluetoothLEAdvertisementWatcher) RemoveReceived(token foundation.EventRegistrationToken) error {
	hr, _, _ := syscall.SyscallN(
		v.VTable().RemoveReceived,
		uintptr(unsafe.Pointer(v)),      // this
		uintptr(unsafe.Pointer(&token)), // in token
	)

	if hr != 0 {
		return ole.NewError(hr)
	}

	return nil
}

func (v *iBluetoothLEAdvertisementWatcher) AddStopped(handler *foundation.TypedEventHandler) (foundation.EventRegistrationToken, error) {
	var out foundation.EventRegistrationToken
	hr, _, _ := syscall.SyscallN(
		v.VTable().AddStopped,
		uintptr(unsafe.Pointer(v)),       // this
		uintptr(unsafe.Pointer(handler)), // in foundation.TypedEventHandler
		uintptr(unsafe.Pointer(&out)),    // out foundation.EventRegistrationToken
	)

	if hr != 0 {
		return foundation.EventRegistrationToken{}, ole.NewError(hr)
	}

	return out, nil
}

func (v *iBluetoothLEAdvertisementWatcher) RemoveStopped(token foundation.EventRegistrationToken) error {
	hr, _, _ := syscall.SyscallN(
		v.VTable().RemoveStopped,
		uintptr(unsafe.Pointer(v)),      // this
		uintptr(unsafe.Pointer(&token)), // in token
	)

	if hr != 0 {
		return ole.NewError(hr)
	}

	return nil
}

const GUIDiBluetoothLEAdvertisementWatcher2 string = "01bf26bc-b164-5805-90a3-e8a7997ff225"

type iBluetoothLEAdvertisementWatcher2 struct {
	ole.IInspectable
}

type iBluetoothLEAdvertisementWatcher2Vtbl struct {
	ole.IInspectableVtbl

	GetAllowExtendedAdvertisements uintptr
	SetAllowExtendedAdvertisements uintptr
}

func (v *iBluetoothLEAdvertisementWatcher2) VTable() *iBluetoothLEAdvertisementWatcher2Vtbl {
	return (*iBluetoothLEAdvertisementWatcher2Vtbl)(unsafe.Pointer(v.RawVTable))
}
