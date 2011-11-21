//Copyright 2011, Andreas Sinz
// Use of this source code is governed by the GPLv2
// license that can be found in the LICENSE file.

package pop3

type Auth interface {
	Authenticate(client *Client) error
}


type PlainAuthentication struct {
	user, pass string
}

func CreatePlainAuthentication (user, pass string) *PlainAuthentication {
	return &PlainAuthentication{user, pass}
}

func (auth *PlainAuthentication) Authenticate(client *Client) error {
	client.WriteMessage(USER + " " + auth.user)
	_, errUser := client.ReadMessage(false)
	if errUser != nil {
		return errUser
	}

	client.WriteMessage(PASSWORD + " " + auth.pass)
	_, errPass := client.ReadMessage(false)
	if errPass != nil {
		return errPass
	}

	return nil
}
