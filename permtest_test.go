package permtest

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func mkTestDir(prefix string, perm os.FileMode, t *testing.T) string {
	d, err := ioutil.TempDir(prefix, "permtest")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if err := os.Remove(d); err != nil {
		t.Fatalf("err: %s", err)
	}
	if err := os.MkdirAll(d, perm); err != nil {
		t.Fatalf("err: %s", err)
	}
	return d
}

func mkTestFile(prefix string, perm os.FileMode, t *testing.T) string {
	f, err := ioutil.TempFile(prefix, "permtest")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer f.Close()
	if err := f.Chmod(perm); err != nil {
		t.Fatalf("err: %s", err)
	}
	return f.Name()
}

func TestPermtest_WriteDir(t *testing.T) {
	path1 := mkTestDir("", 0700, t)
	defer os.Remove(path1)

	path2 := mkTestDir("", 0000, t)
	defer os.Remove(path2)

	if _, err := WriteDir(path1); err != nil {
		t.Fatalf("err: %s", err)
	}

	if _, err := WriteDir(path2); err == nil {
		t.Fatalf("expected permission error")
	}

	path3 := mkTestDir("", 0700, t)
	path4 := mkTestDir(path3, 0000, t)
	defer os.RemoveAll(path3)

	if _, err := WriteDir(path3); err != nil {
		t.Fatalf("err: %s", err)
	}

	last, err := WriteDir(path4)
	if err == nil {
		t.Fatalf("expected permission error")
	}
	if last != path4 {
		t.Fatalf("bad: %s", last)
	}
}

func TestPermtest_WriteFile(t *testing.T) {
	path1 := mkTestFile("", 0600, t)
	defer os.Remove(path1)

	path2 := mkTestFile("", 0000, t)
	defer os.Remove(path2)

	if _, err := WriteFile(path1); err != nil {
		t.Fatalf("err: %s", err)
	}

	if _, err := WriteFile(path2); err == nil {
		t.Fatalf("expected permission error")
	}

	path3 := mkTestFile("", 0600, t)
	defer os.Remove(path3)

	path4 := mkTestFile("", 0000, t)
	defer os.Remove(path4)

	dir := mkTestDir("", 0700, t)
	path5 := filepath.Join(dir, "test_file")
	path6 := filepath.Join(dir, "subdir", "test_file")

	if _, err := WriteFile(path3); err != nil {
		t.Fatalf("err: %s", err)
	}

	if _, err := WriteFile(path4); err == nil {
		t.Fatalf("expected permission error")
	}

	if _, err := WriteFile(path5); err != nil {
		t.Fatalf("err: %s", err)
	}

	if _, err := WriteFile(path6); err == nil {
		t.Fatalf("expected nonexistent dir error")
	}
}
