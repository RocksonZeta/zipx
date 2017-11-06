package zipx

import (
	"os"
	"testing"
)

func TestZipSuite(t *testing.T) {
	zipFile := "test.zip"
	zn := New(zipFile, ZIP_DEFAULT_COMPRESSION_LEVEL, ZIP_MODE_WRITE)
	ok := zn.EntryOpen("test.txt")
	if ok != 0 {
		t.Errorf("ZipEntryOpen error , status=%d", ok)
	}
	ok = zn.EntryFWrite("test.txt")
	if ok != 0 {
		t.Errorf("ZipEntryFWrite error , status=%d", ok)
	}
	zn.EntryClose()
	zn.Close()
	defer os.Remove(zipFile)
	//read test
	zr := New(zipFile, ZIP_DEFAULT_COMPRESSION_LEVEL, ZIP_MODE_READ)
	ok = zr.EntryOpen("test.txt")
	if ok != 0 {
		t.Errorf("ZipEntryOpen error , status=%d", ok)
	}
	bs, ok := zr.EntryReadCopy()
	if ok != 0 {
		t.Errorf("ZipEntryReadCopy error , status=%d", ok)
	}
	if len(bs) <= 0 {
		t.Errorf("ZipEntryReadCopy read no bytes , bs=%s", string(bs))
	}
	ok = zr.EntryClose()
	if ok != 0 {
		t.Errorf("ZipEntryClose error , status=%d", ok)
	}
	zr.Close()

	// //write test
	// zw := ZipNew(zipFile, ZIP_DEFAULT_COMPRESSION_LEVEL, ZIP_MODE_WRITE)
	// defer zw.ZipClose()
	// zw.ZipEntryOpen("new.txt")
	// defer zw.ZipEntryClose()
	// ok = zw.ZipEntryWrite([]byte("new"))
	// if ok != 0 {
	// 	t.Errorf("ZipEntryWrite error , status=%d", ok)
	// }

	//append test

	za := New(zipFile, ZIP_DEFAULT_COMPRESSION_LEVEL, ZIP_MODE_APPEND)
	za.EntryOpen("append.txt")
	ok = za.EntryWrite([]byte("append"))
	if ok != 0 {
		t.Errorf("append ZipEntryWrite  error , status=%d", ok)
	}
	za.EntryClose()
	za.Close()

	zar := New(zipFile, ZIP_DEFAULT_COMPRESSION_LEVEL, ZIP_MODE_READ)
	zar.EntryOpen("append.txt")
	bs, _ = zar.EntryReadCopy()
	if string(bs) != "append" {
		t.Errorf("Zip Append failed")
	}
	zar.EntryClose()

	zar.Close()
	bs, exists := GetCopy("test.zip", "test.txt")
	if !exists || len(bs) <= 0 {
		t.Errorf("ZipGetCopy failed")
	}
}
