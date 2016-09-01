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

// A Mention is a wrapper for TextSpan objects.
type Mention struct {
	Text TextSpan `json:"text"`
}

// A TextSpan specifies the offset in the document where an entity was found.
type TextSpan struct {
	Content     string `json:"content"`
	BeginOffset int    `json:"beginOffset"`
}

// A Request represents the JSON object sent to the entities API.
type request struct {
	Document gcnl.Document `json:"document"`
	Encoding gcnl.Encoding `json:"encodingType"`

	key string
}

// NewRequest returns a Request object with the given API key.
func NewRequest(key string) *request {
	return &request{
		Encoding: gcnl.EncodingDefault,
		key:      key,
	}
}

// FromURL returns a slice of entities retrieved using a given a URL. It expects
// the content retrieved from URL to be valid HTML.
func (req *request) FromURL(url string) (entities []Entity, err error) {
	doc, err := gcnl.NewHTMLDocument(url)
	if err != nil {
		return
	}
	req.Document = doc
	return req.do()
}

// FromPlainText returns a slice of entities retrieved using a given plain text.
func (req *request) FromPlainText(content string) (entities []Entity, err error) {
	panic(errors.New("TODO: FromPlainText"))
}

// Do makes the actual API request for a given Request.
func (req *request) do() (entities []Entity, err error) {
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

	entities = jsonResp.Entities
	return
}
