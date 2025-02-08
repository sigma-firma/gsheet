[//]: # (Copyright [c] 2025 Σfirma)

[//]: # (Permission is hereby granted, free of charge, to any person obtaining a copy)
[//]: # (of this software and associated documentation files [the "Software"], to deal)
[//]: # (in the Software without restriction, including without limitation the rights)
[//]: # (to use, copy, modify, merge, publish, distribute, sublicense, and/or sell)
[//]: # (copies of the Software, and to permit persons to whom the Software is)
[//]: # (furnished to do so, subject to the following conditions:)

[//]: # (The above copyright notice and this permission notice shall be included in all)
[//]: # (copies or substantial portions of the Software.)

[//]: # (THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR)
[//]: # (IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,)
[//]: # (FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE)
[//]: # (AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER)
[//]: # (LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,)
[//]: # (OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE)
[//]: # (SOFTWARE.)

# DEPRECATED See: [gsheet](https://github.com/sigma-firma/gsheet) for the same functionality but for human beings

# INBOXER

Package inboxer is a Go package for checking your gmail inbox, it has the following features:

  - Send emails
  - Mark emails (read/unread/important/etc)
  - Get labels used in inbox
  - Get emails by query (eg "in:sent after:2017/01/01 before:2017/01/30")
  - Get email metadata
  - Get email main body ("text/plain", "text/html")
  - Get the number of unread messages
  - Convert email dates to human readable format

#  USE

## CREDENTIALS:

For inboxer to work you must have a gmail account and a file named "client_secret.json" containing your authorization info in the root directory of your project. To obtain credentials please see step one of this guide: https://developers.google.com/gmail/api/quickstart/go

 Turning on the gmail API

 - Use this wizard (https://console.developers.google.com/start/api?id=gmail) to create or select a project in the Google Developers Console and automatically turn on the API. Click Continue, then Go to credentials.

 - On the Add credentials to your project page, click the Cancel button.

 - At the top of the page, select the OAuth consent screen tab. Select an Email address, enter a Product name if not already set, and click the Save button.

 - Select the Credentials tab, click the Create credentials button and select OAuth client ID.

 - Select the application type Other, enter the name "Gmail API Quickstart", and click the Create button.

 - Click OK to dismiss the resulting dialog.

 - Click the file_download (Download JSON) button to the right of the client ID.

 - Move this file to your working directory and rename it client_secret.json.


```go
package main

import (
        "context"
        "fmt"

        "github.com/sigma-firma/gmailAPI"
        "github.com/sigma-firma/inboxer"
        gmail "google.golang.org/api/gmail/v1"
)

func main() {
        // Connect to the gmail API service.
        ctx := context.Background()
        srv := gmailAPI.ConnectToService(ctx, gmail.MailGoogleComScope)

        msgs, err := inboxer.Query(srv, "category:forums after:2017/01/01 before:2017/01/30")
        if err != nil {
                fmt.Println(err)
        }

        // Range over the messages
        for _, msg := range msgs {
                fmt.Println("========================================================")
                time, err := inboxer.ReceivedTime(msg.InternalDate)
                if err != nil {
                        fmt.Println(err)
                }
                fmt.Println("Date: ", time)
                md := inboxer.GetPartialMetadata(msg)
                fmt.Println("From: ", md.From)
                fmt.Println("Sender: ", md.Sender)
                fmt.Println("Subject: ", md.Subject)
                fmt.Println("Delivered To: ", md.DeliveredTo)
                fmt.Println("To: ", md.To)
                fmt.Println("CC: ", md.CC)
                fmt.Println("Mailing List: ", md.MailingList)
                fmt.Println("Thread-Topic: ", md.ThreadTopic)
                fmt.Println("Snippet: ", msg.Snippet)
                body, err := inboxer.GetBody(msg, "text/plain")
                if err != nil {
                        fmt.Println(err)
                }
                fmt.Println(body)
        }
}

```
## SENDING MAIL
```go
package main

import (
	"context"
	"log"

	"github.com/sigma-firma/gmailAPI"
	"github.com/sigma-firma/inboxer"
	gmail "google.golang.org/api/gmail/v1"
)

func main() {
    // Connect to the gmail API service.
	ctx := context.Background()
	srv := gmailAPI.ConnectToService(ctx, gmail.MailGoogleComScope)

    // Create a message
	var msg *inboxer.Msg = &inboxer.Msg{
		From:    "me",  // the authenticated user
		To:      "leadership@firma.com",
		Subject: "testing",
		Body:    "testing gmail api. lmk if you get this scott",
	}

    // send the email with the message
	err := msg.Send(srv)
	if err != nil {
		log.Println(err)
	}
}
```
## QUERIES

```go
func main() {
        // Connect to the gmail API service.
        ctx := context.Background()
        srv := gmailAPI.ConnectToService(ctx, gmail.GmailComposeScope)
        msgs, err := inboxer.Query(srv, "category:forums after:2017/01/01 before:2017/01/30")
        if err != nil {
                fmt.Println(err)
        }

        // Range over the messages
        for _, msg := range msgs {
                // do stuff
        }
}

```
## MARKING EMAILS

```go
func main() {
        // Connect to the gmail API service.
        ctx := context.Background()
        srv := gmailAPI.ConnectToService(ctx, gmail.GmailComposeScope)

        msgs, err := inboxer.Query(srv, "category:forums after:2017/01/01 before:2017/01/30")
        if err != nil {
                fmt.Println(err)
        }

        req := &gmail.ModifyMessageRequest{
                RemoveLabelIds: []string{"UNREAD"},
                AddLabelIds: []string{"OLD"}
        }

        // Range over the messages
        for _, msg := range msgs {
                msg, err := inboxer.MarkAs(srv, msg, req)
        }
}

```
## MARK ALL "UNREAD" EMAILS AS "READ"

```go
func main() {
        // Connect to the gmail API service.
        ctx := context.Background()
        srv := gmailAPI.ConnectToService(ctx, gmail.GmailComposeScope)

        inboxer.MarkAllAsRead(srv)
}
```
## GETTING LABELS

```go
func main() {
        // Connect to the gmail API service.
        ctx := context.Background()
        srv := gmailAPI.ConnectToService(ctx, gmail.GmailComposeScope)

        labels, err := inboxer.GetLabels(srv)
        if err != nil {
                fmt.Println(err)
        }

        for _, label := range labels {
                fmt.Println(label)
        }
}

```
## METADATA

```go
func main() {
        // Connect to the gmail API service.
        ctx := context.Background()
        srv := gmailAPI.ConnectToService(ctx, gmail.MailGoogleComScope)

        msgs, err := inboxer.Query(srv, "category:forums after:2017/01/01 before:2017/01/30")
        if err != nil {
                fmt.Println(err)
        }

        // Range over the messages
        for _, msg := range msgs {
                fmt.Println("========================================================")
                md := inboxer.GetPartialMetadata(msg)
                fmt.Println("From: ", md.From)
                fmt.Println("Sender: ", md.Sender)
                fmt.Println("Subject: ", md.Subject)
                fmt.Println("Delivered To: ", md.DeliveredTo)
                fmt.Println("To: ", md.To)
                fmt.Println("CC: ", md.CC)
                fmt.Println("Mailing List: ", md.MailingList)
                fmt.Println("Thread-Topic: ", md.ThreadTopic)
        }
}

```
## GETTING THE EMAIL BODY

```go
func main() {
        // Connect to the gmail API service.
        ctx := context.Background()
        srv := gmailAPI.ConnectToService(ctx, gmail.GmailComposeScope)
        msgs, err := inboxer.Query(srv, "category:forums after:2017/01/01 before:2017/01/30")
        if err != nil {
                fmt.Println(err)
        }

        // Range over the messages
        for _, msg := range msgs {
                body, err := inboxer.GetBody(msg, "text/plain")
                if err != nil {
                        fmt.Println(err)
                }
                fmt.Println(body)
        }
}

```
## GETTING THE NUMBER OF UNREAD MESSAGES

```go
// NOTE: to actually view the email text use inboxer.Query and query for unread
// emails.
func main() {
        // Connect to the gmail API service.
        ctx := context.Background()
        srv := gmailAPI.ConnectToService(ctx, gmail.GmailComposeScope)

        // num will be -1 on err
        num, err :=	inboxer.CheckForUnread(srv)
        if err != nil {
                fmt.Println(err)
        }
        fmt.Printf("You have %s unread emails.", num)
}


```
## CONVERTING DATES

```go
// Convert UNIX time stamps to human readable format
func main() {
        // Connect to the gmail API service.
        ctx := context.Background()
        srv := gmailAPI.ConnectToService(ctx, gmail.GmailComposeScope)

        msgs, err := inboxer.Query(srv, "category:forums after:2017/01/01 before:2017/01/30")
        if err != nil {
                fmt.Println(err)
        }

        // Range over the messages
        for _, msg := range msgs {
                // Convert the date
                time, err := inboxer.ReceivedTime(msg.InternalDate)
                if err != nil {
                        fmt.Println(err)
                }
                fmt.Println("Date: ", time)
        }
}

```

## SNIPPET

```go
// Snippets are not really part of the package but I'm including them in the doc
// because they'll likely be useful to anyone working with this package.
func main() {
        // Connect to the gmail API service.
        ctx := context.Background()
        srv := gmailAPI.ConnectToService(ctx, gmail.GmailComposeScope)

        msgs, err := inboxer.Query(srv, "category:forums after:2017/01/01 before:2017/01/30")
        if err != nil {
                fmt.Println(err)
        }

        // Range over the messages
        for _, msg := range msgs {
                // this one is part of the api
                fmt.Println(msg.Snippet)
        }
}
```
