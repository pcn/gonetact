package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"code.google.com/p/goauth2/oauth"
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

type contacts_response struct {
	Feed struct {
		Entries entry_list `json:"entry"`
	} `json:"feed"`
}

type entry_list []contact_entry

func (list *entry_list) UnmarshalJSON(b []byte) error {
	// Unmarshal all entries
	var entries []contact_entry
	err := json.Unmarshal(b, &entries)
	if err != nil {
		return err
	}

	// Keep only the entries with email
	for _, v := range entries {
		if v.Email != "" && v.Name != "" {
			*list = append(*list, v)
		}
	}

	return nil
}

type contact_entry struct {
	Email, Name string
}

func (ce *contact_entry) UnmarshalJSON(b []byte) error {
	data := make(map[string]interface{})

	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	gdEmail := data["gd$email"]
	if gdEmail == nil {
		ce = nil
		return nil
	}

	email, err := get_primary_email(gdEmail.([]interface{}))
	if err != nil {
		return err
	}

	title := data["title"].(map[string]interface{})["$t"]
	if title != nil && title != "" {
		ce.Name = title.(string)
	}

	ce.Email = email
	return nil
}

func print_all_contacts(transport *oauth.Transport) {
	request_url := fmt.Sprintf("https://www.google.com/m8/feeds/contacts/default/thin?alt=json&max-results=10000")
	// fmt.Println("request_url is", request_url)
	resp, _ := transport.Client().Get(request_url)
	defer resp.Body.Close()

	var result contacts_response
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range result.Feed.Entries {
		fmt.Printf("%s\t%s\t\n", v.Email, v.Name)
	}
}
