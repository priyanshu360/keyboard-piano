package repository

import (
	"fmt"
	"os"
)

type MediaStore interface {
	Fetch(string) ([]byte, error)
}

type LyricsStore interface {
	Fetch(string) ([][]string, error)
}

type localMediaStore struct {
	basePath string
}

func (m *localMediaStore) Fetch(note string) ([]byte, error) {
	file := fmt.Sprintf("%s/%s.mp3", m.basePath, note)
	return os.ReadFile(file)
}

func NewLocalMediaStore(basePath string) *localMediaStore {
	return &localMediaStore{
		basePath: basePath,
	}
}

type localLyricsStore struct {
	store map[string][][]string
}

func (s *localLyricsStore) Fetch(note string) ([][]string, error) {
	return s.store[note], nil
}

func NewLocalLyricsStore() *localLyricsStore {
	store := map[string][][]string{}

	store["Treat You Better"] = [][]string{
		{"F3", "Db3", "Eb3", "Db3", "", "Db3", "Db3", "Eb3", "", "F3", "F3", "", "Db3", "Eb3", "Db3", "", "Db3"},
		{"I", "won't", "lie", "to", "", "you", "", "I", "know", "he's", "", "just", "not", "", "right", "for", "you"},

		{"Db3", "Db3", "", "Db3", "F3", "", "F3", "", "F3", "", "F3", "Ab3", ""},
		{"And", "you", "", "can", "tell", "", "me", "", "if", "", "I'm", "off", ""},

		{"Db3", "", "Db3", "Eb3", "", "Eb3", "Eb3", "F3", "", "Eb3", ""},
		{"But", "", "I", "see", "", "it", "on", "your", "", "face", ""},

		{"Db3", "Db3", "", "Db3", "", "Db3", "Db3", "", "Db3", "Eb3", "F3", "F3", ""},
		{"When", "you", "", "say", "", "that", "he's", "", "the", "one", "that", "you", "want"},

		{"Db3", "Db3", "F3", "F3", "", "F3", "F3", "Ab3", "", "Db3", "Db3", "", "Eb3", "Eb3", "", "F3", "Eb3", "Db3"},
		{"And", "you're", "spending", "", "all", "your", "", "time", "in", "this", "", "wrong", "situation"},

		{"Db3", "", "Db3", "", "Db3", "", "Db3", "Db3", "", "Eb3", "F3", "", "F3", ""},
		{"And", "", "anytime", "", "you", "want", "it", "to", "", "stop", ""},
	}

	return &localLyricsStore{
		store: store,
	}
}
