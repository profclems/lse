package ls

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pee2pee/lse/ls/color"
)

const dotCharacter = 46

var dotFiles = []string{".", ".."}

type Flags struct {
	A bool // ls -a
	D bool // ls -d
	G bool // ls --group
	L bool // ls -l
	Q bool // ls --quote
	R bool // ls -R
	T bool // ls -t
}

type LS struct {
	Dir string

	Stderr io.Writer
	StdOut io.Writer
	Color  *color.Palette

	Flags
}

type Dir struct {
	Path string
	Info fs.FileInfo
}

func (l *LS) ListDir() error {
	if l.D {
		return l.showDirStructure()
	}

	if l.R {
		return l.listDirRecursively()
	}
	return l.nonRecursiveListing()
}

func (l *LS) nonRecursiveListing() error {
	dirs, err := os.ReadDir(l.Dir)
	if err != nil {
		return err
	}
	var d []Dir

	// list dotfile if -a is specified
	if l.A {
		for _, file := range dotFiles {
			stat, err := os.Stat(file)
			if err != nil {
				return err
			}
			d = append(d, Dir{
				Info: stat,
				Path: file,
			})
		}
	}

	var modTimes []string

	for _, entry := range dirs {
		if !isHiddenPath(entry.Name(), l.A) {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			d = append(d, Dir{
				Info: info,
				Path: filepath.Join(l.Dir, info.Name()),
			})
			modTimes = append(modTimes, info.ModTime().String())
		}
	}

	if l.T {
		var sortedFiles []Dir
		sort.Sort(sort.Reverse(sort.StringSlice(modTimes)))

		for _, v := range modTimes {
			for _, dir := range d {
				if dir.Info.ModTime().String() == v {
					sortedFiles = append(sortedFiles, dir)
				}
			}
		}

		d = sortedFiles

	}

	if l.G && !l.T {
		var dirs []Dir
		var fileDirs []Dir
		for _, file := range d {
			if file.Info.IsDir() {
				dirs = append(dirs, file)
			} else {
				fileDirs = append(fileDirs, file)
			}
		}
		d = append(dirs, fileDirs...)
	}

	return l.display(d)
}

// listDirRecursively list all subdirectories encountered from the folder
func (l *LS) listDirRecursively() error {
	err := filepath.Walk(l.Dir,

		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if !isHiddenPath(path, l.A) {
				fmt.Fprint(l.StdOut, path, "  ")
				return err
			}
			return nil
		})

	if err != nil {
		return err
	}
	return nil
}

func isHiddenPath(path string, forceHidden bool) bool {
	if forceHidden {
		return false
	}
	return path[0] == dotCharacter
}

func (l *LS) showDirStructure() error {
	file, err := os.Stat(l.Dir)
	if err != nil {
		return err
	}

	p := strings.TrimSuffix(l.Dir, "/")
	if file.IsDir() {
		p = p + "/"
	}

	fmt.Fprintln(l.StdOut, p)
	return nil
}

func (l *LS) listFileIndex() error {
	dirs, err := os.ReadDir(l.Dir)
	if err != nil {
		return err
	}

	for _, file := range dirs {
		fileInfo, err := os.Stat(file.Name())
		if err != nil {
			return err
		}

		fmt.Fprintln(l.StdOut, fileInfo)
	}

	return nil
}
