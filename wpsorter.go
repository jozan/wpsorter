package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

/**
 * --------------------------
 * HAIKU:
 *
 * A thousand words and me
 * Move is about to come, Bob
 * Low resolution
 * --------------------------
 */

const (
	dir_images string = "images"
	dir_wp     string = "wp"
	dir_lowres string = "lowres"
	minWidth   int    = 1800
	minHeight  int    = 900
)

func createPaths(paths []string) error {
	for _, path := range paths {
		p := "./" + path
		err := os.Mkdir(p, 0777)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Directory created: ", path)
	}
	return nil
}

func validateResolution(i image.Config) bool {
	if i.Width >= minWidth && i.Height >= minHeight {
		return true
	}
	return false
}

func move(path string, fname string, highres bool) error {
	sfx := "/" + fname
	newPath := dir_lowres

	if highres {
		newPath = dir_wp
	}

	err := os.Rename(path, newPath+sfx)
	if err != nil {
		return err
	}
	return nil
}

func visit(path string, f os.FileInfo, err error) error {
	i, ierr := os.Open(path)
	defer i.Close()
	if ierr != nil {
		return err
	}

	// check if file jpg or png
	// get dimensions of images
	// lowres images to other directory and highres to another

	switch strings.ToLower(filepath.Ext(path)) {

	case ".jpg", ".jpeg":
		j, jerr := jpeg.DecodeConfig(i)
		i.Close()
		if jerr != nil {
			return jerr
		}
		v, fname := validateResolution(j), f.Name()
		err := move(path, fname, v)
		if err != nil {
			fmt.Println(err)
		}
		return nil

	case ".png":
		p, perr := png.DecodeConfig(i)
		i.Close()
		if perr != nil {
			return perr
		}
		v, fname := validateResolution(p), f.Name()
		err := move(path, fname, v)
		if err != nil {
			fmt.Println(err)
		}
		return nil
	}

	return nil
}

func main() {

	// check if directories exist and create if necessary
	paths := []string{dir_images, dir_wp, dir_lowres}

	err := createPaths(paths)
	if err != nil {
		fmt.Println(err)
	}

	// loop trough files in the directory
	if err := filepath.Walk(dir_images, visit); err != nil {
		fmt.Println("Walk returned: ", err)
		return
	}

}
