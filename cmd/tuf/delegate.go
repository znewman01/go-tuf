package main

import (
	// "encoding/json"

	"fmt"
	"strconv"

	"github.com/flynn/go-docopt"
	"github.com/theupdateframework/go-tuf"
)

func init() {
	register("delegate", cmdDelegate, `
usage: tuf delegate --role=<role> --name=<name> --threshold=<threshold> [<path>...]

Add a target file delegation.

Options:
  <path>                   The target path pattern(s) to delegate.
  --role=<role>            Delegate key(s) corresponding to this role.
  --name=<name>            The role name for the delegation (must be unique).
  --threshold=<threshold>  Number of signatures required for a file to be valid.
`)
}

func cmdDelegate(args *docopt.Args, repo *tuf.Repo) error {
	paths := args.All["<path>"].([]string)
	threshold, err := strconv.Atoi(args.String["--threshold"])
	if err != nil {
		return fmt.Errorf("could not parse --threshold as integer: %w", err)
	}
	return repo.Delegate(args.String["--role"], args.String["--name"], paths, threshold)
}
