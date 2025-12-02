// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"fmt"
	"net/url"

	cbor "github.com/fxamacker/cbor/v2"
)

type Dloa struct {
	Registrar        url.URL
	PlatformLabel    string
	ApplicationLabel *string
}

//nolint:gocritic
func (d Dloa) MarshalJSON() ([]byte, error) {
	tmp := []interface{}{d.Registrar.String(), d.PlatformLabel}
	if d.ApplicationLabel != nil {
		tmp = append(tmp, *d.ApplicationLabel)
	}
	return json.Marshal(tmp)
}

//nolint:dupl
func (d *Dloa) UnmarshalJSON(data []byte) error {
	var tmp []json.RawMessage
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	if len(tmp) < 2 || len(tmp) > 3 {
		return fmt.Errorf("invalid array length: %d", len(tmp))
	}

	var regStr string
	if err := json.Unmarshal(tmp[0], &regStr); err != nil {
		return err
	}
	regURL, err := url.Parse(regStr)
	if err != nil {
		return err
	}
	d.Registrar = *regURL

	if err := json.Unmarshal(tmp[1], &d.PlatformLabel); err != nil {
		return err
	}

	if len(tmp) == 3 {
		var appLabel string
		if err := json.Unmarshal(tmp[2], &appLabel); err != nil {
			return err
		}
		d.ApplicationLabel = &appLabel
	}

	return nil
}

//nolint:gocritic
func (d Dloa) MarshalCBOR() ([]byte, error) {
	tmp := []interface{}{d.Registrar.String(), d.PlatformLabel}
	if d.ApplicationLabel != nil {
		tmp = append(tmp, *d.ApplicationLabel)
	}
	return em.Marshal(tmp)
}

//nolint:dupl
func (d *Dloa) UnmarshalCBOR(data []byte) error {
	var tmp []cbor.RawMessage
	if err := dm.Unmarshal(data, &tmp); err != nil {
		return err
	}
	if len(tmp) < 2 || len(tmp) > 3 {
		return fmt.Errorf("invalid array length: %d", len(tmp))
	}

	var regStr string
	if err := dm.Unmarshal(tmp[0], &regStr); err != nil {
		return err
	}
	regURL, err := url.Parse(regStr)
	if err != nil {
		return err
	}
	d.Registrar = *regURL

	if err := dm.Unmarshal(tmp[1], &d.PlatformLabel); err != nil {
		return err
	}

	if len(tmp) == 3 {
		var appLabel string
		if err := dm.Unmarshal(tmp[2], &appLabel); err != nil {
			return err
		}
		d.ApplicationLabel = &appLabel
	}

	return nil
}
