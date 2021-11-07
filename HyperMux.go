package hyper_mux

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

const (
	MethodGet     = "GET"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodDelete  = "DELETE"
	MethodHead    = "HEAD"
	MethodOptions = "OPTIONS"
	MethodPatch   = "PATCH"
)

// NewHyperMux creates new instance of HyperMux
func NewHyperMux() *HyperMux {
	return &HyperMux{
		endPoints:   make([]*endpoint, 0),
		middlewares: make([]func(next http.Handler) http.Handler, 0),
	}
}

// HyperMux holds all the end-point routings and middlewares.
type HyperMux struct {
	endPoints   []*endpoint
	middlewares []func(next http.Handler) http.Handler
}

// UseMiddleware append a middleware to the end of middleware chain in this Mux
func (m *HyperMux) UseMiddleware(mw func(next http.Handler) http.Handler) {
	m.middlewares = append(m.middlewares, mw)
}

// ServeHTTP serve the HTTP request, its an implementation of http.Handler.ServeHTTP
func (m *HyperMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hf, pattern := m.handleFuncForRequest(r)
		if hf == nil {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("not found"))
			return
		} else {
			mp, _ := parsePathParams(pattern, r.URL.Path)
			if mp != nil {
				for k, v := range mp {
					r.Header.Add(k, v)
				}
			}
			hf.ServeHTTP(w, r)
		}
	})
	if m.middlewares != nil && len(m.middlewares) > 0 {
		for i := len(m.middlewares) - 1; i >= 0; i-- {
			mw := m.middlewares[i]
			h = mw(h)
		}
	}

	h.ServeHTTP(w, r)
}

// AddRoute add a routing to a HTTP resource
func (m *HyperMux) AddRoute(pattern, method string, hFunc http.HandlerFunc) {
	ep := &endpoint{
		pathMethod: &pathMethod{
			pathPattern: pattern,
			method:      method,
		},
		handleFunc: hFunc,
	}
	if len(m.endPoints) > 1 {
		sort.Slice(m.endPoints, func(i, j int) bool {
			return len(m.endPoints[j].pathMethod.String()) > len(m.endPoints[i].pathMethod.String())
		})
	}
	m.endPoints = append(m.endPoints, ep)
}

func (m *HyperMux) handleFuncForRequest(r *http.Request) (http.HandlerFunc, string) {
	for _, ep := range m.endPoints {
		if good, pattern := ep.pathMethod.matchRequest(r); good {
			return ep.handleFunc, pattern
		}
	}
	return nil, ""
}

type endpoint struct {
	pathMethod *pathMethod
	handleFunc http.HandlerFunc
}

type pathMethod struct {
	pathPattern string
	method      string
}

func (pm *pathMethod) matchRequest(r *http.Request) (bool, string) {
	good := isTemplateCompatible(pm.pathPattern, r.URL.Path)
	if good && r.Method == pm.method {
		return true, pm.pathPattern
	}
	return false, ""
}

func (pm *pathMethod) String() string {
	return fmt.Sprintf("[%s]%s", pm.method, pm.pathPattern)
}

func isTemplateCompatible(template, path string) bool {
	if template == path {
		return true
	}
	if !strings.Contains(template, "{") {
		return false
	}
	templatePaths := strings.Split(template, "/")
	pathPaths := strings.Split(path, "/")
	if len(templatePaths) != len(pathPaths) {
		return false
	}
	for idx, templateElement := range templatePaths {
		pathElement := pathPaths[idx]
		if len(templateElement) > 0 && len(pathElement) > 0 {
			if templateElement[:1] == "{" && templateElement[len(templateElement)-1:] == "}" {
				continue
			} else if templateElement != pathElement {
				return false
			}
		}
	}
	return true
}

func parsePathParams(template, path string) (map[string]string, error) {
	templatePaths := strings.Split(template, "/")
	pathPaths := strings.Split(path, "/")
	if len(templatePaths) != len(pathPaths) {
		return nil, fmt.Errorf("pathElement length not equals to templateElement length")
	}
	ret := make(map[string]string)
	for idx, templateElement := range templatePaths {
		pathElement := pathPaths[idx]
		if len(templateElement) > 0 && len(pathElement) > 0 {
			if templateElement[:1] == "{" && templateElement[len(templateElement)-1:] == "}" {
				tKey := templateElement[1 : len(templateElement)-1]
				ret[tKey] = pathElement
			} else if templateElement != pathElement {
				return nil, fmt.Errorf("template %s not compatible with path %s", template, path)
			}
		}
	}
	return ret, nil
}

// InternalServerError writes error message to the http.ResponseWriter
func InternalServerError(w http.ResponseWriter, err error) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf("error while serving request. got %s", err.Error())))
}

// WriteString simply write to the http.ResponseWriter a text of type text/plain
func WriteString(w http.ResponseWriter, code int, text string) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(code)
	w.Write([]byte(text))
}

// WriteJson simply writes a JSON to the http.ResponseWriter of type application/json
// if the marshaled data is not marshallable to json, it write an internal server error.
func WriteJson(w http.ResponseWriter, code int, data interface{}) {
	byteArray, err := json.Marshal(data)
	if err != nil {
		WriteString(w, http.StatusInternalServerError, "error marshaling json. got "+err.Error())
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(byteArray)
	}
}
