package main

import (
	"fmt"

	"github.com/flynn/go-docopt"
	tuf "github.com/theupdateframework/go-tuf/client"
	"github.com/theupdateframework/go-tuf/verify"
	"time"
)

func init() {
	register("who-can-sign", cmdWhoCanSign, `
usage: tuf-client who-can-sign [-s|--store=<path>] <url> <target>

Options:
  -s <path>    The path to the local file store [default: tuf.db]

Who can sign the given target.
  `)
}

func cmdWhoCanSign(args *docopt.Args, client *tuf.Client) error {
	verify.IsExpired = func(t time.Time) bool { return false }
	if _, err := client.Update(); err != nil && !tuf.IsLatestSnapshot(err) {
		return err
	}
	people, err := client.PrintPeople(args.String["<target>"])
	if err != nil {
		return err
	}
	fmt.Println("people:", people)
	return nil
}
