package util

import (
	"io"
	"os"
)

func CopyFile(f, t string) error {
	from, err := os.Open(f)
	if err != nil {
		return err
	}
	defer from.Close()
	to, err := os.OpenFile(t, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer to.Close()
	_, err = io.Copy(to, from)
	return err
}
