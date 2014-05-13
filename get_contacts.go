package main

import (
	"fmt"
	"errors"
	"io/ioutil"
	"encoding/json"
	"code.google.com/p/goauth2/oauth"
	"log"
)

func get_primary_email(gd_email_list []interface{}) (string, error) {
	if len(gd_email_list) < 1 {
		return "", errors.New("Less than one element in the email list")
	}
	for i := range gd_email_list {
		email_mapping := gd_email_list[i].(map[string]interface{})
		if email_mapping["primary"] == "true" {
			return email_mapping["address"].(string), nil
		}
	}
	email_mapping := gd_email_list[0].(map[string]interface{})
	email := email_mapping["address"].(string)
	return email, nil
}

func print_all_contacts(transport *oauth.Transport) {
	request_url := fmt.Sprintf("https://www.google.com/m8/feeds/contacts/default/thin?alt=json&max-results=10000")
	// fmt.Println("request_url is", request_url)
	resp, _ := transport.Client().Get(request_url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))
	data := make(map[string]interface{})
	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal("Couldn't Unmarshal the response body: ", err)
	}
	if data["feed"] == nil {
		return
	}
	feed := data["feed"].(map[string]interface{})
	if feed["entry"] == nil {
		fmt.Println("Entry was nil")
		return
	}
	entries := feed["entry"].([]interface{})
	for i := range entries {
		this_entry := entries[i].(map[string]interface{})
		// fmt.Println(this_entry)
		// fmt.Println()
		if this_entry["gd$email"] == nil {
			continue
		}
		email, err := get_primary_email(this_entry["gd$email"].([]interface{}))
		if err != nil {
			continue
		}
		title := this_entry["title"].(map[string]interface{})
		if title["$t"] != nil && title["$t"] != "" {
			name := title["$t"].(string)
			fmt.Printf("%s\t%s\t\n", email, name)
		}
	}
}
