package fsclient

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/DAv10195/submit_commons/encryption"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	url2 "net/url"
	"os"
	"path/filepath"
)

type FileServerClient struct {
	adr *url2.URL
	username string
	password string
	logger   *logrus.Entry
	encryption encryption.Encryption
	tlsConf	 *tls.Config
}

// create a new fs client
func NewFileServerClient(adr string, username string, password string, logger *logrus.Entry, encryption encryption.Encryption, tlsConf *tls.Config) (*FileServerClient, error){
	if username == "" {
		return nil, errors.New("failed to initialize file server client, username was not given")
	}
	if  password == ""{
		return nil, errors.New("failed to initialize file server client, password was not given")
	}
	if adr == "" {
		return nil, errors.New("failed to initialize file server client, fs address was not given")
	}
	if encryption == nil {
		return nil, errors.New("failed to initialize file server client, encryption was not initialized")
	}
	adrUrl,err := url2.Parse(adr)
	if err != nil {
		return nil, err
	}
	return &FileServerClient{
		adr:        adrUrl,
		username:   username,
		password:   password,
		logger:     logger,
		encryption: encryption,
		tlsConf: tlsConf,
	},nil
}

// upload a file to the given path (url). If isFolder is true then the given files will all be treated
// as archives which will be extracted in the file server
func (fsc *FileServerClient) UploadFile (url string, isFolder bool, reader ...*os.File) error {

	fullURL := fsc.adr
	fullURL.Path = url
	url = fullURL.String()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	var part io.Writer
	var err error
	for i:=0;i<len(reader);i++ {
		part, err = writer.CreateFormFile("file", filepath.Base(reader[i].Name()))
		if err != nil {
			err2 := writer.Close()
			if err2 != nil {
				if fsc.logger != nil {
					fsc.logger.WithError(err).Error("error closing the multipart writer while uploading")
				}
			}
			return err
		}
		_, err = io.Copy(part, reader[i])
		if err != nil {
			err2 := writer.Close()
			if err2 != nil {
				if fsc.logger != nil {
					fsc.logger.WithError(err).Error("error closing the multipart writer while uploading")
				}
			}
			return err
		}
	}
	err = writer.Close()
	if err != nil {
		if fsc.logger != nil {
			fsc.logger.WithError(err).Error("error closing the multipart writer while uploading")
		}
	}
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}
	if isFolder {
		q := request.URL.Query()
		q.Add("isFolder", "true")
		request.URL.RawQuery = q.Encode()
	}
	decryptedPass, err := fsc.encryption.Decrypt(fsc.password)
	if err != nil {
		return err
	}
	request.SetBasicAuth(fsc.username, decryptedPass)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	if fsc.tlsConf != nil {
		client.Transport = &http.Transport{TLSClientConfig: fsc.tlsConf}
	}
	response, err := client.Do(request)

	if err != nil {
		return err
	}
	defer func() {
		err =  response.Body.Close()
		if err != nil {
			if fsc.logger != nil {
				fsc.logger.WithError(err).Error("error closing the resp body while uploading")
			}
		}
	}()
	if response.StatusCode != http.StatusAccepted {
		msg := fmt.Sprintf("Upload request failed for file %s. status code is %d", url ,response.StatusCode)
		return errors.New(msg)
	}
	return nil
}

// download the resource at the given path (url) and write it to given writer (e.g. a file)
func (fsc *FileServerClient) DownloadFile(url string, writer io.Writer) (http.Header, error) {

	fullURL := fsc.adr
	fullURL.Path = url
	url = fullURL.String()
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	decryptedPass, err := fsc.encryption.Decrypt(fsc.password)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(fsc.username, decryptedPass)
	client := &http.Client{}
	if fsc.tlsConf != nil {
		client.Transport = &http.Transport{TLSClientConfig: fsc.tlsConf}
	}
	// get the response from file server.
	resp,err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		err =  resp.Body.Close()
		if err != nil {
			if fsc.logger != nil {
				fsc.logger.WithError(err).Error("error closing the resp body while downloading")
			}
		}
	}()
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Downloading request for file %s failed. status code is %d", url ,resp.StatusCode)
		return nil, errors.New(msg)
	}

	// copy the body to writer and return it.
	if _, err = io.Copy(writer, resp.Body); err != nil {
		return nil, err
	}
	return resp.Header, nil
}

// upload the given bytes to the given path (url)
func (fsc *FileServerClient) UploadTextToFS(url string, data []byte) error {
	fullURL := fsc.adr
	fullURL.Path = url
	url = fullURL.String()

	body := bytes.NewBuffer(data)
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}
	decryptedPass, err := fsc.encryption.Decrypt(fsc.password)
	if err != nil {
		return err
	}
	request.SetBasicAuth(fsc.username, decryptedPass)
	client := &http.Client{}
	if fsc.tlsConf != nil {
		client.Transport = &http.Transport{TLSClientConfig: fsc.tlsConf}
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		err =  response.Body.Close()
		if err != nil {
			if fsc.logger != nil {
				fsc.logger.WithError(err).Error("error closing the resp body while uploading")
			}
		}
	}()
	if response.StatusCode != http.StatusAccepted {
		msg := fmt.Sprintf("Upload request failed for file %s. status code is %d", url ,response.StatusCode)
		return errors.New(msg)
	}

	return nil

}

// delete the resource in the given path
func (fsc *FileServerClient) Delete(path string) error {
	url := fsc.adr
	url.Path = path
	password, err := fsc.encryption.Decrypt(fsc.password)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodDelete, url.String(), nil)
	if err != nil {
		return err
	}
	request.SetBasicAuth(fsc.username, password)
	client := &http.Client{}
	if fsc.tlsConf != nil {
		client.Transport = &http.Transport{TLSClientConfig: fsc.tlsConf}
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		err =  response.Body.Close()
		if err != nil {
			if fsc.logger != nil {
				fsc.logger.WithError(err).Error("error closing the resp body while deleting")
			}
		}
	}()
	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("error deleting path (%s). Status code is %d", path ,response.StatusCode)
	}
	return nil
}

// forward the data in the given reader (b, e.g. an http request body) to the given path with the given content type
func (fsc *FileServerClient) ForwardBody(path, contentType string, b io.Reader) error {
	url := fsc.adr
	url.Path = path
	password, err := fsc.encryption.Decrypt(fsc.password)
	if err != nil {
		return err
	}
	r, err := http.NewRequest(http.MethodPost, url.String(), b)
	if err != nil {
		return err
	}
	r.SetBasicAuth(fsc.username, password)
	r.Header.Add("Content-Type", contentType)
	client := &http.Client{}
	if fsc.tlsConf != nil {
		client.Transport = &http.Transport{TLSClientConfig: fsc.tlsConf}
	}
	response, err := client.Do(r)
	defer func() {
		err =  response.Body.Close()
		if err != nil {
			if fsc.logger != nil {
				fsc.logger.WithError(err).Error("error closing the resp body while forwarding the given body")
			}
		}
	}()
	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("error forwarding body (%s). Status code is %d", path, response.StatusCode)
	}
	return nil
}
