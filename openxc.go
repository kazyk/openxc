package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

//greater index is prior
var extensions = []string{
	".xcodeproj",
	".xcworkspace",
}

func searchDir(dir string) (string, bool) {
	fdir, err := os.Open(dir)
	if err != nil {
		fmt.Errorf("%v", err)
		return "", false
	}

	var path string
	maxPriority := -1

	files, _ := fdir.Readdir(0)
	for _, fi := range files {
		for priority, ext := range extensions {
			if priority > maxPriority && filepath.Ext(fi.Name()) == ext {
				maxPriority = priority
				path = filepath.Join(dir, fi.Name())
			}
		}
	}

	if path != "" {
		return path, true
	}

	return searchDir(filepath.Dir(dir))
}

func main() {
	app := flag.String("a", "", "Application name to open files")
	flag.Parse()

	var dir string

	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	} else {
		dir, _ = os.Getwd()
	}

	result, found := searchDir(dir)

	if !found {
		fmt.Println("not found")
		return
	}

	args := []string{result}
	if *app != "" {
		args = append(args, "-a", *app)
	}
	cmd := exec.Command("open", args...)
	out, err := cmd.CombinedOutput()
	fmt.Print(string(out))
	if err != nil {
		fmt.Errorf("%v", err)
	}
}
