package util

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Zip ...
func Zip(dst, src string) error {
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	zw := zip.NewWriter(f)
	defer func() {
		if err = zw.Close(); err != nil {
			panic(err)
		}
	}()
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		var (
			fh *zip.FileHeader
			w  io.Writer
			_f *os.File
		)
		if fh, err = zip.FileInfoHeader(info); err != nil {
			return err
		}
		fh.Name = strings.TrimPrefix(path, string(filepath.Separator))
		if info.IsDir() {
			fh.Name += "/"
		}
		if w, err = zw.CreateHeader(fh); err != nil {
			return err
		}
		if !fh.Mode().IsRegular() {
			return nil
		}

		if _f, err = os.Open(path); err != nil {
			return err
		}
		defer func(_f *os.File) {
			err := _f.Close()
			if err != nil {

			}
		}(_f)

		if _, err = io.Copy(w, _f); err != nil {
			return err
		}
		return nil
	})
}

// UnZip ...
func UnZip(dst, src string) error {
	zr, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func(zr *zip.ReadCloser) {
		err := zr.Close()
		if err != nil {

		}
	}(zr)
	if dst != "" {
		if err = os.MkdirAll(dst, 0755); err != nil {
			return err
		}
	}
	unZipFunc := func(f *zip.File) error {
		path := filepath.Join(dst, f.Name)
		if f.FileInfo().IsDir() {
			if err = os.MkdirAll(path, f.Mode()); err != nil {
				return err
			}
			return nil
		}
		var (
			fr io.ReadCloser
			_f *os.File
		)
		if fr, err = f.Open(); err != nil {
			return err
		}
		defer func(fr io.ReadCloser) {
			err := fr.Close()
			if err != nil {

			}
		}(fr)

		if _f, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, f.Mode()); err != nil {
			return err
		}
		defer func(_f *os.File) {
			err := _f.Close()
			if err != nil {

			}
		}(_f)

		if _, err = io.Copy(_f, fr); err != nil {
			return err
		}
		return nil
	}
	for _, file := range zr.File {
		if err = unZipFunc(file); err != nil {
			return err
		}
	}
	return nil
}
