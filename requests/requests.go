// Http request interface.
package requests

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/dropbox/godropbox/errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	Get     = "GET"
	Put     = "PUT"
	Post    = "POST"
	Delete  = "DELETE"
	Trace   = "TRACE"
	Options = "OPTIONS"
	Connect = "CONNECT"
	Patch   = "PATCH"
)

var (
	DefaultTimeout      = 10 * time.Second
	DefaultWatchTimeout = 2 * time.Minute
)

type Request struct {
	Method    string
	Url       string
	Timeout   time.Duration
	Params    map[string]string
	Header    http.Header
	Json      interface{}
	Body      io.Reader
	Transport *http.Transport
	Client    *http.Client
}

func (r *Request) client() (client *http.Client) {
	timeout := r.Timeout
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	transport := r.Transport
	if transport == nil {
		transport = &http.Transport{
			Dial: func(network string, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, r.Timeout)
			},
		}
	}

	client = r.Client
	if client == nil {
		client = &http.Client{}
	}
	client.Transport = transport

	return
}

func (r *Request) Do() (resp *Response, err error) {
	client := r.client()

	body := r.Body
	if r.Json != nil {
		data, e := json.Marshal(r.Json)
		if e != nil {
			err = &ParseError{
				errors.Wrap(e, "requests: Json marshal error"),
			}
			return
		}

		body = bytes.NewReader(data)
	}

	req, err := http.NewRequest(r.Method, r.Url, body)
	if err != nil {
		err = &ParseError{
			errors.Wrap(err, "requests: Invalid request"),
		}
		return
	}

	if r.Header != nil {
		req.Header = r.Header
	}

	res, err := client.Do(req)
	if err != nil {
		err = &RequestError{
			errors.Wrap(err, "requests: Request failed"),
		}
		return
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		body, _ := ioutil.ReadAll(res.Body)

		err = &ResponseError{
			errors.Newf("requests: %d Bad response:\n\n%s",
				res.StatusCode, body),
		}
		return
	}

	resp = &Response{
		res,
	}

	return
}

func (r *Request) Watch() (watcher *Watcher, err error) {
	client := r.client()

	req, err := http.NewRequest(r.Method, r.Url, nil)
	if err != nil {
		err = &ParseError{
			errors.Wrap(err, "requests: Invalid request"),
		}
		return
	}

	if r.Header != nil {
		req.Header = r.Header
	}

	req.Header.Add("Accept", "text/event-stream")

	res, err := client.Do(req)
	if err != nil {
		err = &RequestError{
			errors.Wrap(err, "requests: Request failed"),
		}
		return
	}

	timeout := r.Timeout
	if timeout == 0 {
		timeout = DefaultWatchTimeout
	}

	watcher = &Watcher{
		res:     res,
		scanner: bufio.NewScanner(res.Body),
		timeout: timeout,
		Stream:  make(chan *Event),
	}

	watcher.scanner.Split(bufio.ScanLines)

	go watcher.scan()

	return
}

type Response struct {
	*http.Response
}

func (r *Response) Json(v interface{}) (err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = &ParseError{
			errors.Wrap(err, "requests: Io error"),
		}
		return
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		err = &ParseError{
			errors.Wrap(err, "requests: Json unmarshal error"),
		}
	}

	return
}
