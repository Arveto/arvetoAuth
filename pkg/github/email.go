package github

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

var emailNotFound = errors.New("Email not found")

type emails []struct {
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
	Primary  bool   `json:"primary"`
}

// Get the Emails.
func getMail(client *http.Client) (string, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var list emails
	err = json.Unmarshal(data, &list)

	if len(list) == 0 {
		return "", emailNotFound
	}
	email := list[0].Email
	for _, e := range list {
		if e.Primary {
			email = e.Email
			continue
		}
	}

	return email, err
}
