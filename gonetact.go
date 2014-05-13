// -*- golang -*-

// The idea is to provide a daemon that will run alongside my
// mutt instance, and talk to google to get contacts from my google account.
// The interface to it will be command-line based, similar to
// aboot or something similar to that.
//
// The workflow will be:
// * execute
// ** if there is a token that exists on disk already
//    in ~/.gonetact-tokencache, use it.
// *** Connect to google's contacts and query for contacts
// *** Provide commands to
// **** Get all contacts
// **** add a contact (deletes will be up to the web UI
// ** if not then:
// *** determine the hostname via either:
// **** .gonetactrc
// **** failing the config file, hostname := os.Hostname()
// *** put an http server on a listening socket
// *** display a URL to go to which will then ask you to auth to google.
// *** get the oauth token
// *** save it to ~/.gonetact-token

package main

import (
	"fmt"
	"log"
	"os/user"
	"github.com/docopt/docopt-go"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}

	rc_dir := fmt.Sprintf("%s/.gonetact", usr.HomeDir)
	docstring := fmt.Sprintf(`Limited interaction with google contacts

Usage:
  gonetact [-o] [--client-id=<filename>] [--cache=<cache-file>] [--query=<query>]

Options:
  -o                        use oauth2 [default: true]
  --client-id=<filename>    file containing a json client_id [default: %[1]s/client.json]
  --cache=<filename>        file to cache the access token[default: %[1]s/cache.json]
  --user=<gmail address>    user whose contacts will be authenticated
  --query=<query>           a search term, currently doesn't filter anything
  -h --help                 Show this message

The client_id is a file containing a json document per
https://code.google.com/apis/console#access`, rc_dir)

	args, _ := docopt.Parse(docstring, nil, true, "goneact 0.1", false)

	// fmt.Println(args)
	transport := get_oauth_token(string(args["--client-id"].(string)), string(args["--cache"].(string)))
	// fmt.Println(transport)
	print_all_contacts(transport)
}
