// Package maili provides a quick wrapper around the Mailinator API.
package maili

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	// InboxURL is the Mailinator URL to send inbox API requests.
	InboxURL = "https://api.mailinator.com/api/inbox?token=%s&to=%s"

	// EmailURL is the Mailinator URL to send email API requests.
	EmailURL = "https://api.mailinator.com/api/email?token=%s&msgid=%s"
)

// doReq issues a request to url, after fmt.Sprintf(url, tok, id), decoding
// response into obj and returning any errors encountered.
func doReq(url, tok, id string, obj interface{}) error {
	u := fmt.Sprintf(url, tok, id)
	res, err := http.Get(u)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error: could not retrieve %s", u)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, obj)
}

// InboxMsg is the type for inbox messages.
type InboxMsg struct {
	Fromfull   string `json:"fromfull"`
	Subject    string `json:"subject"`
	From       string `json:"from"`
	ID         string `json:"id"`
	To         string `json:"to"`
	SecondsAgo int    `json:"seconds_ago"`
}

// MsgPart is a message part.
type MsgPart struct {
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// A Msg is an email message.
type Msg struct {
	Fromfull   string            `json:"fromfull"`
	Headers    map[string]string `json:"headers"`
	Subject    string            `json:"subject"`
	RequestID  string            `json:"requestId"`
	IP         string            `json:"ip"`
	Parts      []MsgPart         `json:"parts"`
	From       string            `json:"from"`
	BeenRead   bool              `json:"been_read"`
	To         string            `json:"to"`
	ID         string            `json:"id"`
	Time       int               `json:"time"`
	SecondsAgo int               `json:"seconds_ago"`
}

// GetInbox retrieves the list of emails from the Mailinator API.
func GetInbox(tok, email string) ([]InboxMsg, error) {
	msgs := struct {
		Messages []InboxMsg `json:"messages"`
	}{}

	err := doReq(InboxURL, tok, email, &msgs)
	if err != nil {
		return nil, err
	}

	return msgs.Messages, nil
}

// GetEmail retrieves a email message from the Mailinator API.
//
// The returned int is the numebr of API requests left.
func GetEmail(tok, msgid string) (int, *Msg, error) {
	env := struct {
		Data      Msg `json:"data"`
		Remaining int `json:"apiEmailFetchesLeft"`
	}{}

	err := doReq(EmailURL, tok, msgid, &env)
	if err != nil {
		return 0, nil, err
	}

	return env.Remaining, &env.Data, nil
}
