package provider

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Babepedia struct{}

func (p *Babepedia) SlugGet() string {
	return "babepedia"
}

func (p *Babepedia) NiceName() string {
	return "Babepedia"
}

func (p *Babepedia) ActorSearch(actorName string) (url string, err error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	caser := cases.Title(language.English)

	url = "https://www.babepedia.com/babe/" + strings.ReplaceAll(caser.String(actorName), " ", "_")
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

func (p *Babepedia) ActorGetThumb(actorName, url string) (thumb []byte, err error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	caser := cases.Title(language.English)
	_url := "https://www.babepedia.com/pics/" + caser.String(actorName) + ".jpg"
	req, err := http.NewRequest("GET", _url, nil)
	if err != nil {
		return thumb, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0")
	resp, err := client.Do(req)
	if err != nil {
		return thumb, err
	}

	if resp.StatusCode != 200 {
		return thumb, err
	}

	// process thumb
	thumbRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		return thumb, err
	}

	return thumbRaw, nil
}
