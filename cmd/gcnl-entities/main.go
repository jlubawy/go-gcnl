package main

//go:generate go-bindata-assetfs data/...

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/jlubawy/go-gcnl"
	"github.com/jlubawy/go-gcnl/entities"
)

var Options struct {
	Port string
}

var apiKey string
var t = make(map[string]*template.Template)

func init() {
	flag.StringVar(&Options.Port, "port", ":8080", "TCP port to listen on")
	flag.Parse()

	// Get the API key
	apiKey = os.Getenv("GOOGLE_API_KEY")
	if len(apiKey) == 0 {
		fmt.Fprintln(os.Stderr, "must set GOOGLE_API_KEY environment variable")
		os.Exit(1)
	}

	// Initialize templates
	path := "data/templ"
	names, err := AssetDir(path)
	if err != nil {
		panic(err)
	}

	for _, name := range names {
		data, err := Asset(path + "/" + name)
		if err != nil {
			panic(err)
		}

		t[name] = template.Must(template.New(name).Parse(string(data)))
	}
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static").Handler(http.FileServer(assetFS()))
	r.HandleFunc("/", HandleIndex)

	fmt.Fprintln(os.Stderr, "Listening on port", Options.Port)
	if err := http.ListenAndServe(Options.Port, r); err != nil {
		panic(err)
	}
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	var data interface{}

	switch r.Method {
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return

	case "GET":

	case "POST":
		// Sanitize content before using
		content := template.HTMLEscapeString(r.FormValue("content"))
		if len(content) == 0 {
			data = "Error: Must provide content."
			goto HANDLE_GET
		}

		req := entities.NewRequest(apiKey)
		entityMap, err := req.FromPlainText(content)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResp := struct {
			HTML string `json:"html"`
		}{
			AnnotateDocument(req.Document(), entityMap),
		}

		// Return JSON response
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&jsonResp); err != nil {
			data = "Error: Must provide URL."
			goto HANDLE_GET
		}
		return
	}

HANDLE_GET:
	if err := t["index.html"].Execute(w, data); err != nil {
		panic(err)
	}
}

// AnnotateDocument highlights mentions within a document.
func AnnotateDocument(doc gcnl.Document, entityMap entities.Map) string {
	w := bytes.Buffer{}
	textSpanMap := make(map[int]entities.Entity)

	for t, _ := range entityMap {
		for _, e := range entityMap[t] {
			for _, m := range e.Mentions {
				textSpanMap[m.TextSpan.BeginOffset] = e
			}
		}
	}

	l := 0
	for i, b := range []byte(doc.Content()) {
		if l > 0 {
			l -= 1
			if l == 0 {
				w.WriteString("</span>")
			}
		} else if l == 0 {
			if e, ok := textSpanMap[i]; ok {
				for _, m := range e.Mentions {
					if m.TextSpan.BeginOffset == i {
						l = len(m.TextSpan.Content)
						break
					}
				}

				w.WriteString(fmt.Sprintf(`<span class="type-%s" data-toggle="tooltip" title="%s (%f)">`, e.Type, e.Type, e.Salience))
			}
		}
		w.WriteByte(b)
	}

	return w.String()
}
