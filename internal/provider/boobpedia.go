package provider

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Boobpedia struct{}

func (p *Boobpedia) SlugGet() string {
	return "boobpedia"
}

func (p *Boobpedia) NiceName() string {
	return "Boobpedia"
}

func (p *Boobpedia) CapabilitySearchActor() bool {
	return true
}

func (p *Boobpedia) CapabilityScrapePicture() bool {
	return true
}

func boobpediaGet(client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT x.y; Win64; x64; rv:10.0) Gecko/20100101 Firefox/10.0")
	return client.Do(req)
}

func (p *Boobpedia) ActorSearch(offlineMode bool, actorName string) (url string, err error) {
	if offlineMode {
		return url, ErrOfflineMode
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	caser := cases.Title(language.English)
	url = "https://www.boobpedia.com/boobs/" + strings.ReplaceAll(caser.String(actorName), " ", "_")
	resp, err := boobpediaGet(client, url)
	if err != nil {
		return url, err
	}

	if resp.StatusCode == 200 {
		return url, nil
	}

	return url, errors.New("provider did not find actor")
}

func (p *Boobpedia) ActorGetThumb(offlineMode bool, actorName, url string) (thumb []byte, err error) {
	if offlineMode {
		return thumb, ErrOfflineMode
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := boobpediaGet(client, url)
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
	rLegacy := regexp.MustCompile(`class=\"mw-file-description\"[\s\w="<>]* src=\"([^\"]*)`)
	thumbURLMatches := rLegacy.FindStringSubmatch(string(pageData))

	if len(thumbURLMatches) != 2 || thumbURLMatches[1] == "" {
		return thumb, errors.New("provider did not return a thumbnail")
	}

	// set found url
	url = thumbURLMatches[1]

	// retrieve thumb
	url = "https://www.boobpedia.com/" + url
	resp, err = boobpediaGet(client, url)
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
