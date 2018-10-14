package main

import (
	"net/http"

	"github.com/favclip/ucon"
	"github.com/favclip/ucon/swagger"
	"google.golang.org/appengine/log"
)

func setupOrganizationAPI(swPlugin *swagger.Plugin) {
	api := &OrganizationAPI{}
	tag := swPlugin.AddTag(&swagger.Tag{Name: "Organization", Description: "Organization API list"})
	var hInfo *swagger.HandlerInfo

	hInfo = swagger.NewHandlerInfo(api.List)
	ucon.Handle(http.MethodGet, "/api/1/organization", hInfo)
	hInfo.Description, hInfo.Tags = "get to organization", []string{tag.Name}
}

// OrganizationAPI is Organizationに関するAPI
type OrganizationAPI struct{}

// OrganizationAPIListResponse is OrganizationAPI ListのResponse
type OrganizationAPIListResponse struct {
	List    []*Organization `json:"list"`
	HasNext bool            `json:"hasNext"`
	Cursor  string          `json:"cursor"`
}

// List is OrganizationをOrder順に全件返す API
// Interfaceを他のListと合わせてHasNext, Cursorを返すようにしているが、Organizationがそんなに増えることはないので、使ってはない
func (api *OrganizationAPI) List(r *http.Request) (*OrganizationAPIListResponse, error) {
	ctx := r.Context()

	store, err := NewOrganizationStore(ctx)
	if err != nil {
		log.Errorf(ctx, "failed NewOrganizationStore: %+v", err)
		return nil, &HTTPError{Code: http.StatusInternalServerError, Message: "error"}
	}
	ol, err := store.ListAll(ctx)
	if err != nil {
		log.Errorf(ctx, "failed OrganizationStore.List: %+v", err)
		return nil, &HTTPError{Code: http.StatusInternalServerError, Message: "error"}
	}

	return &OrganizationAPIListResponse{
		List:    ol,
		HasNext: false,
		Cursor:  "",
	}, nil
}
