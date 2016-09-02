// go-gcnl - Golang library for accessing the Google Cloud Natural Language API
// Copyright (C) 2016 Josh Lubawy <jlubawy@gmail.com>
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, write to the Free Software Foundation, Inc.,
// 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

package entities

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jlubawy/go-gcnl"
)

const Endpoint = "https://language.googleapis.com/v1beta1/documents:analyzeEntities"

var ErrMissingKey = errors.New("must provide an API key")

// An Entity represents a phrase is the text that is a known entity of a given Type.
type Entity struct {
	Name     string            `json:"name"`
	Type     Type              `json:"type"`
	Metadata map[string]string `json:"metadata"`
	Salience float64           `json:"salience"`
	Mentions []Mention         `json:"mentions"`
}

// A Type specifies the valid entity types returned by the API.
type Type string

const (
	TypeUnknown      Type = "UNKNOWN"
	TypePerson            = "PERSON"
	TypeLocation          = "LOCATION"
	TypeOrganization      = "ORGANIZATION"
	TypeEvent             = "EVENT"
	TypeWorkOfArt         = "WORK_OF_ART"
	TypeConsumerGood      = "CONSUMER_GOOD"
	TypeOther             = "OTHER"
)

// A Map is a map of Types to Entities.
type Map map[Type][]Entity

// A Mention is a wrapper for TextSpan objects.
type Mention struct {
	TextSpan TextSpan `json:"text"`
}

// A TextSpan specifies the offset in the document where an entity was found.
type TextSpan struct {
	Content     string `json:"content"`
	BeginOffset int    `json:"beginOffset"`
}

// A Request represents the JSON object sent to the entities API.
type request struct {
	Doc gcnl.Document `json:"document"`
	Enc gcnl.Encoding `json:"encodingType"`
	key string
}

// NewRequest returns a Request object with the given API key.
func NewRequest(key string) *request {
	return &request{
		Enc: gcnl.EncodingDefault,
		key: key,
	}
}

// Document returns the document used in the request.
func (req *request) Document() gcnl.Document {
	return req.Doc
}

// FromURL returns a slice of entities retrieved using a given a URL. It expects
// the content retrieved from URL to be valid HTML.
func (req *request) FromURL(url string) (entityMap Map, err error) {
	doc, err := gcnl.NewHTMLDocument(url)
	if err != nil {
		return
	}
	req.Doc = doc
	return req.do()
}

// FromPlainText returns a slice of entities retrieved using a given plain text.
func (req *request) FromPlainText(content string) (entityMap Map, err error) {
	req.Doc = gcnl.NewPlainTextDocument(content)
	return req.do()
}

// Do makes the actual API request for a given Request.
func (req *request) do() (entityMap Map, err error) {
	if len(req.key) == 0 {
		err = ErrMissingKey
		return
	}

	d, err := json.Marshal(req)
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(d)

	r, err := http.NewRequest("POST", fmt.Sprintf("%s?key=%s", Endpoint, req.key), buf)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("entities.do: returned %s", resp.Status)
		return
	}

	var jsonResp struct {
		Entities []Entity `json:"entities"`
	}

	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		return
	}

	entityMap = make(Map)

	for _, e := range jsonResp.Entities {
		if entityMap[e.Type] == nil {
			entityMap[e.Type] = make([]Entity, 0)
		}
		entityMap[e.Type] = append(entityMap[e.Type], e)
	}

	return
}
