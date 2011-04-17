package pop3

import "os"

type Auth interface {
	Authenticate(client *Client) os.Error
}


type PlainAuthentication struct {
	user, pass string
}

func CreatePlainAuthentication (user, pass string) *PlainAuthentication {
	return &PlainAuthentication{user, pass}
}

func (auth *PlainAuthentication) Authenticate(client *Client) os.Error {
	client.WriteMessage("USER " + auth.user)
	_, errUser := client.ReadMessage(false)
	if errUser != nil {
		return errUser
	}
		
	client.WriteMessage("PASS " + auth.pass)
	_, errPass := client.ReadMessage(false)
	if errPass != nil {
		return errPass
	}
	
	return nil
}