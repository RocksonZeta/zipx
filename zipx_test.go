package zipx

import (
	"os"
	"testing"
)

func TestZipSuite(t *testing.T) {
	zipFile := "test.zip"
	zn := ZipNew(zipFile, ZIP_DEFAULT_COMPRESSION_LEVEL, ZIP_MODE_WRITE)
	ok := zn.ZipEntryOpen("test.txt")
	if ok != 0 {
		t.Errorf("ZipEntryOpen error , status=%d", ok)
	}
	ok = zn.ZipEntryFWrite("test.txt")
	if ok != 0 {
		t.Errorf("ZipEntryFWrite error , status=%d", ok)
	}
	zn.ZipEntryClose()
	zn.ZipClose()
	defer os.Remove(zipFile)
	//read test
	zr := ZipNew(zipFile, ZIP_DEFAULT_COMPRESSION_LEVEL, ZIP_MODE_READ)
	ok = zr.ZipEntryOpen("test.txt")
	if ok != 0 {
		t.Errorf("ZipEntryOpen error , status=%d", ok)
	}
	bs, ok := zr.ZipEntryReadCopy()
	if ok != 0 {
		t.Errorf("ZipEntryReadCopy error , status=%d", ok)
	}
	if len(bs) <= 0 {
		t.Errorf("ZipEntryReadCopy read no bytes , bs=%s", string(bs))
	}
	ok = zr.ZipEntryClose()
	if ok != 0 {
		t.Errorf("ZipEntryClose error , status=%d", ok)
	}
	zr.ZipClose()

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

	za := ZipNew(zipFile, ZIP_DEFAULT_COMPRESSION_LEVEL, ZIP_MODE_APPEND)
	za.ZipEntryOpen("append.txt")
	ok = za.ZipEntryWrite([]byte("append"))
	if ok != 0 {
		t.Errorf("append ZipEntryWrite  error , status=%d", ok)
	}
	za.ZipEntryClose()
	za.ZipClose()

	zar := ZipNew(zipFile, ZIP_DEFAULT_COMPRESSION_LEVEL, ZIP_MODE_READ)
	zar.ZipEntryOpen("append.txt")
	bs, _ = zar.ZipEntryReadCopy()
	if string(bs) != "append" {
		t.Errorf("Zip Append failed")
	}
	zar.ZipEntryClose()

	zar.ZipClose()
	bs, exists := ZipGetCopy("test.zip", "test.txt")
	if !exists || len(bs) <= 0 {
		t.Errorf("ZipGetCopy failed")
	}
}
