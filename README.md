# zipx
golang wrap for github.com/kuba--/zip


## installation
```
go get github.com/RocksonZeta/zipx^
```
## Example
```go
//Create zip file
zn := New(zipFile, ZIP_DEFAULT_COMPRESSION_LEVEL, ZIP_MODE_WRITE)
zn.EntryOpen("test.txt")
zn.EntryWrite([]byte("this is content"))
//zn.ZipEntryFWrite("test.txt")
zn.EntryClose()
zn.Close()

//append file to zip
za := New(zipFile, ZIP_DEFAULT_COMPRESSION_LEVEL, ZIP_MODE_APPEND)
za.EntryOpen("append.txt")
za.EntryWrite([]byte("append"))
za.EntryClose()
za.Close()


//read a file
zr := New(zipFile, ZIP_DEFAULT_COMPRESSION_LEVEL, ZIP_MODE_READ)
zr.EntryOpen("dir/test.txt")
bs, ok := zr.EntryReadCopy()
zr.EntryClose()
zr.Close()

//read a file directly
bs, exists := GetCopy("test.zip", "dir/test.txt")
fmt.Println(string(bs))
```
