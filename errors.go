// Copyright (c) 2022, Direct Decisions Go client AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package directdecisions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Errors that are returned by the API.
var (
	ErrHTTPStatusBadRequest          = errors.New("http status: " + http.StatusText(http.StatusBadRequest))
	ErrHTTPStatusUnauthorized        = errors.New("http status: " + http.StatusText(http.StatusUnauthorized))
	ErrHTTPStatusForbidden           = errors.New("http status: " + http.StatusText(http.StatusForbidden))
	ErrHTTPStatusNotFound            = errors.New("http status: " + http.StatusText(http.StatusNotFound))
	ErrHTTPStatusMethodNotAllowed    = errors.New("http status: " + http.StatusText(http.StatusMethodNotAllowed))
	ErrHTTPStatusTooManyRequests     = errors.New("http status: " + http.StatusText(http.StatusTooManyRequests))
	ErrHTTPStatusInternalServerError = errors.New("http status: " + http.StatusText(http.StatusInternalServerError))
	ErrHTTPStatusServiceUnavailable  = errors.New("http status: " + http.StatusText(http.StatusServiceUnavailable))
	ErrHTTPStatusBadGateway          = errors.New("http status: " + http.StatusText(http.StatusBadGateway))

	ErrInvalidData    = errors.New("Invalid Data")
	ErrMissingChoices = errors.New("Missing Choices")
	ErrChoiceRequired = errors.New("Choice Required")
	ErrChoiceTooLong  = errors.New("Choice Too Long")
	ErrTooManyChoices = errors.New("Too Many Choices")
	ErrBallotRequired = errors.New("Ballot Required")
	ErrVoterIDTooLong = errors.New("Voter ID Too Long")
	ErrInvalidVoterID = errors.New("Invalid Voter ID")
)

var statusToError = map[int]error{
	http.StatusBadRequest:          ErrHTTPStatusBadRequest,
	http.StatusUnauthorized:        ErrHTTPStatusUnauthorized,
	http.StatusForbidden:           ErrHTTPStatusForbidden,
	http.StatusNotFound:            ErrHTTPStatusNotFound,
	http.StatusMethodNotAllowed:    ErrHTTPStatusMethodNotAllowed,
	http.StatusTooManyRequests:     ErrHTTPStatusTooManyRequests,
	http.StatusInternalServerError: ErrHTTPStatusInternalServerError,
	http.StatusServiceUnavailable:  ErrHTTPStatusServiceUnavailable,
	http.StatusBadGateway:          ErrHTTPStatusBadGateway,
}

var messageToError = map[string]error{
	"Invalid Data":      ErrInvalidData,
	"Missing Choices":   ErrMissingChoices,
	"Choice Required":   ErrChoiceRequired,
	"Choice Too Long":   ErrChoiceTooLong,
	"Too Many Choices":  ErrTooManyChoices,
	"Ballot Required":   ErrBallotRequired,
	"Voter ID Too Long": ErrVoterIDTooLong,
	"Invalid Voter ID":  ErrInvalidVoterID,
}

type messageResponse struct {
	Message string   `json:"message,omitempty"`
	Code    int      `json:"code,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// responseErrorHandler returns an error based on the HTTP status code or nil if
// the status code is from 200 to 299.
func responseErrorHandler(r *http.Response) error {
	if r.StatusCode/100 == 2 {
		return nil
	}

	statusErr, ok := statusToError[r.StatusCode]
	if !ok {
		statusErr = errors.New("http status: " + http.StatusText(r.StatusCode))
	}

	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return statusErr
	}

	errs := []error{
		statusErr,
	}

	var e messageResponse
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		if errors.Is(err, io.EOF) { // empty body
			return statusErr
		}
		errs = append(errs, fmt.Errorf("json decode: %w", err))
	} else {
		for _, e := range e.Errors {
			if err, ok := messageToError[e]; ok {
				errs = append(errs, err)
			} else {
				errs = append(errs, errors.New(e))
			}
		}
	}

	return errors.Join(errs...)
}
