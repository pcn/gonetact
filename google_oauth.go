package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"code.google.com/p/goauth2/oauth"
)

var requestURL = "https://www.googleapis.com/oauth2/v1/contacts"

// To obtain a request token you must specify both -id and -secret.
//
// To obtain Client ID and Secret, see the "OAuth 2 Credentials" section under
// the "API Access" tab on this page: https://code.google.com/apis/console/
//
// Once you have completed the OAuth flow, the credentials should be stored inside
// the file specified by -cache and you may run without the -id and -secret flags.

// Accepts a filename - this is the json file that contains the native client
// client id, which is the obvious way I've seen so far of enabling this.
// the file looks like this:
// {
//     "installed": {
//         "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
//         "auth_uri": "https://accounts.google.com/o/oauth2/auth",
//         "client_email": "",
//         "client_id": "824874312315-ab5qku68sd6cfmh0cdcdhdjqu9m76tfo.apps.googleusercontent.com",
//         "client_secret": "QIMOeESlrv25CUEIGCq_69qt",
//         "client_x509_cert_url": "",
//         "redirect_uris": [
//             "urn:ietf:wg:oauth:2.0:oob",
//             "oob"
//         ],
//         "token_uri": "https://accounts.google.com/o/oauth2/token"
//     }
// }
//
// TODO: unit test

type client_json struct {
	Installed client_info `json:"installed"`
}

type client_info struct {
	Id                      string   `json:"client_id"`
	Secret                  string   `json:"client_secret"`
	Email                   string   `json:"client_email"`
	AuthUri                 string   `json:"auth_uri"`
	TokenUri                string   `json:"token_uri"`
	RedirectUris            []string `json:"redirect_uris"`
	ClientX509CertUrl       string   `json:"client_x509_cert_url"`
	AuthProviderX509CertUrl string   `json:"auth_provider_x509_cert_url"`

	Code string `json:"code"` // Not present?
}

func get_native_app_client_id(filename string) (*client_info, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var c client_json
	if err = json.NewDecoder(f).Decode(&c); err != nil {
		return nil, err
	}

	return &c.Installed, nil
}

// Internal function to get a token, using vars defined in get_oauth_token.
// I need better terminology.
func get_token(client *client_info, cachefile_name string) (*oauth.Transport, *oauth.Token, *oauth.Config, error) {
	var (
		scope       = "https://www.google.com/m8/feeds"
		authURL     = "https://accounts.google.com/o/oauth2/auth"
		redirectURL = "urn:ietf:wg:oauth:2.0:oob"
		tokenURL    = "https://accounts.google.com/o/oauth2/token"
		cachefile   = cachefile_name
	)
	config := &oauth.Config{
		ClientId:     client.Id,
		ClientSecret: client.Secret,
		RedirectURL:  redirectURL,
		Scope:        scope,
		AuthURL:      authURL,
		TokenURL:     tokenURL,
		TokenCache:   oauth.CacheFile(cachefile),
	}
	// Set up a Transport using the config.
	transport := &oauth.Transport{Config: config}

	// Try to pull the token from the cache; if this fails, we need to get one.
	token, err := config.TokenCache.Token()
	if err != nil {
		if client.Id == "" || client.Secret == "" {
			log.Printf("Error in obtaining a token: %s\n", err)
			log.Fatal("cachefile is:  %s\n", cachefile)
		}
		log.Println("Err is not nil")
		log.Println("token is ", token)
		log.Println("transport is ", transport)
		log.Println("config is ", config)
	}
	return transport, token, config, err
}

func get_oauth_token(filename string, cachefile_name string) *oauth.Transport {
	client, err := get_native_app_client_id(filename)
	if err != nil {
		return nil
	}

	transport, token, config, err := get_token(client, cachefile_name)
	if err != nil {
		if client.Code == "" {
			// Get an authorization code from the data provider, then continue
			// ("Please ask the user if I can access this resource.")
			url := config.AuthCodeURL("")
			fmt.Println("Visit this URL to get a code, then paste the code here\n")
			fmt.Println(url)
			bio := bufio.NewReader(os.Stdin)
			line, _, _ := bio.ReadLine() // TODO: check err.  Don't worry about hasmoreinline
			code := string(line)
			// Exchange the authorization code for an access token.
			// ("Here's the code you gave the user, now give me a token!")
			token, err = transport.Exchange(code)
			if err != nil {
				log.Fatal("Exchange:", err)
			}
			// (The Exchange method will automatically cache the token.)
			fmt.Printf("Token is cached in %v\n", config.TokenCache)

			transport, token, config, err = get_token(client, cachefile_name)
			if err != nil {
				log.Fatal("Trying to get the token with a passed in code: ", err)
			}
		} else {
			log.Fatal(err)
		}
	}

	// Make the actual request using the cached token to authenticate.
	// ("Here's the token, let me in!")
	transport.Token = token
	return transport
}
