package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kjk/u"
)

/*
To make it easy to use from shell to detect if content changed,
returns exit code 1 when something changed and 0 if didn't
*/

var (
	contentChanged = false
)

func calcMainCSSSha1Short() string {
	pattern := filepath.Join("www", "css", "main*.css")
	matches, err := filepath.Glob(pattern)
	u.PanicIfErr(err)
	u.PanicIf(len(matches) != 1)
	path := matches[0]
	sha1, err := u.Sha1HexOfFile(path)
	u.PanicIfErr(err)
	sha1 = sha1[:8]
	dstPath := filepath.Join("www", "css", "main-"+sha1+".css")
	if path == dstPath {
		return ""
	}
	cmd := exec.Command("git", "mv", path, dstPath)
	_, err = cmd.CombinedOutput()
	u.PanicIfErr(err)
	return sha1
}

var (
	rxMainCSS = regexp.MustCompile("/css/main.*.css")
)

// return true if file was updated
func updateCSSSha1InFile(path, sha1 string) bool {
	d, err := ioutil.ReadFile(path)
	u.PanicIfErr(err)
	replacement := "/css/main-" + sha1 + ".css"
	s := string(d)
	s2 := rxMainCSS.ReplaceAllString(s, replacement)
	if s == s2 {
		return false
	}
	err = ioutil.WriteFile(path, []byte(s2), 0644)
	u.PanicIfErr(err)
	return true
}

// rename www/css/main.css to www/css/main-${sha1}.css and update references
func updateMainCSSSha1() {
	sha1 := calcMainCSSSha1Short()
	if sha1 == "" {
		// no change
		fmt.Printf("css/main.css didn't change\n")
		return
	}
	contentChanged = true
	filepath.Walk("www", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".html" {
			return nil
		}
		if updateCSSSha1InFile(path, sha1) {
			fmt.Printf("Updated %s\n", path)
		}
		return nil
	})
}

func main() {
	updateMainCSSSha1()
	if contentChanged {
		os.Exit(1)
	}
}
