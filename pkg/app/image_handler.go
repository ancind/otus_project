package app

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handlers) ImageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var pu string
	parsedURL, err := url.Parse(vars["imageUrl"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("validation error"))
		h.logger.Err(err).Msg(err.Error())

		return
	}

	pu = parsedURL.String()
	if parsedURL.Scheme == "" {
		pu = fmt.Sprintf("http://%s", parsedURL.String())
	}

	width, err := strconv.Atoi(vars["width"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("validation error"))
		h.logger.Err(err).Msg(err.Error())

		return
	}

	height, err := strconv.Atoi(vars["height"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("validation error"))
		h.logger.Err(err).Msg(err.Error())

		return
	}

	img, err := h.svc.ResizeImage(r.Context(), pu, r.Header, width, height)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("resize error"))
		h.logger.Err(err).Msg(err.Error())

		return
	}

	w.Header().Add("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(img)))

	if _, err := w.Write(img); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}
}
