// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	cbor "github.com/fxamacker/cbor/v2"
)

// StringOrURI is either an arbitrary text string or a RFC3986 compliant URI.
//    string-or-uri = tstr / uri
type StringOrURI struct {
	text *string
	uri  *url.URL
}

// FromString initializes the StringOrURI value from the specified string,
// overwriting any existing value. If the value contains a colon (":"), then an
// attempt will be made to parse it as a URI (see RFC7519, section 2),
// otherwise, the value is assumed to be a non-URI string.
func (s *StringOrURI) FromString(value string) error {
	if strings.Contains(value, ":") {
		u, err := url.Parse(value)
		if err != nil {
			return err
		}

		s.uri = u
		s.text = nil
	} else {
		s.uri = nil
		s.text = &value
	}

	return nil
}

// FromURL initializes the StringOrURI value from the specified url.URL,
// overwriting any existing value.
func (s *StringOrURI) FromURL(value *url.URL) {
	s.text = nil
	s.uri = value
}

// IsURI returns true iff the underlying value is a URI and not a string.
// NOTE: this only indicates whether the value was set as such -- it possible
//       that the arbitrary string value happens to be a valid URI, however, if
//       it was not set as such, this will return false.
func (s StringOrURI) IsURI() bool {
	return s.uri != nil
}

// String returns the string representation of the StringOrURI value.
func (s StringOrURI) String() string {
	if s.uri != nil {
		return s.uri.String()
	}

	if s.text != nil {
		return *s.text
	}

	return ""
}

// ToURL will return the url.URL representation of the underlying value, if
// possible. This will attempt to parse the underlying string value as a URL if
// it isn't one already.
func (s StringOrURI) ToURL() (*url.URL, error) {
	if s.IsURI() {
		return s.uri, nil
	}

	if s.text != nil {
		return url.Parse(*s.text)
	}

	return nil, nil
}

// MarshalCBOR will encode the StringOrURI value as a CBOR text string,
// wrapping it in Tag 32, if it's a URI. See RFC7049, Section 2.4.4.3.
func (s StringOrURI) MarshalCBOR() ([]byte, error) {
	if s.IsURI() {
		tag := cbor.Tag{Number: 32, Content: s.uri.String()}
		return em.Marshal(tag)
	}

	if s.text != nil {
		return em.Marshal(s.text)
	}

	return []byte{}, nil
}

// UnmarshalCBOR attempts to initializes the StringOrURI from the presented
// CBOR data. The data must be a text string, possibly wrapped in a Tag with
// the value 32 (URI). See RFC7049, Section 2.4.4.3.
func (s *StringOrURI) UnmarshalCBOR(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	if isCBORTextString(data) {
		var value string
		err := dm.Unmarshal(data, &value)
		if err != nil {
			return err
		}

		err = s.FromString(value)
		if err != nil {
			return err
		}
	} else if isCBORTag(data) {
		var tag cbor.Tag
		err := dm.Unmarshal(data, &tag)
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

// MarshalJSON encodes the receiver StringOrURI into a JSON string
func (s StringOrURI) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON attempts at decoding the supplied JSON data into the receiver
// StringOrURI
func (s *StringOrURI) UnmarshalJSON(data []byte) error {
	var v string

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if err := s.FromString(v); err != nil {
		return err
	}

	return nil
}
