package main

import (
	"fmt"
	"time"

	"github.com/flynn/go-docopt"
	"github.com/theupdateframework/go-tuf"
)

func init() {
	register("gen-key", cmdGenKey, `
usage: tuf gen-key [--expires=<days>] [--delegated] <role>

Generate a new signing key for the given role.

The key will be serialized to JSON and written to the "keys" directory with
filename pattern "ROLE-KEYID.json". If "--delegated" is not given, the root
metadata file will also be staged with the addition of the key's ID to the
role's list of key IDs.

Alternatively, passphrases can be set via environment variables in the
form of TUF_{{ROLE}}_PASSPHRASE

Options:
  --expires=<days>   Set the root metadata file to expire <days> days from now.
  --delegated        This role is delegated; do not add it to the root metadata file.
`)
}

func cmdGenKey(args *docopt.Args, repo *tuf.Repo) error {
	role := args.String["<role>"]
	var keyids []string
	var err error
	opts := tuf.NewGenKeyOpts().SetDelegated(args.Bool["--delegated"])
	if arg := args.String["--expires"]; arg != "" {
		var expires time.Time
		expires, err = parseExpires(arg)
		if err != nil {
			return err
		}
		opts.SetExpires(&expires)
	}
	keyids, err = repo.GenKeyWithOpts(role, *opts)
	if err != nil {
		return err
	}
	for _, id := range keyids {
		fmt.Println("Generated", role, "key with ID", id)
	}
	return nil
}
