package server

import (
	"bytes"
	"fmt"
	"html"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/zaydek/go-ssr/pkg/pretty"
	"github.com/zaydek/go-ssr/pkg/terminal"
)

type Head struct {
	URL         string
	Title       string
	Description string
	ImageURL    string
}

type MetaWriter struct {
	buf bytes.Buffer
}

func (w *MetaWriter) Write(meta string) {
	if w.buf.Len() > 0 {
		w.buf.WriteString("\n" + strings.Repeat(" ", 4))
	}
	w.buf.WriteString(meta)
}

func (w *MetaWriter) Writeformat(metaf, v string) {
	if v == "" {
		return
	}

	if w.buf.Len() > 0 {
		w.buf.WriteString("\n" + strings.Repeat(" ", 4))
	}
	str := fmt.Sprintf(metaf, html.EscapeString(v))
	w.buf.WriteString(str)
}

func (w MetaWriter) String() string {
	return w.buf.String()
}

func (h Head) String() string {
	var meta MetaWriter

	meta.Writeformat(`<title>%s</title>`, h.Title)
	meta.Writeformat(`<meta name="title" content="%s">`, h.Title)
	meta.Writeformat(`<meta name="description" content="%s">`, h.Description)

	meta.Write(`<meta property="og:type" content="website">`)
	meta.Writeformat(`<meta property="og:url" content="%s">`, h.URL)
	meta.Writeformat(`<meta property="og:title" content="%s">`, h.Title)
	meta.Writeformat(`<meta property="og:description" content="%s">`, h.Description)
	meta.Writeformat(`<meta property="og:image" content="%s">`, h.ImageURL)

	meta.Write(`<meta property="twitter:card" content="summary_large_image">`)
	meta.Writeformat(`<meta property="twitter:url" content="%s">`, h.URL)
	meta.Writeformat(`<meta property="twitter:title" content="%s">`, h.Title)
	meta.Writeformat(`<meta property="twitter:description" content="%s">`, h.Description)
	meta.Writeformat(`<meta property="twitter:image" content="%s">`, h.ImageURL)

	return meta.String()
}

func observe(w http.ResponseWriter, r *http.Request) func() {
	start := time.Now()
	return func() {
		fmt.Fprintf(os.Stderr, "%s %s %s\n",
			r.Method, r.URL.Path, terminal.Dimf("(%s)", pretty.Duration(time.Since(start))))
	}
}

func Run() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer observe(w, r)()
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
  <head>
    `+Head{
			Title:       "Hello, world! (/)",
			Description: "Welcome to my wonderful site.",
		}.String()+`
  </head>
  <body>
    <h1>Hello, world! (/)</h1>
    <script src="net/vendor.js"></script>
    <script src="net/client.js"></script>
  </body>
</html>
`)
	})

	http.HandleFunc("/pokemon/", func(w http.ResponseWriter, r *http.Request) {
		defer observe(w, r)()
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
  <head>
    `+Head{
			Title:       "Hello, world! (/pokemon/)",
			Description: "Welcome to my wonderful site.",
		}.String()+`
  </head>
  <body>
    <h1>Hello, world! (/pokemon/)</h1>
    <script src="net/vendor.js"></script>
    <script src="net/client.js"></script>
  </body>
</html>
`)
	})

	http.HandleFunc("/nested/pokemon/", func(w http.ResponseWriter, r *http.Request) {
		defer observe(w, r)()
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
  <head>
    `+Head{
			Title:       "Hello, world! (/nested/pokemon/)",
			Description: "Welcome to my wonderful site.",
		}.String()+`
  </head>
  <body>
    <h1>Hello, world! (/nested/pokemon/)</h1>
    <script src="net/vendor.js"></script>
    <script src="net/client.js"></script>
  </body>
</html>
`)
	})

	http.HandleFunc("/net/", func(w http.ResponseWriter, r *http.Request) {
		defer observe(w, r)()
		http.ServeFile(w, r, filepath.Join("www", r.URL.Path))
	})

	var port = 8000
	if envPort := os.Getenv("PORT"); envPort != "" {
		port, _ = strconv.Atoi(envPort)
	}

	fmt.Printf("ready on port %d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
