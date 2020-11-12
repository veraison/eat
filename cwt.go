package eat

// Defined by RFC8392
type CWTClaims struct {
	Issuer     *string      `cbor:"1,keyasint,omitempty" json:"iss,omitempty"`
	Subject    *string      `cbor:"2,keyasint,omitempty" json:"sub,omitempty"`
	Audience   *Audience    `cbor:"3,keyasint,omitempty" json:"aud,omitempty"`
	Expiration *NumericDate `cbor:"4,keyasint,omitempty" json:"exp,omitempty"`
	NotBefore  *NumericDate `cbor:"5,keyasint,omitempty" json:"nbf,omitempty"`
	IssuedAt   *NumericDate `cbor:"6,keyasint,omitempty" json:"iat,omitempty"`
	CwtID      *[]byte      `cbor:"7,keyasint,omitempty" json:"cti,omitempty"`
}
