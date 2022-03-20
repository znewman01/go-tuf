package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/flynn/go-docopt"
	"github.com/theupdateframework/go-tuf"
	"github.com/theupdateframework/go-tuf/pkg/keys"
)

func init() {
	register("import-key", cmdImportKey, `
usage: tuf import-key --role=<role> <file>

Import an ECDSA P256 public key from a PEM file.

Assumes no passphrase for now, I'm lazy.
`)
}

func cmdImportKey(args *docopt.Args, repo *tuf.Repo) error {
	pemBytes, err := os.ReadFile(args.String["<file>"])
	if err != nil {
		return err
	}
	derBytes, _ := pem.Decode(pemBytes)
	if derBytes == nil {
		return errors.New("PEM decoding failed")
	}
	pubKey, err := x509.ParsePKIXPublicKey(derBytes.Bytes)
	if err != nil {
		return err
	}
	ecdsaKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("couldn't convert: %v", pubKey)
	}
	tufPubKey := keys.Import(*ecdsaKey)
	role := args.String["--role"]
	if err = repo.ImportPubKey(tufPubKey, role); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Imported key with ID(s) %v to role %v", tufPubKey.IDs(), role)
	return nil
}
