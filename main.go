// Package gsheet is a Go library for reading and writing data to/from Google
// Sheets using the Google Sheets API, and also send email using the Gmail API.
//
// Copyright (c) 2025 Sigma-Firma
//
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

// Almost everything in this file (main.go) and auth.go was provided in the
// sheetsAPI example code, and I, the author of this module, do not take credit
// for the code found in these file(s), but give credit to the authors at
// Google.
// (Everything else in this module is written by the author).
import (
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/sheets/v4"
)

// Access contains values used for accessing Googles API services, including
// access token, API credentials, and scopes
type Access struct {
	Context context.Context
	// TokenPath is the path to your token.json file. If you're having trouble
	// authenticating, try deleting this file and running the program
	// again. This should renew your token. If you've never run this
	// program, you may not have a token. This program will generate a
	// token for you. Also see: Access.CredentialsPath
	TokenPath string
	// CredentialsPath is the path to your credentials.json file. This file can
	// be obtained from the API Keys section of Google Cloud Platform. You
	// may need to generate the file and enable the API you're interfacing with
	// from within Google Cloud Platform.
	CredentialsPath string
	// Scopes define what level(s) of access we'll have to the API service.
	// If modifying these scopes, delete your previously saved token.json.
	Scopes []string
	// Config is mostly generated and used by the API.
	Config *oauth2.Config
	// The token is the token used to authenticate with the API, and will need
	// to be refreshed using *Access.Cycle()
	Token *oauth2.Token
	// GmailAPI is used as an alternative access point to the Gmail API
	// service, if you don't wish to use the Gmailer struct. see: gmail.go
	GmailAPI *gmail.Service
	// SheetsAPI is used as an alternative access point to the Sheets API
	// service, if you don't wish to use the Sheeter struct. see: sheets.go
	SheetsAPI *sheets.Service
	// RefreshRate
	RefreshRate *time.Ticker
	// LastRefreshed
	LastRefreshed time.Time
	Client        *http.Client
}

// NewAccess() instantiates a new *Access struct, initializing it with default
// values most people will probably need, and authenticated the user, taking
// them through the manual authentication process as necessary (see READEME.md)
func NewAccess(credentialsPath, tokenPath string, scopes []string) *Access {
	var a *Access = &Access{
		Context:         context.Background(),
		CredentialsPath: credentialsPath,
		TokenPath:       tokenPath,
		Scopes:          scopes,
		Config:          &oauth2.Config{},
		Token:           &oauth2.Token{},
		RefreshRate:     time.NewTicker(23 * time.Hour),
	}
	// see: auth.go
	a.ReadCredentials()

	a.Cycle(0)
	return a
}

// *Access.Connect(any) connects us to a service. At this time gsheet only
// supports the Gmail API, the Google Sheets API, or both. Pass either of the
// following empty (as empty structs):
//
// *gmail.Sevice{}
// *sheets.Service{}
func (a *Access) Connect(service any) {
	switch service.(type) {
	case *gmail.Service:
		a.Gmail() // see: gmail.go
	case *sheets.Service:
		a.Sheets() // see: sheets.go
	}
}
