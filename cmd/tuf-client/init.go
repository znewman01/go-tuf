package main

import (
	"io"
	"os"

	"github.com/flynn/go-docopt"
	tuf "github.com/znewman01/go-tuf/client"
)

func init() {
	register("init", cmdInit, `
usage: tuf-client init [-s|--store=<path>] --tofu <url> [<root-metadata-file>]

Options:
  -s <path>    The path to the local file store [default: tuf.db]
  --tofu       Trust on first use, probably insecure

Initialize the local file store with root metadata.
  `)
}

func cmdInit(args *docopt.Args, client *tuf.Client) error {
	var bytes []byte
	var err error

	tofu := args.Bool["--tofu"]
	if tofu {
		bytes, err = client.DownloadMetaUnsafe("root.json", tuf.DefaultRootDownloadLimit)
		if err != nil {
			return err
		}
	} else {
		file := args.String["<root-metadata-file>"]
		var in io.Reader
		if file == "" || file == "-" {
			in = os.Stdin
		} else {
			in, err = os.Open(file)
			if err != nil {
				return err
			}
		}
		bytes, err = io.ReadAll(in)
		if err != nil {
			return err
		}
	}

	return client.InitLocal(bytes)
}
