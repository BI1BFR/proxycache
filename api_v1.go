package proxycache

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func (p *ProxyCache) HTTPHandlerV1() http.Handler {
	return &handlerV1{p}
}

type handlerV1 struct {
	p *ProxyCache
}

func (h *handlerV1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "v1/keys/") {
		h.keys(w, r)
	} else if strings.HasPrefix(r.URL.Path, "v1/status") {
		h.status(w)
	} else {
		http.NotFound(w, r)
	}
}

func (h *handlerV1) keys(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "v1/keys/")
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
