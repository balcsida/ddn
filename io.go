package main

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

func generateProps(filename string) error {
	filename, conf := generateInteractive(filename)

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("couldn't create file: %s", err.Error())
	}
	defer file.Close()

	prop := `vendor="{{.Vendor}}"
version="{{.Version}}"
executable="{{.Exec}}"
dbport="{{.DBPort}}"
dbAddress="{{.DBAddress}}"
connectorPort="{{.ConnectorPort}}"
username="{{.User}}"
password="{{.Password}}"
masterAddress="{{.MasterAddress}}"
`

	if conf.SID != "" {
		prop += "oracle-sid=\"{{.SID}}\"\n"
	}

	if conf.DefaultTablespace != "" {
		prop += "default-tablespace=\"{{.DefaultTablespace}}\"\n"
	}

	tmpl, err := template.New("props").Parse(prop)
	if err != nil {
		return fmt.Errorf("couldn't parse template: %s", err.Error())
	}

	err = tmpl.Execute(file, conf)
	if err != nil {
		return fmt.Errorf("couldn't execute template: %s", err.Error())
	}

	file.Sync()

	return nil
}

func unzip(filepath string) ([]string, error) {
	r, err := zip.OpenReader(filepath)
	if err != nil {
		return nil, fmt.Errorf("creating zip reader failed: %s", err.Error())
	}
	defer r.Close()

	var files []string
	for _, f := range r.File {
		name, err := unzipFile(f)
		if err != nil {
			return nil, fmt.Errorf("extracting zip file failed: %s", err.Error())
		}

		files = append(files, name)
	}

	return files, nil
}

func unzipFile(f *zip.File) (string, error) {
	src, err := f.Open()
	if err != nil {
		return "", fmt.Errorf("opening zipfile failed: %s", err.Error())
	}
	defer src.Close()

	dst, err := os.OpenFile(f.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return "", fmt.Errorf("opening destination file failed: %s", err.Error())
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", fmt.Errorf("copying from archive failed: %s", err.Error())
	}

	return f.Name, nil
}

func ungzip(path string) ([]string, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening gzipfile failed: %s", err.Error())
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return nil, fmt.Errorf("creating gzip reader failed: %s", err.Error())
	}
	defer archive.Close()

	dst, err := os.Create(archive.Header.Name)
	if err != nil {
		return nil, fmt.Errorf("could not create output file: %s", err.Error())
	}
	defer dst.Close()

	_, err = io.Copy(dst, archive)

	e := []string{dst.Name()}

	return e, nil
}

func isArchive(path string) bool {
	switch filepath.Ext(path) {
	case ".zip", ".tar", ".gz":
		return true
	}

	return false
}
