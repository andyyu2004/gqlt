package parser

import (
	"errors"
	"os"
	"reflect"
	"strings"

	"andyyu2004/gqlt/lex"
	"andyyu2004/gqlt/syn"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vektah/gqlparser/v2/parser"
)

type Parser struct {
	lexer  lex.Lexer
	stmts  []syn.Stmt
	errors []error
}

func NewFromPath(path string) (*Parser, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lexer, err := lex.New(&ast.Source{Name: path, Input: string(bytes)})
	return &Parser{lexer: lexer}, err
}

func New(src *ast.Source) (*Parser, error) {
	lexer, err := lex.New(src)
	return &Parser{lexer: lexer}, err
}

func (p *Parser) Parse() (syn.File, error) {
	for !p.at(lex.EOF) {
		stmt := p.parseStmt()
		// could do some error recovery here
		if !isNil(stmt) {
			p.stmts = append(p.stmts, stmt)
		}

		// error recovery by skipping tokens until the next start of statement
		for !p.at(lex.EOF) && !p.at(lex.Let) {
			p.lexer.Next()
		}

	}

	return syn.File{Stmts: p.stmts}, errors.Join(p.errors...)
}

func (p *Parser) peek() lex.Token {
	return p.lexer.Peek()
}

func (p *Parser) at(kind lex.TokenKind) bool {
	tok := p.lexer.Peek()
	return tok.Kind == kind
}

func (p *Parser) eat(kind lex.TokenKind) bool {
	tok := p.lexer.Peek()
	if tok.Kind == kind {
		p.lexer.Next()
		return true
	}
	return false
}

func (p *Parser) bump(kind lex.TokenKind) {
	assert(p.eat(kind))
}

func (p *Parser) expect(kind lex.TokenKind) bool {
	if p.eat(kind) {
		return true
	}

	tok := p.lexer.Peek()
	p.errors = append(p.errors, mkError(tok, "expected token %s, found %s `%s`", kind.Name(), tok.Kind.Name(), tok.String()))
	return false
}

func mkError(tok lex.Token, msg string, args ...any) *gqlerror.Error {
	return gqlerror.ErrorLocf(tok.Pos.Src.Name, tok.Pos.Line, tok.Pos.Column, msg, args...)
}

func (p *Parser) parseStmt() syn.Stmt {
	tok := p.peek()

	switch tok.Kind {
	case lex.Let:
		return p.parseLetStmt()
	default:
		p.errors = append(p.errors, mkError(tok, "expected statement, found %s `%s`", tok.Kind.Name(), tok.String()))
		return nil
	}
}

func (p *Parser) parseLetStmt() *syn.LetStmt {
	assert(p.expect(lex.Let))
	pat := p.parsePat()
	if pat == nil {
		return nil
	}

	if !p.expect(lex.Equals) {
		return nil
	}

	expr := p.parseExpr()
	if expr == nil {
		return nil
	}

	return &syn.LetStmt{Pat: pat, Expr: expr}
}

func (p *Parser) parsePat() syn.Pat {
	tok := p.lexer.Peek()
	switch tok.Kind {
	case lex.Name:
		p.bump(lex.Name)
		return &syn.NamePat{Name: tok.Value}
	default:
		p.errors = append(p.errors, mkError(tok, "expected pattern, found %s `%s`", tok.Kind.Name(), tok.String()))
		return nil
	}
}

func (p *Parser) parseExpr() syn.Expr {
	tok := p.peek()
	switch tok.Kind {
	case lex.Query, lex.Mutation:
		return p.parseQueryExpr()
	default:
		p.errors = append(p.errors, mkError(tok, "expected expression, found %s `%s`", tok.Kind.Name(), tok.String()))
		return nil
	}
}

func (p *Parser) parseQueryExpr() *syn.OperationExpr {
	parser := parser.New(&p.lexer)
	startPos := p.lexer.Peek().Pos
	operation := parser.ParseOperationDefinition()
	if err := parser.Err(); err != nil {
		p.errors = append(p.errors, err)
		return nil
	}

	endTok := p.lexer.Peek()
	query := strings.TrimRight(endTok.Pos.Src.Input[startPos.Start:endTok.Pos.Start], "\n")

	return &syn.OperationExpr{Query: query, Operation: operation}
}

func assert(cond bool) {
	if !cond {
		panic("assertion failed")
	}
}

func isNil(i any) bool {
	if i == nil {
		return true
	}

	val := reflect.ValueOf(i)
	return val.Kind() == reflect.Ptr && val.IsNil()
}
