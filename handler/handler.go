package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dhucsik/daribar-test/middleware"
	"github.com/dhucsik/daribar-test/storage"
	"github.com/dhucsik/daribar-test/utils"
)

type Handler struct {
	DataUpdates *utils.Fanout
	Repo        *storage.Storage
	Auth        *middleware.JWTMiddleware
}

func NewHandler(dataUpd *utils.Fanout, repo *storage.Storage, auth *middleware.JWTMiddleware) *Handler {
	return &Handler{
		DataUpdates: dataUpd,
		Repo:        repo,
		Auth:        auth,
	}
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.DataUpdates.Connected()
	phone, err := h.Auth.GetPhoneFromHeader(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	quantity := h.Repo.GetOpenOrdersQuantity(phone)

	updates := make(chan int, 100)

	h.DataUpdates.Subscribe(phone, quantity, updates)
	defer h.DataUpdates.Unsubscribe(phone)

	h.SetSseHeaders(rw)

	for {
		data, ok := <-updates
		log.Println(data, ok)
		if !ok {
			return
		}

		h.pingAndFlush(rw, strconv.FormatInt(int64(data), 10))
	}
}

func (h *Handler) SetSseHeaders(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
}

func (h *Handler) pingAndFlush(rw http.ResponseWriter, data string) {
	_, err := fmt.Fprintf(rw, "data: %s\n\n", data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.(http.Flusher).Flush()
}
