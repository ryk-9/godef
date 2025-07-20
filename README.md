# godef cli dictionary and thesaurus

## install

```bash
git clone [https://github.com/ryk-9/godef.git](https://github.com/ryk-9/godef.git)
cd godef
go mod init
go mod tidy
go build -o godef
sudo mv godef /usr/local/bin
```

should now be runnable with syntax `godef <word>`

outputs full thesaurus and dictionary entry.
