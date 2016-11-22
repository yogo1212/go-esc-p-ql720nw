package ql720nw

import (
	"errors"
	"fmt"
	"net"
)

var boolOneZeroMap = map[bool]string {
	false: "0",
	true: "1",
}

type Ql720nw struct {
	conn net.Conn
}

func (q *Ql720nw) ReInit() (err error) {
	_, err = q.conn.Write([]byte("\x1B@"))
	return
}

func (q *Ql720nw) PageFeed() (err error) {
	_, err = q.conn.Write([]byte{12})
	return
}

func (q *Ql720nw) SetCut(b bool) (err error) {
	_, err = q.conn.Write([]byte("\x1BiC" + boolOneZeroMap[b]))
	return
}

func (q *Ql720nw) commandMode() (err error) {
	_, err = q.conn.Write([]byte("\x1Bia0"))
	return
}

type CodeSet byte

const (
	Standard CodeSet = 0x0
	Eastern          = 0x1
	Western          = 0x2
)

type Font byte

const (
	Brougham Font       = 0x0
	LetterGothicBold    = 0x1
	Brussels            = 0x2
	Helsinki            = 0x3
	SanDiego            = 0x4
	LetterGothicOutline = 0x9
	BrusselsOutline     = 0xA
	HelsinkiOutline     = 0xB
)

func(q *Ql720nw) SetFont(f Font) (err error) {
	_, err = q.conn.Write(append([]byte("\x1Bk"), byte(f)))
	return
}

func (q *Ql720nw) SetDefaultFont(f Font) (err error) {
	_, err = q.conn.Write(append([]byte("\x1BiXk210"), byte(f)))
	return
}

type CharacterStyle byte

const (
	None CharacterStyle = 0
	Bold                = 1
	Outline             = 2
	Shadow              = 3
	ShadowAndOutline    = 4
)

func (q *Ql720nw) SetCharacterStyle(cs CharacterStyle) (err error) {
	_, err = q.conn.Write(append([]byte("\x1bq"), byte(cs)))
	return
}

func (q *Ql720nw) SetDefaultCharacterStyle(cs CharacterStyle) (err error) {
	_, err = q.conn.Write(append([]byte("\x1BXQ210"), byte(cs)))
	return
}

func (q *Ql720nw) SetLandScape(b bool) (err error) {
	_, err = q.conn.Write([]byte("\x1BiL" + boolOneZeroMap[b]))
	return
}

func (q *Ql720nw) ApplyItalic() (err error) {
	_, err = q.conn.Write([]byte("\x1B4"))
	return
}

func (q *Ql720nw) CancelItalic() (err error) {
	_, err = q.conn.Write([]byte("\x1B5"))
	return
}

func (q *Ql720nw) ApplyBold() (err error) {
	_, err = q.conn.Write([]byte("\x1BE"))
	return
}

func (q *Ql720nw) CancelBold() (err error) {
	_, err = q.conn.Write([]byte("\x1BF"))
	return
}

func (q *Ql720nw) ApplyDoubleStrike() (err error) {
	_, err = q.conn.Write([]byte("\x1BG"))
	return
}

func (q *Ql720nw) CancelDoubleStrike() (err error) {
	_, err = q.conn.Write([]byte("\x1BH"))
	return
}

func (q *Ql720nw) SetProportional(b bool) (err error) {
	_, err = q.conn.Write([]byte("\x1Bp" + boolOneZeroMap[b]))
	return
}

func (q *Ql720nw) SetDoubleWidth(b bool) (err error) {
	_, err = q.conn.Write([]byte("\x1BW" + boolOneZeroMap[b]))
	return
}

func (q *Ql720nw) SetPageFormat(top, bottom uint16) (err error) {
	_, err = q.conn.Write(append([]byte("\x1B(c"), 4, 0, byte(top), byte(top >> 8), byte(bottom), byte(bottom >> 8)))
	return
}

// look up m in thq ql720nw esc/p spec
/*
 * func (q *Ql720nw) PrintBitmap(m byte, n uint16) (err error) {
 * 	_, err = q.conn.Write(append([]byte("\x1B(c"), 4, 0, byte(top), byte(top >> 8), byte(bottom), byte(bottom >> 8)))
 * 	return
 * }
*/

type QRCellSize byte
const (
	QR3Dot QRCellSize = 3
	QR4Dot            = 4
	QR5Dot            = 5
	QR6Dot            = 6
	QR8Dot            = 8
	QR10Dot           = 10
)

type QRSymbolType byte
const (
	QRModel1 QRSymbolType = 1
	QRModel2              = 2
	QRMicro               = 3
)

type QRDataInputMethod byte
const (
	QRMethodAuto QRDataInputMethod = 0
	QRMethodAlphaNum               = 'A'
	QRMethodKanji                  = 'K'
	QRMethodBinary                 = 'B'
)

type QRECLevel byte
const (
	QRECLevelL QRECLevel = 1
	QRECLevelM           = 2
	QRECLevelQ           = 3
	QRECLevelH           = 4
)

func (q *Ql720nw) QRCode(data []byte, cs QRCellSize, sym QRSymbolType, ec QRECLevel, dim QRDataInputMethod) (err error) {
	manualDataInput := byte(0)
	if dim != QRMethodAuto {
		manualDataInput = 1
	}
	buf := append([]byte("\x1BiQ"), byte(cs), byte(sym), 0, 0, 0, 0, byte(ec), manualDataInput)
	if manualDataInput == 1 {
		buf = append(buf, byte(dim))
	}
	if dim == QRMethodBinary {
		decs := fmt.Sprintf("%04d", len(data))
		buf = append(buf, []byte(decs)...)
	}
	buf = append(buf, data...)
	_, err = q.conn.Write(append(buf, '\\', '\\', '\\'))
	return
}

type Underlining byte

const (
	Cancel Underlining = 0
	Underlining1Dot    = 1
	Underlining2Dot    = 2
	Underlining3Dot    = 3
	Underlining4Dot    = 4
)

func (q *Ql720nw) SetUnderlining(u Underlining) (err error) {
	_, err = q.conn.Write(append([]byte("\x1B-"), byte(u)))
	return
}

func (q *Ql720nw) SetCharacterSpacing(dots byte) (err error) {
	if dots > 127 {
		err = errors.New("out of range")
	}
	_, err = q.conn.Write(append([]byte("\x1B "), dots))
	return
}

func (q *Ql720nw) SetCharacterSize(s byte) (err error) {
	// TODO verify sizes?
	_, err = q.conn.Write(append([]byte("\x1BX"), 0, 0, s))
	return
}


func (q *Ql720nw) SetLeftMargin(m byte) (err error) {
	_, err = q.conn.Write(append([]byte("\x1Bl"), m))
	return
}

func (q *Ql720nw) SetRightMargin(m byte) (err error) {
	_, err = q.conn.Write(append([]byte("\x1BQ"), m))
	return
}

type Alignment byte

const (
	AlignLeft Alignment = 0
	AlignCenter         = 1
	AlignRight          = 2
	AlignNothing        = 3
)

func (q *Ql720nw) SetAlignment(a Alignment) (err error) {
	_, err = q.conn.Write(append([]byte("\x1Ba"), byte(a)))
	return
}

func (q *Ql720nw) Close() (err error) {
	return q.conn.Close()
}

func New(host string) (q Ql720nw, err error) {
	q.conn, err = net.Dial("tcp", host + ":9100")
	if err != nil {
		return
	}

	// turns out to be a good idea to do
	_, err = q.conn.Write([]byte{0x1B})
	if err != nil {
		return
	}

	err = q.commandMode()
	if err != nil {
		return
	}

	err = q.ReInit()

	return
}