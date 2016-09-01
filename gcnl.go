package gcnl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Default to English
const LanguageEnglish = "en"

type Document interface {
	Type() Type
	Language() string
	Content() string
	json.Marshaler
}

type Type string

const (
	TypeUnspecified Type = "TYPE_UNSPECIFIED"
	TypePlainText        = "PLAIN_TEXT"
	TypeHTML             = "HTML"
)

// MarshalJSON serializes a Document into JSON.
func MarshalJSON(doc Document) ([]byte, error) {
	s := struct {
		Type     Type   `json:"type"`
		Language string `json:"language"`
		Content  string `json:"content"`
	}{
		doc.Type(),
		doc.Language(),
		doc.Content(),
	}

	return json.Marshal(&s)
}

type PlainTextDocument struct{ content string }

func (doc *PlainTextDocument) Type() Type       { return TypePlainText }
func (doc *PlainTextDocument) Language() string { return LanguageEnglish }
func (doc *PlainTextDocument) Content() string  { return doc.content }

func NewPlainTextDocument(content string) Document {
	return &PlainTextDocument{content}
}

// MarshalJSON satisfies the json.Marshaler interface for PlainTextDocument.
func (doc *PlainTextDocument) MarshalJSON() ([]byte, error) {
	return MarshalJSON(doc)
}

type HTMLDocument struct{ content string }

func (doc *HTMLDocument) Type() Type       { return TypeHTML }
func (doc *HTMLDocument) Language() string { return LanguageEnglish }
func (doc *HTMLDocument) Content() string  { return doc.content }

func NewHTMLDocument(url string) (doc Document, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("NewHTMLDocument: returned %s", resp.Status)
		return
	}

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	doc = &HTMLDocument{string(d)}
	return
}

// MarshalJSON satisfies the json.Marshaler interface for HTMLDocument.
func (doc *HTMLDocument) MarshalJSON() ([]byte, error) {
	return MarshalJSON(doc)
}

type Encoding string

const (
	EncodingNone  Encoding = "NONE"
	EncodingUTF8           = "UTF8"
	EncodingUTF16          = "UTF16"
	EncodingUTF32          = "UTF32"

	// Default to UTF-8
	EncodingDefault = EncodingUTF8
)
