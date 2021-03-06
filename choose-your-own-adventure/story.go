package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strings"
)

var defaultHandlerTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Choose Your Own Adventure</title>
</head>
<body>
  <h1>{{.Title}}</h1>
  {{range .Paragraphs}}
  <p>{{.}}</p>
  {{end}}

  <ul>
  {{range .Options}}
    <li>
      <a href="/{{.Chapter}}">{{.Text}}</a>
    </li>
  {{end}}
  </ul>
</body>
</html>`

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.template = t
	}
}

func WithParseFunc(p func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.parsePathFn = p
	}
}

func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{
		story:       s,
		template:    template.Must(template.New("").Parse(defaultHandlerTemplate)),
		parsePathFn: parsePath,
	}

	for _, opt := range opts {
		opt(&h)
	}

	return h
}

func parsePath(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)

	if path == "" || path == "/" {
		// Act like they went to "/intro"
		path = "/intro"
	}

	// Remove the leading slash
	path = path[1:]

	return path
}

type handler struct {
	story       Story
	template    *template.Template
	parsePathFn func(r *http.Request) string
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.parsePathFn(r)

	// Check if the chapter exists
	if chapter, ok := h.story[path]; ok {
		err := h.template.Execute(w, chapter)

		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}

		return
	}

	http.Error(w, "Chapter not found", http.StatusNotFound)
}

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
