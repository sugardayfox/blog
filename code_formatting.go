package main

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/kjk/u"
)

var langToPrettifyLangMap = map[string]string{
	"c":          "c",
	"c++":        "cc",
	"cpp":        "cpp",
	"python":     "py",
	"html":       "html",
	"xml":        "xml",
	"perl":       "pl",
	"c#":         "cs",
	"javascript": "js",
	"java":       "java",
}

func langToPrettifyLang(lang string) string {
	//from http://google-code-prettify.googlecode.com/svn/trunk/README.html
	//"bsh", "c", "cc", "cpp", "cs", "csh", "cyc", "cv", "htm", "html",
	//"java", "js", "m", "mxml", "perl", "pl", "pm", "py", "rb", "sh",
	//"xhtml", "xml", "xsl"
	if l, ok := langToPrettifyLangMap[lang]; ok {
		return fmt.Sprintf("lang-%s", l)
	}
	return ""
}

func txtCookie(s string) string {
	return u.Sha1HexOfBytes([]byte(s))
}

func extractLang(s []byte) (rest, lang []byte) {
	for len(s) > 0 {
		if s[0] != ' ' {
			break
		}
		s = s[1:]
	}
	i := bytes.IndexByte(s, '>')
	if 0 == i {
		return s[1:], nil
	}
	return s[i+1:], s[:i]
}

// flags:
// s : let . match \n (default false)
// i : case insensitive
// U : ungreedy
var reCode = regexp.MustCompile("(?siU)<code.*>.+</code>")

func txtWithCodeParts(txt []byte) ([]byte, map[string][]byte) {
	codeParts := make(map[string][]byte)
	res := reCode.ReplaceAllFunc(txt, func(s []byte) []byte {
		s = s[len("<code") : len(s)-len("</code>")]
		s, lang := extractLang(s)
		var newCode string
		if lang != nil {
			l := langToPrettifyLang(string(lang))
			newCode = fmt.Sprintf(`<pre class="prettyprint %s">%s</pre>`, l, string(s))
		} else {
			newCode = fmt.Sprintf(`<pre class="prettyprint">%s</pre>`, string(s))
		}
		cookie := txtCookie(newCode)
		codeParts[cookie] = []byte(newCode)
		return []byte(cookie)
	})
	return res, codeParts
}

// flags:
// s : let . match \n (default false)
// i : case insensitive
// U : ungreedy
// m : multi-line mode: ^ and $ match begin/end line in addition to begin/end text (default false)
// https://github.com/google/re2/wiki/Syntax
var reCodeMarkdown = regexp.MustCompile("(?siUm)^```.+\n```")

var (
	htmlFormatter  *html.Formatter
	highlightStyle *chroma.Style
)

// based on https://github.com/alecthomas/chroma/blob/master/quick/quick.go
func htmlHighlight(w io.Writer, source, lexer string) error {
	// Determine lexer.
	l := lexers.Get(lexer)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	if htmlFormatter == nil {
		htmlFormatter = html.New(html.WithClasses(), html.TabWidth(1))
		u.PanicIf(htmlFormatter == nil, "couldn't create html formatter")
		styleName := "monokailight"
		highlightStyle = styles.Get(styleName)
		u.PanicIf(highlightStyle == nil, "didn't find style '%s'", styleName)
	}

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}
	return htmlFormatter.Format(w, highlightStyle, it)
}

var (
	debugMarkdownCodeHighligh = false
)

// extracts code snippets from markdown s, converts them to syntax-higlighted
// html. Replaces snippets with a random string and creates mapping of
// random string => html
func markdownCodeHighligh(d []byte) ([]byte, map[string][]byte) {
	codeParts := make(map[string][]byte)
	d = normalizeNewlines(d)
	res := reCodeMarkdown.ReplaceAllFunc(d, func(d []byte) []byte {
		d = d[:len(d)-4] // remove \n``` from the end
		// first line is ``` + optional lang
		var line string
		line, d = bytesRemoveFirstLine(d)
		// remove ``` from the start which leaves (optional) lang
		lang := string(line[3:])

		var buf bytes.Buffer
		err := htmlHighlight(&buf, string(d), lang)
		u.PanicIfErr(err)
		res := buf.Bytes()
		if debugMarkdownCodeHighligh {
			fmt.Printf("d:\n%s.\n=>\n%s\n\n", string(d), res)
		}
		newCode := string(res)
		cookie := txtCookie(newCode)
		codeParts[cookie] = []byte(newCode)
		return []byte(cookie)
	})
	return res, codeParts
}
