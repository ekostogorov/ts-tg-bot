package file_repo

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
)

type Repo struct{}

func New() *Repo {
	return &Repo{}
}

func (r *Repo) Open(filepath string) ([]byte, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			log.Printf("failed to close file: %s", closeErr.Error())
		}
	}()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	return b, nil
}

func (r *Repo) Save(data []byte, filepath string) error {
	if err := ioutil.WriteFile(filepath, data, os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to save file")
	}

	return nil
}
