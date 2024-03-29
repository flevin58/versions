/*
Copyright © 2024 Fernando Julio Levin

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/flevin58/versions/cfg"
	"github.com/flevin58/versions/model"
)

var (
	ftext bool
	fedit bool
	fcsv  string
	fpath bool
)

func init() {
	flag.BoolVar(&fedit, "e", false, "Edits the config file")
	flag.BoolVar(&fpath, "p", false, "Shows the full path to the config file")
	flag.StringVar(&fcsv, "c", "", "Outputs as CSV file with given name")
	flag.Parse()
}

func Execute() {
	if len(flag.Args()) > 0 {
		log.Fatalf("error: no arguments expected")
	}

	if fpath {
		fmt.Println(cfg.ConfigFile)
		return
	}

	if fedit {
		editor := cfg.Data.Editor
		do := exec.Command(editor, cfg.ConfigFile)
		do.Stdout = os.Stdout
		do.Stdin = os.Stdin
		do.Stderr = os.Stderr
		if err := do.Run(); err != nil {
			log.Fatalln(err)
		}
		return
	}
	for _, cmd := range cfg.Data.Commands {
		model.Add(cmd)
	}
	if ftext {
		model.ToText()
		return
	}
	if fcsv != "" {
		model.ToCSV(fcsv)
		return
	}
	model.ToTable()
}
