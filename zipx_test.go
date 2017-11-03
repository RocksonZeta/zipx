package zipx

import (
	"fmt"
	"testing"
)

func TestZipEntryRead(t *testing.T) {
	zip := ZipNew("test.zip", ZIP_DEFAULT_COMPRESSION_LEVEL, ZIP_MODE_READ)
	defer zip.ZipClose()
	ok := zip.ZipEntryOpen("cmp.html")
	if ok != 0 {
		t.Error("ok ", ok)
	}
	bs, freeFn, ok := zip.ZipEntryRead()
	zip.ZipEntryClose()
	fmt.Println(string(bs))
	freeFn()
	fmt.Println(ok)
}
func TestEntryReadCopy(t *testing.T) {
	zip := ZipNew("test.zip", ZIP_DEFAULT_COMPRESSION_LEVEL, ZIP_MODE_READ)
	defer zip.ZipClose()
	bs, ok := zip.Get("中文/中文.txt")
	fmt.Println(string(bs), ok)
}
func TestCreate(t *testing.T) {
	ok := ZipCreate("test.zip", []string{"zipx.go"})
	fmt.Println("ok", ok)
}

func TestZipGetCopy(t *testing.T) {
	bs, exists := ZipGetCopy("test.zip", "中文/中文.txt1")
	fmt.Println(string(bs), exists)
}
func TestZipAppend(t *testing.T) {
	zip := ZipNew("test.zip", ZIP_DEFAULT_COMPRESSION_LEVEL, ZIP_MODE_APPEND)
	defer zip.ZipClose()
	zip.ZipEntryOpen("entry1.txt")
	zip.ZipEntryWrite0()
	zip.ZipEntryClose()
}
