package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/mazzz1y/go-matrix-webhook/internal/matrix"
	zerolog "github.com/rs/zerolog/log"
)

type WebhookPayload struct {
	Message string `json:"message"`
	RoomID  string `json:"room_id"`
}

type ResponseData struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

const (
	successSent         = "message sent"
	errInvalidSecret    = "invalid secret header"
	errParseBodyError   = "parse body error"
	errJoinRoomError    = "join room error"
	errSendMessageError = "send message error"
	errEmptyMessage     = "empty message"
)

func NewHandler(m matrix.Matrix, secret string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		secretHeader := r.Header.Get("X-Secret")
		userHeader := r.Header.Get("X-Forwarded-User")
		reqIP := getIP(r)

		if secret != "" && secretHeader != secret {
			sendResponse(w, http.StatusBadRequest, errInvalidSecret)
			return
		}

		log := zerolog.With().Str("ip", reqIP).Str("path", r.URL.Path).Logger()
		if userHeader != "" {
			log = log.With().Str("user", userHeader).Logger()
		}

		payload, err := parsePayload(r.Body)
		if err != nil {
			log.Error().Err(err).Msg(errParseBodyError)
			sendResponse(w, http.StatusBadRequest, errParseBodyError)
			return
		}

		err = m.JoinRoom(payload.RoomID)
		if err != nil {
			log.Error().Str("room_id", payload.RoomID).Err(err).Msg(errJoinRoomError)
			sendResponse(w, http.StatusInternalServerError, errJoinRoomError)
			return
		}

		err = m.SendMessage(payload.RoomID, payload.Message)
		if err != nil {
			log.Error().Str("room_id", payload.RoomID).Err(err).Msg(errSendMessageError)
			sendResponse(w, http.StatusOK, errSendMessageError)
			return
		}

		log.Debug().Str("room_id", payload.RoomID).Msg(successSent)
		sendResponse(w, http.StatusOK, "")
	}
}

func parsePayload(body io.ReadCloser) (*WebhookPayload, error) {
	var payload WebhookPayload
	err := json.NewDecoder(body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	if payload.Message == "" {
		return nil, errors.New(errEmptyMessage)
	}

	return &payload, nil
}

func getIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ipParts := strings.Split(xff, ",")
		return strings.TrimSpace(ipParts[0])
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	if cf := r.Header.Get("CF-Connecting-IP"); cf != "" {
		return cf
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func sendResponse(w http.ResponseWriter, code int, msg string) {
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
