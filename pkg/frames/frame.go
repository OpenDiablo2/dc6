package frames

import (
	"fmt"
	"image"
	"image/color"

	"github.com/gravestench/bitstream"
)

const (
	bytesPerInt32  = 4
	terminatorSize = 3
)

func newFrame() *Frame {
	return &Frame{
		FrameData:  make([]byte, 0),
		Terminator: make([]byte, terminatorSize),
	}
}

// Frame represents a single frame in a DC6.
type Frame struct {
	palette *color.Palette

	FrameData  []byte
	Terminator []byte
	Flipped    uint32
	OffsetX    int32
	OffsetY    int32
	Unknown    uint32
	NextBlock  uint32
	Width      uint32
	Height     uint32
}

// Load loads frame data
func (f *Frame) Load(r *bitstream.Reader, palette *color.Palette) error {
	f.palette = palette

	var err error

	r.Next(bytesPerInt32) // set bytes len to uint32

	if f.Flipped, err = r.Bytes().AsUInt32(); err != nil {
		return fmt.Errorf("reading flipped: %w", err)
	}

	if f.Width, err = r.Bytes().AsUInt32(); err != nil {
		return fmt.Errorf("reading width: %w", err)
	}

	if f.Height, err = r.Bytes().AsUInt32(); err != nil {
		return fmt.Errorf("reading height: %w", err)
	}

	if f.OffsetX, err = r.Bytes().AsInt32(); err != nil {
		return fmt.Errorf("reading x-offset: %w", err)
	}

	if f.OffsetY, err = r.Bytes().AsInt32(); err != nil {
		return fmt.Errorf("reading y-offset: %w", err)
	}

	if f.Unknown, err = r.Bytes().AsUInt32(); err != nil {
		return fmt.Errorf("reading frame unknown: %w", err)
	}

	if f.NextBlock, err = r.Bytes().AsUInt32(); err != nil {
		return fmt.Errorf("reading next block: %w", err)
	}

	l, err := r.Bytes().AsUInt32()
	if err != nil {
		return fmt.Errorf("reading length of frame data: %w", err)
	}

	if f.FrameData, err = r.Next(int(l)).Bytes().AsBytes(); err != nil {
		return fmt.Errorf("reading frame data: %w", err)
	}

	if f.Terminator, err = r.Next(terminatorSize).Bytes().AsBytes(); err != nil {
		return fmt.Errorf("reading terminator: %w", err)
	}

	return nil
}

/* TODO: rewrite to use gravestench/bitstream
// Encode encodes frame data into a byte slice
func (f *Frame) Encode() []byte {
	sw := d2datautils.CreateStreamWriter()
	sw.PushUint32(f.Flipped)
	sw.PushUint32(f.Width)
	sw.PushUint32(f.Height)
	sw.PushInt32(f.OffsetX)
	sw.PushInt32(f.OffsetY)
	sw.PushUint32(f.Unknown)
	sw.PushUint32(f.NextBlock)
	sw.PushUint32(uint32(len(f.FrameData)))
	sw.PushBytes(f.FrameData...)
	sw.PushBytes(f.Terminator...)

	return sw.GetBytes()
}
*/

func (f *Frame) ColorIndexAt(x, y int) uint8 {
	idx := (y * int(f.Width)) + x

	return f.FrameData[idx]
}

func (f *Frame) ColorModel() color.Model {
	return color.RGBAModel
}

func (f *Frame) Bounds() image.Rectangle {
	origin := image.Point{X: int(f.OffsetX), Y: int(f.OffsetY)}
	delta := image.Point{X: int(f.Width), Y: int(f.Height)}

	return image.Rectangle{
		Min: origin,
		Max: origin.Add(delta),
	}
}

func (f *Frame) At(x, y int) color.Color {
	cidx := f.ColorIndexAt(x, y)

	return (*f.palette)[cidx]
}

func (f *Frame) ToImageRGBA() *image.RGBA {
	img := image.NewRGBA(image.Rectangle{
		Max: f.Bounds().Size(),
	})

	for py := 0; py < int(f.Height); py++ {
		for px := 0; px < int(f.Width); px++ {
			img.Set(px, py, f.At(px, py))
		}
	}

	return img
}
