package provider

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Boobpedia struct{}

func (p *Boobpedia) SlugGet() string {
	return "boobpedia"
}

func (p *Boobpedia) NiceName() string {
	return "Boobpedia"
}

func (p *Boobpedia) ActorSearch(actorName string) (url string, err error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	url = "https://www.boobpedia.com/boobs/" + strings.ReplaceAll(strings.Title(actorName), " ", "_")
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

	return url, errors.New("provider did not find actor")
}

func (p *Boobpedia) ActorGetThumb(actorName, url string) (thumb []byte, err error) {
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
	rLegacy := regexp.MustCompile("<a.*class=\"mw-file-description\">\n*\\s*<img src=\"([^\"]*)")
	thumbURLMatches := rLegacy.FindStringSubmatch(string(pageData))

	if len(thumbURLMatches) != 2 || thumbURLMatches[1] == "" {
		return thumb, errors.New("provider did not return a thumbnail")
	}

	// set found url
	url = thumbURLMatches[1]

	//TODO: add secondary search with new profiles
	// for reference, regex is: <div class="thumbImage">\n\s*<img src="([^"]*)
	// note: regex is most likely multilined

	// retrieve thumb
	req, err = http.NewRequest("GET", "https://www.boobpedia.com/"+url, nil)
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
