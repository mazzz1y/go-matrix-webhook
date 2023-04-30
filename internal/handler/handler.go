package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mazzz1y/go-matrix-webhook/internal/matrix"
	"github.com/rs/zerolog/log"
)

type WebhookPayload struct {
	Message string `json:"message"`
	RoomID  string `json:"room_id"`
}

type ResponseData struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func NewHandler(m matrix.Matrix, s string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Secret") != s {
			response(w, http.StatusUnauthorized, "")
			log.Debug().Msg("invalid secret header")
			return
		}

		payload, err := parsePayload(r)
		if err != nil {
			response(w, http.StatusBadRequest, "")
			log.Error().Str("room_id", payload.RoomID).Err(err).Msg("parse body error")
			return
		}

		err = m.SendMessage(payload.RoomID, payload.Message)
		if err != nil {
			response(w, http.StatusInternalServerError, "")
			log.Error().Str("room_id", payload.RoomID).Err(err).Msg("send error")
			return
		}

		log.Debug().Str("room_id", payload.RoomID).Msg("sent msg")
		response(w, http.StatusOK, "")
	}
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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonRes)
}

func parsePayload(r *http.Request) (*WebhookPayload, error) {
	var payload WebhookPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
