package matrix

import (
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type Matrix struct {
	client *mautrix.Client
}

func NewMatrix(serverUrl string, userID string, token string) (*Matrix, error) {
	c, err := mautrix.NewClient(serverUrl, id.NewUserID(userID, serverUrl), token)
	if err != nil {
		return &Matrix{}, err
	}

	return &Matrix{c}, nil
}

func (m Matrix) SendMessage(roomID, message string) error {
	_, err := m.client.SendMessageEvent(id.RoomID(roomID), event.EventMessage, &event.MessageEventContent{
		MsgType: event.MsgText,
		Body:    message,
	})
	return err
}
