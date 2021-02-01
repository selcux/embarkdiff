package diff

import (
	"context"
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"
	"sync"
)

const dirKey string = "__dir__"

type FileData struct {
	Path  string
	IsDir bool
}

type ChecksumInfo struct {
	FileData
	Checksum []byte
}

func ExecuteChecksum(ctx context.Context, root string, errCh chan<- error) <-chan ChecksumInfo {
	entitiesCh := entities(ctx, root, errCh)

	return streamChecksum(ctx, entitiesCh, errCh, root)
}

func entities(ctx context.Context, root string, errCh chan<- error) <-chan FileData {
	paths := make(chan FileData)

	go func() {
		defer close(paths)

		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if root == path {
				return nil
			}

			select {
			case paths <- FileData{path, info.IsDir()}:
			case <-ctx.Done():
				return ctx.Err()
			}

			return nil
		})
		if err != nil {
			errCh <- err
		}
	}()

	return paths
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

func createChecksumInfo(entity FileData) (ChecksumInfo, error) {
	if entity.IsDir {
		return ChecksumInfo{FileData: entity}, nil
	}

	sum, err := checksum(entity.Path)
	if err != nil {
		return ChecksumInfo{}, err
	}

	return ChecksumInfo{
		FileData: entity,
		Checksum: sum,
	}, nil
}

func streamChecksum(ctx context.Context, entities <-chan FileData, errCh chan<- error, root string) <-chan ChecksumInfo {
	const numDigesters = 1
	checksumCh := make(chan ChecksumInfo)

	var wg1 sync.WaitGroup
	wg1.Add(numDigesters)

	for i := 0; i < numDigesters; i++ {
		go func(wg2 *sync.WaitGroup) {
			defer wg2.Done()
			for entity := range entities {
				checksumInfo, err := createChecksumInfo(entity)
				if err != nil {
					errCh <- err
					return
				}

				path, err := filepath.Rel(root, checksumInfo.Path)
				if err != nil {
					errCh <- err
					return
				}
				checksumInfo.Path = path // convert to a relative path
				checksumCh <- checksumInfo
			}
		}(&wg1)
	}

	go func() {
		wg1.Wait()
		close(checksumCh)
	}()

	return checksumCh
}
