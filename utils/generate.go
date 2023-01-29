package utils

import (
	"bytes"
	"embed"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
)

// 0: InjectHere
// 1: InjectInterfaceHere
// 2: holding
type InjectFile struct {
	Target string
	Type   uint
}

func Generate(files map[string]string, fs embed.FS, m any) error {
	for tpl, target := range files {
		targetFileDir := filepath.Dir(target)
		if err := os.MkdirAll(targetFileDir, os.ModePerm); err != nil {
			return errors.Wrapf(err, "mkdir %s failed", targetFileDir)
		}

		// get tplFilePath file contents
		b, err := fs.ReadFile(tpl)
		if err != nil {
			return errors.Wrapf(err, "read template file(%s) failed", tpl)
		}

		t, err := template.New("file").Parse(string(b))
		if err != nil {
			return errors.Wrapf(err, "parse template(%s) failed", tpl)
		}

		fd, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "open go file(%s) failed", target)
		}
		defer fd.Close()

		if err := t.Execute(fd, m); err != nil {
			return errors.Wrapf(err, "generate file failed, target: %s", target)
		}
	}
	return nil
}

type InjectInterface interface {
	InjectReset()
	InjectSet(uint, string)
}

func Inject(files map[string]InjectFile, fs embed.FS, m any) error {
	if _, ok := m.(InjectInterface); !ok {
		return errors.New("not inject interface")
	}

	for tpl, target := range files {
		m.(InjectInterface).InjectReset()

		b, err := fs.ReadFile(tpl)
		if err != nil {
			return errors.Wrapf(err, "read template file(%s) failed", tpl)
		}

		t, err := template.New(tpl).Parse(string(b))
		if err != nil {
			return errors.Wrapf(err, "parse template(%s) failed", tpl)
		}

		parsedData := bytes.NewBuffer(nil)
		if err := t.Execute(parsedData, m); err != nil {
			return errors.Wrapf(err, "inject data failed")
		}
		// set data
		m.(InjectInterface).InjectSet(target.Type, parsedData.String())

		b, err = os.ReadFile(target.Target)
		if err != nil {
			return err
		}

		t, err = template.New(target.Target).Parse(string(b))
		if err != nil {
			return errors.Wrapf(err, "parse template(%s) failed", target.Target)
		}

		fd, err := os.OpenFile(target.Target, os.O_RDWR, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "open go file(%s) failed", target.Target)
		}

		if err := t.Execute(fd, m); err != nil {
			fd.Close()
			return errors.Wrapf(err, "inject file failed, %s", target.Target)
		}
		fd.Close()
	}
	return nil
}
