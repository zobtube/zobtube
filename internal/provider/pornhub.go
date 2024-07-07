package provider

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Pornhub struct{}

func (p *Pornhub) SlugGet() string {
	return "pornhub"
}

func (p *Pornhub) NiceName() string {
	return "PornHub"
}

func (p *Pornhub) ActorSearch(actorName string) (url string, err error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	url = "https://www.pornhub.com/pornstar/" + strings.ReplaceAll(strings.ToLower(actorName), " ", "-")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return url, err
	}

	req.Header.Add("Cookie", "accessAgeDisclaimerPH=1")
	resp, err := client.Do(req)
	if err != nil {
		return url, err
	}

	if resp.StatusCode == 200 {
		return url, nil
	}

	url = "https://www.pornhub.com/model/" + strings.ReplaceAll(strings.ToLower(actorName), " ", "-")
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return url, err
	}

	req.Header.Add("Cookie", "accessAgeDisclaimerPH=1")
	resp, err = client.Do(req)
	if err != nil {
		return url, err
	}

	if resp.StatusCode == 200 {
		return url, nil
	}

	return url, errors.New("provider did not find actor")
}

func (p *Pornhub) ActorGetThumb(actor_name, url string) (thumb []byte, err error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return thumb, err
	}

	req.Header.Add("Cookie", "accessAgeDisclaimerPH=1")
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
	rLegacy := regexp.MustCompile("<img id=\"getAvatar\".*src=\"([^\"]*)")
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
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return thumb, err
	}

	req.Header.Add("Cookie", "accessAgeDisclaimerPH=1")
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
