package otto

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"io"

	"golang.org/x/crypto/ssh"
)

// Identity is DER encoded private key
type Identity struct {
	data []byte
}

// NewIdentity will generate a new ed25519 identity
func NewIdentity() (*Identity, error) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	// confidence check
	if _, err := ssh.NewSignerFromKey(privateKey); err != nil {
		return nil, err
	}

	return &Identity{privateKeyBytes}, nil
}

// ParseIdentity will parse the data as an identity
func ParseIdentity(data []byte) (*Identity, error) {
	pkey, err := x509.ParsePKCS8PrivateKey(data)
	if err != nil {
		return nil, err
	}
	// confidence check
	if _, err := ssh.NewSignerFromKey(pkey); err != nil {
		return nil, err
	}
	return &Identity{data}, nil
}

// Write will write the identity to the given writer
func (i *Identity) Write(w io.Writer) (int, error) {
	return w.Write(i.data)
}

// Signer return the SSH signer for the identity
func (i *Identity) Signer() ssh.Signer {
	privateKey, err := x509.ParsePKCS8PrivateKey(i.data)
	if err != nil {
		panic("x509.ParsePKCS8PrivateKey: " + err.Error())
	}

	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		panic("ssh.NewSignerFromKey: " + err.Error())
	}

	return signer
}

// PublicKey will return a DER-encoded representation of the public key for this identity
func (i *Identity) PublicKey() ssh.PublicKey {
	return i.Signer().PublicKey()
}

// String will return a base64-encoded representation of the identity
func (i *Identity) String() string {
	return base64.StdEncoding.EncodeToString(i.data)
}

// PublicKey will return a base64-encoded representation of the public key for this identity
func (i *Identity) PublicKeyString() string {
	return base64.StdEncoding.EncodeToString(i.PublicKey().Marshal())
}
