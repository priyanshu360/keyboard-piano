package main

import (
	"fmt"

	"github.com/priyanshu360/keyboard-piano/internal/repository"
	"github.com/priyanshu360/keyboard-piano/internal/service/notes"
	"github.com/priyanshu360/keyboard-piano/internal/service/player"
)

func main() {
	mediaStore := repository.NewLocalMediaStore("./media")
	lyricStore := repository.NewLocalLyricsStore()

	notesSvc := notes.NewNoteFetcher(mediaStore, lyricStore)
	playerSvc := player.NewPlayer(notesSvc)
	lyrics, _ := lyricStore.Fetch("Treat You Better")
	for idx := range lyrics {
		if idx%2 == 1 {
			fmt.Println(lyrics[idx])
		} else {
			for col := range lyrics[idx] {
				if lyrics[idx][col] == "" {
					fmt.Print("  ")
					continue
				}
				val, _ := notes.GetKey(lyrics[idx][col])
				fmt.Printf("%s  ", val)
			}
			fmt.Println()
		}
	}

	playerSvc.Play("Treat You Better")
}
