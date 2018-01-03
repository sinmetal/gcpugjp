package gcpug

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mjibson/goon"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func init() {
	api := ConnpassAPI{}
	http.HandleFunc("/cron/connpass", api.handler)
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

// PugEvent is Event Model
type PugEvent struct {
	Id             string    `datastore:"-" goon:"id" json:"id"`       // UUID
	OrganizationId string    `json:"organizationId"`                   // 支部Id
	Title          string    `json:"title" datastore:",noindex"`       // イベントタイトル
	Description    string    `json:"description" datastore:",noindex"` // イベント説明
	Url            string    `json:"url"`                              // イベント募集URL
	StartAt        time.Time `json:"startAt"`                          // 開催日時
	CreatedAt      time.Time `json:"createdAt"`                        // 作成日時
	UpdatedAt      time.Time `json:"updatedAt"`                        // 更新日時
}

// Create is Connpassから取得した情報を元にEventをDatastoreに登録する
func (pe *PugEvent) Create(ctx context.Context, g *goon.Goon) error {
	q := datastore.NewQuery("PugEvent")
	q = q.Filter("Url=", pe.Url)

	var pes []PugEvent
	_, err := q.GetAll(ctx, &pes)
	if err != nil {
		return err
	}
	if len(pes) > 0 {
		return errors.New("exists event")
	}

	return g.RunInTransaction(func(g *goon.Goon) error {
		stored := &PugEvent{
			Id: pe.Id,
		}
		err := g.Get(stored)
		if err == nil {
			return errors.New("conflict key") // TODO 専用Errorにする
		}
		if err != datastore.ErrNoSuchEntity {
			return err
		}

		_, err = g.Put(pe)
		if err != nil {
			return err
		}

		return nil
	}, nil)
}

// Load is PugEvent Load
func (pe *PugEvent) Load(ps []datastore.Property) error {
	if err := datastore.LoadStruct(pe, ps); err != nil {
		return err
	}
	return nil
}

// Save is PugEvent Save
func (pe *PugEvent) Save() ([]datastore.Property, error) {
	now := time.Now()
	pe.UpdatedAt = now

	if pe.CreatedAt.IsZero() {
		pe.CreatedAt = now
	}

	return datastore.SaveStruct(pe)
}

// ConnpassAPI is connpass API Func Collection
// ConnpassAPIは https://connpass.com/about/api/ を実行してイベントを拾ってくる機能を持つ
type ConnpassAPI struct{}

func (api *ConnpassAPI) handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	result, err := api.getConnpassEvents(ctx)
	if err != nil {
		log.Errorf(ctx, "failed connpass events api: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sm := api.getSeriesIDMap()
	for _, v := range result.Events {
		j, err := json.Marshal(v)
		if err != nil {
			log.Errorf(ctx, "failed json.Marshal: %+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Infof(ctx, "%s\n", string(j))

		pe := PugEvent{}
		pe.Id = uuid.New()
		pe.StartAt = v.StartedAt
		pe.Title = v.Title
		pe.Url = v.URL
		pe.OrganizationId = sm[v.Series.ID]

		g := goon.FromContext(ctx)
		if err := pe.Create(ctx, g); err != nil {
			// 重複エラーもあるので、失敗しても気にしない
			log.Warningf(ctx, "failed put PugEvent title=%s. err:%+v", pe.Title, err)
		}
	}
	j, err := json.Marshal(result.Events)
	if err != nil {
		log.Errorf(ctx, "failed json.Marshal: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (api *ConnpassAPI) getConnpassEvents(ctx context.Context) (*ConnpassResult, error) {
	client := urlfetch.Client(ctx)
	resp, err := client.Get(fmt.Sprintf("https://connpass.com/api/v1/event/?series_id=%s", api.getSeriesIDParam()))
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

func (api *ConnpassAPI) getSeriesIDParam() string {
	a := []string{}
	m := api.getSeriesIDMap()
	for k := range m {
		a = append(a, strconv.Itoa(k))
	}
	return strings.Join(a, ",")
}

func (api *ConnpassAPI) getSeriesIDMap() map[int]string {
	return map[int]string{
		1898: "tokyo",
		2239: "nagoya",
		1658: "shonan",
		1422: "osaka",
		4758: "kagoshima",
	}

}
