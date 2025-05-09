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

func (p *BabesDirectory) ActorSearch(actorName string) (url string, err error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	url = "https://babesdirectory.online/profile/" + strings.ReplaceAll(strings.ToLower(actorName), " ", "-")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return url, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT x.y; Win64; x64; rv:10.0) Gecko/20100101 Firefox/10.0")
	resp, err := client.Do(req)
	if err != nil {
		return url, err
	}

	if resp.StatusCode == 200 {
		return url, nil
	}

	url = url + "-pornstar"
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return url, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT x.y; Win64; x64; rv:10.0) Gecko/20100101 Firefox/10.0")
	resp, err = client.Do(req)
	if err != nil {
		return url, err
	}

	if resp.StatusCode == 200 {
		return url, nil
	}

	return url, errors.New("provider did not find actor")
}

func (p *BabesDirectory) ActorGetThumb(actorName, url string) (thumb []byte, err error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return thumb, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT x.y; Win64; x64; rv:10.0) Gecko/20100101 Firefox/10.0")
	resp, err := client.Do(req)
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
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return thumb, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT x.y; Win64; x64; rv:10.0) Gecko/20100101 Firefox/10.0")
	resp, err = client.Do(req)
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
