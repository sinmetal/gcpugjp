package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/favclip/testerator"
)

// TestOrganizationAPI_List is とりあえず 200 OKが返ってくることを確認
func TestOrganizationAPI_List(t *testing.T) {
	inst, _, err := testerator.SpinUp()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer testerator.SpinDown()

	r, err := inst.NewRequest(http.MethodGet, "/api/1/organization", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		b, _ := ioutil.ReadAll(w.Body)
		t.Fatalf("unexpected %d, expected 200, body=%s", w.Code, string(b))
	}
}
