package diff

import (
	"fmt"
	"github.com/selcux/embarkdiff/util"
	"log"
	"path"
	"path/filepath"
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
	checksumInfo
}

type DirWithChannel struct {
	Ch  <-chan checksumInfo
	Dir string
}

func Compare(source *DirWithChannel, target *DirWithChannel) {
	var smap sync.Map
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		iterateFiles(source, &smap, Source)
	}()

	go func() {
		defer wg.Done()
		iterateFiles(target, &smap, Target)
	}()

	wg.Wait()

	smap.Range(func(key, value interface{}) bool {
		file := key.(string)
		fType := value.(fileWithType)

		switch fType.Type {
		case Source:
			printOperation(file, Delete)
		case Target:
			absPath := path.Join(target.Dir, file)
			isDir, err := util.DirExists(absPath)
			if err != nil {
				log.Fatalln(err)
			}

			if isDir {
				printOperation(file, Create)
			} else {
				printOperation(file, Copy)
			}
		}

		return true
	})
}

func printOperation(file string, method SyncOperation) {
	fmt.Printf("%s `%s`\n", method.String(), file)
}

func iterateFiles(dc *DirWithChannel, smap *sync.Map, sourceType SourceType) {
	for pair := range dc.Ch {
		file, err := filepath.Rel(dc.Dir, pair.file)
		if err != nil {
			log.Fatalln(err)
		}
		pair.file = file

		fType := fileWithType{
			Type:         sourceType,
			checksumInfo: pair,
		}

		if actual, loaded := smap.LoadOrStore(pair.file, fType); loaded {
			targetChecksum := actual.(fileWithType)

			if targetChecksum.checksum != pair.checksum {
				printOperation(pair.file, Copy)
			}

			smap.Delete(pair.file)
		}
	}
}
