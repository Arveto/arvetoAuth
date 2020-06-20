package github

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

var emailNotFound = errors.New("Email not found")

type Emails []struct {
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
	Primary  bool   `json:"primary"`
}

// Get the Emails.
func getMail(token string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "token "+token)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var list Emails
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
