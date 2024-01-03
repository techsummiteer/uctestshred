package main

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestShredText(t *testing.T) {
	testFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(testFile.Name())

	_, err = testFile.WriteString("temporary file content")
	if err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	testFile.Close()

	err = Shred(testFile.Name())
	if err != nil {
		t.Errorf("Shred failed: %v", err)
	}

	if _, err := os.Stat(testFile.Name()); !os.IsNotExist(err) {
		t.Errorf("File was not deleted")
	}
}
func TestShredBinFileUnder1Block(t *testing.T) {
	testFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(testFile.Name())

	cmd := exec.Command("dd", "if=/dev/zero", "of="+testFile.Name(), "bs=1k", "count=127")
	err = cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	err = Shred(testFile.Name())
	if err != nil {
		t.Errorf("Shred failed: %v", err)
	}

	if _, err := os.Stat(testFile.Name()); !os.IsNotExist(err) {
		t.Errorf("File was not deleted")
	}
}

func TestShredBinFileOver1Block(t *testing.T) {
	testFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(testFile.Name())

	cmd := exec.Command("dd", "if=/dev/zero", "of="+testFile.Name(), "bs=1k", "count=512")
	err = cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
	err = Shred(testFile.Name())
	if err != nil {
		t.Errorf("Shred failed: %v", err)
	}

	if _, err := os.Stat(testFile.Name()); !os.IsNotExist(err) {
		t.Errorf("File was not deleted")
	}
}

func TestShredFileWithUnicode(t *testing.T) {
	testFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file with Unicode: %v", err)
	}
	_, err = testFile.WriteString("世界こんにちは")
	if err != nil {
		t.Fatalf("Failed to write Unicode to temporary file: %v", err)
	}
	testFileName := testFile.Name()
	testFile.Close()
	defer os.Remove(testFileName)

	err = Shred(testFileName)
	if err != nil {
		t.Errorf("Shred failed on file with Unicode : %v", err)
	}

	if _, err := os.Stat(testFileName); !os.IsNotExist(err) {
		t.Errorf("File was not deleted")
	}
}

func TestShredNonExistentFile(t *testing.T) {
	err := Shred("missingfile.txt")
	if err == nil {
		t.Errorf("Expected an error for non-existent file, but got none")
	}
}

func TestShredReadOnlyFile(t *testing.T) {
	testFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	testFileName := testFile.Name()
	testFile.Close()
	os.Chmod(testFileName, 0400)
	defer os.Remove(testFileName)

	err = Shred(testFileName)
	if err == nil {
		t.Errorf("Expected an error for read-only file, but got none")
	}
}
func TestShredFasterThanUnixCmd(t *testing.T) {
	testFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(testFile.Name())

	cmd := exec.Command("dd", "if=/dev/zero", "of="+testFile.Name(), "bs=1M", "count=1024")
	err = cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
	start_time := time.Now()
	err = Shred(testFile.Name())
	if err != nil {
		t.Errorf("Shred failed: %v", err)
	}
	durationMyShred := time.Since(start_time)

	cmd = exec.Command("dd", "if=/dev/zero", "of="+testFile.Name(), "bs=1M", "count=1024")
	err = cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
	cmd = exec.Command("shred", testFile.Name())
	start_time = time.Now()
	err = cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
	durationUnixShred := time.Since(start_time)
	if durationMyShred > durationUnixShred {
		t.Errorf("My Shred is slower than unix shred")
	}

}

func TestShredEmptyFile(t *testing.T) {
	testFile, err := os.CreateTemp("", "emptytestfile")
	if err != nil {
		t.Fatalf("Failed to create empty temporary file: %v", err)
	}
	testFileName := testFile.Name()
	testFile.Close()
	defer os.Remove(testFileName)

	err = Shred(testFileName)
	if err != nil {
		t.Errorf("Shred failed on empty file: %v", err)
	}

	if _, err := os.Stat(testFileName); !os.IsNotExist(err) {
		t.Errorf("File was not deleted")
	}
}

func TestShredDirectoryPath(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	err = Shred(tempDir)
	if err == nil {
		t.Errorf("Expected an error for directory path, but got none")
	}
}
