package samplesrc

import (
	"errors"
	"github.com/madhanganesh/goerr"
)

func Controller() error {
	err := Service()
	if err != nil {
		return goerr.New(err, "controller failed")
	}
	return nil
}

func Service() error {
	err := Repository()
	if err != nil {
		return goerr.New(err, "service failed")
	}
	return err
}

func Repository() error {
	err := errors.New("error from database")
	return goerr.New(nil, err.Error())
}
