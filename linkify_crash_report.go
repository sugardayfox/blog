package main

import (
	"bytes"
	"fmt"
	"strings"
)

// given a path in the form foo\bar.cpp, return
// foo\Bar.cpp (i.e. uppercase the first letter
// of the file name)
func uppercaseFile(path string) string {
	i := strings.LastIndex(path, "\\")
	if i == -1 {
		return path
	}
	// don't crash on foo\ and foo\p
	if i+2 >= len(path) {
		return path
	}
	letter := strings.ToUpper(path[i+1 : i+2])
	return path[:i+1] + letter + path[i+2:]
}

var validSrcExts = []string{".cpp", ".c", ".h"}

func isSrcFile(s string) bool {
	if !strings.HasPrefix(s, "src") {
		return false
	}
	for _, ext := range validSrcExts {
		if strings.HasSuffix(s, ext) {
			return true
		}
	}
	return false
}

func findFileFixes(dir string) {
	n := len(dir)
	allFiles := ListFilesInDir(dir, true)
	files := make([]string, 0)
	for _, f := range allFiles {
		f = f[n+1:]
		if isSrcFile(f) {
			f = strings.Replace(f, "/", "\\", -1)
			files = append(files, f)
		}
	}
	needFix := make([]string, 0)
	for _, f := range files {
		tmp := strings.ToLower(f)
		tmp = uppercaseFile(tmp)
		if f != tmp {
			needFix = append(needFix, f)
			fmt.Printf("%s => %s\n", f, tmp)
		}
	}
}

// 0114C072 01:0004B072 sumatrapdf.exe!CrashMe+0x2
// c:\users\kkowalczyk\src\sumatrapdf\src\utils\baseutil.cpp+14

func linkifyCrashReportLine(l []byte) []byte {
	i := bytes.Index(l, []byte("\\sumatrapdf"))
	if -1 == i {
		return l
	}
	found := false
	for ; i >= 0; i-- {
		if l[i] == ' ' {
			found = true
			break
		}
	}
	if !found {
		return l
	}
	before := l[:i]
	file := l[i+1:]
	// by convention, sumatra files are built from a directory
	// whose name starts with \sumatrapdf (e.g. \sumatrapdf\ or
	// \sumatrapdf-2.1.1\)
	i = bytes.Index(file, []byte("\\sumatrapdf")) + 1
	found = false
	for ; i < len(file); i++ {
		if file[i] == '\\' {
			found = true
			break
		}
	}
	if !found {
		return l
	}
	file = file[i+1:]
	fileStr := string(file)
	// at this point file is in the form src\print+420
	parts := strings.Split(fileStr, "+")
	if len(parts) != 2 {
		return l
	}
	fileStr = parts[0]
	lineStr := parts[1]

	res := make([]byte, 0)
	res = append(res, before...)
	res = append(res, '+')
	res = append(res, []byte(fileStr)...)
	res = append(res, '#')
	res = append(res, []byte(lineStr)...)
	return res
}

func linkifyCrashReport(s []byte) []byte {
	res := make([]byte, 0)
	var l []byte
	for {
		l, s = ExtractLine(s)
		if l == nil {
			break
		}
		l = linkifyCrashReportLine(l)
		res = append(res, l...)
		res = append(res, '\n')
	}
	return res
}
