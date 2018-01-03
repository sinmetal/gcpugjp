package gcpug

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"google.golang.org/appengine"
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

// ConnpassAPI is connpass API Func Collection
// ConnpassAPIは https://connpass.com/about/api/ を実行してイベントを拾ってくる機能を持つ
type ConnpassAPI struct{}

func (api *ConnpassAPI) handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	client := urlfetch.Client(ctx)
	resp, err := client.Get(fmt.Sprintf("https://connpass.com/api/v1/event/?series_id=%d", 1898))
	if err != nil {
		log.Errorf(ctx, "failed connpass event query: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Errorf(ctx, "failed connpass event query: status = %s", resp.Status)
		http.Error(w, fmt.Sprintf("failed connpass event query: status = %s", resp.Status), http.StatusInternalServerError)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf(ctx, "failed connpass api result body: %+v", err)
		http.Error(w, "failed connpass event api result body", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	result := ConnpassResult{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Errorf(ctx, "failed connpass api result body to json.Unmarshal: %+v", err)
		http.Error(w, "failed connpass api result body to json.Unmarshal", http.StatusInternalServerError)
		return
	}
	for _, v := range result.Events {
		j, err := json.Marshal(v)
		if err != nil {
			log.Errorf(ctx, "failed json.Marshal: %+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Infof(ctx, "%s\n", string(j))
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
