// Copyright (c) 2022, Direct Decisions Go client AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package directdecisions

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	version = "0.1.0"

	userAgent   = "directdecisions-go/" + version
	contentType = "application/json; charset=utf-8"

	defaultBaseURL = "https://api.directdecisions.com/"
)

// Client manages communication with the Direct Decisions API.
type Client struct {
	httpClient *http.Client // HTTP client must handle authentication implicitly.
	service    service      // Reuse a single struct instead of allocating one for each service on the heap.

	// rate contains the current rate limit for the client as determined
	// by the most recent API call.
	rate   Rate
	rateMu sync.RWMutex

	// Services that API provides.
	Votings *VotingsService
}

// ClientOptions holds optional parameters for the Client.
type ClientOptions struct {
	HTTPClient *http.Client
	BaseURL    *url.URL
}

// NewClient constructs a new Client that uses API key authentication.
func NewClient(key string, o *ClientOptions) (c *Client) {
	if o == nil {
		o = new(ClientOptions)
	}
	return newClient(httpClientWithTransport(o.HTTPClient, o.BaseURL, func(r *http.Request) {
		r.Header.Set("X-Key", key)
	}))
}

// newClient constructs a new *Client with the provided http Client, which
// should handle authentication implicitly, and sets all API services.
func newClient(httpClient *http.Client) (c *Client) {
	c = &Client{httpClient: httpClient}
	c.service.client = c
	c.Votings = (*VotingsService)(&c.service)
	return c
}

func httpClientWithTransport(c *http.Client, baseURL *url.URL, authFunc func(r *http.Request)) *http.Client {
	if c == nil {
		c = new(http.Client)
	}

	transport := c.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	if baseURL == nil {
		baseURL, _ = url.Parse(defaultBaseURL)
	}
	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path += "/"
	}

	c.Transport = roundTripperFunc(func(r *http.Request) (resp *http.Response, err error) {
		r.Header.Set("User-Agent", userAgent)
		authFunc(r)
		u, err := baseURL.Parse(r.URL.String())
		if err != nil {
			return nil, err
		}
		r.URL = u
		return transport.RoundTrip(r)
	})
	return c
}

// request handles the HTTP request response cycle. It JSON encodes the request
// body, creates an HTTP request with provided method on a path with required
// headers, sets current request rate information to the Client and decodes
// request body if the v argument is not nil and content type is
// application/json.
func (c *Client) request(ctx context.Context, method, path string, body, v interface{}) (err error) {
	var bodyBuffer io.ReadWriter
	if body != nil {
		bodyBuffer = new(bytes.Buffer)
		if err = encodeJSON(bodyBuffer, body); err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, path, bodyBuffer)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", contentType)

	r, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer drain(r.Body)

	c.setRate(r)

	if err := responseErrorHandler(r); err != nil {
		return err
	}

	if v != nil && strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return json.NewDecoder(r.Body).Decode(&v)
	}
	return nil
}

// encodeJSON writes a JSON-encoded v object to the provided writer with
// SetEscapeHTML set to false.
func encodeJSON(w io.Writer, v interface{}) (err error) {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

// drain discards all of the remaining data from the reader and closes it,
// asynchronously.
func drain(r io.ReadCloser) {
	go func() {
		// Panicking here does not put data in
		// an inconsistent state.
		defer func() {
			_ = recover()
		}()

		_, _ = io.Copy(io.Discard, r)
		r.Close()
	}()
}

// service is the base type for all API service providing the Client instance
// for them to use.
type service struct {
	client *Client
}

// roundTripperFunc type is an adapter to allow the use of ordinary functions as
// http.RoundTripper interfaces. If f is a function with the appropriate
// signature, roundTripperFunc(f) is a http.RoundTripper that calls f.
type roundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip calls f(r).
func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
