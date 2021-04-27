package fsclient

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func UploadFile (url string, srcPath string) error{

	file, err := os.Open(srcPath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()


	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(filepath.Ext(filepath.Base(url)), filepath.Base(file.Name()))

	if err != nil {
		log.Fatal(err)
		return err
	}

	io.Copy(part, file)
	writer.Close()
	request, err := http.NewRequest("POST", url, body)
	request.SetBasicAuth("admin", "admin")

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

func DownloadFile(url string, destPath string) error{
	request, err := http.NewRequest("GET", url, nil)
	request.SetBasicAuth("admin", "admin")
	client := &http.Client{}
	// get the response from file server
	resp,err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	destPath = filepath.Join(destPath, filepath.Base(url))
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

