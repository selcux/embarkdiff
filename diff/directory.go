package diff

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const dirKey string = "__dir__"

type checksumInfo struct {
	file     string
	checksum string
}

func ExecuteChecksum(dir string) (<-chan checksumInfo, error) {
	entities, err := getEntities(dir)
	if err != nil {
		return nil, err
	}

	return streamChecksum(entities)
}

func getEntities(dir string) (map[string]string, error) {
	entities := make(map[string]string)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if dir == path {
			return nil
		}

		entities[path] = ""

		if info.IsDir() {
			entities[path] = dirKey
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return entities, nil
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

func streamChecksum(entities map[string]string) (<-chan checksumInfo, error) {
	//This is the maximum number of concurrent file operations
	var MaxCount = 50
	var maxGoroutines = make(chan struct{}, MaxCount)

	chksum := make(chan checksumInfo)
	var wg sync.WaitGroup
	wg.Add(len(entities))

	for k, v := range entities {
		p1 := &checksumInfo{k, v}

		go func(p2 *checksumInfo) {
			defer wg.Done()

			//This limits the maximum number of concurrent file operations
			// to revent `too many open files` error
			maxGoroutines <- struct{}{}
			defer func() { <-maxGoroutines }()

			if p2.checksum == dirKey {
				chksum <- *p2
				return
			}

			check, err := checksum(p2.file)
			if err != nil {
				log.Fatalln(err)
			}

			chksum <- checksumInfo{
				file:     p2.file,
				checksum: hex.EncodeToString(check),
			}
		}(p1)
	}
	go func() {
		wg.Wait()
		close(chksum)
	}()

	return chksum, nil
}

func batchChecksum(entities map[string]string) (map[string]string, error) {
	//This is the maximum number of concurrent file operations
	var MaxCount = 50
	var maxGoroutines = make(chan struct{}, MaxCount)

	chksum := make(chan checksumInfo)
	quit := make(chan struct{})
	fmap := make(map[string]string)
	var smap sync.Map
	var wg sync.WaitGroup

	computeChecksum(chksum, maxGoroutines, &wg, entities)
	go waitChecksum(chksum, quit, &smap)

	wg.Wait()
	close(quit)

	smap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		fmap[k] = string(v)
		return true
	})

	return fmap, nil
}

func computeChecksum(chksum chan checksumInfo, maxGoroutines chan struct{}, wg *sync.WaitGroup, entities map[string]string) {
	for k, v := range entities {
		wg.Add(1)

		p1 := &checksumInfo{k, v}

		go func(p2 *checksumInfo) {
			defer wg.Done()

			//This limits the maximum number of concurrent file operations
			// to revent `too many open files` error
			maxGoroutines <- struct{}{}
			defer func() { <-maxGoroutines }()

			if p2.checksum == dirKey {
				chksum <- *p2
				return
			}

			check, err := checksum(p2.file)
			if err != nil {
				log.Fatalln(err)
			}

			chksum <- checksumInfo{
				file:     p2.file,
				checksum: hex.EncodeToString(check),
			}
		}(p1)
	}
}

func waitChecksum(chksum chan checksumInfo, quit chan struct{}, smap *sync.Map) {
	for {
		select {
		case p := <-chksum:
			smap.Store(p.file, p.checksum)
		case <-quit:
			return
		}
	}
}
