package markdown

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/rwxrob/pegn/ast"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	gmast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/pegn"
)

var (
	MD goldmark.Markdown
)

func init() {
	MD = goldmark.New(
		goldmark.WithExtensions(extension.GFM, meta.Meta),
	)
}

const (
	Untyped int = iota
	Title
	FrontMater
)

// ------------------------------- Title ------------------------------

func ScanTitle(s pegn.Scanner, buf *[]rune) bool {
	m := s.Mark()
	if !s.Scan() || s.Rune() != '#' {
		return s.Revert(m, Title)
	}
	if !s.Scan() || s.Rune() != ' ' {
		return s.Revert(m, Title)
	}
	var count int
	for s.Scan() {
		if count > 100 {
			return s.Revert(m, Title)
		}
		r := s.Rune()
		if r == '\n' {
			if count > 0 {
				return true
			} else {
				return s.Revert(m, Title)
			}
		}
		if buf != nil {
			*buf = append(*buf, r)
		}
		count++
	}
	return true
}

func ParseTitle(s pegn.Scanner) *ast.Node {
	buf := make([]rune, 0, 100)
	if !ScanTitle(s, &buf) {
		return nil
	}
	return &ast.Node{T: Title, V: string(buf)}
}

//func ScanFrontMater(s pegn.Scanner, buf *[]byte) bool {
//	p := MD.Parser()
//	p.Parse()
//	//p.Parse(buf)
//	//m := s.Mark()
//	//if s.Peek(`---`) {
//	//	for s.Peek() {
//	//	}
//	//}
//	//s.Mark()
//	//start := s.Mark()
//	//if !s.Scan() || s.Peek(`---`) {
//	//}
//	return false
//}

func readNodeData(path string) ([]byte, error) {
	if !strings.HasSuffix(path, `README.md`) {
		path = filepath.Join(path, `README.md`)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func ParseFrontMater(data []byte) (map[string]interface{}, error) {
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := MD.Convert(data, &buf, parser.WithContext(context)); err != nil {
		return nil, err
	}

	metaData := meta.Get(context)
	return metaData, nil
}

//func ParseTitle(data bytes.Reader) (string, error) {
//	bytes.NewReader(data)
//	var buf bytes.Buffer
//	context := parser.NewContext(parser.With)
//	p := MD.Parser()
//	if err := MD.Convert(data, &buf, parser.WithContext(context)); err != nil {
//		return "", err
//	}
//
//	p.Parse(data)
//	metaData := meta.Get(context)
//	context.Get()
//	return metaData, nil
//}

func ReadFrontMater(path string) (interface{}, error) {
	data, err := readNodeData(path)
	if err != nil {
		return "", err
	}
	return ParseFrontMater(data)
}

func findFirstHeading(root gmast.Node) *gmast.Heading {
	var node *gmast.Heading = nil
	err := gmast.Walk(root, func(n gmast.Node, entering bool) (gmast.WalkStatus, error) {
		t := n.Type()
		_ = t
		if n.Kind() == gmast.KindHeading {
			h, ok := n.(*gmast.Heading)
			if !ok {
				return gmast.WalkStop, fmt.Errorf("not a header")
			}
			if h.Level == 1 {
				node = h
				return gmast.WalkStop, nil
			}
		}

		return gmast.WalkContinue, nil
	})
	if err != nil {
		return nil
	}

	return node
}

// ReadTitle reads a KEG node title from KEGML file.
func ReadTitle(path string) (string, error) {
	// Check to see if there is a title in the metadata
	titleKeys := []string{"Title", "title"}
	data, err := readNodeData(path)
	if err != nil {
		return "", err
	}
	metadata, err := ParseFrontMater(data)
	if err != nil {
		return "", err
	}

	for _, k := range titleKeys {
		title, exists := metadata[k].(string)
		if exists {
			return fmt.Sprint(title), nil
		}
	}

	// Parse the title from the first H1 found
	r := text.NewReader(data)
	document := MD.Parser().Parse(r)
	var heading *gmast.Heading = nil
	err = gmast.Walk(document, func(n gmast.Node, entering bool) (gmast.WalkStatus, error) {
		t := n.Type()
		_ = t
		if n.Kind() == gmast.KindHeading {
			h, ok := n.(*gmast.Heading)
			if !ok {
				return gmast.WalkStop, fmt.Errorf("invariant: expected a header")
			}
			if h.Level == 1 {
				heading = h
				return gmast.WalkStop, nil
			}
		}

		return gmast.WalkContinue, nil
	})
	if err != nil {
		return "", nil
	}
	if heading == nil {
		return "", nil
	}
	lines := heading.Lines()
	var buf bytes.Buffer
	for i := 0; i < lines.Len(); i++ {
		segment := lines.At(i)
		segment.Value(data)
		buf.Write(segment.Value(data))
	}
	title := buf.String()
	return title, nil
}
