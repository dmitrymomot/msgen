package src

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/dmitrymomot/msgen/statik" // import templates
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
)

// Generate project
func Generate(cmd *cobra.Command, args []string) {
	sfs, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	opt := parseOptions(cmd, args)
	files := getFilesList(opt)

	if err := handleFiles(sfs, files, opt); err != nil {
		log.Fatal(err)
	}
}

func handleFiles(sfs http.FileSystem, files []file, opt options) error {
	for _, file := range files {
		if err := render(sfs, file, opt); err != nil {
			return fmt.Errorf("file: %+v, options: %+v, error: %v", file, opt, err)
		}
	}
	return nil
}
