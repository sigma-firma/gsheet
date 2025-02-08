package gsheet

import (
	"os"

	"testing"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/sheets/v4"
)

func TestAuth(t *testing.T) {
	var access *Access = NewAccess(
		os.Getenv("HOME")+"/credentials/credentials.json",
		os.Getenv("HOME")+"/credentials/token.json",
		[]string{
			sheets.SpreadsheetsScope,
			gmail.GmailComposeScope,
			gmail.GmailLabelsScope,
			gmail.GmailSendScope,
			gmail.GmailModifyScope,
		},
	)

	var gm = access.Gmail()
	msgs, err := gm.Query("category:forums after:2025/02/01 before:2025/02/30")
	if err != nil {
		t.Error("connect", err)
	}

	req := &gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{"UNREAD"},
		AddLabelIds:    []string{"OLD"},
	}

	// Range over the messages
	for _, msg := range msgs {
		_, err := gm.MarkAs(msg, req)
		if err != nil {
			t.Error("connect", err)
		}
	}
}
