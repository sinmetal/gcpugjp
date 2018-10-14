package backend

import (
	"context"
	"net/http"

	"github.com/favclip/ucon"
	"github.com/favclip/ucon/swagger"
	"google.golang.org/appengine/log"
)

func SetupOrganizationAdminAPI(swPlugin *swagger.Plugin) {
	api := &OrganizationAdminAPI{}

	tag := swPlugin.AddTag(&swagger.Tag{Name: "Organization", Description: "Organization Admin API List"})
	var hInfo *swagger.HandlerInfo

	hInfo = swagger.NewHandlerInfo(api.Post)
	ucon.Handle(http.MethodPost, "/api/admin/1/organization", hInfo)
	hInfo.Description, hInfo.Tags = "post to organization", []string{tag.Name}
}

// OrganizationAdminAPI is Organization Admin API Functions
type OrganizationAdminAPI struct{}

// OrganizationAdminAPIPostRequest is Organization Admin Post API Request
type OrganizationAdminAPIPostRequest struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	LogoURL string `json:"logoUrl"`
	Order   int    `json:"order"`
}

// OrganizationAdminAPIPostResponse is Organization Admin Post API Response
type OrganizationAdminAPIPostResponse struct {
	*Organization
}

// Post is Organizationを登録する
func (api *OrganizationAdminAPI) Post(ctx context.Context, form *OrganizationAdminAPIPostRequest) (*OrganizationAdminAPIPostResponse, error) {
	store, err := NewOrganizationStore(ctx)
	if err != nil {
		log.Errorf(ctx, "failed NewOrganizationStore: %+v", err)
		return nil, &HTTPError{Code: http.StatusInternalServerError, Message: "error"}
	}

	// TODO Validation

	e := Organization{
		Name:    form.Name,
		URL:     form.URL,
		LogoURL: form.LogoURL,
		Order:   form.Order,
	}
	key := store.NameKey(ctx, form.ID)

	o, err := store.Create(ctx, key, &e)
	if err != nil {
		log.Errorf(ctx, "failed OrganizationStore.Create: %+v", err)
		return nil, &HTTPError{Code: http.StatusInternalServerError, Message: "error"}
	}

	return &OrganizationAdminAPIPostResponse{o}, nil
}
