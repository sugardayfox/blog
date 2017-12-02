package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func panicIf(cond bool) {
	if cond {
		panic("condition failed")
	}
}

func mkdirForFile(filePath string) error {
	dir := filepath.Dir(filePath)
	return os.MkdirAll(dir, 0755)
}

func copyFile(dst string, src string) error {
	err := mkdirForFile(dst)
	if err != nil {
		return err
	}
	fin, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fin.Close()
	fout, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, err = io.Copy(fout, fin)
	err2 := fout.Close()
	if err != nil || err2 != nil {
		os.Remove(dst)
	}

	return err
}

func skipFile(path string) bool {
	if strings.Contains(path, ".tmpl.") {
		return true
	}
	return false
}

func dirCopyRecur(dst string, src string) (int, error) {
	fmt.Printf("dirCopyRecur: %s => %s\n", src, dst)
	nFilesCopied := 0
	dirsToVisit := []string{src}
	for len(dirsToVisit) > 0 {
		n := len(dirsToVisit)
		dir := dirsToVisit[n-1]
		dirsToVisit = dirsToVisit[:n-1]
		fileInfos, err := ioutil.ReadDir(dir)
		if err != nil {
			return nFilesCopied, err
		}
		for _, fi := range fileInfos {
			path := filepath.Join(dir, fi.Name())
			if fi.IsDir() {
				dirsToVisit = append(dirsToVisit, path)
				continue
			}
			if skipFile(path) {
				continue
			}
			dstPath := dst + path[len(src):]
			err := copyFile(dstPath, path)
			if err != nil {
				return nFilesCopied, err
			}
			nFilesCopied++
		}
	}
	return nFilesCopied, nil
}

func netlifyBuild() {
	// verify we're in the right directory
	_, err := os.Stat("netlify_static")
	panicIfErr(err)
	outDir := filepath.Join("netlify_static", "www")
	err = os.RemoveAll(outDir)
	panicIfErr(err)
	err = os.MkdirAll(outDir, 0755)
	panicIfErr(err)
	nCopied, err := dirCopyRecur(outDir, "www")
	panicIfErr(err)
	fmt.Printf("Copied %d files\n", nCopied)
}
