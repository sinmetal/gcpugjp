package main

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/favclip/ucon"
	"github.com/favclip/ucon/swagger"
	"github.com/sinmetal/gcpugjp/backend"
	"google.golang.org/appengine"
)

func main() {
	ucon.Middleware(UseAppengineContext)
	ucon.Orthodox()
	ucon.Middleware(swagger.RequestValidator())

	swPlugin := swagger.NewPlugin(&swagger.Options{
		Object: &swagger.Object{
			Info: &swagger.Info{
				Title:   "GCPUG",
				Version: "v1",
			},
			Schemes: []string{"http", "https"},
		},
		DefinitionNameModifier: func(refT reflect.Type, defName string) string {
			if strings.HasSuffix(defName, "JSON") {
				return defName[:len(defName)-4]
			}
			return defName
		},
	})
	ucon.Plugin(swPlugin)

	backend.SetupOrderP1(swPlugin)
	backend.SetupConnpassAPI(swPlugin)
	backend.SetupPugEventAPI(swPlugin)
	backend.SetupOrganizationAPI(swPlugin)
	backend.SetupOrganizationAdminAPI(swPlugin)
	backend.SetupAppEngineEnvAPI(swPlugin)

	ucon.DefaultMux.Prepare()
	http.Handle("/api/", ucon.DefaultMux)
	http.HandleFunc("/", backend.StaticContentsHandler)

	appengine.Main()
}

// UseAppengineContext is UseAppengineContext
func UseAppengineContext(b *ucon.Bubble) error {
	if b.Context == nil {
		b.Context = b.R.Context()
	} else {
		// TODO contextのnestってできるんだっけ？
	}

	return b.Next()
}
