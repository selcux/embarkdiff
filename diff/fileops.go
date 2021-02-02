package diff

import (
	"os"
	"path"
	"sort"
	"strings"
)

//FileOperation defines the operation type for the related file
type FileOperation struct {
	Operation SyncOperation
	FileData
}

//PathDepth returns the hierarchical length of the file or directory
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

//SortByPathLen sorts the files according to the increasing hierarchical depth
func (fo FileOps) SortByPathLen() {
	sort.Slice(fo, func(i, j int) bool {
		return fo[i].PathDepth() < fo[j].PathDepth()
	})
}

//GroupByOps groups the files by the operation type
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

//Reverse the order of files in the slice
func (fo FileOps) Reverse() {
	for i, j := 0, len(fo)-1; i < j; i, j = i+1, j-1 {
		fo[i], fo[j] = fo[j], fo[i]
	}
}

/*
Transform the file array into required operation order.
The required order is defined by the order of `SyncOperation`

1. Delete the files and folders in a bottom up hierarchy
2. Create the missing folders
3. Copy the changed or the missing files
*/
func (fo FileOps) Transform() FileOps {
	return flatten(fo.GroupByOps())
}

//flatten the grouped elements into an array
func flatten(groupMap map[SyncOperation]FileOps) FileOps {
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
