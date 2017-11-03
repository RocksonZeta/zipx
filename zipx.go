package zipx

/**
install :
gcc -c zip.c
ar cru libzip.a zip.o
*/

/*
#cgo CFLAGS: -Izip/src
#cgo LDFLAGS: -L zip/src -lzip

#include "zip.h"
#include <stdlib.h>
#include <stdio.h>

struct entry_data {
	void* p;
	size_t len;
	int status;
	int exists;
};

struct entry_data zip_get(char * file ,char* entry) {
	// void *buf = NULL;
	struct entry_data buf;
	buf.len = 0;
	buf.status = 0;
	buf.exists = 1;
	buf.p = NULL;
    // size_t bufsize;
	struct zip_t *zip = zip_open(file, 0, 'r');
	if(NULL == zip) {
		buf.status=-1;
		return buf;
	}
	buf.status = zip_entry_open(zip, entry);
    if(0 == buf.status){
		buf.status = zip_entry_read(zip, &buf.p, &buf.len);
		if(0 == buf.status){
			zip_entry_close(zip);
		}
	}else {
		buf.exists = 0;
		buf.status = 1;
	}
	zip_close(zip);
	return buf;
}

*/
import "C"
import (
	"unsafe"
)

const (
	ZIP_DEFAULT_COMPRESSION_LEVEL = 6
	ZIP_MODE_READ                 = 'r'
	ZIP_MODE_WRITE                = 'w'
	ZIP_MODE_APPEND               = 'a'
)

type Zip struct {
	zip     *C.struct_zip_t
	ZipName string
	Level   int
	Mode    byte
}

//ZipOen extern struct zip_t *zip_open(const char *zipname, int level, char mode);
func ZipNew(zipName string, level int, mode byte) *Zip {
	z := &Zip{
		ZipName: zipName,
		Level:   level,
		Mode:    mode,
	}
	c_zipName := C.CString(zipName)
	defer C.free(unsafe.Pointer(c_zipName))
	z.zip = C.zip_open(c_zipName, C.int(level), C.char(mode))
	return z
}

//ZipClose extern void zip_close(struct zip_t *zip);
func (z *Zip) ZipClose() {
	C.zip_close(z.zip)
}

//ZipEntryOpen extern int zip_entry_open(struct zip_t *zip, const char *entryname);
func (z *Zip) ZipEntryOpen(entryName string) int {
	c_entryName := C.CString(entryName)
	defer C.free(unsafe.Pointer(c_entryName))
	return int(C.zip_entry_open(z.zip, c_entryName))
}

//extern int zip_entry_close(struct zip_t *zip);
func (z *Zip) ZipEntryClose() int {
	return int(C.zip_entry_close(z.zip))
}

//extern int zip_entry_write(struct zip_t *zip, const void *buf, size_t bufsize);
func (z *Zip) ZipEntryWrite(buf []byte) int {
	if 0 == len(buf) {
		return int(C.zip_entry_write(z.zip, unsafe.Pointer(nil), C.size_t(0)))
	} else {
		return int(C.zip_entry_write(z.zip, unsafe.Pointer(&buf[0]), C.size_t(len(buf))))
	}
}

//extern int zip_entry_fwrite(struct zip_t *zip, const char *filename);
func (z *Zip) ZipEntryFWrite(fileName string) int {
	c_fileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(c_fileName))
	return int(C.zip_entry_fwrite(z.zip, c_fileName))
}

//extern int zip_entry_read(struct zip_t *zip, void **buf, size_t *bufsize);
func (z *Zip) ZipEntryRead() ([]byte, func(), int) {
	var p unsafe.Pointer
	var len C.size_t
	ok := int(C.zip_entry_read(z.zip, &p, &len))
	bs := C.GoBytes(p, C.int(len))
	freeFn := func() {
		C.free(p)
	}
	return bs, freeFn, ok
}
func (z *Zip) ZipEntryReadCopy() ([]byte, int) {
	var p unsafe.Pointer
	var size C.size_t
	ok := int(C.zip_entry_read(z.zip, &p, &size))
	if ok == 0 {
		bs := C.GoBytes(p, C.int(size))
		copyBs := make([]byte, len(bs))
		copy(copyBs, bs)
		C.free(p)
		return copyBs, ok
	}
	return nil, ok
}

//extern int zip_entry_fread(struct zip_t *zip, const char *filename);
func (z *Zip) ZipEntryFRead(fileName string) int {
	c_fileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(c_fileName))
	return int(C.zip_entry_fread(z.zip, c_fileName))
}

func (z *Zip) Get(entryName string) ([]byte, int) {
	ok := z.ZipEntryOpen(entryName)
	if 0 != ok {
		return nil, ok
	}
	defer z.ZipEntryClose()
	return z.ZipEntryReadCopy()
}

//extern int zip_create(const char *zipname, const char *filenames[], size_t len);
func ZipCreate(zipName string, fileNames []string) int {
	c_zipName := C.CString(zipName)
	C.free(unsafe.Pointer(c_zipName))
	length := len(fileNames)
	fs := make([]*C.char, length, length)
	for i, f := range fileNames {
		fs[i] = C.CString(f)
		defer C.free(unsafe.Pointer(fs[i]))
	}
	return int(C.zip_create(c_zipName, (**C.char)(&fs[0]), C.size_t(length)))

}

func ZipGetCopy(name, entry string) ([]byte, bool) {
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))
	c_entry := C.CString(entry)
	defer C.free(unsafe.Pointer(c_entry))
	var buf C.struct_entry_data = C.zip_get(c_name, c_entry)
	bs := C.GoBytes(unsafe.Pointer(buf.p), C.int(buf.len))
	if int(buf.status) == 0 {
		copyBs := make([]byte, len(bs))
		copy(copyBs, bs)
		C.free(unsafe.Pointer(buf.p))
		return copyBs, 1 == int(buf.exists)
	}
	return nil, 1 == int(buf.exists)
}
func ZipGet(name, entry string) ([]byte, bool, func()) {
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))
	c_entry := C.CString(entry)
	defer C.free(unsafe.Pointer(c_entry))
	var buf C.struct_entry_data = C.zip_get(c_name, c_entry)
	bs := C.GoBytes(unsafe.Pointer(buf.p), C.int(buf.len))
	freeFn := func() {
		if int(buf.status) == 0 {
			C.free(unsafe.Pointer(buf.p))
		}
	}
	return bs, (1 == int(buf.exists)), freeFn
}

//unimplement
/*
extern int zip_entry_extract(struct zip_t *zip,
                             size_t (*on_extract)(void *arg,
                                                  unsigned long long offset,
                                                  const void *data,
                                                  size_t size),
							 void *arg);
*/

/*
extern int zip_extract(const char *zipname, const char *dir,
                       int (*on_extract_entry)(const char *filename, void *arg),
					   void *arg);
*/
