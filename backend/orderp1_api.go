package backend

import (
	"context"
	"net/http"

	"github.com/favclip/ucon"
	"github.com/favclip/ucon/swagger"
)

func setupOrderP1(swPlugin *swagger.Plugin) {
	api := &OrderP1API{}
	tag := swPlugin.AddTag(&swagger.Tag{Name: "OrderP1", Description: "OrderP1 list"})
	var hInfo *swagger.HandlerInfo

	hInfo = swagger.NewHandlerInfo(api.Post)
	ucon.Handle(http.MethodPost, "/orderP1", hInfo)
	hInfo.Description, hInfo.Tags = "post to orderP1", []string{tag.Name}

	hInfo = swagger.NewHandlerInfo(api.Get)
	ucon.Handle(http.MethodGet, "/orderP1", hInfo)
	hInfo.Description, hInfo.Tags = "get to orderP1", []string{tag.Name}
}

// OrderP1API is OrderP1API
type OrderP1API struct{}

// OrderP1APIPostRequest is Request form
type OrderP1APIPostRequest struct {
}

// OrderP1APIPostResponse is Response Body
type OrderP1APIPostResponse struct {
}

// Post is OrderP1 Post Handler
func (api *OrderP1API) Post(c context.Context, form *OrderP1APIPostRequest) (*OrderP1APIPostResponse, error) {
	return nil, nil
}

// OrderP1APIGetRequest is Request form
type OrderP1APIGetRequest struct {
}

// OrderP1APIGetResponse is Response Body
type OrderP1APIGetResponse struct {
}

// Get is OrderP1 Get Handler
func (api *OrderP1API) Get(c context.Context, form *OrderP1APIGetRequest) (*OrderP1APIGetResponse, error) {
	return &OrderP1APIGetResponse{}, nil
}
