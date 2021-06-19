package sound

import (
	"errors"
	"io/ioutil"
	"log"
	"path"
	"time"

	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/utility"
	"github.com/ted537/go-openal/openal"
)

/* in OpenAL, there are 3 types of entities: Listener, Source and Buffer
Listener is singleton and defines where the listener positioned
Sources may be multiple and are the objects that emit the sound
Buffers are the sound the sources emit



*/

type SoundSource struct {
	Source *openal.Source
	// if its two dimensional, update its position
	//otherwise, set position equal to listener position, so the volume will be constant
	TwoDimensional bool
	Alive          bool
}

type SoundOptions struct {
	IsStatic bool
	Looped   bool
	Gain     float32
	Pitch    float32
}

const (
	max_sources int     = 10
	z           float32 = 0
	basePath    string  = "./assets/sound"
)

var (
	device   *openal.Device
	context  *openal.Context
	Listener *openal.Listener
	DEBUG    bool = false

	events  map[string]*openal.Buffer
	sources []*openal.Source

	soundSources []*SoundSource

	muted bool = false
)

//initialize openal
func Init() {

	device = openal.OpenDevice("")
	context = device.CreateContext()
	context.Activate()
	if checkError() {
		log.Fatal("Could not init sound")
	}

	Listener = &openal.Listener{}

	for i := 0; i < max_sources; i++ {
		source := openal.NewSource()
		source.SetLooping(false)

		sources = append(sources, &source)
	}

	events = make(map[string]*openal.Buffer)
}

//close openal
func Close() {
	for _, source := range sources {
		source.Delete()
	}

	for _, buffer := range events {
		buffer.Delete()
	}

	context.Destroy()
	device.CloseDevice()
}

// update sound source  positions
func Update() {
	var listenerPos openal.Vector
	Listener.GetPosition(&listenerPos)

	for _, soundSource := range soundSources {
		vec := openal.Vector{}
		vec2 := openal.Vector{}
		soundSource.Source.GetPosition(&vec)
		Listener.GetPosition(&vec2)
		// fmt.Println(vec, "    ", vec2)
		if soundSource.Source.State() == openal.Playing {
			if !soundSource.TwoDimensional {
				soundSource.Source.SetPosition(&listenerPos)
			}
		} else if soundSource.Source.State() == openal.Stopped {
			soundSource.Alive = false
		}
	}

	newSoundSources := make([]*SoundSource, 0)
	for _, soundSource := range soundSources {
		if soundSource.Alive {
			newSoundSources = append(newSoundSources, soundSource)
		}
	}

	soundSources = newSoundSources
}

//play simple sound
func PlaySound(event_name string, options ...SoundOptions) uint32 {
	source := getFreeSource()
	if source == nil {
		printDebug("NO AVAILABLE SOURCES")
		return 0
	}
	buffer, ok := events[event_name]
	if !ok {
		printDebug("NO SUCH EVENT")
		return 0
	}
	source.SetBuffer(*buffer)
	if len(options) > 0 {
		if options[0].Gain != 0 {
			source.SetGain(options[0].Gain)
		}
		if options[0].Pitch != 0 {
			source.SetPitch(options[0].Pitch)
		}
		source.SetLooping(options[0].Looped)
	} else { //restore default
		source.SetGain(1.0)
		source.SetPitch(1.0)
		source.SetLooping(false)
	}

	var listenerPos openal.Vector
	Listener.GetPosition(&listenerPos)
	source.SetPosition(&listenerPos)

	source.Play()

	soundSources = append(soundSources, &SoundSource{
		TwoDimensional: false,
		Source:         source,
		Alive:          true,
	})

	return uint32(*source)
}

//Play sound placed in space. Pass in optional sound options struct
func Play2DSound(event_name string, position *physics.Vec2, options ...SoundOptions) uint32 {
	source := getFreeSource()
	if source == nil {
		printDebug("NO AVAILABLE SOURCES")
		return 0
	}

	buffer, ok := events[event_name]
	if !ok {
		printDebug("NO SUCH EVENT")
		return 0
	}

	source.SetBuffer(*buffer)
	source.Set3f(openal.AlPosition, position.X, position.Y, z)

	isStatic := false
	if len(options) > 0 {
		if options[0].Gain != 0 {
			source.SetGain(options[0].Gain)
		}
		if options[0].Pitch != 0 {
			source.SetPitch(options[0].Pitch)
		}
		source.SetLooping(options[0].Looped)
		isStatic = options[0].IsStatic
	} else { //restore default
		source.SetGain(1.0)
		source.SetPitch(1.0)
		source.SetLooping(false)
	}

	if !isStatic {
		go func() {
			for source.State().String() == "Playing" {
				time.Sleep(50 * time.Millisecond)
				source.Set3f(openal.AlPosition, position.X, position.Y, z)
			}
		}()
	}

	// if not static, update sources position in gorouting for the duration of the sound

	source.Play()

	soundSources = append(soundSources, &SoundSource{
		TwoDimensional: true,
		Source:         source,
		Alive:          true,
	})

	return uint32(*source)
}

func StopSound(sourceInt uint32) {
	if sourceInt == 0 {
		printDebug("Wrong source ID!")
	}
	source := openal.Source(sourceInt)
	source.Stop()
	checkError()
}

//register sound event. If relative path is given, its base is considered ./assets/sound
func LoadSound(event_name string, filename string) error {
	_, ok := events[event_name]
	if ok {
		return errors.New("already registered")
	}
	if path.Dir(filename) == "." {
		filename = path.Join(basePath, filename)
	}

	buffer, err := NewBuffer(filename)
	if err != nil {
		return err
	}

	events[event_name] = buffer

	return nil
}

//set speaker position
func SetListenerPosition(position physics.Vec2) {
	Listener.Set3f(openal.AlPosition, position.X, position.Y, z)
}

//get free sound source, and if not, return nil
func getFreeSource() *openal.Source {
	for _, source := range sources {
		state := source.State().String()
		if state == "Initial" || state == "Stopped" {
			return source
		}
	}
	return nil
}

func printDebug(messages ...interface{}) {
	if !DEBUG {
		log.Println(messages...)
	}
}

//create new buffer from file
func NewBuffer(filename string) (*openal.Buffer, error) {

	if path.Dir(filename) == "." {
		filename = path.Join(basePath, filename)
	}
	wavInfo, err := getWavInfo(filename)

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	buffer := openal.NewBuffer()
	audioFormat := determineFormat(wavInfo)
	buffer.SetData(audioFormat, data, int32(wavInfo.SampleRate))

	return &buffer, nil
}

// from 0 to 100
func SetVolume(value float32) {
	clampedValue := utility.ClampF(value, 0, 100)

	Listener.SetGain(clampedValue / 100)
}

func ToggleMute() {
	if !muted {
		muted = true
		SetVolume(0)
	} else {
		muted = false
		SetVolume(100)
	}
}

func checkError() bool {
	err := openal.Err()
	switch err {
	case openal.ErrInvalidContext:
		printDebug("[OPENAL ERROR]", "Invalid context")
	case openal.ErrInvalidDevice:
		printDebug("[OPENAL ERROR]", "Invalid debug")
	case openal.ErrInvalidEnum:
		printDebug("[OPENAL ERROR]", "Invalid enum value passed as parameter")
	case openal.ErrInvalidName:
		printDebug("[OPENAL ERROR]", "Bad id passed as a parameter")
	case openal.ErrInvalidOperation:
		printDebug("[OPENAL ERROR]", "Invalid operation")
	case openal.ErrInvalidValue:
		printDebug("[OPENAL ERROR]", "Invalid value passed as a parameter")
	case openal.ErrOutOfMemory:
		printDebug("[OPENAL ERROR]", "Run out of memory")
	default:
		return false
	}
	return true
}
