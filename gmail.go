// Copyright (c) 2025 Sigma-Firma

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package gsheet

import (
	"encoding/base64"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

///////////////////////////////////////////////////////////////////////////////
///////////////////     Access the Google Gmail API      //////////////////////
///////////////////////////////////////////////////////////////////////////////

// *Access.Gmail() gives usaccess to the Google Gmail API via *Gmailer.Service
func (a *Access) Gmail() *Gmailer {
	service, err := gmail.NewService(
		a.Context,
		option.WithHTTPClient(a.GetClient()),
	)
	if err != nil {
		log.Println(err)
	}
	a.GmailAPI = service
	return &Gmailer{Service: a.GmailAPI}
}

///////////////////////////////////////////////////////////////////////////////
///////////////////       Use the Google Gmail API       //////////////////////
///////////////////////////////////////////////////////////////////////////////

// Gmailer is a wrapper around the *gmail.Service type, giving us access to
// the Google Gmail API service.
type Gmailer struct {
	Service *gmail.Service
	Msgs    []*Msg
}

// Gmailer.Msg is an email message, self explanatory
type Msg struct {
	From      string
	To        string
	Subject   string
	Body      string
	ImagePath string
	MimeType  string
	Markup    string
	Bytes     []byte
	Formed    *gmail.Message
	PartialMetadata
}

// *Msg.Form() forms the message into a proper *gmail.Message type.
func (m *Msg) Form() *gmail.Message {
	var gm *gmail.Message = &gmail.Message{}
	var m_b []byte = []byte(
		"From: " + m.From + "\r\n" +
			"To: " + m.To + "\r\n" +
			"Subject: " + m.Subject + "\r\n" +
			"Content-Type: text/html; charset=\"utf-8\"\r\n\r\n" +
			m.Body)
	gm.Raw = base64.URLEncoding.EncodeToString(m_b)
	m.Formed = gm
	return m.Formed
}

// *Msg.Send() allows us to send mail
func (g *Gmailer) Send(m *Msg) error {
	m.Form()
	sendCall := g.Service.Users.Messages.Send(m.From, m.Formed)
	_, err := sendCall.Do()
	if err != nil {
		return err
	}

	return nil
}

// *Gmailer.MarkAs() allows you to mark an email with a specific label using the
// gmail.ModifyMessageRequest struct.
func (g *Gmailer) MarkAs(msg *gmail.Message, req *gmail.ModifyMessageRequest) (*gmail.Message, error) {
	return g.Service.Users.Messages.Modify("me", msg.Id, req).Do()
}

// *Gmailer.MarkAllAsRead() removes the UNREAD label from all emails.
func (g *Gmailer) MarkAllAsRead() error {
	// Request to remove the label ID "UNREAD"
	req := &gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{"UNREAD"},
	}

	// Get the messages labeled "UNREAD"
	msgs, err := g.Query("label:UNREAD")
	if err != nil {
		return err
	}

	// For each UNREAD message, request to remove the "UNREAD" label (thus
	// maring it as "READ").
	for _, v := range msgs {
		_, err := g.MarkAs(v, req)
		if err != nil {
			return err
		}
	}

	return nil
}

// *Gmailer.GetBody() gets, decodes, and returns the body of the email. It
// returns an error if decoding goes wrong. mimeType is used to indicate
// whether you want the plain text or html encoding ("text/html", "text/plain").
func (g *Gmailer) GetBody(msg *gmail.Message, mimeType string) (string, error) {
	// Loop through the message payload parts to find the parts with the
	// mimetypes we want.
	for _, v := range msg.Payload.Parts {
		if v.MimeType == "multipart/alternative" {
			for _, l := range v.Parts {
				if l.MimeType == mimeType && l.Body.Size >= 1 {
					dec, err := g.DecodeEmailBody(l.Body.Data)
					if err != nil {
						return "", err
					}
					return dec, nil
				}
			}
		}
		if v.MimeType == mimeType && v.Body.Size >= 1 {
			dec, err := g.DecodeEmailBody(v.Body.Data)
			if err != nil {
				return "", err
			}
			return dec, nil
		}
	}
	return "", errors.New("Couldn't Read Body")
}

// *Gmailer.DecodeEmailBody() is used to decode the email body by converting
// from URLEncoded base64 to a string.
func (g *Gmailer) DecodeEmailBody(data string) (string, error) {
	decoded, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// *Gmailer.ReceivedTime() parses and converts a Unix time stamp into a human
// readable format ().
func (g *Gmailer) ReceivedTime(datetime int64) (time.Time, error) {
	conv := strconv.FormatInt(datetime, 10)
	// Remove trailing zeros.
	conv = conv[:len(conv)-3]
	tc, err := strconv.ParseInt(conv, 10, 64)
	if err != nil {
		return time.Unix(0, 0), err
	}
	return time.Unix(tc, 0), nil
}

// *Gmailer.Query() queries the inbox for a string following the search style
// of the gmail online mailbox.
// example:
// "in:sent after:2017/01/01 before:2017/01/30"
func (g *Gmailer) Query(query string) ([]*gmail.Message, error) {
	inbox, err := g.Service.Users.Messages.List("me").Q(query).Do()
	if err != nil {

		return []*gmail.Message{}, err
	}
	msgs, err := g.GetByID(inbox)
	if err != nil {
		return msgs, err
	}
	return msgs, nil
}

// *Gmailer.GetByID() gets emails individually by ID. This is necessary because this is
// how the gmail API is set [0][1] up apparently (but why?).
// [0] https://developers.google.com/gmail/api/v1/reference/users/messages/get
// [1] https://stackoverflow.com/questions/36365172/message-payload-is-always-null-for-all-messages-how-do-i-get-this-data
func (g *Gmailer) GetByID(msgs *gmail.ListMessagesResponse) ([]*gmail.Message, error) {
	var msgSlice []*gmail.Message
	for _, v := range msgs.Messages {
		msg, err := g.Service.Users.Messages.Get("me", v.Id).Do()
		if err != nil {
			return msgSlice, err
		}
		msgSlice = append(msgSlice, msg)
	}
	return msgSlice, nil
}

// *Gmailer.GetMessages() gets and returns gmail messages
func (g *Gmailer) GetMessages(howMany uint) ([]*gmail.Message, error) {
	var msgSlice []*gmail.Message

	// Get the messages
	inbox, err := g.Service.Users.Messages.List("me").MaxResults(int64(howMany)).Do()
	if err != nil {
		return msgSlice, err
	}

	msgs, err := g.GetByID(inbox)
	if err != nil {
		return msgs, err
	}
	return msgs, nil
}

// *Gmailer.CheckFByLabel() checks for mail matching the specified label.
func (g *Gmailer) CheckByLabel(label string) (int64, error) {
	inbox, err := g.Service.Users.Labels.Get("me", label).Do()
	if err != nil {
		return -1, err
	}

	if inbox.MessagesUnread == 0 && inbox.ThreadsUnread == 0 {
		return 0, nil
	}

	return inbox.MessagesUnread + inbox.ThreadsUnread, nil
}

// *Gmailer.CheckForUnread() checks for mail labeled "UNREAD".
// NOTE: When checking your inbox for unread messages, it's not uncommon for
// it to return thousands of unread messages that you don't know about. To see
// them in gmail, query your mail for "label:unread". For CheckForUnread to
// work properly you need to mark all mail as read either through gmail or
// through the MarkAllAsRead() function found in this library.
func (g *Gmailer) CheckForUnread() (int64, error) {
	inbox, err := g.Service.Users.Labels.Get("me", "UNREAD").Do()
	if err != nil {
		return -1, err
	}

	if inbox.MessagesUnread == 0 && inbox.ThreadsUnread == 0 {
		return 0, nil
	}

	return inbox.MessagesUnread + inbox.ThreadsUnread, nil
}

// *Gmailer.GetLabels() gets a list of the labels used in the users inbox.
func (g *Gmailer) GetLabels() (*gmail.ListLabelsResponse, error) {
	return g.Service.Users.Labels.List("me").Do()
}

// PartialMetadata stores email metadata. Some fields may sound redundant, but
// in fact have different contexts. Some are slices of string because the ones
// that have multiple values are still being sorted from those that don't.
type PartialMetadata struct {
	// Sender is the entity that originally created and sent the message
	Sender string
	// From is the entity that sent the message to you (e.g. googlegroups). Most
	// of the time this information is only relevant to mailing lists.
	From string
	// Subject is the email subject
	Subject string
	// Mailing list contains the name of the mailing list that the email was
	// posted to, if any.
	MailingList string
	// CC is the "carbon copy" list of addresses
	CC []string
	// To is the recipient of the email.
	To []string
	// ThreadTopic contains the topic of the thread (e.g. google groups threads)
	ThreadTopic []string
	// DeliveredTo is who the email was sent to. This can contain multiple
	// addresses if the email was forwarded.
	DeliveredTo []string
}

// *Gmailer.GetPartialMetadata() gets some of the useful metadata from the headers.
func (g *Gmailer) GetPartialMetadata(msg *gmail.Message) *PartialMetadata {
	info := &PartialMetadata{}
	for _, v := range msg.Payload.Headers {
		switch strings.ToLower(v.Name) {
		case "sender":
			info.Sender = v.Value
		case "from":
			info.From = v.Value
		case "subject":
			info.Subject = v.Value
		case "mailing-list":
			info.MailingList = v.Value
		case "cc":
			info.CC = append(info.CC, v.Value)
		case "to":
			info.To = append(info.To, v.Value)
		case "thread-Topic":
			info.ThreadTopic = append(info.ThreadTopic, v.Value)
		case "delivered-To":
			info.DeliveredTo = append(info.DeliveredTo, v.Value)
		}
	}
	return info
}
