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

// Almost everything in this file (auth.go) and main.go was provided in the
// sheets/gmailAPI example code, and I, the author of this module, do not take
// credit for the code found in these file(s), but give credit to the authors
// at Google. (Everything else in this module is written by the author).

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

///////////////////////////////////////////////////////////////////////////////
//////////////////////            CREDENTIALS           ///////////////////////
///////////////////////////////////////////////////////////////////////////////

func (a *Access) ReadCredentials() {
	cred, err := os.ReadFile(a.CredentialsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(cred, a.Scopes...)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	a.Config = config
}

///////////////////////////////////////////////////////////////////////////////
//////////////////////              TOKEN               ///////////////////////
///////////////////////////////////////////////////////////////////////////////

// Retrieve a token, saves the token, then returns the generated client.
func (a *Access) GetClient() *http.Client {
	// The file [filename] stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	err := a.TokenFromFile()
	if err != nil {
		a.TokenFromWeb()
	}

	a.Client = a.Config.Client(a.Context, a.Token)
	return a.Client
}

// Retrieves a token from a local file.
func (a *Access) TokenFromFile() error {
	f, err := os.Open(a.TokenPath)
	if err != nil {
		return err
	}
	defer f.Close()
	var tok *oauth2.Token = &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	if err != nil {
		return err
	}
	a.Token = tok
	return nil
}

// Request a token from the web, then returns the retrieved token.
func (a *Access) TokenFromWeb() {
	authURL := a.Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := a.Config.Exchange(a.Context, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	a.Token = tok
	err = a.SaveToken()
	if err != nil {
		log.Fatal(err)
	}
}

// Saves a token to a file path.
func (a *Access) SaveToken() error {
	fmt.Printf("Saving credential file to: %s\n", a.TokenPath)
	f, err := os.OpenFile(a.TokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	a.LastRefreshed = time.Now()
	return json.NewEncoder(f).Encode(a.Token)
}

// *Access.Cycle(time.Duration) is used to cycle (refresh) the token, if
// ran with rate = 0 it'll refresh every 23 hours as default.
func (a *Access) Cycle(rate time.Duration) {
	if rate != 0 {
		a.RefreshRate = time.NewTicker(rate)
	}
	go func() {
		for {
			<-a.RefreshRate.C
			t, err := a.Config.TokenSource(a.Context, a.Token).Token()
			if err != nil {
				log.Println(err)
			}
			a.Token = t
			err = a.SaveToken()
			if err != nil {
				log.Println(err)
			}
			a.RefreshRate.Reset(rate)
		}
	}()
}
