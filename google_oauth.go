package main

import (
	"fmt"
	//	"io"
	"bufio"
	"io/ioutil"
	"log"
	"os"
	// "strings"
	"encoding/json"
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
func get_native_app_client_id(client_id_file string) map[string]interface{} {
	data := make(map[string]interface{})
	contents := []byte{}
	f, err := os.Open(client_id_file)
	if err != nil {
		log.Print("While trying to open the application client_id: %s", client_id_file)
		log.Fatal(err)
	}
	contents, err = ioutil.ReadAll(f)
	if err != nil {
		log.Print("While reader data from the file %s: ", client_id_file)
		log.Fatal(err)
	}
	err = json.Unmarshal(contents, &data)
	if err != nil {
		log.Print("While trying to Unmarshall json from the application client_id file %s: ", client_id_file)
		log.Fatal(err)
	}
	// TODO: add json schema checking
	return data
}

// Internal function to get a token, using vars defined in get_oauth_token.
// I need better terminology.
func get_token(installed map[string]interface{}, cachefile_name string, code string) (*oauth.Transport, *oauth.Token, *oauth.Config, error) {
	clientId, _     := installed["client_id"].(string)
	clientSecret, _ := installed["client_secret"].(string)

	var (
		scope      = "https://www.google.com/m8/feeds"
		authURL    = "https://accounts.google.com/o/oauth2/auth"
		redirectURL= "urn:ietf:wg:oauth:2.0:oob"
		tokenURL   = "https://accounts.google.com/o/oauth2/token"
		cachefile  = cachefile_name
	)
        config := &oauth.Config{
                ClientId    :     clientId,
                ClientSecret: clientSecret,
                RedirectURL :  redirectURL,
                Scope	    :        scope,
                AuthURL	    :      authURL,
                TokenURL    :     tokenURL,
                TokenCache  :   oauth.CacheFile(cachefile),
        }
        // Set up a Transport using the config.
        transport := &oauth.Transport{Config: config}

        // Try to pull the token from the cache; if this fails, we need to get one.
        token, err := config.TokenCache.Token()
        if err != nil {
                if clientId == "" || clientSecret == "" {
                        fmt.Fprint(os.Stderr, "Error in obtaining a token: ", err )
			fmt.Print("\n")
			fmt.Fprint(os.Stderr, "cachefile is: ", cachefile, "\n")
                        os.Exit(2)
                }
		fmt.Println("Err is not nil")
		fmt.Println("token is ", token)
		fmt.Println("transport is ", transport)
		fmt.Println("config is ", config)
	}
	return transport, token, config, err
}

func get_oauth_token(client_json string, cachefile_name string) *oauth.Transport {
	client_info := get_native_app_client_id(client_json)
	// fmt.Println(client_info)

	installed := client_info["installed"].(map[string]interface{})
	code, _   := installed["code"].(string)

	transport, token, config, err := get_token(installed, cachefile_name, code)
        if err != nil {
                if code == "" {
                        // Get an authorization code from the data provider, then continue
                        // ("Please ask the user if I can access this resource.")
                        url := config.AuthCodeURL("")
                        fmt.Println("Visit this URL to get a code, then paste the code here\n")
                        fmt.Println(url)
			bio := bufio.NewReader(os.Stdin)
			line, _, _ := bio.ReadLine() // TODO: check err.  Don't worry about hasmoreinline
			code = string(line)
			// Exchange the authorization code for an access token.
			// ("Here's the code you gave the user, now give me a token!")
			token, err = transport.Exchange(code)
			if err != nil {
				log.Fatal("Exchange:", err)
			}
			// (The Exchange method will automatically cache the token.)
			fmt.Printf("Token is cached in %v\n", config.TokenCache)

			transport, token, config, err = get_token(installed, cachefile_name, code)
			if err != nil {
				log.Fatal("Trying to get the token with a passed in code: ", err)
			}
                }
        }

        // Make the actual request using the cached token to authenticate.
        // ("Here's the token, let me in!")
        transport.Token = token
	return transport
}
