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
	_, userErr := client.Command(USER + " " + auth.user, false)
	if userErr != nil {
		return userErr
	}

	_, pwErr := client.Command(PASSWORD + " " + auth.pass, false)
	if pwErr != nil {
		return pwErr
	}

	return nil
}
