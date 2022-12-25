// Copyright (c) 2022, Direct Decisions Go client AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package directdecisions_test

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestRate(t *testing.T) {
	client, mux, _ := newClient(t, "")

	mux.HandleFunc("/v1/votings", requireMethod("POST", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Limit", "10")
		w.Header().Set("X-RateLimit-Remaining", "20")
		w.Header().Set("X-RateLimit-Reset", "500")
		w.Header().Set("Retry-After", "30")
	}))

	now := time.Now()
	_, err := client.Votings.Create(context.Background(), []string{"Something", "Anything", "Nothing"})
	if err != nil {
		t.Fatal(err)
	}

	got := client.Rate()

	assertEqual(t, "limit", got.Limit, 10)
	assertEqual(t, "remaining", got.Remaining, 20)
	if got.Reset.Round(5*time.Second) != now.Add(500*time.Second).Round(5*time.Second) {
		t.Errorf("got reset %s, want %s", got.Reset.Round(time.Second), now.Add(500*time.Second).Round(time.Second))
	}
	if got.Retry.Round(time.Second) != now.Add(30*time.Second).Round(time.Second) {
		t.Errorf("got retry %s, want %s", got.Retry.Round(time.Second), now.Add(30*time.Second).Round(time.Second))
	}
}
