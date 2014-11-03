package engine

import "io"

type Engine struct {
	Logger Logger
}

func (e Engine) Build(location string) {

}

func (e Engine) Install(Package io.Reader) error {
	return nil
}
