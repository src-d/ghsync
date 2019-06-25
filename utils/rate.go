package utils

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/go-github/github"
)

const (
	ctxEtag    = "etag"
	ctxId      = "id"
	writeDelay = 1 * time.Second
)

// rateLimitTransport implements GitHub's best practices
// for avoiding rate limits
// https://developer.github.com/v3/guides/best-practices-for-integrators/#dealing-with-abuse-rate-limits
type rateLimitTransport struct {
	transport        http.RoundTripper
	delayNextRequest bool
	responseBody     []byte

	m sync.Mutex
}

func NewRateLimitTransport(rt http.RoundTripper) *rateLimitTransport {
	return &rateLimitTransport{transport: rt}
}

func (rlt *rateLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Make requests for a single user or client ID serially
	// This is also necessary for safely saving
	// and restoring bodies between retries below
	rlt.lock(req)

	// If you're making a large number of POST, PATCH, PUT, or DELETE requests
	// for a single user or client ID, wait at least one second between each request.
	if rlt.delayNextRequest {
		log.Printf("[DEBUG] Sleeping %s between write operations", writeDelay)
		time.Sleep(writeDelay)
	}

	rlt.delayNextRequest = isWriteMethod(req.Method)

	resp, err := rlt.transport.RoundTrip(req)
	if err != nil {
		rlt.unlock(req)
		return resp, err
	}

	if resp.Header.Get("X-RateLimit-Remaining") == "0" {
		rlt.delayNextRequest = false

		var limit int
		if limitHeader := resp.Header.Get("X-RateLimit-Limit"); limitHeader != "" {
			limit, _ = strconv.Atoi(limitHeader)
		}

		var reset github.Timestamp
		if resetHeader := resp.Header.Get("X-RateLimit-Reset"); resetHeader != "" {
			if v, _ := strconv.ParseInt(resetHeader, 10, 64); v != 0 {
				reset = github.Timestamp{time.Unix(v, 0)}
			}
		}

		retryAfter := reset.Sub(time.Now())

		log.Printf("[DEBUG] Rate limit %d reached, sleeping for %s before retrying",
			limit, retryAfter)
		if retryAfter < 0 {
			log.Printf("[WARN] retryAfter < 0. reset: %v | now: %v",
				reset, time.Now())
		} else {
			time.Sleep(retryAfter)
		}

		rlt.unlock(req)
		return rlt.RoundTrip(req)
	}

	rlt.unlock(req)
	return resp, nil
}

func (rlt *rateLimitTransport) lock(req *http.Request) {
	rlt.m.Lock()
}

func (rlt *rateLimitTransport) unlock(req *http.Request) {
	rlt.m.Unlock()
}

func isWriteMethod(method string) bool {
	switch method {
	case "POST", "PATCH", "PUT", "DELETE":
		return true
	}
	return false
}
