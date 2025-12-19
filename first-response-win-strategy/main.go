package main

import (
	"errors"
)

var ErrorNorFound = errors.New("Not Found")

func DistributedQuery(query string, replicas []DatabaseHost) (string, error) {
	return "", nil
}
