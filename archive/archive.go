package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

// compress the contents of the given path (src) and write the archive to the given writers
func Compress(src string, logger *logrus.Entry, writers ...io.Writer) error {

	// ensure the src actually exists before trying to tar it
	if _, err := os.Stat(src); err != nil {
		return err
	}
	mw := io.MultiWriter(writers...)
	gzw := gzip.NewWriter(mw)
	defer func() {
		err := gzw.Close()
		if err != nil {
			if logger != nil {
				logger.WithError(err).Error("closing failed for gzip writer")
			}
		}
	}()

	tw := tar.NewWriter(gzw)
	defer func() {
		err := tw.Close()
		if err != nil {
			if logger != nil {
				logger.WithError(err).Error("closing failed for tar writer")
			}
		}
	}()

	// walk path
	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {

		// return on any error
		if err != nil {
			return err
		}

		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		// update the name to correctly reflect the desired destination when untaring
		header.Name = filepath.ToSlash(file)
		if fi.IsDir() {
			header.Name = fmt.Sprintf("%s/", header.Name)
		}

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !fi.IsDir() {
			// open files for taring
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			// copy file data into tar writer
			if _, err = io.Copy(tw, f); err != nil {
				return err
			}
			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			err = f.Close()
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// extract the contents of the given reader (r) to the given destination (dst)
func Extract(dst string, logger *logrus.Entry, r io.Reader) error {

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer func() {
		err := gzr.Close()
		if err != nil {
			if logger != nil {
				logger.WithError(err).Error("closing gzip writer failed")
			}
		}
	}()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()

		switch {
		// if no more files are found return
		case err == io.EOF:
			return nil
		// return any other error
		case err != nil:
			return err
		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}
		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if !os.IsNotExist(err){
					return  err
				}
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			err = f.Close()
			if err != nil {
				return err
			}
		}
	}
}
