package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mazzz1y/go-matrix-webhook/internal/matrix"
	"github.com/rs/zerolog/log"
)

type WebhookPayload struct {
	Message string `json:"message"`
	RoomID  string `json:"room_id"`
}

type ResponseData struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewHandler(m matrix.Matrix, s string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Secret") != s {
			response(w, http.StatusBadRequest, "invalid secret header")
			return
		}

		payload, err := parsePayload(r)
		if err != nil {
			log.Error().Err(err).Msg("")
			response(w, http.StatusBadRequest, "parse body error")
			return
		}

		err = m.JoinRoom(payload.RoomID)
		if err != nil {
			log.Error().Str("room_id", payload.RoomID).Err(err).Msg("")
			response(w, http.StatusInternalServerError, "join room error")
			return
		}

		err = m.SendMessage(payload.RoomID, payload.Message)
		if err != nil {
			log.Error().Str("room_id", payload.RoomID).Err(err).Msg("")
			response(w, http.StatusOK, "send message error")
			return
		}

		log.Debug().Str("room_id", payload.RoomID).Msg("message sent")
		response(w, http.StatusOK, "")
	}
}

func parsePayload(r *http.Request) (*WebhookPayload, error) {
	var payload WebhookPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	if payload.Message == "" {
		return nil, errors.New("empty message")
	}

	return &payload, nil
}

func response(w http.ResponseWriter, code int, msg string) {
	if msg == "" {
		msg = http.StatusText(code)
	}

	res := ResponseData{
		Message: msg,
		Status:  code,
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonRes)
}
