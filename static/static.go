package static

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"crypto/md5"
	"encoding/base32"
	"strings"
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
	}
)

type File struct {
	Type string
	Data []byte
}

type Store struct {
	files map[string]*File
	lookup map[string]string
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

		baseName := name[:len(name) - len(ext)]

		data, e := ioutil.ReadFile(fullPath)
		if e != nil {
			err = e
			return
		}

		file := &File{
			Type: mimeTypes[ext],
			Data: data,
		}

		hash := md5.Sum(data)
		hashStr := base32.StdEncoding.EncodeToString(hash[:])
		hashStr = strings.Replace(hashStr, "/", "", -1)
		hashStr = strings.Replace(hashStr, "+", "", -1)
		hashStr = strings.Replace(hashStr, "=", "", -1)
		hashStr = strings.ToLower(hashStr[:10])

		hashName := baseName + "." + hashStr + ext

		s.files[hashName] = file
		s.lookup[name] = hashName
	}

	return
}

func NewStore(root string) (store *Store, err error) {
	store = &Store{
		files: map[string]*File{},
		lookup: map[string]string{},
	}

	err = store.addDir(root)
	if err != nil {
		return
	}

	return
}