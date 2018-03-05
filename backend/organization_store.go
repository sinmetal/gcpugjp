package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.mercari.io/datastore"
	"google.golang.org/appengine/log"
)

var _ datastore.PropertyLoadSaver = &Organization{}
var _ datastore.KeyLoader = &Organization{}

// Organization is 支部 Datastore Entity Model
// qbg
type Organization struct {
	Key           datastore.Key `datastore:"-" json:"-"`
	KeyStr        string        `datastore:"-" json:"key"` // Json変換用 Encode Key
	Name          string        // 支部名 : example GCPUG東京
	URL           string        // connpassなどのURL
	LogoURL       string        // 支部のLogoURL
	Order         int           // Sort順 市区町村コードを入れている
	CreatedAt     time.Time     `json:"createdAt"` // 作成日時
	UpdatedAt     time.Time     `json:"updatedAt"` // 更新日時
	SchemaVersion int           `json:"-"`
}

// LoadKey is Entity Load時にKeyを設定する
func (e *Organization) LoadKey(ctx context.Context, k datastore.Key) error {
	e.Key = k

	return nil
}

// Load is Entity Load時に呼ばれる
func (e *Organization) Load(ctx context.Context, ps []datastore.Property) error {
	err := datastore.LoadStruct(ctx, e, ps)
	if err != nil {
		return err
	}

	return nil
}

// Save is Entity Save時に呼ばれる
func (e *Organization) Save(ctx context.Context) ([]datastore.Property, error) {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
	e.UpdatedAt = time.Now()
	e.SchemaVersion = 1

	return datastore.SaveStruct(ctx, e)
}

// OrganizationStore is OrganizationのDatastoreの操作を司る
type OrganizationStore struct {
	DatastoreClient datastore.Client
}

// NewOrganizationStore is OrganizationStoreを作成
func NewOrganizationStore(ctx context.Context) (*OrganizationStore, error) {
	ds, err := FromContext(ctx)
	if err != nil {
		log.Errorf(ctx, "failed Datastore New Client: %+v", err)
		return nil, err
	}
	return &OrganizationStore{ds}, nil
}

// Kind is OrganizationのKindを返す
func (store *OrganizationStore) Kind() string {
	return "Organization"
}

// NameKey is Organizationの指定したNameを利用したKeyを生成する
func (store *OrganizationStore) NameKey(ctx context.Context, name string) datastore.Key {
	return store.DatastoreClient.NameKey(store.Kind(), name, nil)
}

// Create is OrganizationをDatastoreにputする
func (store *OrganizationStore) Create(ctx context.Context, key datastore.Key, e *Organization) (*Organization, error) {
	ds := store.DatastoreClient

	_, err := ds.Put(ctx, key, e)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed put Organization to Datastore. key=%v", key))
	}
	e.Key = key
	e.KeyStr = key.Encode()
	return e, nil
}

// Get is OrganizationをDatastoreからgetする
func (store *OrganizationStore) Get(ctx context.Context, key datastore.Key) (*Organization, error) {
	ds := store.DatastoreClient

	var e Organization
	err := ds.Get(ctx, key, &e)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed get Organization from Datastore. key=%v", key))
	}
	e.KeyStr = key.Encode()

	return &e, nil
}
