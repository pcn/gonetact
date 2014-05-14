package main

import (
	"testing"
	"os/user"
	"fmt"
)

func Test_readCommandLine_1(t *testing.T) {
	// Test the no_browser flag
	usr, _ := user.Current()
	rc_dir := fmt.Sprintf("%s/.gonetact", usr.HomeDir)
	client_id  := fmt.Sprintf("%s/client.json", rc_dir)
	cache_file := fmt.Sprintf("%s/cache.json", rc_dir)

	input_argv := []string{"gonetacts", "--no-browser"}
	expected_output := CmdLineOpts {
		client_id : client_id,
		cache_file : cache_file,
		no_browser : true,
		username: "",
		query: "",
	}
	parsed := readCommandLine(input_argv)
	if expected_output != *parsed {
		fmt.Println(expected_output)
		fmt.Println(*parsed)
		t.Error("The returned flags do not look like the flags we specified.")
	}
}

func Test_readCommandLine_2(t *testing.T) {
	// Test the query argument
	usr, _ := user.Current()
	rc_dir := fmt.Sprintf("%s/.gonetact", usr.HomeDir)
	client_id  := fmt.Sprintf("%s/client.json", rc_dir)
	cache_file := fmt.Sprintf("%s/cache.json", rc_dir)

	input_argv := []string{"gonetacts", "--query=foo"}
	expected_output := CmdLineOpts {
		client_id : client_id,
		cache_file : cache_file,
		no_browser : false,
		username: "",
		query: "foo",
	}
	parsed := readCommandLine(input_argv)
	if expected_output != *parsed {
		fmt.Println(expected_output)
		fmt.Println(*parsed)
		t.Error("The returned flags do not look like the flags we specified.")
	}
}

func Test_readCommandLine_3(t *testing.T) {
	// Test the client-id argument
	usr, _ := user.Current()
	rc_dir := fmt.Sprintf("%s/.gonetact", usr.HomeDir)
	client_id  := "/something/foo/client_id.json"
	cache_file := fmt.Sprintf("%s/cache.json", rc_dir)

	input_argv := []string{"gonetacts", fmt.Sprintf("--client-id=%s", client_id)}
	expected_output := CmdLineOpts {
		client_id : client_id,
		cache_file : cache_file,
		no_browser : false,
		username: "",
		query: "",
	}
	parsed := readCommandLine(input_argv)
	if expected_output != *parsed {
		fmt.Println(expected_output)
		fmt.Println(*parsed)
		t.Error("The returned flags do not look like the flags we specified.")
	}
}

func Test_readCommandLine_4(t *testing.T) {
	// Test the cache-file argument
	usr, _ := user.Current()
	rc_dir := fmt.Sprintf("%s/.gonetact", usr.HomeDir)
	client_id  := fmt.Sprintf("%s/client.json", rc_dir)
	cache_file := "/something/foo/client_id.json"

	input_argv := []string{"gonetacts", fmt.Sprintf("--cache=%s", cache_file)}
	expected_output := CmdLineOpts {
		client_id : client_id,
		cache_file : cache_file,
		no_browser : false,
		username: "",
		query: "",
	}
	parsed := readCommandLine(input_argv)
	if expected_output != *parsed {
		fmt.Println(expected_output)
		fmt.Println(*parsed)
		t.Error("The cache-file argument didn't come through per the args we specified.")
	}
}
