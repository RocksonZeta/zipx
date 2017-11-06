# zipx
golang wrap for github.com/kuba--/zip


## installation
```
go get github.com/RocksonZeta/zipx
```
## Example
```go
//Create zip file
zn := ZipNew(zipFile, ZIP_DEFAULT_COMPRESSION_LEVEL, ZIP_MODE_WRITE)
zn.ZipEntryOpen("test.txt")
zn.ZipEntryWrite([]byte("this is content"))
//zn.ZipEntryFWrite("test.txt")
zn.ZipEntryClose()
zn.ZipClose()

//append file to zip
za := ZipNew(zipFile, ZIP_DEFAULT_COMPRESSION_LEVEL, ZIP_MODE_APPEND)
za.ZipEntryOpen("append.txt")
za.ZipEntryWrite([]byte("append"))
za.ZipEntryClose()
za.ZipClose()


//read a file
zr := ZipNew(zipFile, ZIP_DEFAULT_COMPRESSION_LEVEL, ZIP_MODE_READ)
zr.ZipEntryOpen("dir/test.txt")
bs, ok := zr.ZipEntryReadCopy()
zr.ZipEntryClose()
zr.ZipClose()

//read a file directly
bs, exists := ZipGetCopy("test.zip", "dir/test.txt")
fmt.Println(string(bs))
```
