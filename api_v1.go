package proxycache

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// HTTPHandlerV1 create HTTP handler (version 1).
// To serve on a sub URI, don't forget to use http.StripPrefix().
// Check example/server for more details.
func (p *ProxyCache) HTTPHandlerV1() http.Handler {
	return &handlerV1{p}
}

type handlerV1 struct {
	p *ProxyCache
}

func (h *handlerV1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/v1/keys/") {
		h.keys(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/v1/status") {
		h.status(w)
	} else if strings.HasPrefix(r.URL.Path, "/v1/config") {
		h.config(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func (h *handlerV1) keys(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/v1/keys/")
	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		h.get(w, r.URL.Path)
	case "PUT":
		value, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		ttw, _ := strconv.Atoi(r.URL.Query().Get("ttw"))
		h.put(w, r.URL.Path, value, int64(ttw))
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *handlerV1) get(w http.ResponseWriter, key string) {
	b := h.p.Get(key)
	if b == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Write(b)
	}
}

func (h *handlerV1) put(w http.ResponseWriter, key string, value []byte, ttw int64) {
	h.p.Put(key, value, ttw)
}

func (h *handlerV1) status(w http.ResponseWriter) {
	w.Write(h.p.Status())
}

func (h *handlerV1) config(w http.ResponseWriter, r *http.Request) {
	loader, _ := strconv.Atoi(r.URL.Query().Get("loader"))
	saver, _ := strconv.Atoi(r.URL.Query().Get("saver"))
	if loader > 0 {
		h.p.SetLoadMaxProc(loader)
	}
	if saver > 0 {
		h.p.SetSaveProc(saver)
	}
	w.WriteHeader(http.StatusOK)
}
