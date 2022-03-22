package main

import (
	"fmt"
	"os"

	"github.com/flynn/go-docopt"
	"github.com/znewman01/go-tuf"
	"github.com/znewman01/go-tuf/pkg/keys"
)

func init() {
	register("import-sigstore", cmdImportSigstore, `
usage: tuf import-sigstore --role=<role> <issuer> <subject>

Import a sigstore ID.
`)
}

func cmdImportSigstore(args *docopt.Args, repo *tuf.Repo) error {
	sigstore := keys.SigstoreIdentity{
		Issuer:  args.String["<issuer>"],
		Subject: args.String["<subject>"],
	}
	key, err := keys.NewSigstorePublicKey(&sigstore)
	if err != nil {
		return err
	}

	role := args.String["--role"]
	if err = repo.ImportPubKey(key, role); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Imported sigstore identity to role %v\n", role)
	return nil
}
