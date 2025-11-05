# Entity Attestation Token

The `eat` package provides a golang API for manipulating Entity Attestation
Tokens defined in
[RFC 9711](https://datatracker.ietf.org/doc/rfc9711/).

## Supported EAT Features

See [eat.go](eat.go) for more detail.

claim | CBOR key id | Supported?
--|--|--
eat_nonce | 10 | ✅
ueid | 256 | ✅
sueids | 257 | :x:
oemid | 258 | :warning: currently `oemid-pem (int)` is not supported
hwmodel | 259 | ✅
hwversion | 260 | ✅
uptime | 261 | ✅
oemboot | 262 | ✅
dbgstat | 263 | ✅
location | 264 | ✅
eat_profile | 265 | ✅
submods | 266 | ✅
bootcount | 267 | ✅
bootseed | 268 | ✅
dloas | 269 | :x:
swname | 270 | ✅
swversion | 271 | ✅
manifests | 272 | ⚠️ (see [Supported Type for Manifests and Measurements](#supported-type-for-manifests-and-measurements))
measurements | 273 | ⚠️ (see [Supported Type for Manifests and Measurements](#supported-type-for-manifests-and-measurements))
measres | 274 | :x:
intuse | 275 | :x:

## Supported CWT Features

See [cwt.go](cwt.go) for more detail.

claim | CBOR key id | Supported?
--|--|--
iss | 1 | ✅
sub | 2 | ✅
aud | 3 | ✅
exp | 4 | ✅
nbf | 5 | ✅
iat | 6 | ✅
cti | 7 | ⚠️ no jti support
cnf | 8 | ⚠️ supports only OKP and EC2 COSE_Key, no EncryptedKey support

## Supported Type for Manifests and Measurements

[RFC 9711](https://www.rfc-editor.org/rfc/rfc9711.html#name-payload-cddl) defines extensible Manifests and Measurements.

The [measured_component.go](./measured_component.go) experimentally provides encoding/decoding feature for [EAT Measured Component (v05)](https://datatracker.ietf.org/doc/draft-ietf-rats-eat-measured-component/05/).
The `untagged-coswid` encoder/decoder is provided by [veraison/swid](https://github.com/veraison/swid).

coap-conent-type | id | Supported?
--|--|--
`application/swid+cbor` (untagged-coswid) | 258 | ✅
`application/measured-component+cbor` | TBD1 in [draft-ietf-rats-eat-measured-component](https://datatracker.ietf.org/doc/draft-ietf-rats-eat-measured-component/) | ✅ e.g. `cbor.Unmarshal(measurement.Format, &mc)`
`application/measured-component+json` | TBD2 in [draft-ietf-rats-eat-measured-component](https://datatracker.ietf.org/doc/draft-ietf-rats-eat-measured-component/) | ✅ e.g. `json.Unmarshal(measurement.Format, &mc)`
