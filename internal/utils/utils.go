package utils

import (
	"os/user"
	"path/filepath"
)

func ExpandPath(path string) (string, error) {
	if len(path) >= 2 && path[:2] == "~/" {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		path = filepath.Join(usr.HomeDir, path[2:])
	}
	return path, nil
}