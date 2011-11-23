//Copyright 2011, Andreas Sinz
// Use of this source code is governed by the GPLv2
// license that can be found in the LICENSE file.

//Implements the Post Office Protocol 3 as defined in RFC 1939
package pop3

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	//Carriage Return + Line Feed
	//CRLF is appended at the end of each commands
	CRLF = "\r\n"

	//POP3-Commands
	USER     = "USER"
	PASSWORD = "PASS"
	NOOP     = "NOOP"
	RESET    = "RSET"
	DELETE   = "DELE"
	QUIT     = "QUIT"
	STATUS   = "STAT"
	LIST     = "LIST"
	RETRIEVE = "RETR"
)

var (
	IndexERR          = errors.New("Index must be greater than zero")
	UnkownResponseERR = errors.New("Unkown Response received")
)

type Client struct {
	conn       net.Conn
	stream     *bufio.ReadWriter
	ServerName string
	Greeting   string
}

//Returns a new Client connected to a POP3 server at addr.
//The format of addr is "ip:port" or "hostname:port"
func Dial(addr string) (client *Client, err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	host := addr[:strings.Index(addr, ":")]
	return NewClient(conn, host)

}

//NewClient returns a new Client using an existing connection
//name is used as the Servername
func NewClient(conn net.Conn, name string) (*Client, error) {
	client := new(Client)

	//Create a new ReadWriter and store it
	stream := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	client.stream = stream

	//Download the greeting from the POP3 server
	msg, err := client.readMessage(false)

	if err != nil {
		return nil, err
	}

	client.ServerName = name
	client.Greeting = msg[4:]

	return client, nil
}

//Sends a command to the POP3-Server and returns the response or an error
func (client *Client) Command(command string, isResponseMultiLine bool) (string, error) {

	//Check, whether the client connection has already been
	if client == nil {
		return "", errors.New("Connection hasn't been established")
	}

	//Send the command to the server
	tmp := command + CRLF
	_, writeErr := client.stream.WriteString(tmp)
	if writeErr != nil {
		return "", writeErr
	}
	client.stream.Flush()

	return client.readMessage(isResponseMultiLine)
}

//Returns the response of the pop3 server, or an error if any
func (client *Client) readMessage(isResponseMultiLine bool) (string, error) {
	//Get first line of the response
	msg, err := client.stream.ReadString('\n')

	if err != nil {
		return "", err
	}

	//Check, whether the response starts with "+OK" or "-ERR", otherwise return an error
	if strings.HasPrefix(msg, "+OK") {
		msg = msg[4:]

		if isResponseMultiLine {

			for true {
				line, err1 := client.stream.ReadString('\n')

				if err1 != nil {
					return "", err1
				}

				if line == "."+CRLF {
					break
				}

				msg += line
			}
		}

	} else if strings.HasPrefix(msg, "-ERR") {
		return "", errors.New(msg[5:])
	} else {
		return "", UnkownResponseERR
	}

	return msg, nil

}

//Authenticates with the Server
//It uses the provided Authenticationtype auth
func (client *Client) Authenticate(auth Auth) (string, error) {
	return "", auth.Authenticate(client)
}

//Sends a "NOOP" command and the server will just reply with a positive repsonse
func (client *Client) Ping() (err error) {
	_, err = client.Command(NOOP, false)
	return
}

//Messages that have been marked as "deleted" will be unmarked after this command
func (client *Client) Reset() (string, error) {
	return client.Command(RESET, false)
}

//Mark a mail as "deleted"
//All marked mails will be deleted, when you close the connection with "QUIT"
func (client *Client) MarkMailAsDeleted(index int) (string, error) {
	if index < 0 {
		return "", IndexERR
	}

	return client.Command(fmt.Sprintf("%s %d", DELETE, index), false)
}

//Issues the Quit-Command, so the POP3 session enters the UPDATE state
//All mails, which are marked as "deleted", are going to be removed now
func (client *Client) Quit() (string, error) {
	return client.Command(QUIT, false)
}

//Retrieves the count of mails and the size of all those mails in the mailbox
//Mails, which are marked as "deleted", won't show up
func (client *Client) GetStatus() (mailCount, mailBoxSize int, err error) {
	response, cmdErr := client.Command(STATUS, false)
	if cmdErr != nil {
		return -1, -1, cmdErr
	}

	digits := getDigitsFromLine(response)

	if len(digits) == 2 {
		mailCount = digits[0]
		mailBoxSize = digits[1]
	} else {
		err = UnkownResponseERR
	}

	return
}

//Returns a list of mails
//First digit is the index of the mail, then a whitespace and the size in octets
func (client *Client) GetRawMailList() (response string, err error) {
	if response, err = client.Command(LIST, true); err != nil {
		return
	}

	//Just return the list of mails, not the header
	response = response[strings.Index(response, "\n"):]
	return
}

//Returns the index and the size of the mail at index
func (client *Client) GetMailStatus(index int) (mailIndex, mailSize int, err error) {
	cmdString := fmt.Sprintf("%s %d", LIST, index)

	response, cmdErr := client.Command(cmdString, false)
	if cmdErr != nil {
		return -1, -1, cmdErr
	}

	digits := getDigitsFromLine(response)

	if len(digits) == 2 {
		mailIndex = digits[0]
		mailSize = digits[1]
	} else {
		err = UnkownResponseERR
	}

	return
}

//Retrieves the raw mail at index as a string
//This string contains the whole header and the body of the mail
func (client *Client) GetRawMail(index int) (mail string, err error) {
	if index < 1 {
		return "", IndexERR
	}

	mail, err = client.Command(fmt.Sprintf("%s %d", RETRIEVE, index), true)

	if err != nil {
		return
	}

	//Remove the first line
	mail = mail[strings.Index(mail, "\n"):]

	return
}

//Returns every digit, which exists in the string
func getDigitsFromLine(line string) (digits []int) {

	for _, part := range strings.Split(line, " ") {
		if tmp, convertErr := strconv.Atoi(strings.TrimSpace(part)); convertErr == nil {
			digits = append(digits, tmp)
		}
	}

	return
}
