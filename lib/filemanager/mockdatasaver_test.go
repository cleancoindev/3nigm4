//
// 3nigm4 crypto package
// Author: Guido Ronchetti <dyst0ni3@gmail.com>
// v1.0 06/03/2016
//
package filemanager

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type localDataSaver struct {
	rootPath string
}

func NewLocalDataSaver(root string) (*localDataSaver, error) {
	if fi, err := os.Stat(root); os.IsNotExist(err) ||
		fi.IsDir() != true {
		return nil, fmt.Errorf("invalid root path")
	}
	return &localDataSaver{
		rootPath: root,
	}, nil
}

func (l *localDataSaver) SaveChunks(filename string, chunks [][]byte, hashedValue []byte) ([]string, error) {
	paths := make([]string, len(chunks))
	for idx, chunk := range chunks {
		id, err := ChunkFileId(filename, idx, hashedValue)
		if err != nil {
			return nil, err
		}
		err = ioutil.WriteFile(filepath.Join(l.rootPath, id), chunk, 0644)
		if err != nil {
			return nil, err
		}
		paths[idx] = id
	}
	return paths, nil
}

func (l *localDataSaver) RetrieveChunks(files []string) ([][]byte, error) {
	chunks := make([][]byte, len(files))
	for idx, file := range files {
		data, err := ioutil.ReadFile(filepath.Join(l.rootPath, file))
		if err != nil {
			return nil, err
		}
		chunks[idx] = data
	}
	return chunks, nil
}

func (l *localDataSaver) Cleanup() error {
	return os.RemoveAll(l.rootPath)
}
