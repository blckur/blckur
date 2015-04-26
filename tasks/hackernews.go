package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/blckur/blckur/cache"
	"github.com/blckur/blckur/database"
	"github.com/dropbox/godropbox/errors"
	"io/ioutil"
	"net/http"
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
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = &ApiError{
			errors.Wrap(err, "Hacker News parse error"),
		}
		return
	}

	stories := []int{}
	json.Unmarshal(body, &stories)

	conn := cache.Get()

	for _, storyId := range stories[:30] {
		exists, e := conn.ExistsAll(fmt.Sprintf("story-%d", storyId))
		if e != nil {
			err = e
			return
		}

		if exists {
			continue
		}

		resp, e := http.Get(fmt.Sprintf(
			"https://hacker-news.firebaseio.com/v0/item/%d.json", storyId))
		if e != nil {
			err = &ApiError{
				errors.Wrap(e, "Hacker News api error"),
			}
			return
		}

		body, e = ioutil.ReadAll(resp.Body)
		if e != nil {
			err = &ApiError{
				errors.Wrap(e, "Hacker News parse error"),
			}
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
			err = &ApiError{
				errors.Wrap(err, "Hacker News json error"),
			}
			return
		}

		if data.Type != "story" {
			continue
		}

		dataStore := struct {
			Title string `json:"title"`
			By    string `json:"by"`
			Time  int    `json:"time"`
			Text  string `json:"text"`
			Url   string `json:"url"`
		}{
			Title: data.Title,
			By:    data.By,
			Time:  data.Time,
			Text:  data.Text,
			Url:   data.Url,
		}

		dataByt, e := json.Marshal(dataStore)
		if e != nil {
			err = &ApiError{
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

		time.Sleep(1 * time.Second)
	}

	dataByt, err := json.Marshal(stories[:30])
	if err != nil {
		err = &ApiError{
			errors.Wrap(err, "Hacker News parse error"),
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
	register(-1, -1, &hackerNews{})
}
