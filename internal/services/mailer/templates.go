package mailer

type Renderable interface {
	Render() []byte
}

//renderable is the renderable struct
type renderable struct {
	body []byte
}

//Render renders the email body
func (r renderable) Render() []byte {
	return r.body
}
