package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/constant"
	"go/format"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

var errCodeDocPrefix = `
# 错误码
## 功能说明
如果返回结果中存在 {{.}}code{{.}} 字段，则表示调用 API 接口失败。例如：
{{.}}{{.}}{{.}}json
{
"code": 100101,
"message": "Database error"
}
{{.}}{{.}}{{.}}
上述返回中 {{.}}code{{.}} 表示错误码，{{.}}message{{.}} 表示该错误的具体信息。每个错误同时也对应一个 HTTP 状态码，比如上述错误码对应了 HTTP 状态码 500(Internal Server Error)。
## 错误码列表
IAM 系统支持的错误码列表如下：
| Identifier | Code | HTTP Code | Description |
| ---------- | ---- | --------- | ----------- |
`
var (
	typeNames  = flag.String("type", "", "comma-separated list of type names; must be set")
	output     = flag.String("output", "", "output file name; default srcdir/<type>_string.go")
	trimprefix = flag.String("trimprefix", "", "trim the `prefix` from the generated constant names")
	buildTags  = flag.String("tags", "", "comma-separated list of build tags to apply")
	doc        = flag.Bool("doc", true, "if false only generate error code documentation in markdown format")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "error用法\n")
	fmt.Fprintf(os.Stderr, "\tgo run code.go -type T [directory]\n")
	fmt.Fprintf(os.Stderr, "\tcodegen [flags] -type T files... # Must be a single package\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("[logusage]")
	log.SetFlags(log.Ldate | log.Llongfile)
	flag.Usage = Usage
	flag.Parse()
	if len(*typeNames) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	types := strings.Split(*typeNames, ",")
	var tags []string
	if len(*buildTags) > 0 {
		tags = strings.Split(*buildTags, ",")
	}
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	var dir string
	g := Generator{
		trimPrefix: *trimprefix,
	}
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
	} else {
		if len(tags) != 0 {
			log.Fatal("-tags option applies only to directories, not when files are specified")
		}
		dir = filepath.Dir(args[0])
	}
	g.parsePackage(args, tags)
	var src []byte
	for _, typeName := range types {
		if *doc {
			g.generate(typeName)
			// Format the output.
			src = g.format()
		} else {
			g.generateDocs(typeName)
			src = g.buf.Bytes()
		}
	}
	outputName := *output
	if outputName == "" {
		absDir, _ := filepath.Abs(dir)
		baseName := fmt.Sprintf("%s_err.go", strings.ReplaceAll(filepath.Base(absDir), "-", "_"))
		if len(flag.Args()) == 1 {
			baseName = fmt.Sprintf(
				"%s_err.go",
				strings.ReplaceAll(filepath.Base(strings.TrimSuffix(flag.Args()[0], ".go")), "-", "_"),
			)
		}
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}
	err := ioutil.WriteFile(outputName, src, 0o600)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
	fmt.Println("output success ╮(￣▽ ￣)╭")
}

type Generator struct {
	buf bytes.Buffer
	pkg *Package

	trimPrefix string
}

type File struct {
	pkg        *Package
	file       *ast.File
	typeName   string
	values     []Value
	trimPrefix string
}

type Package struct {
	name  string
	defs  map[*ast.Ident]types.Object
	files []*File
}

type Value struct {
	comment      string
	originalName string
	name         string
	value        uint64
	signed       bool
	str          string
}

func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}

	return info.IsDir()
}

func (g *Generator) parsePackage(patterns []string, tags []string) {
	cfg := &packages.Config{
		Mode:       packages.LoadSyntax,
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		panic(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	g.addPackage(pkgs[0])
}

func (g *Generator) addPackage(pkg *packages.Package) {
	g.pkg = &Package{
		name:  pkg.Name,
		defs:  pkg.TypesInfo.Defs,
		files: make([]*File, len(pkg.Syntax)),
	}

	for i, file := range pkg.Syntax {
		g.pkg.files[i] = &File{
			file:       file,
			pkg:        g.pkg,
			trimPrefix: g.trimPrefix,
		}
	}
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

func (g *Generator) generate(typeName string) {
	values := make([]Value, 0, 100)
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.typeName = typeName
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			values = append(values, file.values...)
		}
	}

	if len(values) == 0 {
		log.Fatalf("no values defined for type %s", typeName)
	}
	// Generate code that will fail if the constants change value.
	sort.Slice(values, func(i, j int) bool {
		code1, _ := values[i].ParseComment()
		code2, _ := values[j].ParseComment()
		if code1 >= code2 {
			return false
		} else {
			return true
		}
	})
	g.Printf("package errno\n\nvar(\n")
	for _, v := range values {
		code, description := v.ParseComment()
		g.Printf("\tHttp%s = NewHttpErr(%d ,%s ,\"%s\")\n", v.originalName, v.value, code, description)
	}
	g.Printf(")\n\n")

	g.Printf("var(\n")
	for _, v := range values {
		_, description := v.ParseComment()
		g.Printf("\t%s = NewErrNo(%d ,\"%s\")\n", v.originalName, v.value, description)
	}
	g.Printf(")\n")
}

func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")

		return g.buf.Bytes()
	}

	return src
}

func (g *Generator) generateDocs(typeName string) {
	values := make([]Value, 0, 100)
	for _, file := range g.pkg.files {
		file.typeName = typeName
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			values = append(values, file.values...)
		}
	}

	if len(values) == 0 {
		log.Fatalf("no values defined for type %s", typeName)
	}

	tmpl, _ := template.New("doc").Parse(errCodeDocPrefix)
	var buf bytes.Buffer
	_ = tmpl.Execute(&buf, "`")
	g.Printf(buf.String())
	for _, v := range values {
		code, description := v.ParseComment()
		g.Printf("| %s | %d | %s | %s |\n", v.originalName, v.value, code, description)
	}
	g.Printf("\n")
}

func (v *Value) ParseComment() (string, string) {
	reg := regexp.MustCompile(`\w\s*-\s*(\d{3})\s*:\s*([A-Z].*)\s*\.\n*`)
	if !reg.MatchString(v.comment) {
		return "500", "Internal server error"
	}

	groups := reg.FindStringSubmatch(v.comment)
	if len(groups) != 3 {
		return "500", "Internal server error"
	}

	return groups[1], groups[2]
}

func (f *File) genDecl(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.CONST {
		return true
	}
	typ := ""
	for _, spec := range decl.Specs {
		vspec, _ := spec.(*ast.ValueSpec)
		if vspec.Type == nil && len(vspec.Values) > 0 {
			typ = ""
			ce, ok := vspec.Values[0].(*ast.CallExpr)
			if !ok {
				continue
			}
			id, ok := ce.Fun.(*ast.Ident)
			if !ok {
				continue
			}
			typ = id.Name
		}
		if vspec.Type != nil {
			ident, ok := vspec.Type.(*ast.Ident)
			if !ok {
				continue
			}
			typ = ident.Name
		}
		if typ != f.typeName {
			continue
		}
		for _, name := range vspec.Names {
			if name.Name == "_" {
				continue
			}

			obj, ok := f.pkg.defs[name]
			if !ok {
				log.Fatalf("no value for constant %s", name)
			}
			info := obj.Type().Underlying().(*types.Basic).Info()
			if info&types.IsInteger == 0 {
				log.Fatalf("can't handle non-integer constant type %s", typ)
			}
			value := obj.(*types.Const).Val()
			if value.Kind() != constant.Int {
				log.Fatalf("can't happen: constant is not an integer %s", name)
			}
			i64, isInt := constant.Int64Val(value)
			u64, isUint := constant.Uint64Val(value)
			if !isInt && !isUint {
				log.Fatalf("internal error: value of %s is not an integer: %s", name, value.String())
			}
			if !isInt {
				u64 = uint64(i64)
			}
			v := Value{
				originalName: name.Name,
				value:        u64,
				signed:       info&types.IsUnsigned == 0,
				str:          value.String(),
			}
			if vspec.Doc != nil && vspec.Doc.Text() != "" {
				v.comment = vspec.Doc.Text()
			} else if c := vspec.Comment; c != nil && len(c.List) == 1 {
				v.comment = c.Text()
			}

			v.name = strings.TrimPrefix(v.originalName, f.trimPrefix)
			f.values = append(f.values, v)
		}
	}

	return false
}
