// Copyright (c) 2022, Direct Decisions Go client AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package directdecisions_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"directdecisions.com/directdecisions"
)

func TestClient_key(t *testing.T) {
	client, mux, _ := newClient(t, "xapp-1-my-authkey")

	mux.HandleFunc("/v1/votings", requireMethod("POST", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer xapp-1-my-authkey" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", jsonContentType)
		fmt.Fprintln(w, struct{}{})
	}))

	_, err := client.Votings.Create(context.Background(), []string{"Something", "Anything", "Nothing"})
	assertErrors(t, err, nil)
}

func TestClient_userAgent(t *testing.T) {
	client, mux, _ := newClient(t, "")

	mux.HandleFunc("/v1/votings", requireMethod("POST", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != directdecisions.UserAgent {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", jsonContentType)
		fmt.Fprintln(w, struct{}{})
	}))

	_, err := client.Votings.Create(context.Background(), []string{"Something", "Anything", "Nothing"})
	assertErrors(t, err, nil)

}

const jsonContentType = "application/json; charset=utf-8"

func newClient(t testing.TB, key string) (client *directdecisions.Client, mux *http.ServeMux, baseURL *url.URL) {
	t.Helper()

	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	baseURL, err := url.Parse(server.URL)
	assertErrors(t, err, nil)

	client = directdecisions.NewClient(key, &directdecisions.ClientOptions{
		BaseURL:    baseURL,
		HTTPClient: server.Client(),
	})

	t.Cleanup(server.Close)

	return client, mux, baseURL
}

func newStaticHandler(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", jsonContentType)
		fmt.Fprintln(w, body)
	}
}

func requireMethod(method string, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		f(w, r)
	}
}

func assertEqual(t testing.TB, name string, got, want any) {
	t.Helper()

	if name != "" {
		name += ": "
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%sgot %+v, want %+v", name, got, want)
	}
}

func assertErrors(t testing.TB, got error, want ...error) {
	t.Helper()

	for i, w := range want {
		if !errors.Is(got, w) {
			t.Fatalf("got %v error %[2]T %[2]v, want %[3]T %[3]v", i, got, want[i])
		}
	}
}
