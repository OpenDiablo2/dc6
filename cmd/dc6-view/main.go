package main

import (
	"bytes"
	"flag"
	"fmt"
	"image/color"
	"io/ioutil"
	"path"

	"github.com/AllenDang/giu"

	dc6lib "github.com/OpenDiablo2/dc6/pkg"
	dc6widget "github.com/OpenDiablo2/dc6/pkg/giuwidget"
	gpl "github.com/gravestench/gpl/pkg"
)

const (
	title               = "dc6 viewer"
	windowFlags         = giu.MasterWindowFlagsFloating & giu.MasterWindowFlagsNotResizable
	minWidth, minHeight = 1, 1
	padWidth            = 8 // px
)

func main() {
	var o options

	if parseOptions(&o) {
		flag.Usage()
		return
	}

	srcPath := *o.dc6Path

	fileContents, err := ioutil.ReadFile(srcPath)
	if err != nil {
		const fmtErr = "could not read file, %w"

		fmt.Print(fmt.Errorf(fmtErr, err))

		return
	}

	dc6, err := dc6lib.FromBytes(fileContents)
	if err != nil {
		fmt.Print(err)
		return
	}

	if *o.palPath != "" {
		palData, err := ioutil.ReadFile(*o.palPath)
		if err != nil {
			fmt.Println(err)
			return
		}

		gplInstance, err := gpl.Decode(bytes.NewBuffer(palData))
		if err != nil {
			fmt.Println("palette is not a GIMP palette file...")
			return
		}

		dc6.SetPalette(color.Palette(*gplInstance))
	}

	f0 := dc6.Frames.Direction(0).Frame(0)

	imgW := int(float64(f0.Width) * *o.scale)
	imgH := int(float64(f0.Height) * *o.scale)

	w, h := imgW+padWidth<<1, imgH+padWidth<<1

	if w < minWidth {
		w = minWidth
	}

	if h < minHeight {
		h = minHeight
	}

	windowTitle := fmt.Sprintf("%s - %s", title, path.Base(srcPath))

	window := giu.NewMasterWindow(windowTitle, w, h, windowFlags, nil)
	id := fmt.Sprintf("%s##%s", windowTitle, "dc6")

	viewer := dc6widget.FrameViewer(id, dc6)
	viewer.SetScale(*o.scale)

	window.Run(func() {
		giu.SingleWindow(windowTitle).Layout(viewer)
	})
}

type options struct {
	dc6Path *string
	palPath *string
	pngPath *string
	scale   *float64
}

func parseOptions(o *options) (terminate bool) {
	o.dc6Path = flag.String("dc6", "", "input dc6 file (required)")
	o.palPath = flag.String("pal", "", "input pal file (optional)")
	o.pngPath = flag.String("png", "", "path to png file (optional)")
	o.scale = flag.Float64("scale", 1.0, "scale")

	flag.Parse()

	return *o.dc6Path == ""
}
