package notes

import (
	"fmt"

	"github.com/priyanshu360/keyboard-piano/internal/repository"
)

type Notes map[string][]byte

type NoteFetcher interface {
	Fetch(string) (Notes, error)
}

type noteFetcher struct {
	media  repository.MediaStore
	lyrics repository.LyricsStore
}

func NewNoteFetcher(media repository.MediaStore, lyrics repository.LyricsStore) *noteFetcher {
	return &noteFetcher{
		lyrics: lyrics,
		media:  media,
	}
}

func (nf *noteFetcher) Fetch(song string) (Notes, error) {
	notes := make(Notes)
	lyrics, err := nf.lyrics.Fetch(song)
	if err != nil {
		return notes, err
	}

	for idx := range lyrics {
		if idx%2 == 1 {
			continue
		}
		for col := range lyrics[idx] {
			if lyrics[idx][col] == "" {
				continue
			}

			key, err := GetKey(lyrics[idx][col])
			if err != nil {
				return notes, err
			}

			bytes, err := nf.media.Fetch(lyrics[idx][col])
			if err != nil {
				return notes, err
			}

			notes[key] = bytes
		}
	}

	return notes, nil
}

// getKey maps musical note names to specific keys using simple string manipulation and a switch statement.
func GetKey(lyric string) (string, error) {
	if len(lyric) == 0 {
		return "", fmt.Errorf("invalid note")
	}

	// Remove the last character (octave number) from the note
	baseNote := lyric[:len(lyric)-1]

	// Determine the key based on the base note
	var key string
	switch baseNote {
	case "C":
		key = "a"
	case "D":
		key = "s"
	case "E":
		key = "d"
	case "F":
		key = "f"
	case "G":
		key = "j"
	case "A":
		key = "k"
	case "B":
		key = "l"
	case "Db":
		key = "w"
	case "Eb":
		key = "e"
	case "Gb":
		key = "t"
	case "Ab":
		key = "i"
	case "Bb":
		key = "o"
	default:
		return "", fmt.Errorf("note not found")
	}

	return key, nil
}
