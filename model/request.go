package model

type StayerGetRequest struct {
	ID     int64            `json:"id"`
	Name   string           `json:"name"`
	Room   string           `json:"room"`
	RoomID int              `json:"roomId"`
}