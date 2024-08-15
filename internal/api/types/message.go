package types

import "time"

type Msg struct {
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
