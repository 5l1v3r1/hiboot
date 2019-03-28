package httpclient

import (
	"io"
	"net/http"
)

// Doer interface has the method required to use a type as custom http client.
// The net/*http.client type satisfies this interface.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// Client Is a generic HTTP client interface
type Client interface {
	Get(url string, headers http.Header, callbacks ...func(req *http.Request)) (*http.Response, error)
	Post(url string, body io.Reader, headers http.Header, callbacks ...func(req *http.Request)) (*http.Response, error)
	Put(url string, body io.Reader, headers http.Header, callbacks ...func(req *http.Request)) (*http.Response, error)
	Patch(url string, body io.Reader, headers http.Header, callbacks ...func(req *http.Request)) (*http.Response, error)
	Delete(url string, headers http.Header, callbacks ...func(req *http.Request)) (*http.Response, error)
	Do(req *http.Request) (*http.Response, error)
}
