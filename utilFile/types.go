package utilFile

import "io/fs"

type FilterFunc func(file fs.FileInfo, path string) bool
