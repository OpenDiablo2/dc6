package giuwidget

import (
	"fmt"
	"image"
	"log"

	"github.com/AllenDang/giu"

	dc6 "github.com/OpenDiablo2/dc6/pkg"
)

// FrameViewer creates frame viewer
func FrameViewer(id string, d *dc6.DC6) *FrameViewerDC6 {
	return &FrameViewerDC6{
		id:            id,
		dc6:           d,
		textureLoader: newTextureLoader(),
	}
}

var _ giu.Widget = &FrameViewerDC6{}

type frameViewerState struct {
	images   []*image.RGBA
	textures []*giu.Texture

	frame     int32
	direction int32

	scale float64
}

func (fvs *frameViewerState) Dispose() {
	// noop
}

// FrameViewerDC6 represents a dc6 frame viewer
type FrameViewerDC6 struct {
	textureLoader TextureLoader
	dc6           *dc6.DC6
	id            string
}

// Build implements giu.Widget
func (fv *FrameViewerDC6) Build() {
	const (
		imageW, imageH = 10, 10
	)

	fv.textureLoader.ResumeLoadingTextures()
	fv.textureLoader.ProcessTextureLoadRequests()

	viewerState := fv.getState()

	imageScale := viewerState.scale

	dirIdx := int(viewerState.direction)
	frameIdx := int(viewerState.frame)

	textureIdx := dirIdx*fv.dc6.Frames.FramesPerDirection() + frameIdx

	err := giu.Context.GetRenderer().SetTextureMagFilter(giu.TextureFilterNearest)
	if err != nil {
		log.Print(err)
	}

	var frameImage *giu.ImageWidget

	if viewerState.textures == nil || len(viewerState.textures) <= frameIdx || viewerState.textures[frameIdx] == nil {
		frameImage = giu.Image(nil).Size(imageW, imageH)
	} else {
		bw := fv.dc6.Frames.Direction(dirIdx).Frame(frameIdx).Width
		bh := fv.dc6.Frames.Direction(dirIdx).Frame(frameIdx).Height
		w := float32(float64(bw) * imageScale)
		h := float32(float64(bh) * imageScale)
		frameImage = giu.Image(viewerState.textures[textureIdx]).Size(w, h)
	}

	numDirections := fv.dc6.Frames.NumberOfDirections()
	numFrames := fv.dc6.Frames.FramesPerDirection()

	giu.Layout{
		giu.Custom(func() {
			if numDirections > 1 {
				giu.SliderInt("direction", &viewerState.direction, 0, int32(numDirections-1)).Build()
			}
		}),
		giu.Custom(func() {
			if numFrames > 1 {
				giu.SliderInt("frame", &viewerState.frame, 0, int32(numFrames-1)).Build()
			}
		}),
		frameImage,
	}.Build()
}

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

// SetScale sets image scale
func (fv *FrameViewerDC6) SetScale(scale float64) *FrameViewerDC6 {
	s := fv.getState()

	if scale <= 0 {
		scale = 1.0
	}

	s.scale = scale

	return fv
}

func (fv *FrameViewerDC6) setState(s giu.Disposable) {
	giu.Context.SetState(fv.getStateID(), s)
}
