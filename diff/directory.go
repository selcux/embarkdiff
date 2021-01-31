package diff

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const dirKey string = "__dir__"

type pair struct {
	file     string
	checksum string
}

func ExecuteChecksum(dir string) error {
	entities, err := getEntities(dir)
	if err != nil {
		return err
	}

	files, err := batchChecksum(entities)
	if err != nil {
		return err
	}

	for k, v := range files {
		fmt.Println(k, v)
	}

	return nil
}

func getEntities(dir string) (map[string]string, error) {
	entities := make(map[string]string)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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

func batchChecksum(entities map[string]string) (map[string]string, error) {
	chksum := make(chan pair, 4)
	quit := make(chan struct{}, 0)
	fmap := make(map[string]string)
	var smap sync.Map
	var wg sync.WaitGroup

	computeChecksum(chksum, &wg, entities)
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

func computeChecksum(chksum chan pair, wg *sync.WaitGroup, entities map[string]string) {
	for k, v := range entities {
		wg.Add(1)

		p1 := &pair{k, v}

		go func(p2 *pair, sum chan pair, wgroup *sync.WaitGroup) {
			defer wg.Done()

			if p2.checksum == dirKey {
				sum <- *p2
				return
			}

			check, err := checksum(p2.file)
			if err != nil {
				log.Fatalln(err)
			}

			sum <- pair{
				file:     p2.file,
				checksum: hex.EncodeToString(check),
			}
		}(p1, chksum, wg)
	}
}

func waitChecksum(chksum chan pair, quit chan struct{}, smap *sync.Map) {
	for {
		select {
		case p := <-chksum:
			fmt.Println(p.file, p.checksum)
			smap.Store(p.file, p.checksum)
		case <-quit:
			fmt.Println("goto exit")
			return
		}
	}
}
