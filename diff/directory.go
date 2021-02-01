package diff

import (
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"
	"sync"
)

//FileData keeps the path and whether the path is a directory or not
type FileData struct {
	Path  string
	IsDir bool
}

//ChecksumInfo embeds the FileData and includes the checksum data for the related file.
//If the path is a directory Checksum is nil.
type ChecksumInfo struct {
	FileData
	Checksum []byte
}

//ExecuteChecksum takes the rood directory and returns the checksum info of its files
func ExecuteChecksum(root string, errCh chan<- error) <-chan ChecksumInfo {
	entitiesCh := entities(root, errCh)

	return streamChecksum(entitiesCh, errCh, root)
}

//entities traverse through the given directory and returns a channel where the path of the every file/subfolder is sent
func entities(root string, errCh chan<- error) <-chan FileData {
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

			paths <- FileData{path, info.IsDir()}

			return nil
		})
		if err != nil {
			errCh <- err
		}
	}()

	return paths
}

//checksum calculates the checksum of the given file
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

//createChecksumInfo calculates the checksum of the given FileData and coverts it into ChecksumInfo
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

//streamChecksum retrieves the entities, calculates their checksum and sends the checksum information through a channel
func streamChecksum(entities <-chan FileData, errCh chan<- error, root string) <-chan ChecksumInfo {
	const numDigesters = 50
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
