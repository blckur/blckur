package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/cache"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/requests"
	"github.com/blckur/blckur/stream"
	"github.com/dropbox/godropbox/errors"
	"labix.org/v2/mgo/bson"
	"net/url"
	"strings"
	"time"
)

var (
	hackerNewsAcctId = bson.ObjectIdHex("5542b6ceb0e7304e38000001")
	hackerNewsUserId = bson.ObjectIdHex("5542b6ceb0e7304e38000002")
)

type hackerNewsBackend struct {
	conn    *cache.ClusterConn
	stop    bool
	watcher *requests.Watcher
}

func (b *hackerNewsBackend) parseStory(storyId int) (err error) {
	exists, err := b.conn.ExistsAll(fmt.Sprintf("story-%d", storyId))
	if exists || err != nil {
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
	req := requests.Request{
		Method: requests.Get,
		Url: fmt.Sprintf(
			"https://hacker-news.firebaseio.com/v0/item/%d.json", storyId),
		Timeout: 5 * time.Second,
	}

	resp, err := req.Do()
	if err != nil {
		return
	}

	err = resp.Json(&data)
	if err != nil {
		return
	}

	if data.Type != "story" {
		return
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

	dataByt, err := json.Marshal(dataStore)
	if err != nil {
		err = &ParseError{
			errors.Wrap(err, "tasks.hackernews: Json error"),
		}
		return
	}

	e = b.conn.SetString(fmt.Sprintf("story-%d", data.Id), string(dataByt),
		6*time.Hour)
	if e != nil {
		err = e
		return
	}

	err = b.conn.Publish("hacker_news", "story", string(dataByt))
	if err != nil {
		return
	}

	return
}

func (b *hackerNewsBackend) parse(evt *requests.Event) (err error) {
	if evt.Type != "put" {
		return
	}

	stories := struct {
		Data []int `json:"data"`
	}{}
	err = json.Unmarshal([]byte(evt.Data), &stories)
	if err != nil {
		err = &ParseError{
			errors.Wrap(err, "task.hackernews: Failed to parse event data"),
		}
		return
	}

	if len(stories.Data) > 10 {
		stories.Data = stories.Data[:10]
	}

	for _, storyId := range stories.Data {
		err = b.parseStory(storyId)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"story_id": storyId,
				"error":    err,
			}).Error("tasks.hackernews: Failed to parse story")
			err = nil
		}
	}

	return
}

func (b *hackerNewsBackend) Run() {
	b.conn = cache.Get()

	req := requests.Request{
		Method: requests.Get,
		Url:    "https://hacker-news.firebaseio.com/v0/newstories.json",
	}

	watcher, err := req.Watch()
	if err != nil {
		time.Sleep(constants.RetryDelay)
		b.Run()
		return
	}
	defer watcher.Close()
	b.watcher = watcher

	for evt := range watcher.Stream {
		err := b.parse(evt)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("tasks.hackernews: Failed to parse event")
		}
	}

	if !b.stop {
		b.Run()
		return
	}
}

func (b *hackerNewsBackend) Stop() {
	b.stop = true
	b.watcher.Close()
}

type hackerNews struct {
}

func (h *hackerNews) Type() string {
	return "hacker_news"
}

func (h *hackerNews) Run(db *database.Database) (err error) {
	stream := stream.NewStream(db, hackerNewsAcctId,
		hackerNewsUserId, &hackerNewsBackend{})

	err = stream.Start()
	if err != nil {
		return
	}

	return
}

func init() {
	register(-1, -1, 45, &hackerNews{})
}
