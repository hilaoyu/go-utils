package utilFile

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/hilaoyu/go-utils/utilCmd"
	"github.com/hilaoyu/go-utils/utilStr"
	"github.com/hilaoyu/go-utils/utils"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type ByModTimeAsc []os.FileInfo
type ByModTimeDesc []os.FileInfo

var fileWalkBreakErr = errors.New("file walk break")

func (fis ByModTimeAsc) Len() int {
	return len(fis)
}

func (fis ByModTimeAsc) Swap(i, j int) {
	fis[i], fis[j] = fis[j], fis[i]
}

func (fis ByModTimeAsc) Less(i, j int) bool {
	return fis[i].ModTime().Before(fis[j].ModTime())
}

func (fis ByModTimeDesc) Len() int {
	return len(fis)
}

func (fis ByModTimeDesc) Swap(i, j int) {
	fis[i], fis[j] = fis[j], fis[i]
}

func (fis ByModTimeDesc) Less(i, j int) bool {
	return fis[i].ModTime().After(fis[j].ModTime())
}

func ToLinux(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

func SafePath(path string) string {
	path = ToLinux(path)
	path = strings.ReplaceAll(path, "|", "")
	path = strings.ReplaceAll(path, "&", "")
	path = strings.ReplaceAll(path, "'", "")
	path = strings.ReplaceAll(path, "\"", "")
	path = strings.ReplaceAll(path, ">", "")
	path = strings.ReplaceAll(path, "<", "")
	path = strings.ReplaceAll(path, "?", "")
	path = strings.ReplaceAll(path, "*", "")
	reg, _ := regexp.Compile("\\/\\.+\\/")
	return reg.ReplaceAllString(path, "")
}

func isExists(p string) (bool, os.FileInfo) {
	fi, err := os.Stat(p)
	if err != nil {
		return false, nil
	}
	return true, fi
}
func Exists(p string) bool {
	if e, _ := isExists(p); e {
		return e
	}
	return false
}

func IsFile(p string) bool {
	if e, fi := isExists(p); e {
		return fi.Mode().IsRegular()
	}
	return false
}

func IsDir(p string) bool {
	if e, fi := isExists(p); e {
		return fi.Mode().IsDir()
	}
	return false
}
func CheckDir(p string, perm ...os.FileMode) bool {
	if e, fi := isExists(p); e {
		return fi.Mode().IsDir()
	}
	newDirPerm := os.FileMode(0755)
	if len(perm) > 0 {
		newDirPerm = perm[0]
	}
	err := os.MkdirAll(p, newDirPerm)
	if nil != err {
		//fmt.Println(err)
		return false
	}
	return true
}
func GetDirSize(p string) int64 {
	var size int64
	output, err := utilCmd.RunCommand(true, "du", "-sb", p)
	if nil != err {
		size = 0
	}
	output = utilStr.BeforeFirst(output, "\t")
	size, err = strconv.ParseInt(output, 10, 64)
	if nil != err {
		size = 0
	}

	return size
}

func GetModTime(path string) (int64, error) {

	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}

	return fi.ModTime().Unix(), nil
}

func Move(src string, dst string) error {
	dstDir := filepath.Dir(dst)
	if !CheckDir(dstDir) {
		return errors.New("目标目录不可写")
	}
	if "windows" != utils.RunningOs("windows") {
		var cmd *exec.Cmd
		cmd = exec.Command("mv", src, dst)
		_, err := cmd.Output()
		if err != nil {
			return err
		}
		return nil
	}
	return os.Rename(src, dst)
}

func First(root string, findFn FilterFunc, recursion bool) (string, error) {
	//fmt.Println("First",path)
	var fileFind string
	var err error

	rootAbs, err := filepath.Abs(root)
	if nil != err {
		return "", err
	}
	err = filepath.Walk(rootAbs, func(path string, info os.FileInfo, err error) error {
		if nil != err {
			return err
		}
		if !recursion && info.IsDir() && (path != rootAbs) {
			return filepath.SkipDir
		}
		if findFn(info, path) {
			fileFind = path
			return filepath.SkipAll
		}
		return nil
	})

	return fileFind, err
}

func FirstN(root string, limit int, findFn FilterFunc) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if nil != err {
			return err
		}
		if findFn(info, path) {
			files = append(files, path)
			if len(files) >= limit {
				return fileWalkBreakErr
			}
		}
		return nil
	})

	if err == fileWalkBreakErr {
		err = nil
	}

	return files, err
}

func SortByModTimeAsc(fis []os.FileInfo) (sortedFiles ByModTimeAsc) {
	sortedFiles = ByModTimeAsc(fis)
	sort.Sort(sortedFiles)
	return
}
func SortByModTimeDesc(fis []os.FileInfo) (sortedFiles ByModTimeDesc) {
	sortedFiles = ByModTimeDesc(fis)
	sort.Sort(sortedFiles)
	return
}

func ReadLines(file string) (lines []string, err error) {
	return ReadLinesWithFilter(file, func(line string) (matched bool, broken bool) {
		matched = "" != line
		broken = false
		return
	})

	return
}
func ReadLinesWithFilter(file string, handle func(line string) (matched bool, broken bool)) (lines []string, err error) {
	fd, err := os.Open(file)
	if err != nil {
		return
	}
	defer fd.Close()
	br := bufio.NewReader(fd)

	for {
		line, brErr := br.ReadString('\n')
		line = strings.TrimSpace(line)
		m, b := handle(line)
		if m {
			lines = append(lines, line)
		}

		if b {
			break
		}

		if brErr != nil {
			if brErr == io.EOF {
				break
			}
			err = brErr
			return
		}

	}
	/*for line_end := true; ; {
		line_bytes, is_prefix, err1 := br.ReadLine()
		if err1 != nil {
			if err1 != io.EOF {
				err = err1
			}
			break
		}

		line := string(line_bytes)
		if line_end == false {
			lines[len(lines)-1] += line

		} else {
			lines = append(lines, line)
			line_end = !is_prefix
		}
	}*/

	return
}

func FormatSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%c", float64(b)/float64(div), "KMGTPE"[exp])
}

func Md5(file string) (fileMd5 string, err error) {
	pFile, err := os.Open(file)
	if err != nil {
		err = fmt.Errorf("打开文件失败，filename=%v, err=%v", file, err)
		return
	}
	defer pFile.Close()
	md5h := md5.New()
	io.Copy(md5h, pFile)

	return hex.EncodeToString(md5h.Sum(nil)), nil

}

func RemoveWithGlob(path string) (err error) {
	files, err := filepath.Glob(path)
	if err != nil {
		return
	}
	for _, f := range files {
		if err = os.Remove(f); err != nil {
			return
		}
	}
	return
}
