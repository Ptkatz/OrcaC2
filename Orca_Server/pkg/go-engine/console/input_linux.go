package console

import (
	"bytes"
	"errors"
	"Orca_Server/pkg/go-engine/common"
	"golang.org/x/sys/unix"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// tKeyCode represents a combination of a key code and modifiers.
type tKeyCode struct {
	key Key
}

func NewConsoleInput() *ConsoleInput {
	return &ConsoleInput{}
}

type ConsoleInput struct {
	in             *os.File
	evch           chan *EventKey
	exit           bool
	workResultLock sync.WaitGroup
	keyexist       map[Key]bool
	keycodes       map[string]*tKeyCode
	ti             *Terminfo
	tiosp          *termiosPrivate
}

type termiosPrivate struct {
	tio *unix.Termios
}

func (ci *ConsoleInput) Init() error {
	ci.evch = make(chan *EventKey, 10)

	var err error

	if ci.in, err = os.OpenFile("/dev/tty", os.O_RDONLY, 0); err != nil {
		return err
	}

	tio, err := unix.IoctlGetTermios(int(ci.in.Fd()), unix.TCGETS)
	if err != nil {
		ci.in.Close()
		return err
	}

	ci.tiosp = &termiosPrivate{tio: tio}

	// make a local copy, to make it raw
	raw := &unix.Termios{
		Cflag: tio.Cflag,
		Oflag: tio.Oflag,
		Iflag: tio.Iflag,
		Lflag: tio.Lflag,
		Cc:    tio.Cc,
	}
	raw.Iflag &^= (unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP |
		unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON)
	raw.Lflag &^= (unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG |
		unix.IEXTEN)
	raw.Cflag &^= (unix.CSIZE | unix.PARENB)
	raw.Cflag |= unix.CS8

	err = unix.IoctlSetTermios(int(ci.in.Fd()), unix.TCSETS, raw)
	if err != nil {
		ci.in.Close()
		return err
	}

	ti, e := loadDynamicTerminfo(os.Getenv("TERM"))
	if e != nil {
		ci.in.Close()
		return err
	}
	ci.ti = ti

	ci.keyexist = make(map[Key]bool)
	ci.keycodes = make(map[string]*tKeyCode)
	ci.prepareKeys()

	go ci.inputLoop()
	return nil
}

func fcntl(fd int, cmd int, arg int) (val int, err error) {
	r, _, e := syscall.Syscall(syscall.SYS_FCNTL, uintptr(fd), uintptr(cmd),
		uintptr(arg))
	val = int(r)
	if e != 0 {
		err = e
	}
	return
}

func (ci *ConsoleInput) Stop() {
	ci.exit = true
	ci.workResultLock.Wait()
	ci.in.Close()
}

func (ci *ConsoleInput) inputLoop() {
	defer common.CrashLog()

	ci.workResultLock.Add(1)
	defer ci.workResultLock.Done()

	buf := &bytes.Buffer{}
	for !ci.exit {
		chunk := make([]byte, 128)
		n, e := ci.in.Read(chunk)
		switch e {
		case io.EOF:
		case nil:
		default:
			return
		}

		buf.Write(chunk[:n])

		ci.scanInput(buf)
	}
}

func (ci *ConsoleInput) scanInput(buf *bytes.Buffer) {
	evs := ci.collectEventsFromInput(buf)

	for _, ev := range evs {
		ci.postEvent(ev)
	}
}

// Return an array of Events extracted from the supplied buffer. This is done
// while holding the screen's lock - the events can then be queued for
// application processing with the lock released.
func (ci *ConsoleInput) collectEventsFromInput(buf *bytes.Buffer) []*EventKey {

	res := make([]*EventKey, 0, 20)

	for {
		b := buf.Bytes()
		if len(b) == 0 {
			buf.Reset()
			return res
		}

		partials := 0

		if part, comp := ci.parseRune(buf, &res); comp {
			continue
		} else if part {
			partials++
		}

		if part, comp := ci.parseFunctionKey(buf, &res); comp {
			continue
		} else if part {
			partials++
		}

		if partials == 0 {
			if b[0] == '\x1b' {
				if len(b) == 1 {
					res = append(res, NewEventKey(KeyEsc, 0))
				}
				buf.ReadByte()
				continue
			}
			// Nothing was going to match, or we timed out
			// waiting for more data -- just deliver the characters
			// to the app & let them sort it out.  Possibly we
			// should only do this for control characters like ESC.
			by, _ := buf.ReadByte()
			res = append(res, NewEventKey(KeyRune, rune(by)))
			continue
		}

		// well we have some partial data, wait until we get
		// some more
		break
	}

	return res
}

func (ci *ConsoleInput) parseFunctionKey(buf *bytes.Buffer, evs *[]*EventKey) (bool, bool) {
	b := buf.Bytes()
	partial := false
	for e, k := range ci.keycodes {
		esc := []byte(e)
		if (len(esc) == 1) && (esc[0] == '\x1b') {
			continue
		}
		if bytes.HasPrefix(b, esc) {
			// matched
			var r rune
			if len(esc) == 1 {
				r = rune(b[0])
			}
			*evs = append(*evs, NewEventKey(k.key, r))
			for i := 0; i < len(esc); i++ {
				buf.ReadByte()
			}
			return true, true
		}
		if bytes.HasPrefix(esc, b) {
			partial = true
		}
	}
	return partial, false
}

func (ci *ConsoleInput) parseRune(buf *bytes.Buffer, evs *[]*EventKey) (bool, bool) {
	b := buf.Bytes()
	if b[0] >= ' ' && b[0] <= 0x7F {
		// printable ASCII easy to deal with -- no encodings
		*evs = append(*evs, NewEventKey(KeyRune, rune(b[0])))
		buf.ReadByte()
		return true, true
	}

	if b[0] < 0x80 {
		// Low numbered values are control keys, not runes.
		return false, false
	}

	// Looks like potential escape
	return true, false
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

func (ci *ConsoleInput) prepareKeyMod(key Key, val string) {
	if val != "" {
		// Do not overrride codes that already exist
		if _, exist := ci.keycodes[val]; !exist {
			ci.keyexist[key] = true
			ci.keycodes[val] = &tKeyCode{key: key}
		}
	}
}

func (ci *ConsoleInput) prepareKey(key Key, val string) {
	ci.prepareKeyMod(key, val)
}

func (ci *ConsoleInput) prepareKeys() {
	ti := ci.ti
	ci.prepareKey(KeyBackspace, ti.KeyBackspace)
	ci.prepareKey(KeyF1, ti.KeyF1)
	ci.prepareKey(KeyF2, ti.KeyF2)
	ci.prepareKey(KeyF3, ti.KeyF3)
	ci.prepareKey(KeyF4, ti.KeyF4)
	ci.prepareKey(KeyF5, ti.KeyF5)
	ci.prepareKey(KeyF6, ti.KeyF6)
	ci.prepareKey(KeyF7, ti.KeyF7)
	ci.prepareKey(KeyF8, ti.KeyF8)
	ci.prepareKey(KeyF9, ti.KeyF9)
	ci.prepareKey(KeyF10, ti.KeyF10)
	ci.prepareKey(KeyF11, ti.KeyF11)
	ci.prepareKey(KeyF12, ti.KeyF12)
	ci.prepareKey(KeyF13, ti.KeyF13)
	ci.prepareKey(KeyF14, ti.KeyF14)
	ci.prepareKey(KeyF15, ti.KeyF15)
	ci.prepareKey(KeyF16, ti.KeyF16)
	ci.prepareKey(KeyF17, ti.KeyF17)
	ci.prepareKey(KeyF18, ti.KeyF18)
	ci.prepareKey(KeyF19, ti.KeyF19)
	ci.prepareKey(KeyF20, ti.KeyF20)
	ci.prepareKey(KeyF21, ti.KeyF21)
	ci.prepareKey(KeyF22, ti.KeyF22)
	ci.prepareKey(KeyF23, ti.KeyF23)
	ci.prepareKey(KeyF24, ti.KeyF24)
	ci.prepareKey(KeyF25, ti.KeyF25)
	ci.prepareKey(KeyF26, ti.KeyF26)
	ci.prepareKey(KeyF27, ti.KeyF27)
	ci.prepareKey(KeyF28, ti.KeyF28)
	ci.prepareKey(KeyF29, ti.KeyF29)
	ci.prepareKey(KeyF30, ti.KeyF30)
	ci.prepareKey(KeyF31, ti.KeyF31)
	ci.prepareKey(KeyF32, ti.KeyF32)
	ci.prepareKey(KeyF33, ti.KeyF33)
	ci.prepareKey(KeyF34, ti.KeyF34)
	ci.prepareKey(KeyF35, ti.KeyF35)
	ci.prepareKey(KeyF36, ti.KeyF36)
	ci.prepareKey(KeyF37, ti.KeyF37)
	ci.prepareKey(KeyF38, ti.KeyF38)
	ci.prepareKey(KeyF39, ti.KeyF39)
	ci.prepareKey(KeyF40, ti.KeyF40)
	ci.prepareKey(KeyF41, ti.KeyF41)
	ci.prepareKey(KeyF42, ti.KeyF42)
	ci.prepareKey(KeyF43, ti.KeyF43)
	ci.prepareKey(KeyF44, ti.KeyF44)
	ci.prepareKey(KeyF45, ti.KeyF45)
	ci.prepareKey(KeyF46, ti.KeyF46)
	ci.prepareKey(KeyF47, ti.KeyF47)
	ci.prepareKey(KeyF48, ti.KeyF48)
	ci.prepareKey(KeyF49, ti.KeyF49)
	ci.prepareKey(KeyF50, ti.KeyF50)
	ci.prepareKey(KeyF51, ti.KeyF51)
	ci.prepareKey(KeyF52, ti.KeyF52)
	ci.prepareKey(KeyF53, ti.KeyF53)
	ci.prepareKey(KeyF54, ti.KeyF54)
	ci.prepareKey(KeyF55, ti.KeyF55)
	ci.prepareKey(KeyF56, ti.KeyF56)
	ci.prepareKey(KeyF57, ti.KeyF57)
	ci.prepareKey(KeyF58, ti.KeyF58)
	ci.prepareKey(KeyF59, ti.KeyF59)
	ci.prepareKey(KeyF60, ti.KeyF60)
	ci.prepareKey(KeyF61, ti.KeyF61)
	ci.prepareKey(KeyF62, ti.KeyF62)
	ci.prepareKey(KeyF63, ti.KeyF63)
	ci.prepareKey(KeyF64, ti.KeyF64)
	ci.prepareKey(KeyInsert, ti.KeyInsert)
	ci.prepareKey(KeyDelete, ti.KeyDelete)
	ci.prepareKey(KeyHome, ti.KeyHome)
	ci.prepareKey(KeyEnd, ti.KeyEnd)
	ci.prepareKey(KeyUp, ti.KeyUp)
	ci.prepareKey(KeyDown, ti.KeyDown)
	ci.prepareKey(KeyLeft, ti.KeyLeft)
	ci.prepareKey(KeyRight, ti.KeyRight)
	ci.prepareKey(KeyPgUp, ti.KeyPgUp)
	ci.prepareKey(KeyPgDn, ti.KeyPgDn)
	ci.prepareKey(KeyHelp, ti.KeyHelp)
	ci.prepareKey(KeyPrint, ti.KeyPrint)
	ci.prepareKey(KeyCancel, ti.KeyCancel)
	ci.prepareKey(KeyExit, ti.KeyExit)
	ci.prepareKey(KeyBacktab, ti.KeyBacktab)

	// Sadly, xterm handling of keycodes is somewhat erratic.  In
	// particular, different codes are sent depending on application
	// mode is in use or not, and the entries for many of these are
	// simply absent from terminfo on many systems.  So we insert
	// a number of escape sequences if they are not already used, in
	// order to have the widest correct usage.  Note that prepareKey
	// will not inject codes if the escape sequence is already known.
	// We also only do this for terminals that have the application
	// mode present.

	// Cursor mode
	if ti.EnterKeypad != "" {
		ci.prepareKey(KeyUp, "\x1b[A")
		ci.prepareKey(KeyDown, "\x1b[B")
		ci.prepareKey(KeyRight, "\x1b[C")
		ci.prepareKey(KeyLeft, "\x1b[D")
		ci.prepareKey(KeyEnd, "\x1b[F")
		ci.prepareKey(KeyHome, "\x1b[H")
		ci.prepareKey(KeyDelete, "\x1b[3~")
		ci.prepareKey(KeyHome, "\x1b[1~")
		ci.prepareKey(KeyEnd, "\x1b[4~")
		ci.prepareKey(KeyPgUp, "\x1b[5~")
		ci.prepareKey(KeyPgDn, "\x1b[6~")

		// Application mode
		ci.prepareKey(KeyUp, "\x1bOA")
		ci.prepareKey(KeyDown, "\x1bOB")
		ci.prepareKey(KeyRight, "\x1bOC")
		ci.prepareKey(KeyLeft, "\x1bOD")
		ci.prepareKey(KeyHome, "\x1bOH")
	}

outer:
	// Add key mappings for control keys.
	for i := 0; i < ' '; i++ {
		// Do not insert direct key codes for ambiguous keys.
		// For example, ESC is used for lots of other keys, so
		// when parsing this we don't want to fast path handling
		// of it, but instead wait a bit before parsing it as in
		// isolation.
		for esc := range ci.keycodes {
			if []byte(esc)[0] == byte(i) {
				continue outer
			}
		}

		ci.keyexist[Key(i)] = true

		ci.keycodes[string(rune(i))] = &tKeyCode{key: Key(i)}
	}
}

// Terminfo represents a terminfo entry.  Note that we use friendly names
// in Go, but when we write out JSON, we use the same names as terminfo.
// The name, aliases and smous, rmous fields do not come from terminfo directly.
type Terminfo struct {
	Name         string
	Aliases      []string
	Columns      int    // cols
	Lines        int    // lines
	Colors       int    // colors
	Bell         string // bell
	Clear        string // clear
	EnterCA      string // smcup
	ExitCA       string // rmcup
	ShowCursor   string // cnorm
	HideCursor   string // civis
	AttrOff      string // sgr0
	Underline    string // smul
	Bold         string // bold
	Blink        string // blink
	Reverse      string // rev
	Dim          string // dim
	EnterKeypad  string // smkx
	ExitKeypad   string // rmkx
	SetFg        string // setaf
	SetBg        string // setab
	SetCursor    string // cup
	CursorBack1  string // cub1
	CursorUp1    string // cuu1
	PadChar      string // pad
	KeyBackspace string // kbs
	KeyF1        string // kf1
	KeyF2        string // kf2
	KeyF3        string // kf3
	KeyF4        string // kf4
	KeyF5        string // kf5
	KeyF6        string // kf6
	KeyF7        string // kf7
	KeyF8        string // kf8
	KeyF9        string // kf9
	KeyF10       string // kf10
	KeyF11       string // kf11
	KeyF12       string // kf12
	KeyF13       string // kf13
	KeyF14       string // kf14
	KeyF15       string // kf15
	KeyF16       string // kf16
	KeyF17       string // kf17
	KeyF18       string // kf18
	KeyF19       string // kf19
	KeyF20       string // kf20
	KeyF21       string // kf21
	KeyF22       string // kf22
	KeyF23       string // kf23
	KeyF24       string // kf24
	KeyF25       string // kf25
	KeyF26       string // kf26
	KeyF27       string // kf27
	KeyF28       string // kf28
	KeyF29       string // kf29
	KeyF30       string // kf30
	KeyF31       string // kf31
	KeyF32       string // kf32
	KeyF33       string // kf33
	KeyF34       string // kf34
	KeyF35       string // kf35
	KeyF36       string // kf36
	KeyF37       string // kf37
	KeyF38       string // kf38
	KeyF39       string // kf39
	KeyF40       string // kf40
	KeyF41       string // kf41
	KeyF42       string // kf42
	KeyF43       string // kf43
	KeyF44       string // kf44
	KeyF45       string // kf45
	KeyF46       string // kf46
	KeyF47       string // kf47
	KeyF48       string // kf48
	KeyF49       string // kf49
	KeyF50       string // kf50
	KeyF51       string // kf51
	KeyF52       string // kf52
	KeyF53       string // kf53
	KeyF54       string // kf54
	KeyF55       string // kf55
	KeyF56       string // kf56
	KeyF57       string // kf57
	KeyF58       string // kf58
	KeyF59       string // kf59
	KeyF60       string // kf60
	KeyF61       string // kf61
	KeyF62       string // kf62
	KeyF63       string // kf63
	KeyF64       string // kf64
	KeyInsert    string // kich1
	KeyDelete    string // kdch1
	KeyHome      string // khome
	KeyEnd       string // kend
	KeyHelp      string // khlp
	KeyPgUp      string // kpp
	KeyPgDn      string // knp
	KeyUp        string // kcuu1
	KeyDown      string // kcud1
	KeyLeft      string // kcub1
	KeyRight     string // kcuf1
	KeyBacktab   string // kcbt
	KeyExit      string // kext
	KeyClear     string // kclr
	KeyPrint     string // kprt
	KeyCancel    string // kcan
	Mouse        string // kmous
	MouseMode    string // XM
	AltChars     string // acsc
	EnterAcs     string // smacs
	ExitAcs      string // rmacs
	EnableAcs    string // enacs
	KeyShfRight  string // kRIT
	KeyShfLeft   string // kLFT
	KeyShfHome   string // kHOM
	KeyShfEnd    string // kEND

	// These are non-standard extensions to terminfo.  This includes
	// true color support, and some additional keys.  Its kind of bizarre
	// that shifted variants of left and right exist, but not up and down.
	// Terminal support for these are going to vary amongst XTerm
	// emulations, so don't depend too much on them in your application.

	SetFgBg         string // setfgbg
	SetFgBgRGB      string // setfgbgrgb
	SetFgRGB        string // setfrgb
	SetBgRGB        string // setbrgb
	KeyShfUp        string // shift-up
	KeyShfDown      string // shift-down
	KeyCtrlUp       string // ctrl-up
	KeyCtrlDown     string // ctrl-left
	KeyCtrlRight    string // ctrl-right
	KeyCtrlLeft     string // ctrl-left
	KeyMetaUp       string // meta-up
	KeyMetaDown     string // meta-left
	KeyMetaRight    string // meta-right
	KeyMetaLeft     string // meta-left
	KeyAltUp        string // alt-up
	KeyAltDown      string // alt-left
	KeyAltRight     string // alt-right
	KeyAltLeft      string // alt-left
	KeyCtrlHome     string
	KeyCtrlEnd      string
	KeyMetaHome     string
	KeyMetaEnd      string
	KeyAltHome      string
	KeyAltEnd       string
	KeyAltShfUp     string
	KeyAltShfDown   string
	KeyAltShfLeft   string
	KeyAltShfRight  string
	KeyMetaShfUp    string
	KeyMetaShfDown  string
	KeyMetaShfLeft  string
	KeyMetaShfRight string
	KeyCtrlShfUp    string
	KeyCtrlShfDown  string
	KeyCtrlShfLeft  string
	KeyCtrlShfRight string
	KeyCtrlShfHome  string
	KeyCtrlShfEnd   string
	KeyAltShfHome   string
	KeyAltShfEnd    string
	KeyMetaShfHome  string
	KeyMetaShfEnd   string
}

func loadDynamicTerminfo(term string) (*Terminfo, error) {
	ti, _, e := LoadTerminfo(term)
	if e != nil {
		return nil, e
	}
	return ti, nil
}

type termcap struct {
	name    string
	desc    string
	aliases []string
	bools   map[string]bool
	nums    map[string]int
	strs    map[string]string
}

func (tc *termcap) getnum(s string) int {
	return (tc.nums[s])
}

func (tc *termcap) getflag(s string) bool {
	return (tc.bools[s])
}

func (tc *termcap) getstr(s string) string {
	return (tc.strs[s])
}

func (tc *termcap) setupterm(name string) error {
	cmd := exec.Command("infocmp", "-1", name)
	output := &bytes.Buffer{}
	cmd.Stdout = output

	tc.strs = make(map[string]string)
	tc.bools = make(map[string]bool)
	tc.nums = make(map[string]int)

	if err := cmd.Run(); err != nil {
		return err
	}

	// Now parse the output.
	// We get comment lines (starting with "#"), followed by
	// a header line that looks like "<name>|<alias>|...|<desc>"
	// then capabilities, one per line, starting with a tab and ending
	// with a comma and newline.
	lines := strings.Split(output.String(), "\n")
	for len(lines) > 0 && strings.HasPrefix(lines[0], "#") {
		lines = lines[1:]
	}

	// Ditch trailing empty last line
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	header := lines[0]
	if strings.HasSuffix(header, ",") {
		header = header[:len(header)-1]
	}
	names := strings.Split(header, "|")
	tc.name = names[0]
	names = names[1:]
	if len(names) > 0 {
		tc.desc = names[len(names)-1]
		names = names[:len(names)-1]
	}
	tc.aliases = names
	for _, val := range lines[1:] {
		if (!strings.HasPrefix(val, "\t")) ||
			(!strings.HasSuffix(val, ",")) {
			return (errors.New("malformed infocmp: " + val))
		}

		val = val[1:]
		val = val[:len(val)-1]

		if k := strings.SplitN(val, "=", 2); len(k) == 2 {
			tc.strs[k[0]] = unescape(k[1])
		} else if k := strings.SplitN(val, "#", 2); len(k) == 2 {
			u, err := strconv.ParseUint(k[1], 0, 0)
			if err != nil {
				return (err)
			}
			tc.nums[k[0]] = int(u)
		} else {
			tc.bools[val] = true
		}
	}
	return nil
}

const (
	none = iota
	control
	escaped
)

func unescape(s string) string {
	// Various escapes are in \x format.  Control codes are
	// encoded as ^M (carat followed by ASCII equivalent).
	// escapes are: \e, \E - escape
	//  \0 NULL, \n \l \r \t \b \f \s for equivalent C escape.
	buf := &bytes.Buffer{}
	esc := none

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch esc {
		case none:
			switch c {
			case '\\':
				esc = escaped
			case '^':
				esc = control
			default:
				buf.WriteByte(c)
			}
		case control:
			buf.WriteByte(c - 0x40)
			esc = none
		case escaped:
			switch c {
			case 'E', 'e':
				buf.WriteByte(0x1b)
			case '0', '1', '2', '3', '4', '5', '6', '7':
				if i+2 < len(s) && s[i+1] >= '0' && s[i+1] <= '7' && s[i+2] >= '0' && s[i+2] <= '7' {
					buf.WriteByte(((c - '0') * 64) + ((s[i+1] - '0') * 8) + (s[i+2] - '0'))
					i = i + 2
				} else if c == '0' {
					buf.WriteByte(0)
				}
			case 'n':
				buf.WriteByte('\n')
			case 'r':
				buf.WriteByte('\r')
			case 't':
				buf.WriteByte('\t')
			case 'b':
				buf.WriteByte('\b')
			case 'f':
				buf.WriteByte('\f')
			case 's':
				buf.WriteByte(' ')
			default:
				buf.WriteByte(c)
			}
			esc = none
		}
	}
	return (buf.String())
}

// LoadTerminfo creates a Terminfo by for named terminal by attempting to parse
// the output from infocmp.  This returns the terminfo entry, a description of
// the terminal, and either nil or an error.
func LoadTerminfo(name string) (*Terminfo, string, error) {
	var tc termcap
	if err := tc.setupterm(name); err != nil {
		if err != nil {
			return nil, "", err
		}
	}
	t := &Terminfo{}
	// If this is an alias record, then just emit the alias
	t.Name = tc.name
	if t.Name != name {
		return t, "", nil
	}
	t.Aliases = tc.aliases
	t.Colors = tc.getnum("colors")
	t.Columns = tc.getnum("cols")
	t.Lines = tc.getnum("lines")
	t.Bell = tc.getstr("bel")
	t.Clear = tc.getstr("clear")
	t.EnterCA = tc.getstr("smcup")
	t.ExitCA = tc.getstr("rmcup")
	t.ShowCursor = tc.getstr("cnorm")
	t.HideCursor = tc.getstr("civis")
	t.AttrOff = tc.getstr("sgr0")
	t.Underline = tc.getstr("smul")
	t.Bold = tc.getstr("bold")
	t.Blink = tc.getstr("blink")
	t.Dim = tc.getstr("dim")
	t.Reverse = tc.getstr("rev")
	t.EnterKeypad = tc.getstr("smkx")
	t.ExitKeypad = tc.getstr("rmkx")
	t.SetFg = tc.getstr("setaf")
	t.SetBg = tc.getstr("setab")
	t.SetCursor = tc.getstr("cup")
	t.CursorBack1 = tc.getstr("cub1")
	t.CursorUp1 = tc.getstr("cuu1")
	t.KeyF1 = tc.getstr("kf1")
	t.KeyF2 = tc.getstr("kf2")
	t.KeyF3 = tc.getstr("kf3")
	t.KeyF4 = tc.getstr("kf4")
	t.KeyF5 = tc.getstr("kf5")
	t.KeyF6 = tc.getstr("kf6")
	t.KeyF7 = tc.getstr("kf7")
	t.KeyF8 = tc.getstr("kf8")
	t.KeyF9 = tc.getstr("kf9")
	t.KeyF10 = tc.getstr("kf10")
	t.KeyF11 = tc.getstr("kf11")
	t.KeyF12 = tc.getstr("kf12")
	t.KeyF13 = tc.getstr("kf13")
	t.KeyF14 = tc.getstr("kf14")
	t.KeyF15 = tc.getstr("kf15")
	t.KeyF16 = tc.getstr("kf16")
	t.KeyF17 = tc.getstr("kf17")
	t.KeyF18 = tc.getstr("kf18")
	t.KeyF19 = tc.getstr("kf19")
	t.KeyF20 = tc.getstr("kf20")
	t.KeyF21 = tc.getstr("kf21")
	t.KeyF22 = tc.getstr("kf22")
	t.KeyF23 = tc.getstr("kf23")
	t.KeyF24 = tc.getstr("kf24")
	t.KeyF25 = tc.getstr("kf25")
	t.KeyF26 = tc.getstr("kf26")
	t.KeyF27 = tc.getstr("kf27")
	t.KeyF28 = tc.getstr("kf28")
	t.KeyF29 = tc.getstr("kf29")
	t.KeyF30 = tc.getstr("kf30")
	t.KeyF31 = tc.getstr("kf31")
	t.KeyF32 = tc.getstr("kf32")
	t.KeyF33 = tc.getstr("kf33")
	t.KeyF34 = tc.getstr("kf34")
	t.KeyF35 = tc.getstr("kf35")
	t.KeyF36 = tc.getstr("kf36")
	t.KeyF37 = tc.getstr("kf37")
	t.KeyF38 = tc.getstr("kf38")
	t.KeyF39 = tc.getstr("kf39")
	t.KeyF40 = tc.getstr("kf40")
	t.KeyF41 = tc.getstr("kf41")
	t.KeyF42 = tc.getstr("kf42")
	t.KeyF43 = tc.getstr("kf43")
	t.KeyF44 = tc.getstr("kf44")
	t.KeyF45 = tc.getstr("kf45")
	t.KeyF46 = tc.getstr("kf46")
	t.KeyF47 = tc.getstr("kf47")
	t.KeyF48 = tc.getstr("kf48")
	t.KeyF49 = tc.getstr("kf49")
	t.KeyF50 = tc.getstr("kf50")
	t.KeyF51 = tc.getstr("kf51")
	t.KeyF52 = tc.getstr("kf52")
	t.KeyF53 = tc.getstr("kf53")
	t.KeyF54 = tc.getstr("kf54")
	t.KeyF55 = tc.getstr("kf55")
	t.KeyF56 = tc.getstr("kf56")
	t.KeyF57 = tc.getstr("kf57")
	t.KeyF58 = tc.getstr("kf58")
	t.KeyF59 = tc.getstr("kf59")
	t.KeyF60 = tc.getstr("kf60")
	t.KeyF61 = tc.getstr("kf61")
	t.KeyF62 = tc.getstr("kf62")
	t.KeyF63 = tc.getstr("kf63")
	t.KeyF64 = tc.getstr("kf64")
	t.KeyInsert = tc.getstr("kich1")
	t.KeyDelete = tc.getstr("kdch1")
	t.KeyBackspace = tc.getstr("kbs")
	t.KeyHome = tc.getstr("khome")
	t.KeyEnd = tc.getstr("kend")
	t.KeyUp = tc.getstr("kcuu1")
	t.KeyDown = tc.getstr("kcud1")
	t.KeyRight = tc.getstr("kcuf1")
	t.KeyLeft = tc.getstr("kcub1")
	t.KeyPgDn = tc.getstr("knp")
	t.KeyPgUp = tc.getstr("kpp")
	t.KeyBacktab = tc.getstr("kcbt")
	t.KeyExit = tc.getstr("kext")
	t.KeyCancel = tc.getstr("kcan")
	t.KeyPrint = tc.getstr("kprt")
	t.KeyHelp = tc.getstr("khlp")
	t.KeyClear = tc.getstr("kclr")
	t.AltChars = tc.getstr("acsc")
	t.EnterAcs = tc.getstr("smacs")
	t.ExitAcs = tc.getstr("rmacs")
	t.EnableAcs = tc.getstr("enacs")
	t.Mouse = tc.getstr("kmous")
	t.KeyShfRight = tc.getstr("kRIT")
	t.KeyShfLeft = tc.getstr("kLFT")
	t.KeyShfHome = tc.getstr("kHOM")
	t.KeyShfEnd = tc.getstr("kEND")

	// Terminfo lacks descriptions for a bunch of modified keys,
	// but modern XTerm and emulators often have them.  Let's add them,
	// if the shifted right and left arrows are defined.
	if t.KeyShfRight == "\x1b[1;2C" && t.KeyShfLeft == "\x1b[1;2D" {
		t.KeyShfUp = "\x1b[1;2A"
		t.KeyShfDown = "\x1b[1;2B"
		t.KeyMetaUp = "\x1b[1;9A"
		t.KeyMetaDown = "\x1b[1;9B"
		t.KeyMetaRight = "\x1b[1;9C"
		t.KeyMetaLeft = "\x1b[1;9D"
		t.KeyAltUp = "\x1b[1;3A"
		t.KeyAltDown = "\x1b[1;3B"
		t.KeyAltRight = "\x1b[1;3C"
		t.KeyAltLeft = "\x1b[1;3D"
		t.KeyCtrlUp = "\x1b[1;5A"
		t.KeyCtrlDown = "\x1b[1;5B"
		t.KeyCtrlRight = "\x1b[1;5C"
		t.KeyCtrlLeft = "\x1b[1;5D"
		t.KeyAltShfUp = "\x1b[1;4A"
		t.KeyAltShfDown = "\x1b[1;4B"
		t.KeyAltShfRight = "\x1b[1;4C"
		t.KeyAltShfLeft = "\x1b[1;4D"

		t.KeyMetaShfUp = "\x1b[1;10A"
		t.KeyMetaShfDown = "\x1b[1;10B"
		t.KeyMetaShfRight = "\x1b[1;10C"
		t.KeyMetaShfLeft = "\x1b[1;10D"

		t.KeyCtrlShfUp = "\x1b[1;6A"
		t.KeyCtrlShfDown = "\x1b[1;6B"
		t.KeyCtrlShfRight = "\x1b[1;6C"
		t.KeyCtrlShfLeft = "\x1b[1;6D"
	}
	// And also for Home and End
	if t.KeyShfHome == "\x1b[1;2H" && t.KeyShfEnd == "\x1b[1;2F" {
		t.KeyCtrlHome = "\x1b[1;5H"
		t.KeyCtrlEnd = "\x1b[1;5F"
		t.KeyAltHome = "\x1b[1;9H"
		t.KeyAltEnd = "\x1b[1;9F"
		t.KeyCtrlShfHome = "\x1b[1;6H"
		t.KeyCtrlShfEnd = "\x1b[1;6F"
		t.KeyAltShfHome = "\x1b[1;4H"
		t.KeyAltShfEnd = "\x1b[1;4F"
		t.KeyMetaShfHome = "\x1b[1;10H"
		t.KeyMetaShfEnd = "\x1b[1;10F"
	}

	// And the same thing for rxvt and workalikes (Eterm, aterm, etc.)
	// It seems that urxvt at least send escaped as ALT prefix for these,
	// although some places seem to indicate a separate ALT key sesquence.
	if t.KeyShfRight == "\x1b[c" && t.KeyShfLeft == "\x1b[d" {
		t.KeyShfUp = "\x1b[a"
		t.KeyShfDown = "\x1b[b"
		t.KeyCtrlUp = "\x1b[Oa"
		t.KeyCtrlDown = "\x1b[Ob"
		t.KeyCtrlRight = "\x1b[Oc"
		t.KeyCtrlLeft = "\x1b[Od"
	}
	if t.KeyShfHome == "\x1b[7$" && t.KeyShfEnd == "\x1b[8$" {
		t.KeyCtrlHome = "\x1b[7^"
		t.KeyCtrlEnd = "\x1b[8^"
	}

	// If the kmous entry is present, then we need to record the
	// the codes to enter and exit mouse mode.  Sadly, this is not
	// part of the terminfo databases anywhere that I've found, but
	// is an extension.  The escapedape codes are documented in the XTerm
	// manual, and all terminals that have kmous are expected to
	// use these same codes, unless explicitly configured otherwise
	// vi XM.  Note that in any event, we only known how to parse either
	// x11 or SGR mouse events -- if your terminal doesn't support one
	// of these two forms, you maybe out of luck.
	t.MouseMode = tc.getstr("XM")
	if t.Mouse != "" && t.MouseMode == "" {
		// we anticipate that all xterm mouse tracking compatible
		// terminals understand mouse tracking (1000), but we hope
		// that those that don't understand any-event tracking (1003)
		// will at least ignore it.  Likewise we hope that terminals
		// that don't understand SGR reporting (1006) just ignore it.
		t.MouseMode = "%?%p1%{1}%=%t%'h'%Pa%e%'l'%Pa%;" +
			"\x1b[?1000%ga%c\x1b[?1002%ga%c\x1b[?1003%ga%c\x1b[?1006%ga%c"
	}

	// We only support colors in ANSI 8 or 256 color mode.
	if t.Colors < 8 || t.SetFg == "" {
		t.Colors = 0
	}
	if t.SetCursor == "" {
		return nil, "", errors.New("terminal not cursor addressable")
	}

	// For padding, we lookup the pad char.  If that isn't present,
	// and npc is *not* set, then we assume a null byte.
	t.PadChar = tc.getstr("pad")
	if t.PadChar == "" {
		if !tc.getflag("npc") {
			t.PadChar = "\u0000"
		}
	}

	// For terminals that use "standard" SGR sequences, lets combine the
	// foreground and background together.
	if strings.HasPrefix(t.SetFg, "\x1b[") &&
		strings.HasPrefix(t.SetBg, "\x1b[") &&
		strings.HasSuffix(t.SetFg, "m") &&
		strings.HasSuffix(t.SetBg, "m") {
		fg := t.SetFg[:len(t.SetFg)-1]
		r := regexp.MustCompile("%p1")
		bg := r.ReplaceAllString(t.SetBg[2:], "%p2")
		t.SetFgBg = fg + ";" + bg
	}

	return t, tc.desc, nil
}
