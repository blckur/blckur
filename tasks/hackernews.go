package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/blckur/blckur/cache"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/requests"
	"github.com/dropbox/godropbox/errors"
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
	stories := []int{}
	req := requests.Request{
		Method:  requests.Get,
		Url:     "https://hacker-news.firebaseio.com/v0/newstories.json",
		Timeout: 5 * time.Second,
	}

	resp, err := req.Do()
	if err != nil {
		return
	}

	err = resp.Json(&stories)
	if err != nil {
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

		data := struct {
			Id    int    `json:"id"`
			Title string `json:"title"`
			Type  string `json:"type"`
			By    string `json:"by"`
			Time  int    `json:"time"`
			Text  string `json:"text"`
			Url   string `json:"url"`
		}{}
		req := requests.Request{
			Method: requests.Get,
			Url: fmt.Sprintf(
				"https://hacker-news.firebaseio.com/v0/item/%d.json",
				storyId),
			Timeout: 5 * time.Second,
		}

		resp, e := req.Do()
		if e != nil {
			err = e
			return
		}

		err = resp.Json(&data)
		if err != nil {
			return
		}

		if data.Type != "story" {
			continue
		}

		timestamp := time.Unix(int64(data.Time), 0)

		url, e := url.Parse(data.Url)
		if e != nil {
			err = &ParseError{
				errors.Wrap(e, "tasks.hackernews: Url error"),
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
				errors.Wrap(e, "tasks.hackernews: Json error"),
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
			errors.Wrap(err, "tasks.hackernews: Json error"),
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
