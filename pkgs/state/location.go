package state

import (
	"fmt"
	"os"
	"path"
)

type Location struct {
	Origin      string `json:"origin"`
	IsDirectory bool   `json:"is_directory"`
	Destination string `json:"destination"`
}

func (l *Location) CreateLink() error {
	origin := l.Origin
	destination := l.Destination

	stat, err := os.Stat(origin)
	if err != nil {
		return err
	}
	l.IsDirectory = stat.IsDir()
	if !l.IsDirectory {
		destination = path.Dir(destination)
	}

	err = os.MkdirAll(destination, os.ModePerm)
	if err != nil {
		return err
	}

	return os.Symlink(origin, l.Destination)
}

func (l *Location) RemoveLink() error {
	stat, err := os.Lstat(l.Destination)
	if err != nil {
		return err
	}

	if stat.Mode()&os.ModeSymlink == 0 {
		fmt.Println("Not a symlink, skipping removal:", l.Destination)
		return nil
	}

	return os.Remove(l.Destination)
}

func (l *Location) Compare(other Location) bool {
	return l.Origin == other.Origin && l.Destination == other.Destination && l.IsDirectory == other.IsDirectory
}
