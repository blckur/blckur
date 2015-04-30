package logger

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/constants"
	"github.com/blckur/blckur/errortypes"
	"github.com/blckur/blckur/settings"
	"github.com/dropbox/godropbox/errors"
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
)

func init() {
	senders = append(senders, &rollbarSender{})
}

type rollbarMessage struct {
	Body string `json:"body"`
}

type rollbarBody struct {
	Message *rollbarMessage `json:"message"`
}

type rollbarServer struct {
	Host        string `json:"host"`
	CodeVersion string `json:"code_version"`
}

type rollbarData struct {
	Environment string         `json:"environment"`
	Level       string         `json:"level"`
	Timestamp   int64          `json:"timestamp"`
	CodeVersion string         `json:"code_version"`
	Platform    string         `json:"platform"`
	Language    string         `json:"language"`
	Context     string         `json:"context"`
	Fingerprint string         `json:"fingerprint"`
	Uuid        bson.ObjectId  `json:"uuid"`
	Title       string         `json:"title"`
	Server      *rollbarServer `json:"server"`
	Body        *rollbarBody   `json:"body"`
}

type rollbarItem struct {
	AccessToken string       `json:"access_token"`
	Data        *rollbarData `json:"data"`
}

type rollbarSender struct {
	limit  map[string]time.Time
	buffer chan *logrus.Entry
}

func (r *rollbarSender) listen() {
	for {
		entry := <-r.buffer

		// TODO Use r.send(entry)
		err := rollbarSend(entry)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("logger: Rollbar send error")
		}
	}
}

func (r *rollbarSender) Init() {
	r.buffer = make(chan *logrus.Entry, 128)
	go r.listen()
}

func (r *rollbarSender) Limit() time.Duration {
	return time.Duration(settings.Rollbar.RateLimit) * time.Second
}

func (r *rollbarSender) Parse(entry *logrus.Entry) {
	if settings.Rollbar.Token == "" {
		return
	}

	if len(buffer) <= 125 {
		r.buffer <- entry
	}
}

func rollbarSend(entry *logrus.Entry) (err error) {
	token := settings.Rollbar.Token
	if token == "" {
		return
	}

	msg := entry.Message

	var error string
	for key, val := range entry.Data {
		if key == "error" {
			error = fmt.Sprintf("%s", val)
		} else {
			msg += fmt.Sprintf("\n    %s=%v", key, val)
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

	payload := rollbarItem{
		AccessToken: token,
		Data: &rollbarData{
			Environment: settings.Rollbar.Environment,
			Level:       entry.Level.String(),
			Timestamp:   entry.Time.Unix(),
			CodeVersion: constants.Version,
			Platform:    "linux",
			Language:    "go",
			Context:     constants.Node,
			Fingerprint: hash,
			Uuid:        bson.NewObjectId(),
			Title:       entry.Message,
			Server: &rollbarServer{
				Host:        constants.Host,
				CodeVersion: constants.Version,
			},
			Body: &rollbarBody{
				Message: &rollbarMessage{
					Body: msg,
				},
			},
		},
	}

	payloadByt, err := json.Marshal(payload)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "logger.rollbar: Parse error"),
		}
		return
	}

	resp, err := http.Post("https://api.rollbar.com/api/1/item/",
		"application/json", bytes.NewBuffer(payloadByt))
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "logger.rollbar: Api error"),
		}
		return
	}

	if resp.StatusCode != 200 {
		err = &errortypes.UnknownError{
			errors.New("logger.rollbar: Bad response code"),
		}
		return
	}

	return
}
