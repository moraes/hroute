package hrouter

import (
	"net/http"
	"net/url"
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
	route, err := m.root.addStaticPrefix(prefix, pat)
	if err != nil {
		return nil, err
	}
	return route, nil
}

func (m *matcher) Match(r *http.Request) (Handler, Params) {
	// TODO...
	// node, p := lookup(m.root, r.Path)
	return nil, nil
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
		return r.allowHandler(handlers, 200)
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
	return r.allowHandler(handlers, 405)
}

// allowHandler returns a handler that sets a header with the given
// status code and allowed methods.
func allowHandler(handlers map[string]Handler, code int) Handler {
	allowed := []string{"OPTIONS"}
	for m, _ := range handlers {
		if m != "" && m != "OPTIONS" {
			allowed = append(allowed, m)
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Allow", strings.Join(allowed, ", "))
		w.WriteHeader(code)
		fmt.Fprintln(w, code, http.StatusText(code))
	}
}
