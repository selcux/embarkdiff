package diff

import (
	"bytes"
	"fmt"
	"sync"
)

//SourceType defines whether the folder is a source or a target
type SourceType int8

const (
	Source SourceType = iota
	Target
)

//SyncOperation declares what kind of process will be needed to sync the directories
type SyncOperation int8

const (
	Delete SyncOperation = iota
	Create
	Copy
)

func (o SyncOperation) String() string {
	switch o {
	case Delete:
		return "delete"
	case Create:
		return "create"
	case Copy:
		return "copy"
	default:
		return fmt.Sprintf("%d", int(o))
	}
}

//fileWithType is being used to distinct files from source and target paths
type fileWithType struct {
	Type SourceType
	ChecksumInfo
}

//Compare the files in the source and the target paths and return the determined operations in order to sync them
func Compare(source <-chan ChecksumInfo, target <-chan ChecksumInfo) FileOps {
	var smap sync.Map
	fileOps := make(FileOps, 0)

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

//iterateFiles takes the data of the files and stores them into a thread safe map.
//If there is a duplicate key of the map, compares the checksums to understand if the files is changed or not.
//If the files is changed, removes it from the map and sends the copy operation in return.
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

//merge the given channels into one
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

//PrintOperation prints the file operation to be applied in order to sync the files
func PrintOperation(file string, method SyncOperation) {
	fmt.Printf("%s `%s`\n", method.String(), file)
}
