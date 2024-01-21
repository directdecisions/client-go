// Copyright (c) 2022, Direct Decisions Go client AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package directdecisions_test

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"directdecisions.com/directdecisions"
)

var votingsServiceVotingWant = &directdecisions.Voting{
	ID:      "40f80454800b2bd7c172",
	Choices: []string{"Margarita", "Diavola", "Capricciosa"},
}

func TestVotingsService_Voting(t *testing.T) {
	client, mux, _ := newClient(t, "")

	mux.HandleFunc("/v1/votings/40f80454800b2bd7c172", requireMethod("GET", newStaticHandler(`{
		"id": "40f80454800b2bd7c172",
		"choices": ["Margarita", "Diavola", "Capricciosa"]
	}`)))

	got, err := client.Votings.Voting(context.Background(), "40f80454800b2bd7c172")
	assertErrors(t, err, nil)

	assertEqual(t, "", got, votingsServiceVotingWant)
}

func TestVotingsService_Create(t *testing.T) {
	client, mux, _ := newClient(t, "")

	type createVotingRequest struct {
		Choices []string `json:"choices"`
	}

	mux.HandleFunc("/v1/votings", requireMethod("POST", func(w http.ResponseWriter, r *http.Request) {
		var request createVotingRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(request.Choices, []string{"Margarita", "Diavola", "Capricciosa"}) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		d, err := json.Marshal(votingsServiceVotingWant)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", jsonContentType)
		_, _ = w.Write(d)
	}))

	got, err := client.Votings.Create(context.Background(), []string{"Margarita", "Diavola", "Capricciosa"})
	assertErrors(t, err, nil)

	assertEqual(t, "", got, votingsServiceVotingWant)
}

func TestVotingsService_Set(t *testing.T) {
	client, mux, _ := newClient(t, "")

	type setChoiceRequest struct {
		Choice string `json:"choice"`
		Index  int    `json:"index"`
	}

	type setChoiceResponse struct {
		Choices []string `json:"choices"`
	}

	want := setChoiceResponse{
		Choices: []string{"Margarita", "Capricciosa", "Diavola"},
	}

	mux.HandleFunc("/v1/votings/40f80454800b2bd7c172/choices", requireMethod("POST", func(w http.ResponseWriter, r *http.Request) {
		var request *setChoiceRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(&setChoiceRequest{
			Choice: "Diavola",
			Index:  2,
		}, request) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		d, err := json.Marshal(want)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", jsonContentType)
		_, _ = w.Write(d)
	}))

	got, err := client.Votings.Set(context.Background(), "40f80454800b2bd7c172", "Diavola", 2)
	assertErrors(t, err, nil)

	assertEqual(t, "", got, want.Choices)
}

func TestVotingsService_Delete(t *testing.T) {
	client, mux, _ := newClient(t, "")

	mux.HandleFunc("/v1/votings/40f80454800b2bd7c172", requireMethod("DELETE", func(w http.ResponseWriter, r *http.Request) {}))

	err := client.Votings.Delete(context.Background(), "40f80454800b2bd7c172")
	assertErrors(t, err, nil)

}

func TestVotingsService_Vote(t *testing.T) {
	client, mux, _ := newClient(t, "")

	type voteRequest struct {
		Ballot map[string]int `json:"ballot"`
	}

	type voteResponse struct {
		Revoted bool `json:"revoted"`
	}

	want := voteResponse{
		Revoted: true,
	}

	mux.HandleFunc("/v1/votings/40f80454800b2bd7c172/ballots/leonardo", requireMethod("POST", func(w http.ResponseWriter, r *http.Request) {
		var request *voteRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(&voteRequest{
			Ballot: map[string]int{"Diavola": 1},
		}, request) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		d, err := json.Marshal(want)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", jsonContentType)
		_, _ = w.Write(d)
	}))

	revoted, err := client.Votings.Vote(context.Background(), "40f80454800b2bd7c172", "leonardo", map[string]int{
		"Diavola": 1,
	})
	assertErrors(t, err, nil)

	assertEqual(t, "", revoted, true)
}

func TestVotingsService_Unvote(t *testing.T) {
	client, mux, _ := newClient(t, "")

	mux.HandleFunc("/v1/votings/40f80454800b2bd7c172/ballots/leonardo", requireMethod("DELETE", func(w http.ResponseWriter, r *http.Request) {}))

	if err := client.Votings.Unvote(context.Background(), "40f80454800b2bd7c172", "leonardo"); err != nil {
		t.Fatal(err)
	}
}

func TestVotingsService_Ballot(t *testing.T) {
	client, mux, _ := newClient(t, "")

	type voteResponse struct {
		Ballot map[string]int `json:"ballot"`
	}

	mux.HandleFunc("/v1/votings/40f80454800b2bd7c172/ballots/leonardo", requireMethod("GET", func(w http.ResponseWriter, r *http.Request) {
		d, err := json.Marshal(voteResponse{
			Ballot: map[string]int{"Diavola": 1},
		})
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", jsonContentType)
		_, _ = w.Write(d)
	}))

	ballot, err := client.Votings.Ballot(context.Background(), "40f80454800b2bd7c172", "leonardo")
	assertErrors(t, err, nil)

	assertEqual(t, "", ballot, map[string]int{"Diavola": 1})
}

func TestVotingsService_Results(t *testing.T) {
	client, mux, _ := newClient(t, "")

	mux.HandleFunc("/v1/votings/40f80454800b2bd7c172/results", requireMethod("GET", newStaticHandler(`{
		"results": [
			{
				"choice": "Diavola",
				"index": 2,
				"wins": 1,
				"percentage": 50,
				"strength": 10,
				"advantage": 5
			},
			{
				"choice": "Margherita",
				"index": 1,
				"wins": 0,
				"percentage": 0,
				"strength": 3,
				"advantage": 1
			},
			{
				"choice": "Peperoni",
				"index": 0,
				"wins": 1,
				"percentage": 50,
				"strength": 1,
				"advantage": 0
			}
		],
		"tie": true
	}`)))

	results, tie, err := client.Votings.Results(context.Background(), "40f80454800b2bd7c172")
	assertErrors(t, err, nil)

	assertEqual(t, "", results, []directdecisions.Result{
		{
			Choice:     "Diavola",
			Index:      2,
			Wins:       1,
			Percentage: 50,
			Strength:   10,
			Advantage:  5,
		},
		{
			Choice:     "Margherita",
			Index:      1,
			Wins:       0,
			Percentage: 0,
			Strength:   3,
			Advantage:  1,
		},
		{
			Choice:     "Peperoni",
			Index:      0,
			Wins:       1,
			Percentage: 50,
			Strength:   1,
			Advantage:  0,
		},
	})
	assertEqual(t, "", tie, true)
}

func TestVotingsService_Duels(t *testing.T) {
	client, mux, _ := newClient(t, "")

	mux.HandleFunc("/v1/votings/40f80454800b2bd7c172/results/duels", requireMethod("GET", newStaticHandler(`{
		"results": [
			{
				"choice": "Diavola",
				"index": 2,
				"wins": 1,
				"percentage": 50,
				"strength": 10,
				"advantage": 5
			},
			{
				"choice": "Margherita",
				"index": 1,
				"wins": 0,
				"percentage": 0,
				"strength": 3,
				"advantage": 1
			},
			{
				"choice": "Peperoni",
				"index": 0,
				"wins": 1,
				"percentage": 50,
				"strength": 1,
				"advantage": 0
			}
		],
		"tie": true,
		"duels": [
			{
				"left": {
					"choice": "Peperoni",
					"index": 0,
					"strength": 1
				},
				"right": {
					"choice": "Margherita",
					"index": 1,
					"strength": 0
				}
			},
			{
				"left": {
					"choice": "Peperoni",
					"index": 0,
					"strength": 1
				},
				"right": {
					"choice": "Diavola",
					"index": 2,
					"strength": 1
				}
			},
			{
				"left": {
					"choice": "Margherita",
					"index": 1,
					"strength": 0
				},
				"right": {
					"choice": "Diavola",
					"index": 2,
					"strength": 1
				}
			}
		]
	}`)))

	results, duels, tie, err := client.Votings.Duels(context.Background(), "40f80454800b2bd7c172")
	assertErrors(t, err, nil)

	assertEqual(t, "", results, []directdecisions.Result{
		{
			Choice:     "Diavola",
			Index:      2,
			Wins:       1,
			Percentage: 50,
			Strength:   10,
			Advantage:  5,
		},
		{
			Choice:     "Margherita",
			Index:      1,
			Wins:       0,
			Percentage: 0,
			Strength:   3,
			Advantage:  1,
		},
		{
			Choice:     "Peperoni",
			Index:      0,
			Wins:       1,
			Percentage: 50,
			Strength:   1,
			Advantage:  0,
		},
	})
	assertEqual(t, "", tie, true)
	assertEqual(t, "", duels, []directdecisions.Duel{
		{
			Left: directdecisions.ChoiceStrength{
				Choice:   "Peperoni",
				Index:    0,
				Strength: 1,
			},
			Right: directdecisions.ChoiceStrength{
				Choice:   "Margherita",
				Index:    1,
				Strength: 0,
			},
		},
		{
			Left: directdecisions.ChoiceStrength{
				Choice:   "Peperoni",
				Index:    0,
				Strength: 1,
			},
			Right: directdecisions.ChoiceStrength{
				Choice:   "Diavola",
				Index:    2,
				Strength: 1,
			},
		},
		{
			Left: directdecisions.ChoiceStrength{
				Choice:   "Margherita",
				Index:    1,
				Strength: 0,
			},
			Right: directdecisions.ChoiceStrength{
				Choice:   "Diavola",
				Index:    2,
				Strength: 1,
			},
		},
	})
}
