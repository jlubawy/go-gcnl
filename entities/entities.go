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

var errMissingKey = errors.New("must provide and API key")

type Entity struct {
	Name     string            `json:"name"`
	Type     Type              `json:"type"`
	Metadata map[string]string `json:"metadata"`
	Salience float64           `json:"salience"`
	Mentions []Mention         `json:"mentions"`
}

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

type Mention struct {
	Text TextSpan `json:"text"`
}

type TextSpan struct {
	Content     string `json:"content"`
	BeginOffset int    `json:"beginOffset"`
}

type Request struct {
	Document gcnl.Document `json:"document"`
	Encoding gcnl.Encoding `json:"encodingType"`

	key string
}

func NewRequest(key string) *Request {
	return &Request{
		Encoding: gcnl.EncodingDefault,
		key:      key,
	}
}

func (req *Request) FromURL(url string) (entities []Entity, err error) {
	doc, err := gcnl.NewHTMLDocument(url)
	if err != nil {
		return
	}
	req.Document = doc
	return req.do()
}

func (req *Request) FromPlainText(content string) (entities []Entity, err error) {
	panic(errors.New("TODO: FromPlainText"))
}

func (req *Request) do() (entities []Entity, err error) {
	if len(req.key) == 0 {
		err = errMissingKey
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
