// Copyright (c) 2022, Direct Decisions Go client AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package directdecisions_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"directdecisions.com/directdecisions"
)

func TestErrors(t *testing.T) {

	for _, tc := range []struct {
		body   string
		status int
		errors []error
	}{
		{
			status: http.StatusBadRequest,
			errors: []error{directdecisions.ErrHTTPStatusBadRequest},
		},
		{
			status: http.StatusUnauthorized,
			errors: []error{directdecisions.ErrHTTPStatusUnauthorized},
		},
		{
			status: http.StatusForbidden,
			errors: []error{directdecisions.ErrHTTPStatusForbidden},
		},
		{
			status: http.StatusNotFound,
			errors: []error{directdecisions.ErrHTTPStatusNotFound},
		},
		{
			status: http.StatusMethodNotAllowed,
			errors: []error{directdecisions.ErrHTTPStatusMethodNotAllowed},
		},
		{
			status: http.StatusTooManyRequests,
			errors: []error{directdecisions.ErrHTTPStatusTooManyRequests},
		},
		{
			status: http.StatusInternalServerError,
			errors: []error{directdecisions.ErrHTTPStatusInternalServerError},
		},
		{
			status: http.StatusServiceUnavailable,
			errors: []error{directdecisions.ErrHTTPStatusServiceUnavailable},
		},
		{
			status: http.StatusBadGateway,
			errors: []error{directdecisions.ErrHTTPStatusBadGateway},
		},
		{
			body:   `{"message": "Bad Request", "Code": 400, "errors": ["Invalid Data"]}`,
			status: http.StatusBadRequest,
			errors: []error{directdecisions.ErrHTTPStatusBadRequest, directdecisions.ErrInvalidData},
		},
		{
			body:   `{"message": "Bad Request", "Code": 400, "errors": ["Missing Choices"]}`,
			status: http.StatusBadRequest,
			errors: []error{directdecisions.ErrHTTPStatusBadRequest, directdecisions.ErrMissingChoices},
		},
		{
			body:   `{"message": "Bad Request", "Code": 400, "errors": ["Choice Required"]}`,
			status: http.StatusBadRequest,
			errors: []error{directdecisions.ErrHTTPStatusBadRequest, directdecisions.ErrChoiceRequired},
		},
		{
			body:   `{"message": "Bad Request", "Code": 400, "errors": ["Choice Too Long"]}`,
			status: http.StatusBadRequest,
			errors: []error{directdecisions.ErrHTTPStatusBadRequest, directdecisions.ErrChoiceTooLong},
		},
		{
			body:   `{"message": "Bad Request", "Code": 400, "errors": ["Ballot Required", "Voter ID Too Long", "Invalid Voter ID"]}`,
			status: http.StatusBadRequest,
			errors: []error{directdecisions.ErrHTTPStatusBadRequest, directdecisions.ErrBallotRequired, directdecisions.ErrVoterIDTooLong, directdecisions.ErrInvalidVoterID},
		},
	} {
		t.Run(joinErrors(tc.errors), func(t *testing.T) {
			client, mux, _ := newClient(t, "")

			mux.HandleFunc("/v1/votings", requireMethod("POST", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))

			_, err := client.Votings.Create(context.Background(), nil)
			assertErrors(t, err, tc.errors...)
		})
	}
}

func joinErrors(errors []error) string {
	s := make([]string, len(errors))
	for i, err := range errors {
		s[i] = err.Error()
	}
	return strings.Join(s, ", ")
}
