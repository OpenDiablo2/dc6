package giuwidget

import (
	"image"

	"github.com/AllenDang/giu"
)

func (fv *FrameViewerDC6) initState() {
	state := &frameViewerState{
		scale: 1,
	}

	fv.setState(state)

	numDirections := fv.dc6.Frames.NumberOfDirections()
	numFrames := fv.dc6.Frames.FramesPerDirection()
	totalFrames := numDirections * numFrames
	state.images = make([]*image.RGBA, totalFrames)

	for dirIdx := 0; dirIdx < fv.dc6.Frames.NumberOfDirections(); dirIdx++ {
		for frameIdx := 0; frameIdx < fv.dc6.Frames.FramesPerDirection(); frameIdx++ {
			fw := int(fv.dc6.Frames.Direction(dirIdx).Frame(frameIdx).Width)
			fh := int(fv.dc6.Frames.Direction(dirIdx).Frame(frameIdx).Height)

			absoluteFrameIdx := (dirIdx * numFrames) + frameIdx

			frame := fv.dc6.Frames.Direction(dirIdx).Frame(frameIdx)
			pixels := frame.FrameData

			if state.images[absoluteFrameIdx] == nil {
				state.images[absoluteFrameIdx] = frame.ToImageRGBA()
			}

			state.images[absoluteFrameIdx] = image.NewRGBA(image.Rect(0, 0, fw, fh))

			for y := 0; y < fh; y++ {
				for x := 0; x < fw; x++ {
					idx := x + (y * fw)
					if idx >= len(pixels) {
						continue
					}

					paletteIndex := pixels[idx]

					RGBAColor := fv.dc6.Palette()[paletteIndex]
					state.images[absoluteFrameIdx].Set(x, y, RGBAColor)
				}
			}
		}
	}

	go func() {
		textures := make([]*giu.Texture, totalFrames)

		for frameIndex := 0; frameIndex < totalFrames; frameIndex++ {
			frameIndex := frameIndex
			fv.textureLoader.CreateTextureFromARGB(state.images[frameIndex], func(t *giu.Texture) {
				textures[frameIndex] = t
			})
		}

		s := fv.getState()
		s.textures = textures
		fv.setState(s)
	}()
}
