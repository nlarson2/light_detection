// main.go
package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	root = flag.String("root", "./images", "directory of images")
	addr = flag.String("addr", ":8080", "listen address")
)

type Img struct {
	Name string
	URL  string
}

var page = template.Must(template.New("index").Parse(`
<!doctype html><meta charset="utf-8">
<title>Images</title>
<style>
  body{margin:0;font:14px system-ui;background:#111;color:#eee}
  header{padding:12px 16px;border-bottom:1px solid #222}
  .grid{padding:12px;display:grid;grid-template-columns:repeat(auto-fill,minmax(160px,1fr));gap:10px}
  .card{display:block;background:#161616;border:1px solid #222;border-radius:8px;overflow:hidden;text-decoration:none;color:inherit}
  .card img{width:100%;height:150px;object-fit:cover;display:block}
  .name{padding:8px;font-size:12px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;color:#aaa}
</style>
<header><h1 style="margin:0;font-size:16px">Images in {{.Dir}}</h1></header>
<main class="grid">
  {{range .Images}}
    <a class="card" href="{{.URL}}" target="_blank" title="{{.Name}}">
      <img loading="lazy" src="{{.URL}}" alt="{{.Name}}">
      <div class="name">{{.Name}}</div>
    </a>
  {{end}}
</main>
`))

func isImage(name string) bool {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return true
	default:
		return false
	}
}

func main() {
	flag.Parse()
	if err := os.MkdirAll(*root, 0o755); err != nil {
		log.Fatal(err)
	}

	// Serve raw files under /img/
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(*root))))

	// Simple index that lists images
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		entries, err := os.ReadDir(*root)
		if err != nil {
			http.Error(w, "failed to read dir", 500)
			return
		}
		var imgs []Img
		for _, e := range entries {
			if e.IsDir() || !isImage(e.Name()) {
				continue
			}
			imgs = append(imgs, Img{
				Name: e.Name(),
				URL:  "/img/" + e.Name(),
			})
		}
		w.Header().Set("Cache-Control", "no-store")
		_ = page.Execute(w, map[string]any{"Dir": *root, "Images": imgs})
	})

	log.Printf("serving %q on http://localhost%s\n", *root, *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
