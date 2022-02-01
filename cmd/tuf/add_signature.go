package main

import (
	"os"

	"github.com/flynn/go-docopt"
	"github.com/theupdateframework/go-tuf"
	"github.com/theupdateframework/go-tuf/data"
)

func init() {
	register("add-signature", cmdAddSignature, `
usage: tuf add-signature <role> --key-id <key_id> --signature <sig_file>

Adds a signature (as hex-encoded bytes) generated by an offline tool to the given role.

If the signature does not verify, it will not be added.
`)
}

func cmdAddSignature(args *docopt.Args, repo *tuf.Repo) error {
	role := args.String["<role>"]
	keyID := args.String["<key_id>"]

	f := args.String["<sig_file>"]
	sigBytes, err := os.ReadFile(f)
	if err != nil {
		return err
	}
	sigData := data.HexBytes(sigBytes)

	sig := data.Signature{
		KeyID:     keyID,
		Signature: sigData,
	}
	return repo.AddOrUpdateSignature(role, sig)
}
