package main

import (
	"bytes"
	"fmt"
	"html"
	"net/http"
	"os"
	"strconv"
	"time"
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

func (w *MetaWriter) Write(arg string) {
	if w.buf.Len() > 0 {
		w.buf.Write([]byte("\n\t\t"))
	}
	w.buf.Write([]byte(arg))
}

func (w *MetaWriter) Writeformat(format, arg string) {
	if arg == "" {
		return
	}

	if w.buf.Len() > 0 {
		w.buf.Write([]byte("\n\t\t"))
	}
	str := fmt.Sprintf(format, html.EscapeString(arg))
	w.buf.Write([]byte(str))
}

func (w MetaWriter) String() string {
	return w.buf.String()
}

func (h Head) String() string {
	var meta MetaWriter

	// Web
	meta.Write(`<meta charset="utf-8">`)
	meta.Write(`<meta name="viewport" content="width=device-width, initial-scale=1.0">`)
	meta.Writeformat(`<title>%s</title>`, h.Title)
	meta.Writeformat(`<meta name="title" content="%s">`, h.Title)
	meta.Writeformat(`<meta name="description" content="%s">`, h.Description)

	// og:*
	meta.Write(`<meta property="og:type" content="website">`)
	meta.Writeformat(`<meta property="og:url" content="%s">`, h.URL)
	meta.Writeformat(`<meta property="og:title" content="%s">`, h.Title)
	meta.Writeformat(`<meta property="og:description" content="%s">`, h.Description)
	meta.Writeformat(`<meta property="og:image" content="%s">`, h.ImageURL)

	// twitter:*
	meta.Write(`<meta property="twitter:card" content="summary_large_image">`)
	meta.Writeformat(`<meta property="twitter:url" content="%s">`, h.URL)
	meta.Writeformat(`<meta property="twitter:title" content="%s">`, h.Title)
	meta.Writeformat(`<meta property="twitter:description" content="%s">`, h.Description)
	meta.Writeformat(`<meta property="twitter:image" content="%s">`, h.ImageURL)

	return meta.String()
}

func prettyDur(dur time.Duration) string {
	var out string
	if amount := dur.Nanoseconds(); amount < 1_000 {
		out = strconv.Itoa(int(amount)) + "ns"
	} else if amount := dur.Microseconds(); amount < 1_000 {
		out = strconv.Itoa(int(amount)) + "µs"
	} else if amount := dur.Milliseconds(); amount < 1_000 {
		out = strconv.Itoa(int(amount)) + "ms"
	} else {
		out = strconv.Itoa(int(dur.Seconds())) + "s"
	}
	return out
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			fmt.Printf("%s (%s)\n", r.URL.Path, prettyDur(time.Since(start)))
		}()

		head := Head{
			Title:       "Hello, world! (/)",
			Description: "Welcome to my wonderful site.",
		}
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
	<head>
		`+head.String()+`
	</head>
	<body>
		<h1>Hello, world! (/)</h1>
	</body>
</html>
`)
	})

	http.HandleFunc("/pokemon/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			fmt.Printf("%s (%s)\n", r.URL.Path, prettyDur(time.Since(start)))
		}()

		head := Head{
			Title:       "Hello, world! (/pokemon/)",
			Description: "Welcome to my wonderful site.",
		}
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
	<head>
		`+head.String()+`
	</head>
	<body>
		<h1>Hello, world! (/pokemon/)</h1>
	</body>
</html>
`)
	})

	http.HandleFunc("/nested/pokemon/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			fmt.Printf("%s (%s)\n", r.URL.Path, prettyDur(time.Since(start)))
		}()

		head := Head{
			Title:       "Hello, world! (/nested/pokemon/)",
			Description: "Welcome to my wonderful site.",
		}
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
	<head>
		`+head.String()+`
	</head>
	<body>
		<h1>Hello, world! (/nested/pokemon/)</h1>
	</body>
</html>
`)
	})

	serve := func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, r.URL.Path) }

	http.HandleFunc("/vendor.js", serve)
	http.HandleFunc("/vendor.js.map", serve)
	http.HandleFunc("/client.js", serve)
	http.HandleFunc("/client.js.map", serve)
	http.HandleFunc("/www/", serve)

	var port = 8000
	if envPort := os.Getenv("PORT"); envPort != "" {
		port, _ = strconv.Atoi(envPort)
	}

	fmt.Printf("ready on port %d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
