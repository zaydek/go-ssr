package export

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/zaydek/go-ssr/cmd/server"
	"github.com/zaydek/go-ssr/pkg/pretty"
	"github.com/zaydek/go-ssr/pkg/terminal"
)

const (
	MODE_DIR  = 0755
	MODE_FILE = 0644
)

const (
	EXPORT_DIR = "__export__"
)

var pathnames = []string{
	"/",
	"/pokemon",
	"/nested/pokemon",
	"/404",
}

func getBrowserPath(url string) string {
	out := url
	if strings.HasSuffix(url, "/index.html") {
		out = out[:len(out)-len("index.html")]
	} else if ext := filepath.Ext(url); ext == ".html" {
		out = out[:len(out)-len(".html")]
	}
	return out
}

func getFSPath(url string) string {
	out := url
	if strings.HasSuffix(url, "/") {
		out += "index.html"
	} else if ext := filepath.Ext(url); ext == "" {
		out += ".html"
	}
	return out
}

func dimAll(str string) string {
	arr := strings.Split(str, "\n")
	for x, v := range arr {
		arr[x] = terminal.Dim(v)
	}
	return strings.Join(arr, "\n")
}

func Run() {
	// rm -r __export__
	if err := os.RemoveAll(EXPORT_DIR); err != nil {
		panic(err)
	}

	// cp -r www/net __export__/net
	if err := copyDir("www/net", filepath.Join(EXPORT_DIR, "net"), nil); err != nil {
		panic(err)
	}

	var port = 8000
	if envPort := os.Getenv("PORT"); envPort != "" {
		port, _ = strconv.Atoi(envPort)
	}

	go func() {
		server.Run()
	}()

	start := time.Now()
	for x, pathname := range pathnames {
		reqStart := time.Now()

		if x > 0 {
			fmt.Println()
		}

		res, err := http.Get(fmt.Sprintf("http://localhost:%d/%s", port, pathname))
		if err != nil {
			panic(err)
		} else if res.StatusCode != 200 {
			fmt.Fprintf(os.Stderr, terminal.BoldRedf("'%s %s' - %d\n",
				res.Request.Method, res.Request.URL.Path, res.StatusCode))
			return
		}

		bstr, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		res.Body.Close() // Do not defer
		// fmt.Print(dimAll(strings.ReplaceAll(string(bstr), "\t", "  ")))
		fmt.Print(dimAll(string(bstr)))

		target := filepath.Join(EXPORT_DIR, getFSPath(res.Request.URL.Path))
		if err := os.MkdirAll(filepath.Dir(target), MODE_DIR); err != nil {
			panic(err)
		}

		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf("<!-- %s -->", pretty.Duration(time.Since(reqStart))))
		buf.WriteString("\n")
		buf.Write(bstr)

		if err := ioutil.WriteFile(target, buf.Bytes(), MODE_FILE); err != nil {
			panic(err)
		}
	}

	fmt.Println()
	fmt.Println(terminal.Dimf("(%s)", pretty.Duration(time.Since(start))))
}
