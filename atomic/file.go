package atomic

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
)

func WriteFileReader(filename string, r io.Reader, perm os.FileMode) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return WriteFile(filename, b, perm)
}

// WriteFileAtomic writes the data to a temp file and atomically move if everything else succeeds.
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	dir, name := path.Split(filename)
	f, err := ioutil.TempFile(dir, name)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err == nil {
		err = f.Sync()
	}
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if permErr := os.Chmod(f.Name(), perm); err == nil {
		err = permErr
	}
	if err == nil {
		err = os.Rename(f.Name(), filename)
	}
	// Any err should result in full cleanup.
	if err != nil {
		os.Remove(f.Name())
	}
	return err
}

// MoveFile copies file from src to dst. It only works on single files and
// fails, if dst already exists or if dirname(dst) does not exist. Works across
// partitions.
func MoveFile(src, dst string) error {
	fi, err := os.Stat(dst)
	if err == nil || (fi != nil && fi.IsDir()) {
		return fmt.Errorf("dst already exists: %s", dst)
	}
	fi, err = os.Stat(path.Dir(dst))
	if os.IsNotExist(err) {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("dst directory does not exist")
	}
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()
	dsttmp := fmt.Sprintf("%s-%d", rand.Intn(99999999))
	g, err := os.Create(dsttmp)
	if err != nil {
		return err
	}
	if _, err := io.Copy(g, f); err != nil {
		return err
	}
	if err := g.Close(); err != nil {
		return err
	}
	err = os.Rename(dsttmp, dst)
	if err != nil {
		return err
	}
	return os.Remove(src)
}
