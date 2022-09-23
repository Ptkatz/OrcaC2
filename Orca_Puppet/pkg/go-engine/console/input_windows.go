package console

import (
	"errors"
	"Orca_Puppet/pkg/go-engine/common"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

var k32 = syscall.NewLazyDLL("kernel32.dll")

// We have to bring in the kernel32.dll directly, so we can get access to some
// system calls that the core Go API lacks.
//
// Note that Windows appends some functions with W to indicate that wide
// characters (Unicode) are in use.  The documentation refers to them
// without this suffix, as the resolution is made via preprocessor.
var (
	procReadConsoleInput           = k32.NewProc("ReadConsoleInputW")
	procWaitForMultipleObjects     = k32.NewProc("WaitForMultipleObjects")
	procCreateEvent                = k32.NewProc("CreateEventW")
	procSetEvent                   = k32.NewProc("SetEvent")
	procGetConsoleCursorInfo       = k32.NewProc("GetConsoleCursorInfo")
	procSetConsoleCursorInfo       = k32.NewProc("SetConsoleCursorInfo")
	procSetConsoleCursorPosition   = k32.NewProc("SetConsoleCursorPosition")
	procSetConsoleMode             = k32.NewProc("SetConsoleMode")
	procGetConsoleMode             = k32.NewProc("GetConsoleMode")
	procGetConsoleScreenBufferInfo = k32.NewProc("GetConsoleScreenBufferInfo")
	procFillConsoleOutputAttribute = k32.NewProc("FillConsoleOutputAttribute")
	procFillConsoleOutputCharacter = k32.NewProc("FillConsoleOutputCharacterW")
	procSetConsoleWindowInfo       = k32.NewProc("SetConsoleWindowInfo")
	procSetConsoleScreenBufferSize = k32.NewProc("SetConsoleScreenBufferSize")
	procSetConsoleTextAttribute    = k32.NewProc("SetConsoleTextAttribute")
)

const (
	w32Infinite    = ^uintptr(0)
	w32WaitObject0 = uintptr(0)
)

type ConsoleInput struct {
	in             syscall.Handle
	cancelflag     syscall.Handle
	evch           chan *EventKey
	exit           bool
	workResultLock sync.WaitGroup
}

func NewConsoleInput() *ConsoleInput {
	return &ConsoleInput{}
}

func (ci *ConsoleInput) Init() error {
	ci.evch = make(chan *EventKey, 10)

	in, err := syscall.Open("CONIN$", syscall.O_RDWR, 0)
	if err != nil {
		return err
	}
	ci.in = in

	cf, _, e := procCreateEvent.Call(
		uintptr(0),
		uintptr(1),
		uintptr(0),
		uintptr(0))
	if cf == uintptr(0) {
		return e
	}
	ci.cancelflag = syscall.Handle(cf)

	go ci.scanInput()

	return nil
}

func (ci *ConsoleInput) Stop() {
	ci.exit = true
	procSetEvent.Call(uintptr(ci.cancelflag))
	ci.workResultLock.Wait()
	syscall.Close(ci.in)
}

type inputRecord struct {
	typ  uint16
	_    uint16
	data [16]byte
}

const (
	keyEvent    uint16 = 1
	mouseEvent  uint16 = 2  // don't use
	resizeEvent uint16 = 4  // don't use
	menuEvent   uint16 = 8  // don't use
	focusEvent  uint16 = 16 // don't use
)

type keyRecord struct {
	isdown int32
	repeat uint16
	kcode  uint16
	scode  uint16
	ch     uint16
}

// NB: All Windows platforms are little endian.  We assume this
// never, ever change.  The following code is endian safe. and does
// not use unsafe pointers.
func getu32(v []byte) uint32 {
	return uint32(v[0]) + (uint32(v[1]) << 8) + (uint32(v[2]) << 16) + (uint32(v[3]) << 24)
}
func geti32(v []byte) int32 {
	return int32(getu32(v))
}
func getu16(v []byte) uint16 {
	return uint16(v[0]) + (uint16(v[1]) << 8)
}
func geti16(v []byte) int16 {
	return int16(getu16(v))
}

const (
	// Constants per Microsoft.  We don't put the modifiers
	// here.
	vkCancel = 0x03
	vkBack   = 0x08 // Backspace
	vkTab    = 0x09
	vkClear  = 0x0c
	vkReturn = 0x0d
	vkPause  = 0x13
	vkEscape = 0x1b
	vkSpace  = 0x20
	vkPrior  = 0x21 // PgUp
	vkNext   = 0x22 // PgDn
	vkEnd    = 0x23
	vkHome   = 0x24
	vkLeft   = 0x25
	vkUp     = 0x26
	vkRight  = 0x27
	vkDown   = 0x28
	vkPrint  = 0x2a
	vkPrtScr = 0x2c
	vkInsert = 0x2d
	vkDelete = 0x2e
	vkHelp   = 0x2f
	vkF1     = 0x70
	vkF2     = 0x71
	vkF3     = 0x72
	vkF4     = 0x73
	vkF5     = 0x74
	vkF6     = 0x75
	vkF7     = 0x76
	vkF8     = 0x77
	vkF9     = 0x78
	vkF10    = 0x79
	vkF11    = 0x7a
	vkF12    = 0x7b
	vkF13    = 0x7c
	vkF14    = 0x7d
	vkF15    = 0x7e
	vkF16    = 0x7f
	vkF17    = 0x80
	vkF18    = 0x81
	vkF19    = 0x82
	vkF20    = 0x83
	vkF21    = 0x84
	vkF22    = 0x85
	vkF23    = 0x86
	vkF24    = 0x87
)

var vkKeys = map[uint16]Key{
	vkCancel: KeyCancel,
	vkBack:   KeyBackspace,
	vkTab:    KeyTab,
	vkClear:  KeyClear,
	vkPause:  KeyPause,
	vkPrint:  KeyPrint,
	vkPrtScr: KeyPrint,
	vkPrior:  KeyPgUp,
	vkNext:   KeyPgDn,
	vkReturn: KeyEnter,
	vkEnd:    KeyEnd,
	vkHome:   KeyHome,
	vkLeft:   KeyLeft,
	vkUp:     KeyUp,
	vkRight:  KeyRight,
	vkDown:   KeyDown,
	vkInsert: KeyInsert,
	vkDelete: KeyDelete,
	vkHelp:   KeyHelp,
	vkF1:     KeyF1,
	vkF2:     KeyF2,
	vkF3:     KeyF3,
	vkF4:     KeyF4,
	vkF5:     KeyF5,
	vkF6:     KeyF6,
	vkF7:     KeyF7,
	vkF8:     KeyF8,
	vkF9:     KeyF9,
	vkF10:    KeyF10,
	vkF11:    KeyF11,
	vkF12:    KeyF12,
	vkF13:    KeyF13,
	vkF14:    KeyF14,
	vkF15:    KeyF15,
	vkF16:    KeyF16,
	vkF17:    KeyF17,
	vkF18:    KeyF18,
	vkF19:    KeyF19,
	vkF20:    KeyF20,
	vkF21:    KeyF21,
	vkF22:    KeyF22,
	vkF23:    KeyF23,
	vkF24:    KeyF24,
}

func (ci *ConsoleInput) scanInput() {
	defer common.CrashLog()

	ci.workResultLock.Add(1)
	defer ci.workResultLock.Done()

	for !ci.exit {
		ci.getConsoleInput()
	}
}

func (ci *ConsoleInput) getConsoleInput() error {
	// cancelFlag comes first as WaitForMultipleObjects returns the lowest index
	// in the event that both events are signalled.
	waitObjects := []syscall.Handle{ci.cancelflag, ci.in}
	// As arrays are contiguous in memory, a pointer to the first object is the
	// same as a pointer to the array itself.
	pWaitObjects := unsafe.Pointer(&waitObjects[0])

	rv, _, er := procWaitForMultipleObjects.Call(
		uintptr(len(waitObjects)),
		uintptr(pWaitObjects),
		uintptr(0),
		w32Infinite)
	// WaitForMultipleObjects returns WAIT_OBJECT_0 + the index.
	switch rv {
	case w32WaitObject0: // s.cancelFlag
		return errors.New("cancelled")
	case w32WaitObject0 + 1: // s.in
		rec := &inputRecord{}
		var nrec int32
		rv, _, er := procReadConsoleInput.Call(
			uintptr(ci.in),
			uintptr(unsafe.Pointer(rec)),
			uintptr(1),
			uintptr(unsafe.Pointer(&nrec)))
		if rv == 0 {
			return er
		}
		if nrec != 1 {
			return nil
		}
		switch rec.typ {
		case keyEvent:
			krec := &keyRecord{}
			krec.isdown = geti32(rec.data[0:])
			krec.repeat = getu16(rec.data[4:])
			krec.kcode = getu16(rec.data[6:])
			krec.scode = getu16(rec.data[8:])
			krec.ch = getu16(rec.data[10:])

			if krec.isdown == 0 || krec.repeat < 1 {
				// its a key release event, ignore it
				return nil
			}
			if krec.ch != 0 {
				// synthesized key code
				for krec.repeat > 0 {
					ci.postEvent(NewEventKey(KeyRune, rune(krec.ch)))
					krec.repeat--
				}
				return nil
			}
			key := KeyNUL // impossible on Windows
			ok := false
			if key, ok = vkKeys[krec.kcode]; !ok {
				return nil
			}
			for krec.repeat > 0 {
				ci.postEvent(NewEventKey(key, rune(krec.ch)))
				krec.repeat--
			}

		default:
		}
	default:
		return er
	}

	return nil
}

func (ci *ConsoleInput) postEvent(ev *EventKey) error {
	select {
	case ci.evch <- ev:
		return nil
	default:
		return errors.New("event queue full")
	}
}

func (ci *ConsoleInput) PollEvent() *EventKey {
	select {
	case ev := <-ci.evch:
		return ev
	case <-time.After(100 * time.Millisecond):
		return nil
	}
}
