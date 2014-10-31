package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

//				 Name     Versions
type Index map[string][]string

var index Index

func GetIndex() error {

	if len(index) < 1 {

		repo := path.Join(config.R.String(), "Forkliftindex")

		res, err := http.Get(repo)
		if err != nil {
			return err
		}

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("%d %s", res.StatusCode, http.StatusText(res.StatusCode))
		}

		err = json.NewDecoder(res.Body).Decode(&index)
		if err != nil {
			return err
		}
	}
	return nil
}
