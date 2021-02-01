package diff

import (
	"bytes"
	"fmt"
	"sync"
)

type SourceType int8

const (
	Source SourceType = iota
	Target
)

type SyncOperation int8

const (
	Copy SyncOperation = iota
	Create
	Delete
)

func (o SyncOperation) String() string {
	switch o {
	case Copy:
		return "copy"
	case Create:
		return "create"
	case Delete:
		return "delete"
	default:
		return fmt.Sprintf("%d", int(o))
	}
}

type fileWithType struct {
	Type SourceType
	ChecksumInfo
}

type FileOperation struct {
	Operation SyncOperation
	FileData
}

func Compare(source <-chan ChecksumInfo, target <-chan ChecksumInfo, errCh chan<- error) []FileOperation {
	var smap sync.Map
	fileOps := make([]FileOperation, 0)

	fo1 := iterateFiles(source, &smap, Source)
	fo2 := iterateFiles(target, &smap, Target)

	for fo := range merge(fo1, fo2) {
		fileOps = append(fileOps, fo)
	}

	smap.Range(func(_, value interface{}) bool {
		fType := value.(fileWithType)
		fo := FileOperation{
			FileData: fType.FileData,
		}

		switch fType.Type {
		case Source:
			fo.Operation = Delete
		case Target:
			if fType.IsDir {
				fo.Operation = Create
			} else {
				fo.Operation = Copy
			}
		}

		fileOps = append(fileOps, fo)

		return true
	})

	return fileOps
}

func iterateFiles(dc <-chan ChecksumInfo, smap *sync.Map, sourceType SourceType) <-chan FileOperation {
	out := make(chan FileOperation)

	go func() {
		defer close(out)

		for info := range dc {
			fType := fileWithType{
				Type:         sourceType,
				ChecksumInfo: info,
			}

			if actual, loaded := smap.LoadOrStore(info.Path, fType); loaded {
				targetChecksum := actual.(fileWithType)

				if bytes.Compare(targetChecksum.Checksum, info.Checksum) != 0 {
					out <- FileOperation{
						Operation: Copy,
						FileData:  info.FileData,
					}
				}

				smap.Delete(info.Path)
			}
		}
	}()

	return out
}

func merge(cs ...<-chan FileOperation) <-chan FileOperation {
	var wg sync.WaitGroup
	out := make(chan FileOperation)

	output := func(c <-chan FileOperation) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func PrintOperation(file string, method SyncOperation) {
	fmt.Printf("%s `%s`\n", method.String(), file)
}
