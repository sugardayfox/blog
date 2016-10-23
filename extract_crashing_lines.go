package main

import (
	"bytes"
	"strings"
)

const (
	noCrashingLine = "no crashing line"
)

func ExtractSumatraCrashingLine(d []byte) string {
	if d = SkipPastLine(d, "Crashed thread:"); d == nil {
		return noCrashingLine
	}
	s := findValidCrashLine(d)
	s = shortenCrashLine(s)
	if s == "" {
		return noCrashingLine
	}
	return s
}

// given a line in the format: c:\kjk\src\sumatrapdf-2.0\ext\openjpeg\j2k.c+1369
// try to extract the part relative to sumatra's source dir (i.e. "ext\openjpeg\j2k.c+1369")
func shortenSrcLine(s string) string {
	if i := strings.Index(s, "\\sumatra"); i != -1 {
		rest := s[i+len("\\sumatra"):]
		if i = strings.Index(rest, "\\"); i != -1 {
			return rest[i+1:]
		}
	}
	return s
}

var modulesToSkip = []string{"sumatrapdf-no-mupdf.exe!", "sumatrapdf.exe!", "libmupdf.dll!"}

// s is in the form:
// 771A2CC7 01:00051CC7 ntdll.dll!RtlFreeHeap+0xcd
// or:
// 5F870D6B 01:000EFD6B libmupdf.dll!j2k_read_sot+0x21b c:\kjk\src\sumatrapdf-2.0\ext\openjpeg\j2k.c+1369
func shortenCrashLine(s string) string {
	parts := strings.Split(s, " ")
	if len(parts) < 3 {
		return s
	}
	// verify first part is what we expect ("771A2CC7" for 32-bit and
	// "771A2CC7771A2CC7" for 64bitetc.)
	n := len(parts[0])
	isValidLen := (n == 8) || (n == 16)
	if !isValidLen {
		return s
	}
	addr := parts[2]
	for _, toSkip := range modulesToSkip {
		if strings.HasPrefix(addr, toSkip) {
			addr = addr[len(toSkip):]
		}
	}
	if 3 == len(parts) {
		return addr
	}
	src := shortenSrcLine(parts[3])
	return addr + " " + src
}

func findValidCrashLine(d []byte) string {
	var l []byte
	for {
		l, d = ExtractLine(d)
		if l == nil || len(l) == 0 {
			return ""
		}
		if bytes.IndexByte(l, ' ') == -1 {
			continue
		}
		return string(l)
	}
}
