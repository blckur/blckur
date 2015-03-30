// Versions static files with hash, replaces references and stores in memory.
package static

import (
	"github.com/blckur/blckur/errortypes"
	"github.com/dropbox/godropbox/errors"
	"io/ioutil"
	"path"
	"path/filepath"
	"crypto/md5"
	"encoding/base32"
	"strings"
	"regexp"
)

var (
	mimeTypes = map[string]string{
		".js": "application/javascript",
		".dart": "application/dart",
		".json": "application/json",
		".css": "text/css",
		".html": "text/html",
		".jpg": "image/jpeg",
		".png": "image/png",
		".svg": "image/svg+xml",
		".otf": "application/font-sfnt",
		".ttf": "application/font-sfnt",
		".woff": "application/font-woff",
		".eot": "application/vnd.ms-fontobject",
		".map": "application/json",
	}
	srcReg = regexp.MustCompile(`(src=)('|")(.*?)('|")`)
	pathReg = regexp.MustCompile(`('|")(packages/|/s/img/)(.*?)('|")`)
)

type File struct {
	Type string
	Data []byte
}

type FileName struct {
	Name string
	HashName string
}

type Store struct {
	Files map[string]*File
	HashFiles map[string]*File
	root string
	dirNames map[string][]*FileName
	lookup map[string]*FileName
}

func (s *Store) addDir(dir string) (err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, info := range files {
		name := info.Name()
		fullPath := path.Join(dir, name)

		if info.IsDir() {
			s.addDir(fullPath)
			continue
		}

		ext := filepath.Ext(name)
		if len(ext) == 0 {
			continue
		}

		typ, ok := mimeTypes[ext]
		if !ok {
			continue
		}

		baseName := name[:len(name) - len(ext)]

		data, e := ioutil.ReadFile(fullPath)
		if e != nil {
			err = e
			return
		}

		file := &File{
			Type: typ,
			Data: data,
		}

		s.Files[fullPath] = file

		if name == "index.html" {
			s.HashFiles[fullPath] = file
			continue
		}

		hash := md5.Sum(data)
		hashStr := base32.StdEncoding.EncodeToString(hash[:])
		hashStr = strings.Replace(hashStr, "/", "", -1)
		hashStr = strings.Replace(hashStr, "+", "", -1)
		hashStr = strings.Replace(hashStr, "=", "", -1)
		hashStr = strings.ToLower(hashStr[:10])
		hashName := baseName + "." + hashStr + ext

		s.HashFiles[path.Join(dir, hashName)] = file

		s.lookup[fullPath] = &FileName{
			Name: name,
			HashName: hashName,
		}
	}

	return
}

func (s *Store) parseFiles() {
	for path, file := range s.Files {
		path = filepath.Dir(path)
		dataStr := string(file.Data)

		dataStr = pathReg.ReplaceAllStringFunc(dataStr, func(
		match string) string {
			var matchPath string
			if match[1:4] == "/s/" {
				matchPath = filepath.Join(s.root, match[4:len(match) - 1])
			} else {
				matchPath = filepath.Join(path, match[1:len(match) - 1])
			}

			if name, ok := s.lookup[matchPath]; ok {
				match = strings.Replace(match, name.Name, name.HashName, 1)
			}

			return match
		})

		dataStr = srcReg.ReplaceAllStringFunc(dataStr, func(
				match string) string {
			matchPath := filepath.Join(path, match[5:len(match) - 1])

			if name, ok := s.lookup[matchPath]; ok {
				match = strings.Replace(match, name.Name, name.HashName, 1)
			}

			return match
		})

		file.Data = []byte(dataStr)
	}
}

func NewStore(root string) (store *Store, err error) {
	store = &Store{
		Files: map[string]*File{},
		HashFiles: map[string]*File{},
		root: root,
		lookup: map[string]*FileName{},
	}

	err = store.addDir(root)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "static: Init error"),
		}
		return
	}

	store.parseFiles()

	store.lookup = nil

	return
}
