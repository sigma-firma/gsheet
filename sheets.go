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

import (
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

///////////////////////////////////////////////////////////////////////////////
///////////////////     Access the Google Sheets API     //////////////////////
///////////////////////////////////////////////////////////////////////////////

// *Access.Sheets() gives access to the Google Sheets API via *Sheeter.Service
func (a *Access) Sheets() *Sheeter {
	service, err := sheets.NewService(
		a.Context,
		option.WithHTTPClient(a.GetClient()),
	)
	if err != nil {
		log.Println(err)
	}
	a.SheetsAPI = service
	return &Sheeter{Service: a.SheetsAPI}
}

///////////////////////////////////////////////////////////////////////////////
///////////////////      Use the Google Sheets API       //////////////////////
///////////////////////////////////////////////////////////////////////////////

// gsheet.Sheeter is a wrapper around the *sheets.Service type, giving us access to
// the Google Sheets API service.
type Sheeter struct {
	Service      *sheets.Service
	DefaultSheet *SpreadSheet
	// just guess what these could be used for
	Writables      map[string]*SpreadSheet
	Readables      map[string]*SpreadSheet
	ReadWriteables map[string]*SpreadSheet
}

// SpreadSheet is passed as the argument to the sheets related
// functions found herein.
type SpreadSheet struct {
	ID               string
	WriteRange       string
	Vals             []interface{}
	ReadRange        string
	ValueInputOption string
}

// *Sheeter.AppendRow() is used to append a row to a spreadsheet
func (s *Sheeter) AppendRow(sht *SpreadSheet) (*sheets.AppendValuesResponse, error) {
	var vr sheets.ValueRange
	vr.Values = append(vr.Values, sht.Vals)

	return s.Service.Spreadsheets.Values.Append(
		sht.ID,
		sht.WriteRange,
		&vr,
	).ValueInputOption(sht.ValueInputOption).Do()
}

// *Sheeter.Read() is used to read from a spreadsheet
func (s *Sheeter) Read(sht *SpreadSheet) (*sheets.ValueRange, error) {
	return s.Service.Spreadsheets.Values.Get(sht.ID, sht.ReadRange).Do()
}
