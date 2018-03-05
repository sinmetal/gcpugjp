package backend

import (
	"testing"

	"github.com/favclip/testerator"
)

func TestOrganizationStore_Create(t *testing.T) {
	_, ctx, err := testerator.SpinUp() // gae/pythonのインスタンスが無ければ起動、あれば使いまわす
	if err != nil {
		t.Fatal(err.Error())
	}
	defer testerator.SpinDown() // プロセスをシャットダウンせずに、Datastoreなどの内容をクリアする

	store, err := NewOrganizationStore(ctx)
	if err != nil {
		t.Fatal(err)
	}

	e := Organization{
		Name:    "GCPUG東京",
		URL:     "https://gcpug-tokyo.connpass.com/",
		LogoURL: "//gcpug.jp/images/gcpug_tokyo.png",
		Order:   13103,
	}
	key := store.NameKey(ctx, "tokyo")
	stored, err := store.Create(ctx, key, &e)
	if err != nil {
		t.Fatal(err)
	}
	if stored.Key == nil {
		t.Fatalf("Key is Empty.")
	}
	if stored.KeyStr == "" {
		t.Fatalf("KeyStr is Empty")
	}
	if e, g := e.Name, stored.Name; e != g {
		t.Fatalf("expected Name is %s; got %s", e, g)
	}
	if e, g := e.URL, stored.URL; e != g {
		t.Fatalf("expected URL is %s; got %s", e, g)
	}
	if e, g := e.LogoURL, stored.LogoURL; e != g {
		t.Fatalf("expected LogoURL is %s; got %s", e, g)
	}
	if e, g := e.Order, stored.Order; e != g {
		t.Fatalf("expected Order is %d; got %d", e, g)
	}
	if stored.CreatedAt.IsZero() {
		t.Fatalf("CreatedAt is Zero")
	}
	if stored.UpdatedAt.IsZero() {
		t.Fatalf("UpdatedAt is Zero")
	}
	if stored.SchemaVersion == 0 {
		t.Fatalf("SchemaVersion is Zero")
	}
}

func TestOrganizationStore_Get(t *testing.T) {
	_, ctx, err := testerator.SpinUp() // gae/pythonのインスタンスが無ければ起動、あれば使いまわす
	if err != nil {
		t.Fatal(err.Error())
	}
	defer testerator.SpinDown() // プロセスをシャットダウンせずに、Datastoreなどの内容をクリアする

	store, err := NewOrganizationStore(ctx)
	if err != nil {
		t.Fatal(err)
	}

	e := Organization{
		Name:    "GCPUG東京",
		URL:     "https://gcpug-tokyo.connpass.com/",
		LogoURL: "//gcpug.jp/images/gcpug_tokyo.png",
		Order:   13103,
	}
	key := store.NameKey(ctx, "tokyo")
	_, err = store.Create(ctx, key, &e)
	if err != nil {
		t.Fatal(err)
	}
	stored, err := store.Get(ctx, key)
	if err != nil {
		t.Fatal(err)
	}
	if stored.Key == nil {
		t.Fatalf("Key is Empty.")
	}
	if stored.KeyStr == "" {
		t.Fatalf("KeyStr is Empty")
	}
	if e, g := e.Name, stored.Name; e != g {
		t.Fatalf("expected Name is %s; got %s", e, g)
	}
	if e, g := e.URL, stored.URL; e != g {
		t.Fatalf("expected URL is %s; got %s", e, g)
	}
	if e, g := e.LogoURL, stored.LogoURL; e != g {
		t.Fatalf("expected LogoURL is %s; got %s", e, g)
	}
	if e, g := e.Order, stored.Order; e != g {
		t.Fatalf("expected Order is %d; got %d", e, g)
	}
	if stored.CreatedAt.IsZero() {
		t.Fatalf("CreatedAt is Zero")
	}
	if stored.UpdatedAt.IsZero() {
		t.Fatalf("UpdatedAt is Zero")
	}
	if stored.SchemaVersion == 0 {
		t.Fatalf("SchemaVersion is Zero")
	}
}
