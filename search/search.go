// Interface to elasticsearch.
package search

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/mattbaird/elastigo/lib"
	"strings"
	"time"
)

type Mapping map[string]interface{}

func (m Mapping) Add(field string, enabled bool, typ string, index string) {
	fieldMap := map[string]interface{}{
		"type":  typ,
	}

	if !enabled {
		fieldMap["enabled"] = enabled
	}

	if index != "" {
		fieldMap["index"] = index
	}

	m[field] = fieldMap
}

type Session struct {
	conn *elastigo.Conn
}

func (s *Session) Close() {
	s.conn.Close()
}

func (s *Session) CreateIndex(index string) (err error) {
	_, err = s.conn.CreateIndex("logs")
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			err = nil
		} else {
			err = &UnknownError{
				errors.Wrap(err, "search: Unknown create index error"),
			}
		}
	}

	return
}

func (s *Session) PutMapping(index string, typ string, doc interface{},
	mapping Mapping) (err error) {

	opts := elastigo.MappingOptions{
		Properties: mapping,
	}

	err = s.conn.PutMapping(index, typ, doc, opts)
	if err != nil {
		err = &UnknownError{
			errors.Wrap(err, "search: Unknown mapping error"),
		}
	}

	return
}

func (s *Session) Index(index string, typ string, id string, doc interface{}) (
	err error) {

	_, err = s.conn.Index(index, typ, id, nil, doc)
	if err != nil {
		err = &UnknownError{
			errors.Wrap(err, "search: Unknown index error"),
		}
	}

	return
}

func NewSession() (conn *Session) {
	conn = &Session{
		conn: &elastigo.Conn{
			Protocol:       elastigo.DefaultProtocol,
			Domain:         elastigo.DefaultDomain,
			ClusterDomains: []string{elastigo.DefaultDomain},
			Port:           elastigo.DefaultPort,
			DecayDuration: time.Duration(
				elastigo.DefaultDecayDuration * time.Second),
		},
	}

	return
}
