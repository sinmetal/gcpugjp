package backend

import (
	"context"
	"net/http"

	"github.com/favclip/ucon"
	"github.com/favclip/ucon/swagger"
	"google.golang.org/appengine"
)

func setupAppEngineEnvAPI(swPlugin *swagger.Plugin) {
	api := &AppEngineEnvAPI{}

	tag := swPlugin.AddTag(&swagger.Tag{Name: "AppEngineEnv", Description: "AppEngineEnv API List"})
	var hInfo *swagger.HandlerInfo

	hInfo = swagger.NewHandlerInfo(api.Get)
	ucon.Handle(http.MethodGet, "/api/1/appEngineEnv", hInfo)
	hInfo.Description, hInfo.Tags = "get to appengineenv", []string{tag.Name}
}

// AppEngineEnvAPI is App Engineの環境 API
type AppEngineEnvAPI struct{}

// AppEngineEnvAPIGetResponse is AppEngineEnvAPI Get Response
type AppEngineEnvAPIGetResponse struct {
	Datacenter     string `json:"datacenter"`
	ServerSoftware string `json:"serverSoftware"`
	ModuleName     string `json:"moduleName"`
	VersionID      string `json:"versionId"`
}

// Get is App Engine環境状況を返す API Handler
func (api *AppEngineEnvAPI) Get(ctx context.Context) (*AppEngineEnvAPIGetResponse, error) {
	return &AppEngineEnvAPIGetResponse{
		Datacenter:     appengine.Datacenter(ctx),
		ServerSoftware: appengine.ServerSoftware(),
		ModuleName:     appengine.ModuleName(ctx),
		VersionID:      appengine.VersionID(ctx),
	}, nil
}
