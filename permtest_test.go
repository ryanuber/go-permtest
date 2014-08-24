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

func TestPermtest_Write(t *testing.T) {
	// Write should handle directories properly
	path1 := mkTestDir("", 0700, t)
	defer os.Remove(path1)

	path2 := mkTestDir("", 0000, t)
	defer os.Remove(path2)

	if _, err := Write(path1); err != nil {
		t.Fatalf("err: %s", err)
	}

	if _, err := Write(path2); err == nil {
		t.Fatalf("expected permission error")
	}

	// Write should handle files properly
	path3 := mkTestFile("", 0600, t)
	defer os.Remove(path1)

	path4 := mkTestFile("", 0000, t)
	defer os.Remove(path2)

	if _, err := Write(path3); err != nil {
		t.Fatalf("err: %s", err)
	}

	if _, err := Write(path4); err == nil {
		t.Fatalf("expected permission error")
	}
}

func TestPermtest_WriteFile(t *testing.T) {
	path1 := mkTestFile("", 0600, t)
	defer os.Remove(path1)

	path2 := mkTestFile("", 0000, t)
	defer os.Remove(path2)

	dir := mkTestDir("", 0700, t)
	path3 := filepath.Join(dir, "test_file")
	path4 := filepath.Join(dir, "subdir", "test_file")

	if _, err := WriteFile(path1); err != nil {
		t.Fatalf("err: %s", err)
	}

	if _, err := WriteFile(path2); err == nil {
		t.Fatalf("expected permission error")
	}

	if _, err := WriteFile(path3); err != nil {
		t.Fatalf("err: %s", err)
	}

	if _, err := WriteFile(path4); err == nil {
		t.Fatalf("expected nonexistent dir error")
	}
}

func TestPermtest_WriteDir(t *testing.T) {
	d1 := mkTestDir("", 0700, t)
	d2 := mkTestDir(d1, 0000, t)
	defer os.RemoveAll(d1)

	if _, err := WriteDir(d1); err != nil {
		t.Fatalf("err: %s", err)
	}

	last, err := WriteDir(d2)
	if err == nil {
		t.Fatalf("expected permission error")
	}
	if last != d2 {
		t.Fatalf("bad: %s", last)
	}
}
