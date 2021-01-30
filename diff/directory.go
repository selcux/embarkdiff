package diff

import (
	"crypto/sha256"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type DirInfo struct {
	Directories []string
	Files       map[string][]byte
}

func NewDirInfo(root string) (*DirInfo, error) {
	dirs, files, err := getFiles(root)
	if err != nil {
		return nil, err
	}

	fmap, err := batchChecksum(files)
	if err != nil {
		return nil, err
	}

	return &DirInfo{
		Directories: dirs,
		Files:       fmap,
	}, nil
}

func getFiles(dir string) ([]string, []string, error) {
	files := make([]string, 0)
	dirs := make([]string, 0)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirs = append(dirs, path)
		} else {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return dirs, files, nil
}

func checksum(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha256.New()

	_, err = io.Copy(h, f)

	return h.Sum(nil), nil

}

func batchChecksum(files []string) (map[string][]byte, error) {
	var wg sync.WaitGroup
	fmap := make(map[string][]byte)
	wg.Add(len(files))
	var smap sync.Map

	for _, file := range files {
		go func(f string) {
			sum, err := checksum(f)
			if err != nil {
				log.Fatalln(err)
			}

			// fmap[f] = sum
			smap.Store(f, sum)
			wg.Done()

		}(file)
	}

	wg.Wait()

	smap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.([]byte)
		fmap[k] = v

		return false
	})

	return fmap, nil
}
