package frames

type scanlineState int

const (
	endOfLine scanlineState = iota
	runOfTransparentPixels
	runOfOpaquePixels

	endOfScanLine = 0x80
	maxRunLength  = 0x7f
)

// decodeFrame decodes the given frame to an indexed color texture
func (f *Frame) decodeFrame() {
	indexData := make([]byte, f.Width*f.Height)
	x := 0
	y := int(f.Height) - 1
	offset := 0

loop: // this is a label for the loop, so the switch can break the loop (and not the switch)
	for {
		b := int(f.FrameData[offset])
		offset++

		switch scanlineType(b) {
		case endOfLine:
			if y == 0 {
				break loop
			}

			y--

			x = 0
		case runOfTransparentPixels:
			transparentPixels := b & maxRunLength
			x += transparentPixels
		case runOfOpaquePixels:
			for i := 0; i < b; i++ {
				indexData[x+y*int(f.Width)+i] = f.FrameData[offset]
				offset++
			}

			x += b
		}
	}

	f.IndexData = indexData
}

func scanlineType(b int) scanlineState {
	if b == endOfScanLine {
		return endOfLine
	}

	if (b & endOfScanLine) > 0 {
		return runOfTransparentPixels
	}

	return runOfOpaquePixels
}
