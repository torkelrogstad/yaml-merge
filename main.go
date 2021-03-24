package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/torkelrogstad/yaml-merge/merge"
)

func exit(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("error: %s", msg))
	_, _ = fmt.Fprintln(os.Stderr)
	flag.Usage()
	os.Exit(1)
}

func main() {
	flag.Usage = func() {
		const msg = `yaml-merge: merge YAML files
usage: yaml-merge file1.yaml file2.yaml [...further files]`
		_, _ = fmt.Fprint(os.Stderr, msg+"\n")
	}

	flag.Parse()

	if len(flag.Args()) < 2 {
		exit("needs at least two files")
	}

	var files [][]byte
	for _, arg := range flag.Args() {
		file, err := os.ReadFile(arg)
		if err != nil {
			exit(err.Error())

		}

		files = append(files, file)
	}

	//goland:noinspection GoNilness
	merged, err := merge.Yaml(files[0], files[1], files[2:]...)
	if err != nil {
		exit(fmt.Errorf("could not merge files: %w", err).Error())
	}

	fmt.Println(string(merged))
	os.Exit(0)
}
