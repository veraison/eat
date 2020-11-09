package eat

import (
	"fmt"
	"net/url"

	cbor "github.com/fxamacker/cbor/v2"
)

type StringOrURI struct {
	text *string
	uri  *url.URL
}

func (s *StringOrURI) FromString(value string) {
	s.uri = nil
	s.text = &value
}

func (s *StringOrURI) FromURL(value *url.URL) {
	s.text = nil
	s.uri = value
}

func (s *StringOrURI) IsURI() bool {
	return s.uri != nil
}

func (s *StringOrURI) String() string {
	if s.uri != nil {
		return s.uri.String()
	}

	if s.text != nil {
		return *s.text
	}

	return ""
}

func (s *StringOrURI) ToURL() (*url.URL, error) {
	if s.IsURI() {
		return s.uri, nil
	}

	if s.text != nil {
		return url.Parse(*s.text)
	}

	return nil, nil
}

func (s *StringOrURI) MarshalCBOR() ([]byte, error) {
	if s.IsURI() {
		tag := cbor.Tag{Number: 32, Content: s.uri.String()}
		return cbor.Marshal(tag)
	}

	if s.text != nil {
		return cbor.Marshal(s.text)
	}

	return []byte{}, nil
}

func (s *StringOrURI) UnmarshalCBOR(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	if isCBORTextString(data) {
		var value string
		err := cbor.Unmarshal(data, &value)
		if err != nil {
			return err
		}

		s.FromString(value)
	} else if isCBORTag(data) {
		var tag cbor.Tag
		err := cbor.Unmarshal(data, &tag)
		if err != nil {
			return err
		}

		if tag.Number != 32 {
			return fmt.Errorf("must be URI (tag 32), found: %v", tag.Number)
		}

		switch tag.Content.(type) {
		case string:
			u, err := url.Parse(tag.Content.(string))
			if err != nil {
				return err
			}
			s.FromURL(u)
		default:
			return fmt.Errorf("URI tag value must be a string")
		}
	} else {
		return fmt.Errorf("must be a text string or a URI tag")
	}

	return nil
}

func isCBORTag(data []byte) bool {
	return (data[0] & 0xe0) == 0xc0
}

func isCBORTextString(data []byte) bool {
	return (data[0] & 0xe0) == 0x60
}
