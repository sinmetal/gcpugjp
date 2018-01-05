package backend

import (
	"testing"
	"time"

	"github.com/favclip/testerator"
)

func TestPugEventStore_Create(t *testing.T) {
	_, ctx, err := testerator.SpinUp() // gae/pythonのインスタンスが無ければ起動、あれば使いまわす
	if err != nil {
		t.Fatal(err.Error())
	}
	defer testerator.SpinDown() // プロセスをシャットダウンせずに、Datastoreなどの内容をクリアする

	store, err := NewPugEventStore(ctx)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	e := PugEvent{
		OrganizationID: "tokyo",
		Title:          "GCPUG Day",
		Description:    "GCPUGやるぞー！",
		URL:            "https://gcpug.jp",
		StartAt:        time.Now(),
		EndAt:          time.Now(),
	}

	stored, err := store.Create(ctx, &e)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if stored.Key == nil {
		t.Fatalf("Key is Empty.")
	}
	if e, g := e.OrganizationID, stored.OrganizationID; e != g {
		t.Fatalf("expected OrganizationID %s; got %s", e, g)
	}
	if e, g := e.Title, stored.Title; e != g {
		t.Fatalf("expected Title %s; got %s", e, g)
	}
	if e, g := e.Description, stored.Description; e != g {
		t.Fatalf("expected Description %s; got %s", e, g)
	}
	if e, g := e.URL, stored.URL; e != g {
		t.Fatalf("expected URL %s; got %s", e, g)
	}
	if e, g := e.StartAt, stored.StartAt; !EqualTime(e, g) {
		t.Fatalf("expected StartAt %s; got %s", e, g)
	}
	if e, g := e.EndAt, stored.EndAt; !EqualTime(e, g) {
		t.Fatalf("expected EndAt %s; got %s", e, g)
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

func TestPugEventStore_Get(t *testing.T) {
	_, ctx, err := testerator.SpinUp() // gae/pythonのインスタンスが無ければ起動、あれば使いまわす
	if err != nil {
		t.Fatal(err.Error())
	}
	defer testerator.SpinDown() // プロセスをシャットダウンせずに、Datastoreなどの内容をクリアする

	store, err := NewPugEventStore(ctx)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	e := PugEvent{
		OrganizationID: "tokyo",
		Title:          "GCPUG Day",
		Description:    "GCPUGやるぞー！",
		URL:            "https://gcpug.jp",
		StartAt:        time.Now(),
		EndAt:          time.Now(),
	}

	_, err = store.Create(ctx, &e)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	stored, err := store.Get(ctx, e.Key)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if stored.Key == nil {
		t.Fatalf("Key is Empty.")
	}
	if e, g := e.OrganizationID, stored.OrganizationID; e != g {
		t.Fatalf("expected OrganizationID %s; got %s", e, g)
	}
	if e, g := e.Title, stored.Title; e != g {
		t.Fatalf("expected Title %s; got %s", e, g)
	}
	if e, g := e.Description, stored.Description; e != g {
		t.Fatalf("expected Description %s; got %s", e, g)
	}
	if e, g := e.URL, stored.URL; e != g {
		t.Fatalf("expected URL %s; got %s", e, g)
	}
	if e, g := e.StartAt, stored.StartAt; !EqualTime(e, g) {
		t.Fatalf("expected StartAt %s; got %s", e, g)
	}
	if e, g := e.EndAt, stored.EndAt; !EqualTime(e, g) {
		t.Fatalf("expected EndAt %s; got %s", e, g)
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
