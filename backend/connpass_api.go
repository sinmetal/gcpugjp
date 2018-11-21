package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/favclip/ucon"
	"github.com/favclip/ucon/swagger"
	"github.com/pkg/errors"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
)

func SetupConnpassAPI(swPlugin *swagger.Plugin) {
	api := ConnpassAPI{}
	tag := swPlugin.AddTag(&swagger.Tag{Name: "connpass", Description: "connpass"})
	var hInfo *swagger.HandlerInfo

	hInfo = swagger.NewHandlerInfo(api.HandlerCron)
	ucon.Handle(http.MethodGet, "/api/cron/1/connpass", hInfo)
	hInfo.Description, hInfo.Tags = "get from connpass start", []string{tag.Name}

	hInfo = swagger.NewHandlerInfo(api.Get)
	ucon.Handle(http.MethodGet, "/api/queue/1/connpass/{seriesId}", hInfo)
	hInfo.Description, hInfo.Tags = "get from connpass", []string{tag.Name}
}

// ConnpassResult is Connpass API Result
type ConnpassResult struct {
	Returned  int             `json:"results_returned"`
	Available int             `json:"results_available"`
	Start     int             `json:"results_start"`
	Events    []ConnpassEvent `json:"events"`
}

// ConnpassEvent is Connpass APi Return Event
type ConnpassEvent struct {
	EventID       int            `json:"event_id"`
	Title         string         `json:"title"`
	Catch         string         `json:"catch"`
	Description   string         `json:"description"`
	URL           string         `json:"event_url"`
	Tag           string         `json:"hash_tag"`
	StartedAt     time.Time      `json:"started_at"`
	EndedAt       time.Time      `json:"ended_at"`
	Limit         int            `json:"limit"`
	Etype         string         `json:"event_type"`
	Address       string         `json:"address"`
	Place         string         `json:"place"`
	Lat           string         `json:"lat"`
	Lon           string         `json:"lon"`
	OwnerID       int            `json:"owner_id"`
	OwnerNickname string         `json:"owner_nickname"`
	OwnerName     string         `json:"owner_display_name"`
	Series        ConnpassSeries `json:"series"`
	Accepted      int            `json:"accepted"`
	Waiting       int            `json:"waiting"`
	Updated       time.Time      `json:"updated_at"`
}

// ConnpassSeries is Series
// ConnpassのSeriesはグループを指す
type ConnpassSeries struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

// ConnpassAPI is connpass API Func Collection
// ConnpassAPIは https://connpass.com/about/api/ を実行してイベントを拾ってくる機能を持つ
type ConnpassAPI struct{}

type ConnpassAPIGetForm struct {
	SeriesID int `json:"seriesId" swagger:",in=query"`
}

// HandlerCron is Cron Handler
func (api *ConnpassAPI) HandlerCron(ctx context.Context) error {
	m := api.getSeriesIDMap()
	tl := []*taskqueue.Task{}
	for k, _ := range m {
		tl = append(tl, &taskqueue.Task{
			Path: fmt.Sprintf("/api/queue/1/connpass/%v", k),
		})
	}
	_, err := taskqueue.AddMulti(ctx, tl, "connpass")
	if err != nil {
		log.Errorf(ctx, "failed taskqueue.AddMulti err:%+v", err)
		return err
	}
	return nil
}

// Get is Connpass APIを実行して、新しいイベントがあれば、Datastoreに登録するCronから起動するためのAPI
func (api *ConnpassAPI) Get(ctx context.Context, form *ConnpassAPIGetForm) (*ConnpassResult, error) {
	result, err := api.getConnpassEvents(ctx, form.SeriesID)
	if err != nil {
		log.Errorf(ctx, "failed connpass events api: %+v", err)
		return nil, err
	}

	store, err := NewPugEventStore(ctx)
	if err != nil {
		log.Errorf(ctx, "failed New PugEventStore: %+v", err)
		return nil, err
	}

	sm := api.getSeriesIDMap()
	for _, v := range result.Events {
		j, err := json.Marshal(v)
		if err != nil {
			log.Errorf(ctx, "failed json.Marshal: %+v", err)
			return nil, err
		}
		log.Infof(ctx, "%s\n", string(j))

		pe := PugEvent{}
		pe.Limit = v.Limit
		pe.Waiting = v.Waiting
		pe.Accepted = v.Accepted
		pe.StartAt = v.StartedAt
		pe.EndAt = v.EndedAt
		pe.Title = v.Title
		pe.URL = v.URL
		pe.OrganizationID = sm[v.Series.ID]

		_, err = store.Upsert(ctx, &pe)
		if err != nil {
			// 重複エラーもあるので、失敗しても気にしない
			log.Warningf(ctx, "failed put PugEvent title=%s. err:%+v", pe.Title, err)
		}
	}

	j, err := json.Marshal(result.Events)
	if err != nil {
		log.Warningf(ctx, "failed json.Marshal: %+v", err)
	} else {
		log.Infof(ctx, "%s", string(j))
	}

	return result, nil
}

func (api *ConnpassAPI) getConnpassEvents(ctx context.Context, seriesId int) (*ConnpassResult, error) {
	url := fmt.Sprintf("https://connpass.com/api/v1/event/?series_id=%v", seriesId)
	fmt.Printf("connpass api url = %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed connpass event query")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed connpass event query: status = %s", resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed connpass event api result body")
	}
	defer resp.Body.Close()
	result := ConnpassResult{}
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, errors.New("failed connpass api result body to json.Unmarshal")
	}
	return &result, nil
}

// getSeriesIDMap is Watch対象のconnpass group一覧を取得する
// group idは https://connpass.com/api/v1/event/?keyword=GCPUG を叩くと見える
func (api *ConnpassAPI) getSeriesIDMap() map[int]string {
	return map[int]string{
		1898: "tokyo",
		2774: "beginners-tokyo",
		6273: "yokoyama",
		5270: "ibaraki",
		5424: "fukushima",
		2239: "nagoya",
		1658: "shonan",
		5297: "kyoto",
		5498: "nara",
		1422: "osaka",
		5271: "kobe",
		6362: "okayama",
		5812: "wakayama",
		4086: "hiroshima",
		4609: "shimane",
		6415: "kochi",
		1170: "fukuoka",
		4758: "kagoshima",
		3824: "okinawa",
	}

}
