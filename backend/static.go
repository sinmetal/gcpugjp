package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// ErrDirectory is target path is Directory
var ErrDirectory = errors.New("path is directory")

func handler(w http.ResponseWriter, r *http.Request) {
	fn := r.URL.Path[1:len(r.URL.Path)]
	if fn == "" || fn == "/" {
		fn = "index.html"
	}
	f, err := readFile(fn)
	if err != nil {
		if err == ErrDirectory {
			writeIndexHTML(w)
			return
		}
		if os.IsNotExist(err) {
			writeIndexHTML(w)
			return
		}
		if os.IsPermission(err) {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ext := filepath.Ext(fmt.Sprintf("static/%s", fn))

	ct := ""
	switch ext {
	default:
		ct = "application/octet-stream"
	case ".html", ".htm":
		ct = "text/html;charset=utf-8"
	case ".css":
		ct = "text/css;charset=utf-8"
	case ".js":
		ct = "text/javascript;charset=utf-8"
	case ".jpeg", ".jpg":
		ct = "image/jpeg"
	case ".png":
		ct = "image/png"
	case ".gif":
		ct = "image/gif"
	case ".txt":
		ct = "text/plain;charset=utf-8"
	case ".json":
		ct = "application/json;charset=utf-8"
	case ".pdf":
		ct = "application/pdf"
	case ".ico":
		ct = "image/x-icon"
	}

	w.Header().Set("Content-Type", ct)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(f)))
	w.WriteHeader(http.StatusOK)
	w.Write(f)
}

func readFile(path string) ([]byte, error) {
	p := fmt.Sprintf("static/%s", path)
	stat, err := os.Stat(p)
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, ErrDirectory
	}

	fp, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	return ioutil.ReadAll(fp)
}

func writeIndexHTML(w http.ResponseWriter) {
	fn := "index.html"
	f, err := readFile(fn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(f)))
	w.WriteHeader(http.StatusOK)
	w.Write(f)
}
