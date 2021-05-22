package fsclient

import (
	"bytes"
	"fmt"
	"github.com/DAv10195/submit_commons/encryption"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	url2 "net/url"
)

type FileServerClient struct {
	adr string
	username string
	password string
	logger   *logrus.Entry
	encryption encryption.Encryption
}

func NewFileServerClient(adr string, username string, password string, logger *logrus.Entry, encryption encryption.Encryption) *FileServerClient{
	encryptedPass, err := encryption.Encrypt(password)
	if err != nil {
		panic("error in encryption")
	}
	return &FileServerClient{
		adr: adr,
		username: username,
		password: encryptedPass,
		logger:   logger,
		encryption: encryption,
	}
}

func (fsc *FileServerClient) UploadFile (url string, reader io.Reader, filename string) error {
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
		return err
	}
	_, err = io.Copy(part, reader)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		fsc.logger.WithError(err).Error("error closing multi part writer")
		return err
	}
	request, err := http.NewRequest(http.MethodPost, url, body)
	decryptedPass, err := fsc.encryption.Decrypt(fsc.password)
	if err != nil {
		return err
	}
	request.SetBasicAuth(fsc.username, decryptedPass)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return err
	}
	defer func() {
		err =  response.Body.Close()
		if err != nil {
			fsc.logger.WithError(err).Error("error closing the resp body while uploading")
			return
		}
	}()
	if response.StatusCode != http.StatusAccepted {
		fsc.logger.Error(fmt.Printf("Upload request failed for file %s. status code is %d", url ,response.StatusCode))
	}
	_ , err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return nil
}

func (fsc *FileServerClient) DownloadFile(url string, writer io.Writer) error {
	fullURL,err := url2.Parse(fsc.adr)
	if err != nil {
		return err
	}
	fullURL.Path = url
	url = fullURL.String()
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	decryptedPass, err := fsc.encryption.Decrypt(fsc.password)
	if err != nil {
		return err
	}
	request.SetBasicAuth(fsc.username, decryptedPass)
	client := &http.Client{}
	// get the response from file server.
	resp,err := client.Do(request)
	if err != nil {
		return err
	}

	defer func() {
		err =  resp.Body.Close()
		if err != nil {
			fsc.logger.WithError(err).Error("error closing the resp body while downloading")
			return
		}
	}()
	if resp.StatusCode != http.StatusOK {
		fsc.logger.Error(fmt.Printf("Downloading request for file %s failed. status code is %d",url ,resp.StatusCode))
	}

	// copy the body to writer and return it.
	if _, err = io.Copy(writer, resp.Body); err != nil {
		return err
	}
	return nil
}

