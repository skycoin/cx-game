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
	//HadLoaded
	OpenError
	DecodeError
)

var (
	DEBUG = true
)

type ImgStat struct {
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

func (imgStat *ImgStat) FileSize() int64 {
	return imgStat.fileSize
}

func (imgStat *ImgStat) ConvertTime() time.Duration {
	return imgStat.convertTime
}

func (imgStat *ImgStat) ReadTime() time.Duration {
	return imgStat.readTime
}

func (imgStat *ImgStat) Loaded() bool {
	return imgStat.loaded
}

//func LoadingRate() float32 {
//	var allSize, loadingSize int64
//	for _, ss := range spritesStat {
//		sz := ss.FileSize()
//		allSize += sz
//		if ss.Loaded() {
//			loadingSize += sz
//		}
//	}
//
//	return float32(loadingSize / allSize)
//}

func timeCost(s func(time.Duration)) func() {
	start := time.Now()
	return func() {
		end := time.Since(start)
		s(end)
	}
}

func decodePng(imgStat *ImgStat, imgFile *os.File) (int, *image.RGBA, *ImgStat) {
	defer timeCost(func(end time.Duration) {
		imgStat.convertTime = end
		imgStat.loaded = true
		if DEBUG {
			fmt.Printf("decode png ok \n%v\n", imgStat.name)
		}
	})()
	defer imgFile.Close()

	im, err := png.Decode(imgFile)
	if err != nil {
		if DEBUG {
			fmt.Printf("error decodePng %v %v\n", imgStat, err)
		}
		// log.Fatalln(err)
		return DecodeError, nil, nil
	}

	img := image.NewRGBA(im.Bounds())
	draw.Draw(img, img.Bounds(), im, image.Pt(0, 0), draw.Src)
	return LoadOk, img, imgStat
}

func openPng(imgStat *ImgStat) *os.File {
	defer timeCost(func(end time.Duration) {
		imgStat.readTime = end
		if DEBUG {
			fmt.Printf("open png ok \n%v\n", imgStat.name)
		}
	})()
	imgFile, err := os.Open(imgStat.name)

	if err != nil {
		if DEBUG {
			fmt.Printf("error openPng %v %v\n", imgStat, err)
		}
		// log.Fatalln(err)
		imgFile.Close()
		return nil
	}

	fi, _ := imgFile.Stat()
	imgStat.fileSize = fi.Size()
	return imgFile
}

//please check return value if load failed
func LoadPng(path string) (int, *image.RGBA, *ImgStat) {
	imgStat := ImgStat{
		name:        path,
		readTime:    0,
		convertTime: 0,
		fileSize:    0,
	}
	imgFile := openPng(&imgStat)
	if imgFile == nil {
		return OpenError, nil, nil
	}

	return decodePng(&imgStat, imgFile)
}

func GetCodeString(code int) string {
	var ret string
	switch code {
	case LoadOk:
		ret = "LoadOk"
	//case HadLoaded: ret = "HadLoaded"
	case OpenError:
		ret = "OpenError"
	case DecodeError:
		ret = "DecodeError"
	}
	return ret
}
