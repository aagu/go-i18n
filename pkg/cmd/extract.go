/*
Copyright (c) 2014 Nick Snyder https://github.com/nicksnyder

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
	"github.com/aagu/go-i18n/pkg/translation"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func ExtractMessages(buf []byte) ([]*translation.Message, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", buf, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	extractor := newExtractor(file)
	ast.Walk(extractor, file)
	return extractor.messages, nil
}

func newExtractor(file *ast.File) *extractor {
	return &extractor{i18nPackageName: i18nPackageName(file)}
}

type extractor struct {
	i18nPackageName string
	messages        []*translation.Message
}

func (e *extractor) Visit(node ast.Node) ast.Visitor {
	e.extractMessages(node)

	return e
}

func (e *extractor) extractMessages(node ast.Node) {
	cl, ok := node.(*ast.CompositeLit)
	if !ok {
		return
	}

	switch t := cl.Type.(type) {
	case *ast.SelectorExpr:
		if !e.isMessageType(t) {
			return
		}
		e.extractMessage(cl)
	case *ast.ArrayType:
		if !e.isMessageType(t.Elt) {
			return
		}
		for _, el := range cl.Elts {
			ecl, ok := el.(*ast.CompositeLit)
			if !ok {
				continue
			}
			e.extractMessage(ecl)
		}
	case *ast.MapType:
		if !e.isMessageType(t.Value) {
			return
		}
		for _, el := range cl.Elts {
			kve, ok := el.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			vcl, ok := kve.Value.(*ast.CompositeLit)
			if !ok {
				continue
			}
			e.extractMessage(vcl)
		}
	}
}

func (e *extractor) isMessageType(expr ast.Expr) bool {
	selectorExpr := unwrapSelectorExpr(expr)
	if selectorExpr == nil {
		return false
	}
	if selectorExpr.Sel.Name != "Message" {
		return false
	}
	x, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return false
	}
	return x.Name == e.i18nPackageName
}

func (e *extractor) extractMessage(cl *ast.CompositeLit) {
	data := make(map[string]string)
	for _, elt := range cl.Elts {
		kve, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kve.Key.(*ast.Ident)
		if !ok {
			continue
		}
		value, ok := extractStringLiteral(kve.Value)
		if !ok {
			continue
		}
		data[key.Name] = value
	}
	if !isValidMessage(data) {
		return
	}
	e.messages = append(e.messages, &translation.Message{ID: data["ID"], Text: data["Text"]})
}

func isValidMessage(m map[string]string) bool {
	if len(m) == 0 {
		return false
	}
	if _, exist := m["ID"]; !exist {
		return false
	}
	if _, exist := m["Text"]; !exist {
		return false
	}
	return true
}

func extractStringLiteral(expr ast.Expr) (string, bool) {
	switch v := expr.(type) {
	case *ast.BasicLit:
		if v.Kind != token.STRING {
			return "", false
		}
		s, err := strconv.Unquote(v.Value)
		if err != nil {
			return "", false
		}
		return s, true
	case *ast.BinaryExpr:
		if v.Op != token.ADD {
			return "", false
		}
		x, ok := extractStringLiteral(v.X)
		if !ok {
			return "", false
		}
		y, ok := extractStringLiteral(v.Y)
		if !ok {
			return "", false
		}
		return x + y, true
	case *ast.Ident:
		if v.Obj == nil {
			return "", false
		}
		switch o := v.Obj.Decl.(type) {
		case *ast.ValueSpec:
			if len(o.Values) == 0 {
				return "", false
			}
			s, ok := extractStringLiteral(o.Values[0])
			if !ok {
				return "", false
			}
			return s, true
		}
		return "", false
	default:
		return "", false
	}
}

func unwrapSelectorExpr(expr ast.Expr) *ast.SelectorExpr {
	switch et := expr.(type) {
	case *ast.SelectorExpr:
		return et
	case *ast.StarExpr:
		se, _ := et.X.(*ast.SelectorExpr)
		return se
	default:
		return nil
	}
}

func i18nPackageName(file *ast.File) string {
	for _, i := range file.Imports {
		if i.Path.Kind == token.STRING && i.Path.Value == `"github.com/aagu/go-i18n/pkg/translation"` {
			if i.Name == nil {
				return "translation"
			}
			return i.Name.Name
		}
	}
	return ""
}
