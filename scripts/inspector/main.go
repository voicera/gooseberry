// Package main inspects go source files for the proper use of the log package.
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/voicera/gooseberry/containers/sets"
)

var (
	exitCode              = 0
	sourceFolderPath      = ""
	verbose               = false
	fileSet               = token.NewFileSet()
	funcsToInspect        = sets.NewSetFromStrings("Debug", "Info", "Warn", "Error")
	validFieldNamePattern = regexp.MustCompile(`^"[A-Za-z_][0-9A-Za-z_]+"$`)
)

func init() {
	flag.StringVar(&sourceFolderPath, "i", "", "Path of source folder to inspect.")
	flag.BoolVar(&verbose, "v", false, "Verbose logging")
	flag.Parse()
	if sourceFolderPath == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	panicIfNotNil(filepath.Walk(sourceFolderPath, walkFolder))
	os.Exit(exitCode)
}

func walkFolder(path string, info os.FileInfo, err error) error {
	panicIfNotNil(err)
	if shouldSkipFolder(path, info) {
		return nil
	}
	if verbose {
		fmt.Println("Inspecting log calls for package:", path)
	}
	packages, err := parser.ParseDir(fileSet, path, filter, parser.AllErrors)
	panicIfNotNil(err)
	for _, p := range packages {
		for _, file := range p.Files {
			ast.Inspect(file, inspect)
		}
	}
	return nil
}

func shouldSkipFolder(path string, info os.FileInfo) bool {
	return !info.IsDir() ||
		strings.Contains(path, "/vendor") ||
		strings.ContainsRune(info.Name(), '.') || // hidden (hack; should check first char of each segment)
		info.Name()[0] == '_' // ignored by go build
}

func filter(fileInfo os.FileInfo) bool {
	return true
}

func inspect(node ast.Node) bool {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return true // continue parsing down the AST
	}

	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if !funcsToInspect.Contains(selector.Sel.Name) {
		return false
	}

	packageIdent, ok := selector.X.(*ast.Ident)
	if !ok || packageIdent.Name != "log" {
		return false
	}

	inspectArgs(call.Args)
	return false
}

func inspectArgs(args []ast.Expr) {
	// TODO check that args[0] adheres to log message standards: lower-case start, no delimiter at the end, etc.
	if len(args) < 3 { // allow log.Error(message) and log.Error(message, args...) call patterns
		return
	}
	for i, arg := range args[1:] { // we only care about the varargs
		if i%2 == 0 { // zap field name
			literal, ok := arg.(*ast.BasicLit)
			if !ok { // for now, only allow literal strings (no variables or constants)
				exitCode = 1
				fmt.Printf("%s: Invalid zap field name in log call; must be a literal string!\n",
					fileSet.Position(arg.Pos()))
			} else if !validFieldNamePattern.MatchString(literal.Value) {
				exitCode = 1
				fmt.Printf("%s: Invalid zap field name in log call (%s); should be a quoted variable name\n",
					fileSet.Position(arg.Pos()), literal.Value)
			}
		}
	}
}

func panicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}
