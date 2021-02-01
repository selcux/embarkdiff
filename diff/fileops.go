package diff

import (
	"os"
	"path"
	"sort"
	"strings"
)

type FileOperation struct {
	Operation SyncOperation
	FileData
}

func (f FileOperation) PathDepth() int {
	dir := f.Path
	if !f.IsDir {
		dir = path.Dir(f.Path)
	}

	fileSep := string(os.PathSeparator)
	dir = strings.TrimPrefix(dir, fileSep)
	dirs := strings.Split(dir, fileSep)

	return len(dirs)
}

type FileOps []FileOperation

func (fo FileOps) SortByPathLen() {
	sort.Slice(fo, func(i, j int) bool {
		return fo[i].PathDepth() < fo[j].PathDepth()
	})
}

func (fo FileOps) GroupByOps() map[SyncOperation]FileOps {
	groupMap := make(map[SyncOperation]FileOps)

	for _, o := range fo {
		if _, ok := groupMap[o.Operation]; !ok {
			groupMap[o.Operation] = make(FileOps, 0)
		}

		groupMap[o.Operation] = append(groupMap[o.Operation], o)
	}

	return groupMap
}

func (fo FileOps) Reverse() {
	for i, j := 0, len(fo)-1; i < j; i, j = i+1, j-1 {
		fo[i], fo[j] = fo[j], fo[i]
	}
}

func (fo FileOps) Transform() FileOps {
	return convert(fo.GroupByOps())
}

func convert(groupMap map[SyncOperation]FileOps) FileOps {
	fops := make(FileOps, 0)

	for i := SyncOperation(0); i <= Copy; i++ {
		groupMap[i].SortByPathLen()
		if i == Delete {
			groupMap[i].Reverse()
		}

		fops = append(fops, groupMap[i]...)
	}

	return fops
}
