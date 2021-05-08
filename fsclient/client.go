package fsclient

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	url2 "net/url"
)

type FileServerClient struct {
	adr string
	username string
	password string
}

func NewFileServerClient(adr string, username string, password string) *FileServerClient{
	return &FileServerClient{
		adr: adr,
		username: username,
		password: password,
	}
}

func (fsc *FileServerClient) UploadFile (url string, reader *io.Reader, filename string) error {
	fullURL,err := url2.Parse(fsc.adr)
	if err != nil {
		return err
	}
	fullURL.Path = url
	url = fullURL.String()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)

	if err != nil {
		log.Fatal(err)
		return err
	}

	_,err = io.Copy(part, *reader)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = writer.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}
	request, err := http.NewRequest("POST", url, body)
	request.SetBasicAuth(fsc.username, fsc.password)

	if err != nil {
		log.Fatal(err)
		return err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
		return err
	}
	defer response.Body.Close()

	resp , err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Printf("the response from file server is %s", resp)
	return nil
}

func (fsc *FileServerClient) DownloadFile(url string, writer *io.Writer) error {
	fullURL,err := url2.Parse(fsc.adr)
	if err != nil {
		return err
	}
	fullURL.Path = url
	url = fullURL.String()
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.SetBasicAuth(fsc.username, fsc.password)
	client := &http.Client{}
	// get the response from file server.
	resp,err := client.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		err =  resp.Body.Close()
		if err != nil {
			panic("cannot close body")
			return
		}
	}()

	// copy the body to writer and return it.
	if _, err := io.Copy(*writer, resp.Body); err != nil {
		log.Fatal(err)
	}
	return nil
}

