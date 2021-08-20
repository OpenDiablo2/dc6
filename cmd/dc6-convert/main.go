package main

import (
	"bytes"
	"flag"
	"fmt"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	dc6lib "github.com/OpenDiablo2/dc6/pkg"
	gpl "github.com/gravestench/gpl/pkg"
)

type options struct {
	dc6Path *string
	palPath *string
	pngPath *string
}

func main() {
	var o options

	if parseOptions(&o) {
		flag.Usage()
		return
	}

	dc6Data, err := ioutil.ReadFile(*o.dc6Path)
	if err != nil {
		const fmtErr = "could not read file, %v"

		fmt.Print(fmt.Errorf(fmtErr, err))

		return
	}

	dc6, err := dc6lib.FromBytes(dc6Data)
	if err != nil {
		fmt.Println(err)
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

	numDirections := dc6.Frames.NumberOfDirections()
	framesPerDir := dc6.Frames.FramesPerDirection()
	isMultiFrame := numDirections > 1 || framesPerDir > 1

	outfilePath := *o.pngPath
	if isMultiFrame {
		noExt := fileNameWithoutExt(outfilePath)
		outfilePath = noExt + "_d%v_f%v.png"
	}

	for dirIdx := 0; dirIdx < numDirections; dirIdx++ {
		for frameIdx := 0; frameIdx < framesPerDir; frameIdx++ {
			outPath := outfilePath

			if isMultiFrame {
				outPath = fmt.Sprintf(outfilePath, dirIdx, frameIdx)
			}

			f, err := os.Create(outPath)
			if err != nil {
				log.Fatal(err)
			}

			if err := png.Encode(f, dc6.Frames.Direction(dirIdx).Frame(frameIdx)); err != nil {
				_ = f.Close()

				log.Fatal(err)
			}

			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func parseOptions(o *options) (terminate bool) {
	o.dc6Path = flag.String("dc6", "", "input dc6lib file (required)")
	o.palPath = flag.String("pal", "", "input pal file (optional)")
	o.pngPath = flag.String("png", "", "path to png file (optional)")

	flag.Parse()

	if *o.dc6Path == "" {
		flag.Usage()
		return true
	}

	return false
}

func fileNameWithoutExt(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}
