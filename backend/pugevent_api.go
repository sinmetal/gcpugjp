package backend

import (
	"context"
	"net/http"

	"github.com/favclip/ucon"
	"github.com/favclip/ucon/swagger"
	"google.golang.org/appengine/log"
)

func setupPugEventAPI(swPlugin *swagger.Plugin) {
	api := &PugEventAPI{}
	tag := swPlugin.AddTag(&swagger.Tag{Name: "PugEvent", Description: "PugEvent API list"})
	var hInfo *swagger.HandlerInfo

	hInfo = swagger.NewHandlerInfo(api.List)
	ucon.Handle(http.MethodGet, "/api/1/event", hInfo)
	hInfo.Description, hInfo.Tags = "get to gcpug event", []string{tag.Name}
}

// PugEventAPI is GCPUG Event API
type PugEventAPI struct{}

// PugEventAPIListRequest is GCPUG Event 一覧取得APIのリクエスト
type PugEventAPIListRequest struct {
}

// PugEventAPIListResponse is GCPUG Event 一覧取得APIのレスポンス
type PugEventAPIListResponse struct {
	List    []*PugEvent `json:"list"`
	HasNext bool        `json:"hasNext"`
	Cursor  string      `json:"cursor"`
}

// List is EventをStartAtの降順に100件返す API
func (api *PugEventAPI) List(ctx context.Context, form *PugEventAPIListRequest) (*PugEventAPIListResponse, error) {
	param := PugEventListParam{
		Limit: 100,
	}

	store, err := NewPugEventStore(ctx)
	if err != nil {
		log.Errorf(ctx, "failed NewPugEventStore: %+v", err)
		return nil, &HTTPError{Code: http.StatusInternalServerError, Message: "error"}
	}
	res, err := store.List(ctx, &param)
	if err != nil {
		log.Errorf(ctx, "failed PugEventStore.List: %+v", err)
		return nil, &HTTPError{Code: http.StatusInternalServerError, Message: "error"}
	}

	return &PugEventAPIListResponse{
		List:    res.List,
		HasNext: res.HasNext,
		Cursor:  res.NextCursor.String(),
	}, nil
}
