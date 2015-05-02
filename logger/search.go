package logger

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/requires"
	"github.com/blckur/blckur/search"
	"github.com/blckur/blckur/settings"
	"labix.org/v2/mgo/bson"
	"time"
)

type searchFields map[string]interface{}

type searchEntry struct {
	Environment string       `json:"environment"`
	Level       string       `json:"level"`
	Timestamp   string       `json:"timestamp"`
	Version     string       `json:"version"`
	Fingerprint string       `json:"fingerprint"`
	Server      string       `json:"server"`
	Host        string       `json:"host"`
	Client      string       `json:"client"`
	Message     string       `json:"message"`
	Body        string       `json:"body"`
	Fields      searchFields `json:"fields"`
}

type searchSender struct {
	limit  limiter
	buffer chan *logrus.Entry
}

func (s *searchSender) listen() {
	for {
		entry := <-s.buffer

		conn := search.NewSession()
		if conn == nil {
			continue
		}

		err := s.send(conn, entry)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("logger: Search send error")
		}
	}
}

func (s *searchSender) Init() {
	s.limit = limiter{}
	s.buffer = make(chan *logrus.Entry, 128)
	go s.listen()
}

func (s *searchSender) Parse(entry *logrus.Entry) {
	if !s.limit.Check(entry, 3*time.Second) {
		return
	}

	if len(s.buffer) <= 125 {
		s.buffer <- entry
	}
}

func (s *searchSender) send(conn *search.Session, entry *logrus.Entry) (
	err error) {

	msg := entry.Message
	fields := searchFields{}

	var error string
	for key, val := range entry.Data {
		if key == "error" {
			error = fmt.Sprintf("%s", val)
		} else {
			fields[key] = fmt.Sprintf("%v", val)
		}
	}

	if error != "" {
		msg += "\n\n" + error
	}

	if string(msg[len(msg)-1]) != "\n" {
		msg += "\n"
	}

	hashByt := md5.Sum([]byte(entry.Message))
	hash := hex.EncodeToString(hashByt[:])

	ent := &searchEntry{
		Environment: settings.System.Environment,
		Level:       entry.Level.String(),
		Timestamp:   entry.Time.UTC().Format("2006-01-02T15:04:05"),
		Version:     constants.Version,
		Fingerprint: hash,
		Server:      constants.Node,
		Host:        constants.Host,
		Message:     entry.Message,
		Body:        msg,
		Fields:      fields,
	}

	err = conn.Index("logs", "entry", bson.NewObjectId().Hex(), ent)
	if err != nil {
		return
	}

	return
}

func updateSearch() {
	conn := search.NewSession()
	if conn == nil {
		return
	}

	mapping := search.Mapping{}
	mapping.Add("environment", true, "string", "not_analyzed")
	mapping.Add("level", true, "string", "not_analyzed")
	mapping.Add("timestamp", true, "date", "")
	mapping.Add("version", true, "string", "not_analyzed")
	mapping.Add("fingerprint", true, "string", "not_analyzed")
	mapping.Add("server", true, "string", "not_analyzed")
	mapping.Add("host", true, "string", "not_analyzed")
	mapping.Add("client", true, "string", "not_analyzed")
	mapping.Add("message", true, "string", "not_analyzed")
	mapping.Add("body", true, "string", "analyzed")
	mapping.Add("fields", false, "object", "")

	err := conn.CreateIndex("logs")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("logger.elasticsearch: Update error")
		time.Sleep(constants.RetryDelay)
		updateSearch()
		return
	}

	err = conn.PutMapping("logs", "entry", searchEntry{}, mapping)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("logger.elasticsearch: Update error")
		time.Sleep(constants.RetryDelay)
		updateSearch()
		return
	}
}

func init() {
	senders = append(senders, &searchSender{})

	module := requires.New("logger.search")
	module.After("settings")

	module.Handler = func() {
		messenger.Register("settings", "search", func(_ *messenger.Message) {
			go updateSearch()
		})
		updateSearch()
	}
}
