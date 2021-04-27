package spriteloader

import(
	"fmt"
	"image"
)

type ChanInfo struct {
	name string
	img *image.RGBA
	loadCode int
}

const (
	goroutineCount = 8
)

type ImgLoader struct {
	// the images needed load
	imgList []string
	// input channel
	inChan chan string
	// output channel
	outChan chan *ChanInfo
	//the images
	imgRes map[string] *image.RGBA
	//load count include load failed
	loadCount int
	//loaded flag
	loaded bool
	//callback
	callback func(*ImgLoader)
}

func NewImgLoader(imgList []string, cb func(*ImgLoader)) *ImgLoader {
	count := len(imgList)
	if count == 0 {
		return nil
	}
	loader := ImgLoader{
		imgList: imgList,
		inChan: make(chan string, count),
		outChan: make(chan *ChanInfo, count),
		imgRes: make(map[string] *image.RGBA),
		loadCount:0,
		loaded:false,
		callback:cb,
	}

	return &loader
}

func (il *ImgLoader) GetImg(path string) *image.RGBA {
	return il.imgRes[path]
}

func (il *ImgLoader) Run(){
	go inputChannel(il.imgList, il.inChan)
	for i := 0; i < goroutineCount; i++ {
		go outChannel(il.inChan, il.outChan)
	}
}

func (il *ImgLoader) NeedUpdate() bool {
	return !il.loaded
}

func (il *ImgLoader) Update(){
	//load finish
	if len(il.imgList) == il.loadCount {
		close(il.outChan)
		fmt.Printf("img had loaded finish\n")
		il.loaded = true
		il.callback(il)
		return
	}

	select {
	case ci := <- il.outChan:
		if ci.loadCode != LoadOk {
			fmt.Printf("img load failed %v\n", ci)
			fmt.Printf("img load failed reason %s\n", GetCodeString(ci.loadCode))
		}
		fmt.Printf("img load ok %v\n", ci)
		il.loadCount++
		il.imgRes[ci.name] = ci.img
	}
}

func inputChannel(imgList []string, inChan chan <- string){
	for _, path := range imgList{
		inChan <- path
	}
	close(inChan)
}

func outChannel(inChan <- chan string, outChan chan <- *ChanInfo){
	for path := range inChan{
		loadCode,img := LoadPng(path)
		ci := ChanInfo{
			name:path,
			img:img,
			loadCode: loadCode,
		}
		outChan <- &ci
	}
}