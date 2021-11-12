package models

type PasswordReset struct {
	Receiver string
	Link     string
	Code     string
}

type Signup struct {
	Receiver string
	Link     string
	Code     string
}
