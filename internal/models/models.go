package models

type Email struct {
	Receiver     string `json:"receiver"`
	Sender       string `json:"sender"`
	EmailAddress string `json:"emailAddress"`
}
