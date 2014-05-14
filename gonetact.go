// -*- golang -*-

package main

import (
	"fmt"
	"log"
	"os/user"
	"github.com/docopt/docopt-go"
)

type CmdLineOpts struct {
	client_id  string // filename
	cache_file string // filename
	username   string // gmail email address
	query	   string // substring to match
	no_browser bool   // don't open the browser
}

func readCommandLine() *CmdLineOpts {
	usr, err := user.Current()
	rc_dir := fmt.Sprintf("%s/.gonetact", usr.HomeDir)
	if err != nil {
		log.Fatal( err )
	}
	opts := new(CmdLineOpts)

	docstring := fmt.Sprintf(`Limited interaction with google contacts

Usage:
  gonetact [--client-id=<filename>] [--cache=<cache-file>] [--query=<query>] [--no-browser]

Options:
  --client-id=<filename>    file containing a json client_id [default: %[1]s/client.json]
  --cache=<filename>        file to cache the access token[default: %[1]s/cache.json]
  --user=<gmail address>    user whose contacts will be authenticated
  --query=<query>           match this string in the email and name of the contact
  --no-browser              Don't open the auth link a browser [default: false]
  -h --help                 Show this message

The client_id is a file containing a json document per
https://code.google.com/apis/console#access`, rc_dir)

	args, _ := docopt.Parse(docstring, nil, true, "goneact 0.1", false)

	opts.client_id  = string(args["--client-id"].(string))
	opts.no_browser = bool(args["--no-browser"].(bool))
	opts.cache_file = string(args["--cache"].(string))
	if args["--user"] != nil {
		opts.username   = string(args["--user"].(string))
	}
	if args["--query"] != nil {
		opts.query      = string(args["--query"].(string))
	}
	return opts
}

func main() {
	opts := readCommandLine()

	transport := get_oauth_token(opts.client_id, opts.cache_file, opts.no_browser)

	// fmt.Println(transport)
	if opts.query == "" {
		print_all_contacts(transport)
	} else {
		print_matching_contacts(transport, opts.query)
	}
}
