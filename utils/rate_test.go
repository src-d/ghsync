package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
)

func TestRateLimit(t *testing.T) {
	assert := assert.New(t)

	resetAfter := 2 * time.Second
	mt := &rateTransport{
		Limit: 3,
		Reset: time.Now().Add(resetAfter),
	}
	c := newClient(assert, mt)

	// first 3 requests must be executed almost immediately
	spent := measure(func() {
		c.getSuccess()
		c.getSuccess()
		c.getSuccess()
	})
	assert.True(spent < time.Millisecond)

	// should sleep before success for more than a second
	spent = measure(func() {
		c.getSuccess()
	})
	assert.True(spent > time.Second)

	// next request is quick again
	spent = measure(func() {
		c.getSuccess()
	})
	assert.True(spent < time.Millisecond)
}

func TestAbuseLimit(t *testing.T) {
	assert := assert.New(t)

	mt := &abuseTransport{
		RetryAfter: 1,
	}
	c := newClient(assert, mt)

	// should return only after sleeping ~second
	spent := measure(func() {
		c.getSuccess()
	})
	assert.True(spent >= time.Second)
}

func TestWriteLimit(t *testing.T) {
	assert := assert.New(t)

	resetAfter := time.Second
	mt := &rateTransport{
		Limit: 3,
		Reset: time.Now().Add(resetAfter),
	}
	c := newClient(assert, mt)

	// first requests must be executed almost immediately
	spent := measure(func() {
		c.postSuccess()
	})
	assert.True(spent < time.Millisecond)

	// sequential requests must sleep for 1 second
	spent = measure(func() {
		c.postSuccess()
		c.postSuccess()
	})
	assert.True(spent > time.Second)
}

// helper to mesure time
func measure(fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Now().Sub(start)
}

// test wrapper for github client
type client struct {
	*github.Client
	assert *assert.Assertions
}

// creates new github client with RateLimitTransport
func newClient(a *assert.Assertions, mt http.RoundTripper) *client {
	rt := NewRateLimitTransport(mt)
	hc := &http.Client{Transport: rt}

	return &client{Client: github.NewClient(hc), assert: a}
}

// any GET request
func (c *client) getSuccess() {
	_, _, err := c.Users.Get(context.TODO(), "")
	c.assert.Nil(err)
}

// any POST request
func (c *client) postSuccess() {
	_, err := c.Users.Unfollow(context.TODO(), "")
	c.assert.Nil(err)
}

// emulates responses with X-RateLimit header
type rateTransport struct {
	// number of requests to allow before hitting the limit
	Limit int
	// when to reset the remaining number (mock doesn't actually reset it)
	Reset time.Time

	requests int
}

func (t *rateTransport) RoundTrip(*http.Request) (*http.Response, error) {
	remaining := t.Limit - t.requests
	var resp *http.Response
	if remaining == 0 {
		resp = &http.Response{
			StatusCode: http.StatusForbidden,
			Header:     make(http.Header),
			Body:       ioutil.NopCloser(bytes.NewBuffer(limitPayload)),
		}
	} else {
		resp = &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       http.NoBody,
		}
	}

	resp.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp.Header.Set("X-RateLimit-Limit", strconv.Itoa(t.Limit))
	resp.Header.Set("X-RateLimit-Reset", strconv.FormatInt(t.Reset.Unix(), 10))
	resp.Header.Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

	t.requests++

	return resp, nil
}

// emulates abuse error
type abuseTransport struct {
	RetryAfter int // number of seconds

	returnSuccess bool
}

func (t *abuseTransport) RoundTrip(*http.Request) (*http.Response, error) {
	if t.returnSuccess {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       http.NoBody,
		}
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	}

	t.returnSuccess = true
	resp := &http.Response{
		StatusCode: http.StatusForbidden,
		Header:     make(http.Header),
		Body:       ioutil.NopCloser(bytes.NewBuffer(abusePayload)),
	}

	resp.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp.Header.Set("Retry-After", strconv.Itoa(t.RetryAfter))

	return resp, nil
}

// github library relies on particular error payload from github, not only status
type errorPayload struct {
	Message string `json:"message"`
	DocURL  string `json:"documentation_url"`
}

func (p errorPayload) JSON() []byte {
	b, err := json.Marshal(p)
	if err != nil {
		panic("can't marshal payload")
	}

	return b
}

var limitPayload = errorPayload{
	Message: "API rate limit exceeded for xxx.xxx.xxx.xxx.",
	DocURL:  "https://developer.github.com/v3/#rate-limiting",
}.JSON()

var abusePayload = errorPayload{
	Message: "You have triggered an abuse detection mechanism and have been temporarily blocked from content creation. Please retry your request again later.",
	DocURL:  "https://developer.github.com/v3/#abuse-rate-limits",
}.JSON()
