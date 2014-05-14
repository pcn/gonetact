package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
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

// Returns a list of email_entry structs to be attached to an contact_entry
func get_email(gd_email_list []interface{}) ([]email_entry, error) {
	var email_entry_list []email_entry
	if len(gd_email_list) < 1 {
		return email_entry_list, errors.New("Less than one element in the email list")
	}
	for i := range gd_email_list {
		email_mapping := gd_email_list[i].(map[string]interface{})
		for k, v := range(email_mapping) {
			value := v.(string)
			if k == "rel" {
				idx := strings.Index(value, "#")
				address_association := value[idx+1:]
				email_entry_list = append(email_entry_list, email_entry { email_mapping["address"].(string), address_association})
			}
		}
	}
	return email_entry_list, nil
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
		if len(v.Emails) != 0 && v.Name != "" {
			*list = append(*list, v)
		}
	}

	return nil
}

type contact_entry struct {
	Emails []email_entry
	Name string
}

type email_entry struct {
	Address string
	Association string
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

	emails, err := get_email(gdEmail.([]interface{}))
	if err != nil {
		return err
	}

	title := data["title"].(map[string]interface{})["$t"]
	if title != nil && title != "" {
		ce.Name = title.(string)
	}

	ce.Emails = emails
	return nil
}

// This function does the work of obtaining the contacts from the server
func fetch_all_contacts(transport *oauth.Transport) contacts_response {
	// XXX: increase the max-results
	request_url := fmt.Sprintf("https://www.google.com/m8/feeds/contacts/default/thin?alt=json&max-results=10000")
	// fmt.Println("request_url is", request_url)
	resp, _ := transport.Client().Get(request_url)
	defer resp.Body.Close()

	var result contacts_response
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

// This function turns the response from fetch_all_contacts into a list of strings
func all_contacts(transport *oauth.Transport) []string {
	result := fetch_all_contacts(transport)
	var all_contacts []string
	for _, v := range result.Feed.Entries {
		for i := range v.Emails {
			all_contacts = append(all_contacts, fmt.Sprintf("%s\t%s\t%s", v.Emails[i].Address, v.Name, v.Emails[i].Association))
		}
	}
	return all_contacts
}

func print_all_contacts(transport *oauth.Transport) {
	all := all_contacts(transport)
	for line := range all {
		fmt.Printf("%s\n", all[line])
	}
}

func print_matching_contacts(transport *oauth.Transport, query_str string) {
	all := all_contacts(transport)
	for line := range all {
		if strings.Contains(strings.ToLower(all[line]), strings.ToLower(query_str)) {
			fmt.Printf("%s\n", all[line])
		}
	}
}
