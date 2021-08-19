package giuwidget

import (
	"fmt"
	"image"
	"log"

	"github.com/AllenDang/giu"

	dc6 "github.com/OpenDiablo2/dc6/pkg"
)

func FrameViewer(id string, d *dc6.DC6) *FrameViewerDC6 {
	return &FrameViewerDC6{
		id:            id,
		dc6:           d,
		textureLoader: newTextureLoader(),
	}
}

var _ giu.Widget = &FrameViewerDC6{}

type frameViewerState struct {
	direction, frame int
	scale            float64
	images           []*image.RGBA
	textures         []*giu.Texture
}

func (fvs *frameViewerState) Dispose() {
	// noop
}

type FrameViewerDC6 struct {
	id            string
	dc6           *dc6.DC6
	state         *frameViewerState
	textureLoader TextureLoader
}

func (p *FrameViewerDC6) Build() {
	const (
		imageW, imageH = 10, 10
	)

	p.textureLoader.ResumeLoadingTextures()
	p.textureLoader.ProcessTextureLoadRequests()

	viewerState := p.getState()

	imageScale := viewerState.scale

	dirIdx := 0
	frameIdx := 0

	textureIdx := dirIdx*p.dc6.Frames.FramesPerDirection() + frameIdx

	err := giu.Context.GetRenderer().SetTextureMagFilter(giu.TextureFilterNearest)
	if err != nil {
		log.Print(err)
	}

	var frameImage *giu.ImageWidget

	if viewerState.textures == nil || len(viewerState.textures) <= int(frameIdx) || viewerState.textures[frameIdx] == nil {
		frameImage = giu.Image(nil).Size(imageW, imageH)
	} else {
		bw := p.dc6.Frames.Direction(dirIdx).Frame(frameIdx).Width
		bh := p.dc6.Frames.Direction(dirIdx).Frame(frameIdx).Height
		w := float32(float64(bw) * imageScale)
		h := float32(float64(bh) * imageScale)
		frameImage = giu.Image(viewerState.textures[textureIdx]).Size(w, h)
	}

	// numDirections := len(p.dc6.Directions)
	// numFrames := len(p.dc6.Directions[0].Frames)

	giu.Layout{frameImage}.Build()
}

//func (fv *FrameViewerDC6) Build() {
//	s := fv.getState()
//	if s == nil {
//		return
//	}
//
//	fv.state = s
//
//	absIndex := (len(fv.dc6.Directions) * s.direction) + s.frame
//
//	frame := fv.dc6.Directions[s.direction].Frames[s.frame]
//	w, h := float32(float64(frame.Width)*s.scale), float32(float64(frame.Height)*s.scale)
//
//	layout := giu.Layout{
//		giu.Custom(func() {
//			if s.textures == nil {
//				return
//			}
//
//			if len(s.textures) <= absIndex {
//				return
//			}
//
//			giu.Image(s.textures[absIndex]).Size(w, h)
//		}),
//		giu.Dummy(w, h),
//	}
//
//	fv.ResumeLoadingTextures()
//	fv.ProcessTextureLoadRequests()
//
//	layout.Build()
//}

func (fv *FrameViewerDC6) getStateID() string {
	return fmt.Sprintf("widget_%s", fv.id)
}

func (fv *FrameViewerDC6) getState() *frameViewerState {
	var state *frameViewerState

	s := giu.Context.GetState(fv.getStateID())

	if s != nil {
		state = s.(*frameViewerState)
	} else {
		fv.initState()
		state = fv.getState()
	}

	return state
}

func (fv *FrameViewerDC6) SetScale(scale float64) {
	s := fv.getState()

	if scale <= 0 {
		scale = 1.0
	}

	s.scale = scale

	fv.setState(s)
}

func (fv *FrameViewerDC6) setState(s giu.Disposable) {
	giu.Context.SetState(fv.getStateID(), s)
}

func dirLookup(dir, numDirs int) int {
	d4 := []int{0, 1, 2, 3}
	d8 := []int{0, 5, 1, 6, 2, 7, 3, 4}
	d16 := []int{0, 9, 5, 10, 1, 11, 6, 12, 2, 13, 7, 14, 3, 15, 4, 8}

	lookup := []int{0}

	switch numDirs {
	case 4:
		lookup = d4
	case 8:
		lookup = d8
	case 16:
		lookup = d16
	default:
		dir = 0
	}

	return lookup[dir]
}
