package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/blckur/blckur/cache"
	"github.com/blckur/blckur/database"
	"github.com/dropbox/godropbox/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type hackerNews struct {
}

func (h *hackerNews) Type() string {
	return "hacker_news"
}

func (h *hackerNews) Run(db *database.Database) (err error) {
	resp, err := http.Get(
		"https://hacker-news.firebaseio.com/v0/newstories.json")
	if err != nil {
		err = &ApiError{
			errors.Wrap(err, "Hacker News api error"),
		}
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = &ParseError{
			errors.Wrap(err, "Hacker News io error"),
		}
		return
	}

	stories := []int{}

	err = json.Unmarshal(body, &stories)
	if err != nil {
		err = &ParseError{
			errors.Wrap(err, "Hacker News json error"),
		}
		return
	}

	if len(stories) > 10 {
		stories = stories[:10]
	}

	conn := cache.Get()

	for _, storyId := range stories {
		exists, e := conn.ExistsAll(fmt.Sprintf("story-%d", storyId))
		if e != nil {
			err = e
			return
		}

		if exists {
			continue
		}

		var body []byte

		for i := 0; i < 3; i++ {
			resp, e := http.Get(fmt.Sprintf(
				"https://hacker-news.firebaseio.com/v0/item/%d.json", storyId))
			if e != nil {
				err = &ApiError{
					errors.Wrap(e, "Hacker News api error"),
				}
				time.Sleep(3 * time.Second)
				continue
			}

			if resp.StatusCode != 200 {
				err = &ApiError{
					errors.Newf("Hacker News bad status %d", resp.StatusCode),
				}
				time.Sleep(3 * time.Second)
				continue
			}

			body, e = ioutil.ReadAll(resp.Body)
			if e != nil {
				err = &ParseError{
					errors.Wrap(e, "Hacker News io error"),
				}
				time.Sleep(3 * time.Second)
				continue
			}

			err = nil
			break
		}
		if err != nil {
			return
		}

		data := struct {
			Id    int    `json:"id"`
			Title string `json:"title"`
			Type  string `json:"type"`
			By    string `json:"by"`
			Time  int    `json:"time"`
			Text  string `json:"text"`
			Url   string `json:"url"`
		}{}

		err = json.Unmarshal(body, &data)
		if err != nil {
			err = &ParseError{
				errors.Wrap(err, "Hacker News json error"),
			}
			return
		}

		if data.Type != "story" {
			continue
		}

		timestamp := time.Unix(int64(data.Time), 0)

		url, e := url.Parse(data.Url)
		if e != nil {
			err = &ParseError{
				errors.Wrap(e, "Hacker News url error"),
			}
			return
		}

		storyUrl := fmt.Sprintf("https://news.ycombinator.com/item?id=%d",
			data.Id)

		dataStore := struct {
			Id     int       `json:"id"`
			Title  string    `json:"title"`
			By     string    `json:"by"`
			Time   time.Time `json:"time"`
			Text   string    `json:"text"`
			Url    string    `json:"url"`
			Domain string    `json:"domain"`
		}{
			Id:     data.Id,
			Title:  data.Title,
			By:     data.By,
			Time:   timestamp,
			Text:   data.Text,
			Url:    storyUrl,
			Domain: strings.Replace(url.Host, "www.", "", 1),
		}

		dataByt, e := json.Marshal(dataStore)
		if e != nil {
			err = &ParseError{
				errors.Wrap(e, "Hacker News json error"),
			}
			return
		}

		e = conn.SetString(fmt.Sprintf("story-%d", data.Id), string(dataByt),
			6*time.Hour)
		if e != nil {
			err = e
			return
		}

		err = conn.Publish("hacker_news", "story", string(dataByt))
		if err != nil {
			return
		}

		time.Sleep(1 * time.Second)
	}

	dataByt, err := json.Marshal(stories)
	if err != nil {
		err = &ParseError{
			errors.Wrap(err, "Hacker News json error"),
		}
		return
	}

	err = conn.SetString("new_stories", string(dataByt), 6*time.Hour)
	if err != nil {
		return
	}

	return
}

func init() {
	registerMulti([]int{-1}, []int{-1}, []int{
		3, 13, 23, 33, 43, 53}, &hackerNews{})
}
