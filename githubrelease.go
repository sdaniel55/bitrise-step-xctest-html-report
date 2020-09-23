package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/log"
)

// GithubRelease returned from the Github API
type GithubRelease struct {
	TagName string `json:"tag_name"`
}

// https://api.github.com/repos/TitouanVanBelle/XCTestHTMLReport/releases/latest

// Returns the information of the latest release from the Github API for the provided repository
// More: https://docs.github.com/en/rest/reference/repos#releases
func latestGithubRelease(githubOrg, githubRepository string, accessToken stepconf.Secret) (GithubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", githubOrg, githubRepository)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return GithubRelease{}, fmt.Errorf("failed to create new request for %s, error: %v", url, err)
	}

	if string(accessToken) != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", string(accessToken)))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return GithubRelease{}, fmt.Errorf("failed to call the %s, error: %v", url, err)
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Warnf("Failed to close response body, error: %v", cerr)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GithubRelease{}, fmt.Errorf("failed to read response body from the %s, error: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return GithubRelease{}, fmt.Errorf("response status %v, body: %s", resp.StatusCode, string(body))
	}
	log.Debugf("Response status: %s", resp.Status)

	var githubRelease GithubRelease
	if err := json.Unmarshal([]byte(body), &githubRelease); err != nil {
		return githubRelease, fmt.Errorf("failed to parse latest Github Release JSON, ewrror: %v", err)
	}

	return githubRelease, nil
}
