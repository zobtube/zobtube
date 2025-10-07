package provider

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type BabesDirectory struct{}

func (p *BabesDirectory) SlugGet() string {
	return "babesdirectory"
}

func (p *BabesDirectory) NiceName() string {
	return "Babes Directory"
}

func (p *BabesDirectory) CapabilitySearchActor() bool {
	return true
}

func (p *BabesDirectory) CapabilityScrapePicture() bool {
	return true
}

func babesDirectoryGet(client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT x.y; Win64; x64; rv:10.0) Gecko/20100101 Firefox/10.0")
	return client.Do(req)
}

func (p *BabesDirectory) ActorSearch(offlineMode bool, actorName string) (url string, err error) {
	if offlineMode {
		return url, ErrOfflineMode
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	url = "https://babesdirectory.online/profile/" + strings.ReplaceAll(strings.ToLower(actorName), " ", "-")

	resp, err := babesDirectoryGet(client, url)
	if err != nil {
		return url, err
	}

	if resp.StatusCode == 200 {
		return url, nil
	}

	url += "-pornstar"
	resp, err = babesDirectoryGet(client, url)
	if err != nil {
		return url, err
	}

	if resp.StatusCode == 200 {
		return url, nil
	}

	return url, errors.New("provider did not find actor")
}

func (p *BabesDirectory) ActorGetThumb(offlineMode bool, actorName, url string) (thumb []byte, err error) {
	if offlineMode {
		return thumb, ErrOfflineMode
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := babesDirectoryGet(client, url)
	if err != nil {
		return thumb, err
	}

	if resp.StatusCode != 200 {
		return thumb, errors.New("unable to query provider")
	}

	// process page
	pageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return thumb, err
	}

	// try default avatar position of legacy profiles
	re := regexp.MustCompile("<div class=\"pill-image\">\\n*\\s*<img src=\"([^\"]*)")
	thumbURLMatches := re.FindStringSubmatch(string(pageData))

	if len(thumbURLMatches) != 2 || thumbURLMatches[1] == "" {
		return thumb, errors.New("provider did not return a thumbnail")
	}

	// set found url
	url = thumbURLMatches[1]

	// retrieve thumb
	resp, err = babesDirectoryGet(client, url)
	if err != nil {
		return thumb, err
	}

	if resp.StatusCode != 200 {
		return thumb, errors.New("provider thumb retrieval failed")
	}

	// process thumb
	thumbRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		return thumb, err
	}

	return thumbRaw, nil
}
