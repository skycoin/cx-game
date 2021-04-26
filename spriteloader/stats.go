package spriteloader

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"time"
)

const (
	LoadOk = iota
	HadLoaded
	OpenError
	DecodeError
)

type SpriteStat struct {
	name string
	//how long to read file from disc
	readTime time.Duration
	//how long to convert the .png to internal format
	convertTime time.Duration
	//is loaded
	loaded bool
	//file size
	fileSize int64
}

var spritesStat = make(map[string]*SpriteStat)

func (ss *SpriteStat) FileSize() int64 {
	return ss.fileSize
}

func (ss *SpriteStat) ConvertTime() time.Duration {
	return ss.convertTime
}

func (ss *SpriteStat) ReadTime() time.Duration {
	return ss.readTime
}

func (ss *SpriteStat) Loaded() bool {
	return ss.loaded
}

func LoadingRate() float32 {
	var allSize, loadingSize int64
	for _, ss := range spritesStat {
		sz := ss.FileSize()
		allSize += sz
		if ss.Loaded() {
			loadingSize += sz
		}
	}

	return float32(loadingSize / allSize)
}

func timeCost(s func(time.Duration)) func() {
	start := time.Now()
	return func() {
		end := time.Since(start)
		s(end)
	}
}

func decodePng(ss *SpriteStat, imgFile *os.File) (int, *image.RGBA) {
	defer timeCost(func(end time.Duration) {
		ss.convertTime = end
		ss.loaded = true
		fmt.Printf("load png ok \n%v\n", ss)
	})()
	defer imgFile.Close()

	im, err := png.Decode(imgFile)
	if err != nil {
		fmt.Printf("error decodePng %v %v\n", ss, err)
		// log.Fatalln(err)
		return DecodeError, nil
	}

	img := image.NewRGBA(im.Bounds())
	draw.Draw(img, img.Bounds(), im, image.Pt(0, 0), draw.Src)
	spritesStat[ss.name] = ss
	return LoadOk, img
}

func openPng(ss *SpriteStat) *os.File {
	defer timeCost(func(end time.Duration) {
		ss.readTime = end
		fmt.Printf("open png ok \n%v\n", ss)
	})()
	imgFile, err := os.Open(ss.name)

	if err != nil {
		fmt.Printf("error openPng %v %v\n", ss, err)
		// log.Fatalln(err)
		imgFile.Close()
		return nil
	}

	fi, _ := imgFile.Stat()
	ss.fileSize = fi.Size()
	return imgFile
}

//please check return value if load failed
func LoadPng(path string) (int, *image.RGBA) {
	_, ok := spritesStat[path]
	if ok {
		return HadLoaded, nil
	}

	ss := SpriteStat{
		name:       path,
		readTime:   0,
		convertTime: 0,
		fileSize:   0,
	}
	imgFile := openPng(&ss)
	if imgFile == nil {
		return OpenError, nil
	}

	return decodePng(&ss, imgFile)
}