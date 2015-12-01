package hrouter

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Matcher registers patterns as routes, matches requests and builds URLs.
type Matcher interface {
	Route(pattern string) (*Route, error)
	Match(r *http.Request) (Handler, Params)
	URL(r *Route, p Params) *url.URL
}

// -----------------------------------------------------------------------------

func newMatcher() *matcher {
	return &matcher{
		root: &node{path: "/"},
	}
}

type matcher struct {
	root *node
}

func (m *matcher) Route(pattern string) (*Route, error) {
	pat, err := parsePattern(pattern)
	if err != nil {
		return nil, err
	}
	var prefix string
	prefix, pat.static = pat.static[0], pat.static[1:]
	r, err := m.root.addStaticPrefix(prefix, pat)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (m *matcher) Match(r *http.Request) (Handler, Params) {
	// TODO: handle NotFound and strict slashes (matcher options).
	node, p := lookup(m.root, r.URL.Path)
	if node == nil {
		return nil, nil
	}
	h := methodHandler(node.route.Handlers, r.Method)
	return h, p
}

func (m *matcher) URL(r *Route, p Params) *url.URL {
	// TODO...
	return nil
}

// methodHandler returns the handler registered for the given HTTP method.
func methodHandler(handlers map[string]Handler, method string) Handler {
	if h, ok := handlers[method]; ok {
		return h
	}
	switch method {
	case "OPTIONS":
		return allowHandler(handlers, 200)
	case "HEAD":
		if h, ok := handlers["GET"]; ok {
			return h
		}
		fallthrough
	default:
		if h, ok := handlers[""]; ok {
			return h
		}
	}
	return allowHandler(handlers, 405)
}

// allowHandler returns a handler that sets a header with the given
// status code and allowed methods.
func allowHandler(handlers map[string]Handler, code int) Handler {
	allowed := make([]string, len(handlers)+1)
	allowed[0] = "OPTIONS"
	i := 1
	for m, _ := range handlers {
		if m != "" && m != "OPTIONS" {
			allowed[i] = m
			i++
		}
	}
	return Handler(func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Allow", strings.Join(allowed[:i], ", "))
		w.WriteHeader(code)
		fmt.Fprintln(w, code, http.StatusText(code))
	})
}
