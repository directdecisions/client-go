// Copyright (c) 2022, Direct Decisions Go client AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package directdecisions

import (
	"context"
	"net/http"
	"net/url"
)

// VotingsService provides information and methods to manage votings and ballots.
type VotingsService service

// Voting holds information about votings and ballots.
type Voting struct {
	ID      string   `json:"id"`
	Choices []string `json:"choices,omitempty"`
}

// Voting returns a specific voting referenced by its ID.
func (s *VotingsService) Voting(ctx context.Context, votingID string) (v *Voting, err error) {
	err = s.client.request(ctx, http.MethodGet, "v1/votings/"+url.PathEscape(votingID), nil, &v)
	return v, err
}

// Create adds a new voting with a provided choices.
func (s *VotingsService) Create(ctx context.Context, choices []string) (v *Voting, err error) {

	type createVotingRequest struct {
		Choices []string `json:"choices"`
	}

	err = s.client.request(ctx, http.MethodPost, "v1/votings", createVotingRequest{
		Choices: choices,
	}, &v)
	return v, err
}

// Set adds, moves or removes a choice in a voting.
func (s *VotingsService) Set(ctx context.Context, votingID, choice string, index int) (choices []string, err error) {

	type setChoiceRequest struct {
		Choice string `json:"choice"`
		Index  int    `json:"index"`
	}

	type setChoiceResponse struct {
		Choices []string `json:"choices"`
	}

	var response *setChoiceResponse
	if err = s.client.request(ctx, http.MethodPost, "v1/votings/"+url.PathEscape(votingID)+"/choices", setChoiceRequest{
		Choice: choice,
		Index:  index,
	}, &response); err != nil {
		return nil, err
	}
	return response.Choices, nil
}

// Delete removes a voting referenced by its ID.
func (s *VotingsService) Delete(ctx context.Context, votingID string) (err error) {
	return s.client.request(ctx, http.MethodDelete, "v1/votings/"+url.PathEscape(votingID), nil, nil)
}

func (s *VotingsService) Ballot(ctx context.Context, votingID, voterID string) (ballot map[string]int, err error) {

	type ballotResponse struct {
		Ballot map[string]int `json:"ballot"`
	}

	var response *ballotResponse
	if err = s.client.request(ctx, http.MethodGet, "v1/votings/"+url.PathEscape(votingID)+"/ballots/"+url.PathEscape(voterID), nil, &response); err != nil {
		return nil, err
	}
	return response.Ballot, nil
}

func (s *VotingsService) Vote(ctx context.Context, votingID, voterID string, ballot map[string]int) (revoted bool, err error) {

	type voteRequest struct {
		Ballot map[string]int `json:"ballot"`
	}

	type voteResponse struct {
		Revoted bool `json:"revoted"`
	}

	var response *voteResponse
	if err = s.client.request(ctx, http.MethodPost, "v1/votings/"+url.PathEscape(votingID)+"/ballots/"+url.PathEscape(voterID), voteRequest{
		Ballot: ballot,
	}, &response); err != nil {
		return false, err
	}
	return response.Revoted, nil
}

func (s *VotingsService) Unvote(ctx context.Context, votingID, voterID string) error {
	return s.client.request(ctx, http.MethodDelete, "v1/votings/"+url.PathEscape(votingID)+"/ballots/"+url.PathEscape(voterID), nil, nil)
}

type Result struct {
	Choice     string
	Index      int
	Wins       int
	Percentage float64
}

func (s *VotingsService) Results(ctx context.Context, votingID string) (results []Result, tie bool, err error) {

	type votingResultAPIResponse struct {
		Choice     string  `json:"choice"`
		Index      int     `json:"index"`
		Wins       int     `json:"wins"`
		Percentage float64 `json:"percentage"`
	}

	type computeResultsAPIResponse struct {
		Results []votingResultAPIResponse `json:"results"`
		Tie     bool                      `json:"tie"`
	}

	var response *computeResultsAPIResponse
	if err = s.client.request(ctx, http.MethodGet, "v1/votings/"+url.PathEscape(votingID)+"/results", nil, &response); err != nil {
		return nil, false, err
	}

	results = make([]Result, len(response.Results))
	for i, r := range response.Results {
		results[i] = Result(r)
	}

	return results, response.Tie, nil
}
