package keys

import (
	"encoding/json"

	"github.com/znewman01/go-tuf/data"
)

// func init() {
// 	VerifierMap.Store(data.KeyTypeECDSA_SHA2_P256, NewEcdsaVerifier)
// }
//
// func NewSigstoreVerifier() Verifier {
// 	return &p256Verifier{}
// }
//
// type ecdsaSignature struct {
// 	R, S *big.Int
// }
//
// type p256Verifier struct {
// 	PublicKey data.HexBytes `json:"public"`
// 	key       *data.PublicKey
// }
//
// func (p *p256Verifier) Public() string {
// 	return p.PublicKey.String()
// }
//
// func (p *p256Verifier) Verify(msg, sigBytes []byte) error {
// 	x, y := elliptic.Unmarshal(elliptic.P256(), p.PublicKey)
// 	k := &ecdsa.PublicKey{
// 		Curve: elliptic.P256(),
// 		X:     x,
// 		Y:     y,
// 	}
//
// 	var sig ecdsaSignature
// 	if _, err := asn1.Unmarshal(sigBytes, &sig); err != nil {
// 		return err
// 	}
//
// 	hash := sha256.Sum256(msg)
//
// 	if !ecdsa.Verify(k, hash[:], sig.R, sig.S) {
// 		return errors.New("tuf: ecdsa signature verification failed")
// 	}
// 	return nil
// }
//
// func (p *p256Verifier) MarshalPublicKey() *data.PublicKey {
// 	return p.key
// }
//
// func (p *p256Verifier) UnmarshalPublicKey(key *data.PublicKey) error {
// 	if err := json.Unmarshal(key.Value, p); err != nil {
// 		return err
// 	}
// 	x, _ := elliptic.Unmarshal(elliptic.P256(), p.PublicKey)
// 	if x == nil {
// 		return errors.New("tuf: invalid ecdsa public key point")
// 	}
// 	p.key = key
// 	return nil
// }
//
// type ecdsaPublic struct {
// 	PublicKey data.HexBytes `json:"public"`
// }
//
// func Import(pub ecdsa.PublicKey) *data.PublicKey {
// 	keyValBytes, _ := json.Marshal(ecdsaPublic{PublicKey: elliptic.Marshal(pub.Curve, pub.X, pub.Y)})
// 	return &data.PublicKey{
// 		Type:       data.KeyTypeECDSA_SHA2_P256,
// 		Scheme:     data.KeySchemeECDSA_SHA2_P256,
// 		Algorithms: data.HashAlgorithms,
// 		Value:      keyValBytes,
// 	}
// }

type SigstoreIdentity struct {
	Issuer  string `json:"issuer"`
	Subject string `json:"identity"` // TODO: Subject->Identity
}

// TODO: do this for every PublicKey type
func PublicKeyToSigstoreIdentity(key *data.PublicKey) (*SigstoreIdentity, error) {
	if key.Type != data.KeyTypeSigstore || key.Scheme != data.KeySchemeSigstore {
		panic("uh oh") // TODO: define error
	}
	sigstore := &SigstoreIdentity{}
	if err := json.Unmarshal(key.Value, sigstore); err != nil {
		return nil, err
	}
	return sigstore, nil
}

func NewSigstorePublicKey(sigstore *SigstoreIdentity) (*data.PublicKey, error) {
	bytes, err := json.Marshal(sigstore)
	if err != nil {
		return nil, err
	}
	return &data.PublicKey{
		Type:       data.KeyTypeSigstore,
		Scheme:     data.KeySchemeSigstore,
		Algorithms: data.HashAlgorithms,
		Value:      bytes,
	}, nil
}
