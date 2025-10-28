package cache

import (
	"bytes"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/logger"
	"compress/gzip"
	"encoding/gob"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var MyFileCache *FileCache

type FileCache struct {
	dir string
}

func NewFileCache(dir string) *FileCache {
	return &FileCache{dir: dir}
}

func (fc *FileCache) Set(key string, value interface{}, expire ...time.Duration) error {
	cacheFilePath := filepath.Join(fc.dir, key)
	var cacheBuf bytes.Buffer
	if err := gob.NewEncoder(&cacheBuf).Encode(value); err != nil {
		return err
	}

	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	if _, err := gzWriter.Write(cacheBuf.Bytes()); err != nil {
		return err
	}
	if err := gzWriter.Close(); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(cacheFilePath), 0755); err != nil {
		return err
	}

	if err := ioutil.WriteFile(cacheFilePath, buf.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

func (fc *FileCache) Get(key string) (interface{}, error) {
	cacheFilePath := filepath.Join(fc.dir, key)
	cachedData, err := ioutil.ReadFile(cacheFilePath)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Write(cachedData)
	gzReader, err := gzip.NewReader(&buf)
	if err != nil {
		return nil, err
	}
	defer func(gzReader *gzip.Reader) {
		err := gzReader.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(gzReader)

	decompressedData, err := ioutil.ReadAll(gzReader)
	if err != nil {
		return nil, err
	}

	var cacheData map[string]interface{}
	if err := gob.NewDecoder(bytes.NewReader(decompressedData)).Decode(&cacheData); err != nil {
		return nil, err
	}

	return cacheData, nil
}

func (fc *FileCache) Del(key string) error {
	cacheFilePath := filepath.Join(fc.dir, key)
	return os.Remove(cacheFilePath)
}

func (fc *FileCache) DelByPrefix(prefix string) error {
	files, err := filepath.Glob(filepath.Join(fc.dir, prefix+"*"))
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}

func (fc *FileCache) Update(key string, value map[string]interface{}) error {
	cacheFilePath := filepath.Join(fc.dir, key)

	// Check if the cache file exists
	if _, err := os.Stat(cacheFilePath); err != nil {
		return err // Return error if the file doesn't exist
	}

	// Encode the new value
	var cacheBuf bytes.Buffer
	if err := gob.NewEncoder(&cacheBuf).Encode(value); err != nil {
		return err
	}

	// Compress the encoded value
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	if _, err := gzWriter.Write(cacheBuf.Bytes()); err != nil {
		return err
	}
	if err := gzWriter.Close(); err != nil {
		return err
	}

	// Write the compressed data to the cache file
	if err := ioutil.WriteFile(cacheFilePath, buf.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}
