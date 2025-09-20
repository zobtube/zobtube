package provider

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type IAFD struct {
	client *http.Client
}

func (p *IAFD) SlugGet() string {
	return "iafd"
}

func (p *IAFD) NiceName() string {
	return "IAFD"
}

func (p *IAFD) CapabilitySearchActor() bool {
	return true
}

func (p *IAFD) CapabilityScrapePicture() bool {
	return true
}

func (p *IAFD) IAFDGet(url string) (*http.Response, error) {
	if p.client == nil {
		p.client = &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				DisableCompression: true,
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS13,
				},
			},
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:141.0) Gecko/20100101 Firefox/141.0")
	req.Header.Add("Accept", "*/*")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("provider did not find actor, got status code %d", resp.StatusCode)
	}

	return resp, nil
}

func (p *IAFD) ActorSearch(offlineMode bool, actorName string) (url string, err error) {
	if offlineMode {
		return url, ErrOfflineMode
	}

	// search actor
	caser := cases.Title(language.English)
	url = "https://www.iafd.com/results.asp?searchtype=comprehensive&searchstring=" + strings.ReplaceAll(caser.String(actorName), " ", "+")

	resp, err := p.IAFDGet(url)
	if err != nil {
		return url, err
	}

	// process page
	pageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return url, err
	}

	re := regexp.MustCompile(`<tr><td><a href="(\/person.rme\/[^"]*)`)
	matches := re.FindAllStringSubmatch(string(pageData), -1)

	if len(matches) == 0 {
		return url, errors.New("provider did not find actor")
	}

	if len(matches) > 1 {
		return url, errors.New("provider matches more than one actor")
	}

	return "https://www.iafd.com" + matches[0][1], nil
}

func (p *IAFD) ActorGetThumb(offlineMode bool, actorName, url string) (thumb []byte, err error) {
	if offlineMode {
		return thumb, ErrOfflineMode
	}

	resp, err := p.IAFDGet(url)
	if err != nil {
		return thumb, err
	}

	// process page
	pageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return thumb, err
	}

	divRegex := regexp.MustCompile(`<div id="headshot"><img(.*src="([^\"]*))`)
	thumbURLMatches := divRegex.FindAllStringSubmatch(string(pageData), -1)

	url = ""
	for _, match := range thumbURLMatches {
		if len(match) > 2 {
			url = match[2]
		}
	}

	if url == "" {
		// definitely not found
		return thumb, errors.New("provider did not return a thumbnail")
	}

	// retrieve thumb
	resp, err = p.IAFDGet(url)
	if err != nil {
		return thumb, err
	}

	// process thumb
	thumbRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		return thumb, err
	}

	return thumbRaw, nil
}
