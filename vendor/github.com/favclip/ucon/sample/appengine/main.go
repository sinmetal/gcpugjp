//go:generate jwg -output model_json.go -transcripttag swagger .
//go:generate qbg -output model_query.go .

package example_appengine

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/favclip/ucon"
	"github.com/favclip/ucon/swagger"
	"google.golang.org/appengine"
)

func init() {
	var _ ucon.HTTPErrorResponse = &HttpError{}

	ucon.Middleware(UseAppengineContext)
	ucon.Orthodox()
	ucon.Middleware(swagger.RequestValidator())

	swPlugin := swagger.NewPlugin(&swagger.Options{
		Object: &swagger.Object{
			Info: &swagger.Info{
				Title:   "Todo list",
				Version: "v1",
			},
			Schemes: []string{"http"},
		},
		DefinitionNameModifier: func(refT reflect.Type, defName string) string {
			if strings.HasSuffix(defName, "JSON") {
				return defName[:len(defName)-4]
			}
			return defName
		},
	})
	ucon.Plugin(swPlugin)

	ucon.HandleFunc("GET", "/swagger-ui/", func(w http.ResponseWriter, r *http.Request) {
		localPath := "./node_modules/swagger-ui/dist/" + r.URL.Path[len("/swagger-ui/"):]
		http.ServeFile(w, r, localPath)
	})

	setupTodo(swPlugin)

	ucon.HandleFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		localPath := "./public/" + r.URL.Path[len("/"):]
		http.ServeFile(w, r, localPath)
	})

	ucon.DefaultMux.Prepare()
	http.Handle("/", ucon.DefaultMux)
}

func UseAppengineContext(b *ucon.Bubble) error {
	if b.Context == nil {
		b.Context = appengine.NewContext(b.R)
	} else {
		b.Context = appengine.WithContext(b.Context, b.R)
	}

	return b.Next()
}
