package main

import (
	"encoding/json"
	"os/exec"
	"sync"
)

// Song : Song information fetched from youtube-dl. More fields exist in the JSON response, but this is what we need.
type Song struct {
	Title      string `json:"title"`
	URL        string `json:"url"`
	Thumbnail  string `json:"thumbnail"`
	WebpageURL string `json:"webpage_url"`
}

type SongCache struct {
	cache map[string]Song
	mux   sync.RWMutex
}

var songCache = SongCache{
	cache: make(map[string]Song),
}

func (sc *SongCache) Get(url string) (Song, bool) {
	sc.mux.RLock()
	defer sc.mux.RUnlock()
	song, exists := sc.cache[url]
	return song, exists
}

func (sc *SongCache) Set(url string, song Song) {
	sc.mux.Lock()
	defer sc.mux.Unlock()
	sc.cache[url] = song
}

func fetchSong(url string) (Song, error) {
	// Check cache first
	if song, exists := songCache.Get(url); exists {
		return song, nil
	}

	// If not in cache, fetch and store
	song, err := fetchSongFromYoutubeDL(url)
	if err == nil {
		songCache.Set(url, song)
	}
	return song, err
}

func fetchSongFromYoutubeDL(url string) (Song, error) {
	// Maybe make this function return a promise or something in the future
	var songData Song
	// JSON dump has no extra overhead, and we get more info that we need that might be useful
	cmd := exec.Command("yt-dlp", "--no-playlist", "-J", "--youtube-skip-dash-manifest", "-f bestaudio", url)
	stdOut, err := cmd.StdoutPipe()

	if err != nil {
		return Song{}, err
	}

	if err := cmd.Start(); err != nil {
		return Song{}, err
	}

	if err := json.NewDecoder(stdOut).Decode(&songData); err != nil {
		return Song{}, err
	}

	cmd.Wait()
	// Only using these fields atm
	return songData, nil
}
