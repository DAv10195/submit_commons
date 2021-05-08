package main

import (
	fs_client "github.com/DAv10195/submit_commons/fsclient"
	"io"
	"os"
)

func main(){
	reader,err := os.Open("C:\\ProgramData\\test\\hun.txt")
	if err != nil {
		panic("failed to open yey.txt")
	}
	fsc := fs_client.NewFileServerClient("http://localhost/","admin","admin")
	var ioReader io.Reader = (*os.File)(reader)
	err = fsc.UploadFile("/nikita/", &ioReader,"yey.txt")
	if err != nil {
		panic("error upload")
	}
	fileToCopyTo, err := os.Create("C:\\ProgramData\\test\\yeses.txt")
	if err != nil {
		panic("error open the file downloaded")
	}
	writer := io.MultiWriter(fileToCopyTo)
	err = fsc.DownloadFile("nikita/yey.txt", &writer)
	if err != nil {
		panic("error download")
	}

}
