package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/aymerick/raymond/ast"
	"github.com/aymerick/raymond/lexer"
	"github.com/aymerick/raymond/parser"
)

func main() {
	start := time.Now().UnixNano()
	content, err := ioutil.ReadFile("demo.hbs")

	if err != nil {
		panic("Can not find file `demo.hbs`")
	}

	// source := `
	// <div class={{ klass }}>aaa</div>
	// {{assign "hello" section.settings.blog}}
	// {{assign "world" (append "kk" (split "array"))}}
	// {{#with this as |global|}}
	//     {{#each blogs as |blog|}}
	//         {{snippet "blog" title=blog.title desc=(append blog.title blog.desc)}}
	//     {{/each}}
	// {{/with}}
	// `

	// scan(source)
	parse(string(content))
	fmt.Printf("\n%s", time.Duration(time.Now().UnixNano()-start))
	// fmt.Println("", collect(source))
}

type MyStruct struct {
	helpers map[string]bool
	Name    string
	Type    string
}

func newMyStruct() *MyStruct {
	return &MyStruct{
		helpers: map[string]bool{},
	}
}

func collectHelpers(node ast.Node) []string {
	visitor := newMyStruct()
	node.Accept(visitor)
	return visitor.collect()
}

func (v *MyStruct) collect() []string {
	var helperNames []string

	for name := range v.helpers {
		helperNames = append(helperNames, name)
	}

	return helperNames
}

func (v *MyStruct) VisitBlock(node *ast.BlockStatement) interface{} {
	helperName := node.Expression.HelperName()
	v.helpers[helperName] = true
	node.Expression.Accept(v)

	if node.Program != nil {
		node.Program.Accept(v)
	}

	if node.Inverse != nil {
		node.Inverse.Accept(v)
	}

	return nil
}

func (v *MyStruct) VisitMustache(node *ast.MustacheStatement) interface{} {
	params := node.Expression.Params
	helperName := node.Expression.HelperName()

	if len(params) > 0 {
		v.helpers[helperName] = true
	}
	node.Expression.Accept(v)
	return nil
}

func (v *MyStruct) VisitSubExpression(node *ast.SubExpression) interface{} {
	params := node.Expression.Params
	helperName := node.Expression.HelperName()

	if len(params) > 0 {
		v.helpers[helperName] = true
	}
	node.Expression.Accept(v)
	return nil
}

func (v *MyStruct) VisitBoolean(node *ast.BooleanLiteral) interface{} {
	return nil
}

func (v *MyStruct) VisitComment(node *ast.CommentStatement) interface{} {
	return nil
}

func (v *MyStruct) VisitContent(node *ast.ContentStatement) interface{} {
	return nil
}

func (v *MyStruct) VisitExpression(node *ast.Expression) interface{} {
	// path
	node.Path.Accept(v)

	// params
	for _, n := range node.Params {
		n.Accept(v)
	}

	// hash
	if node.Hash != nil {
		node.Hash.Accept(v)
	}

	return nil
}

func (v *MyStruct) VisitHash(node *ast.Hash) interface{} {
	for _, p := range node.Pairs {
		p.Accept(v)
	}
	return nil
}

func (v *MyStruct) VisitHashPair(node *ast.HashPair) interface{} {
	node.Val.Accept(v)
	return nil
}

func (v *MyStruct) VisitNumber(node *ast.NumberLiteral) interface{} {
	return nil
}

func (v *MyStruct) VisitPartial(node *ast.PartialStatement) interface{} {
	node.Name.Accept(v)

	if len(node.Params) > 0 {
		node.Params[0].Accept(v)
	}

	// hash
	if node.Hash != nil {
		node.Hash.Accept(v)
	}

	return nil
}

func (v *MyStruct) VisitPath(node *ast.PathExpression) interface{} {
	return nil
}

func (v *MyStruct) VisitProgram(node *ast.Program) interface{} {
	for _, n := range node.Body {
		n.Accept(v)
	}

	return nil
}

func (v *MyStruct) VisitString(node *ast.StringLiteral) interface{} {
	return nil
}

func parse(source string) {
	program, err := parser.Parse(source)

	if err != nil {
		panic(err)
	}

	// fmt.Println("%v", program)

	// output := ast.Print(program)

	helpers := collectHelpers(program)

	fmt.Print(helpers, len(helpers))
}

type JsonOutput struct {
	Type string
	Name string
}

func collect(source string) (output []lexer.Token) {
	output = lexer.Collect(source)
	return
}

func scan(source string) {
	output := ""
	lex := lexer.Scan(source)
	for {
		// consume next token
		token := lex.NextToken()

		// token.Kind == lexer
		output += fmt.Sprintf(" %s\n", token)

		// stops when all tokens have been consumed, or on error
		if token.Kind == lexer.TokenEOF || token.Kind == lexer.TokenError {
			break
		}
	}

	fmt.Print(output)
}
