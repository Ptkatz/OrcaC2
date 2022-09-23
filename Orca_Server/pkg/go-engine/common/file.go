package common

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func LoadJson(filename string, conf interface{}) error {
	err := loadJson(filename, conf)
	if err != nil {
		err := loadJson(filename+".back", conf)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadJson(filename string, conf interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(conf)
	if err != nil {
		return err
	}
	return nil
}

func SaveJson(filename string, conf interface{}) error {
	err1 := saveJson(filename, conf)
	err2 := saveJson(filename+".back", conf)
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}

func saveJson(filename string, conf interface{}) error {
	str, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}
	jsonFile, err := os.OpenFile(filename,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	jsonFile.Write(str)
	jsonFile.Close()
	return nil
}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func FileMd5(filename string) (string, error) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	md5 := GetMd5String(string(data))
	return md5, nil
}

func FileReplace(filename string, from string, to string) error {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	str := string(data)
	str = strings.Replace(str, from, to, -1)

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	out.WriteString(str)
	return nil
}

func FileFind(filename string, dst string) int {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0
	}
	str := string(data)

	n := 0
	for _, str := range strings.Split(str, "\n") {
		if strings.Contains(str, dst) {
			n++
		}
	}
	return n
}

func FileLineCount(filename string) int {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0
	}

	lineSep := []byte{'\n'}

	return bytes.Count(data, lineSep) + 1
}

func IsSymlink(filename string) bool {
	fi, err := os.Lstat(filename)
	if err != nil {
		return false
	}
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		return true
	} else {
		return false
	}
}

// symwalkFunc calls the provided WalkFn for regular files.
// However, when it encounters a symbolic link, it resolves the link fully using the
// filepath.EvalSymlinks function and recursively calls symwalk.Walk on the resolved path.
// This ensures that unlink filepath.Walk, traversal does not stop at symbolic links.
//
// Note that symwalk.Walk does not terminate if there are any non-terminating loops in
// the file structure.
func walk(filename string, linkDirname string, walkFn filepath.WalkFunc) error {
	symWalkFunc := func(path string, info os.FileInfo, err error) error {

		if fname, err := filepath.Rel(filename, path); err == nil {
			path = filepath.Join(linkDirname, fname)
		} else {
			return err
		}

		if err == nil && info.Mode()&os.ModeSymlink == os.ModeSymlink {
			finalPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}
			info, err := os.Lstat(finalPath)
			if err != nil {
				return walkFn(path, info, err)
			}
			if info.IsDir() {
				return walk(finalPath, path, walkFn)
			}
		}

		return walkFn(path, info, err)
	}
	return filepath.Walk(filename, symWalkFunc)
}

// Walk extends filepath.Walk to also follow symlinks
func Walk(path string, walkFn filepath.WalkFunc) error {
	return walk(path, path, walkFn)
}
