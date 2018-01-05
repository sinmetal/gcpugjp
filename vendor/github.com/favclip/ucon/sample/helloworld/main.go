package main

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/favclip/ucon"
)

func main() {
	ucon.Orthodox()

	ucon.Middleware(NowInJST)
	ucon.Middleware(Logger)

	ucon.HandleFunc("GET", "/", func(w http.ResponseWriter, r *http.Request, now time.Time) {
		w.Write([]byte(fmt.Sprintf("Hello World! : %s", now.Format("2006/01/02 15:04:05"))))
	})

	ucon.HandleFunc("GET", "/a", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/a"))
	})
	//	ucon.HandleFunc("GET", "/a/", func(w http.ResponseWriter, r *http.Request) {
	//		w.Write([]byte("/a/"))
	//	})
	//	ucon.HandleFunc("GET", "/a/b/", func(w http.ResponseWriter, r *http.Request) {
	//		w.Write([]byte("/a/b"))
	//	})

	ucon.ListenAndServe(":8080")
}

func Logger(b *ucon.Bubble) error {
	fmt.Printf("Received: %s %s\n", b.R.Method, b.R.URL.String())
	return b.Next()
}

func NowInJST(b *ucon.Bubble) error {
	for idx, argT := range b.ArgumentTypes {
		if argT == reflect.TypeOf(time.Time{}) {
			b.Arguments[idx] = reflect.ValueOf(time.Now())
			break
		}
	}
	return b.Next()
}
