package models

type PasswordReset struct {
	Receiver string
	Link     string
	Code     string
}
