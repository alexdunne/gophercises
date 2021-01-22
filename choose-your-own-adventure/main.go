package cyoa

import (
	"encoding/json"
	"io"
)

func JsonStory(r io.Reader) (Story, error) {
	decoder := json.NewDecoder(r)

	var story Story
	if err := decoder.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil
}

type Story map[string]Chapter

type Chapter struct {
	Title      string          `json:"title"`
	Paragraphs []string        `json:"story"`
	Options    []ChapterOption `json:"options"`
}

type ChapterOption struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}
