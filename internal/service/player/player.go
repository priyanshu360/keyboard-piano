package player

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
	"github.com/priyanshu360/keyboard-piano/internal/service/notes"
	"golang.org/x/term"
)

type Player interface {
	Play(string)
}

type player struct {
	notes notes.NoteFetcher
}

func NewPlayer(notes notes.NoteFetcher) *player {
	return &player{
		notes: notes,
	}
}

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

func (p *player) Play(song string) {

	// Get the current terminal state
	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		fmt.Println("Error entering raw mode:", err)
		return
	}
	defer term.Restore(int(syscall.Stdin), oldState) // Restore the terminal when done

	notes, err := p.notes.Fetch(song)
	if err != nil {
		panic(err)
	}

	otoCtx := createContext()
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		r := scanner.Text()
		switch r {
		case "x":
			os.Exit(1)
		default:
			if b, ok := notes[r]; ok {
				go play(otoCtx, b)
			}
		}
	}
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
