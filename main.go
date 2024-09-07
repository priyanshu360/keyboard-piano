package main

import (
	"bytes"
	"fmt"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
	"net"
	"os"
	"time"
)

func createContext() *oto.Context {
	// Prepare an Oto context (this will use your default audio device) that will
	// play all our sounds. Its configuration can't be changed later.

	op := &oto.NewContextOptions{}

	// Usually 44100 or 48000. Other values might cause distortions in Oto
	op.SampleRate = 44100

	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	op.ChannelCount = 2

	// Format of the source. go-mp3's format is signed 16bit integers.
	op.Format = oto.FormatSignedInt16LE

	// Remember that you should **not** create more than one context
	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan

	return otoCtx
}

func play(otoCtx *oto.Context, fileBytes []byte) error {
	// Convert the pure bytes into a reader object that can be used with the mp3 decoder
	fileBytesReader := bytes.NewReader(fileBytes)

	// Decode file
	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	// Create a new 'player' that will handle our sound. Paused by default.
	player := otoCtx.NewPlayer(decodedMp3)

	// Play starts playing the sound and returns without waiting for it (Play() is async).
	player.Play()

	// We can wait for the sound to finish playing using something like this
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	// If you don't want the player/sound anymore simply close
	err = player.Close()
	if err != nil {
		panic("player.Close failed: " + err.Error())
	}
	return nil
}

func readTCP(conn net.Conn) {
	player := NewPlayer()
	for {
		b := make([]byte, 100)
		size, err := conn.Read(b)
		if err != nil {
			panic(err)
		}
		if size != 0 {
			for _, each := range string(b[0:size]) {
				fmt.Print(string(each))
				go player.Play(string(each))
			}
		}
	}
}

func makeTCP() {
	listner, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listner.Accept()
		if err != nil {
			panic(err)
		}
		readTCP(conn)
	}
}

type Player interface {
	Play()
}

type player struct {
	otoCtx *oto.Context
	km     notes
}

func NewPlayer() *player {
	return &player{
		otoCtx: createContext(),
		km:     initNotes(),
	}
}

func (p *player) Play(key string) {
	switch key {
	case "x":
		os.Exit(1)
	default:
		if b, ok := p.km[key]; ok {
			play(p.otoCtx, b)
		}
	}
}

func main() {
	makeTCP()
}

type notes map[string][]byte

func read(file string) []byte {

	fileBytes, err := os.ReadFile(file)
	if err != nil {
		panic("reading my-file.mp3 failed: " + err.Error())
	}
	return fileBytes
}

func initNotes() notes {
	keyMap := make(notes)

	keyMap["a"] = read("./mp3/C3.mp3")
	keyMap["s"] = read("./mp3/D3.mp3")
	keyMap["d"] = read("./mp3/E3.mp3")
	keyMap["f"] = read("./mp3/F3.mp3")
	keyMap["j"] = read("./mp3/G3.mp3")
	keyMap["k"] = read("./mp3/A3.mp3")
	keyMap["l"] = read("./mp3/B3.mp3")
	keyMap["w"] = read("./mp3/Db3.mp3")
	keyMap["e"] = read("./mp3/Eb3.mp3")
	keyMap["t"] = read("./mp3/Gb3.mp3")
	keyMap["i"] = read("./mp3/Ab3.mp3")
	keyMap["o"] = read("./mp3/Bb3.mp3")

	return keyMap
}
