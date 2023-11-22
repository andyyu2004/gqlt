package parser

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"

	"andyyu2004/gqlt/lex"
	"andyyu2004/gqlt/syn"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vektah/gqlparser/v2/parser"
	"github.com/wk8/go-ordered-map/v2"
)

// Implementation notes:
// The error handling invariant is that if you return nil, you must have already emitted an error.
type Parser struct {
	steps  int
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
		if !isNil(stmt) {
			p.expect(lex.Semi)
			p.stmts = append(p.stmts, stmt)
		} else {
			// error recovery by skipping tokens until the next start of statement by searching for a semicolon
			for !p.at(lex.EOF) && !p.eat_(lex.Semi) {
				p.lexer.Next()
			}
		}
	}

	assert(p.at(lex.EOF))

	return syn.File{Stmts: p.stmts}, errors.Join(p.errors...)
}

func (p *Parser) step() {
	p.steps++
	if p.steps > 100000 {
		panic("oops, detected loop")
	}
}

func (p *Parser) peek() lex.Token {
	p.step()
	return p.lexer.Peek()
}

func (p *Parser) at(kind lex.TokenKind) bool {
	tok := p.peek()
	return tok.Kind == kind
}

func (p *Parser) eat_(kind lex.TokenKind) bool {
	_, ok := p.eat(kind)
	return ok
}

func (p *Parser) eat(kind lex.TokenKind) (lex.Token, bool) {
	tok := p.lexer.Peek()
	if tok.Kind == kind {
		tok := p.lexer.Next()
		return tok, true
	}

	return lex.Token{}, false
}

func (p *Parser) bump(kind lex.TokenKind) {
	_, ok := p.eat(kind)
	assert(ok)
}

func (p *Parser) expect_(kind lex.TokenKind) bool {
	_, ok := p.expect(kind)
	return ok
}

func (p *Parser) expect(kind lex.TokenKind) (lex.Token, bool) {
	if tok, ok := p.eat(kind); ok {
		return tok, true
	}

	tok := p.lexer.Peek()
	p.errors = append(p.errors, mkError(tok, "expected `%s`, found `%s`", kind.String(), tok.String()))
	return lex.Token{}, false
}

func mkError(tok lex.Token, msg string, args ...any) *gqlerror.Error {
	return gqlerror.ErrorLocf(tok.Pos.Src.Name, tok.Pos.Line, tok.Pos.Column, msg, args...)
}

func (p *Parser) parseStmt() syn.Stmt {
	tok := p.peek()

	switch tok.Kind {
	case lex.Let:
		return p.parseLetStmt()
	case lex.Assert:
		return p.parseAssertStmt()
	default:
		expr := p.parseExpr()
		if expr == nil {
			return nil
		}

		return &syn.ExprStmt{Expr: expr}
	}
}

func (p *Parser) parseAssertStmt() *syn.AssertStmt {
	p.bump(lex.Assert)
	expr := p.parseExpr()
	if expr == nil {
		return nil
	}

	return &syn.AssertStmt{Expr: expr}
}

func (p *Parser) parseLetStmt() *syn.LetStmt {
	p.bump(lex.Let)
	pat := p.parsePat()
	if pat == nil {
		return nil
	}

	if !p.expect_(lex.Equals) {
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
	case lex.BraceL:
		return p.parseObjectPat()
	case lex.Int, lex.Float, lex.String, lex.BlockString, lex.True, lex.False, lex.Null:
		return p.parseLiteralPat()
	default:
		p.errors = append(p.errors, mkError(tok, "expected pattern, found `%s`", tok.String()))
		return nil
	}
}

func (p *Parser) parseLiteralPat() *syn.LiteralPat {
	return &syn.LiteralPat{Value: p.parseLiteral()}
}

func (p *Parser) parseObjectPat() *syn.ObjectPat {
	p.bump(lex.BraceL)
	fields := orderedmap.New[string, syn.Pat]()
	for !p.at(lex.EOF) && !p.at(lex.BraceR) {
		name, ok := p.expect(lex.Name)
		if !ok {
			return nil
		}

		var pat syn.Pat = &syn.NamePat{Name: name.Value}
		if p.eat_(lex.Colon) {
			pat = p.parsePat()
			if pat == nil {
				return nil
			}
		}

		fields.Set(name.Value, pat)
	}

	p.expect(lex.BraceR)

	return &syn.ObjectPat{Fields: fields}
}

func (p *Parser) parseExpr() syn.Expr {
	expr := p.parseExprBP(1)

	switch p.peek().Kind {
	case lex.ParenL:
		return p.parseCallExpr(expr)
	case lex.Matches:
		return p.parseMatchesExpr(expr)
	default:
		return expr
	}
}

// pratt parser binding power
type bp int

type assoc bool

const (
	left  assoc = false
	right assoc = true
)

func (p *Parser) parseExprBP(minBp bp) syn.Expr {
	expr := p.parseAtomExpr()
	if expr == nil {
		return nil
	}

	for {
		bp, token, assoc := p.currentOp()

		if bp < minBp {
			break
		}

		p.bump(token.Kind)

		if assoc == left {
			bp++
		}

		rhs := p.parseExprBP(bp)
		if rhs == nil {
			return nil
		}

		expr = &syn.BinaryExpr{Left: expr, Op: token.Kind, Right: rhs}
	}

	return expr
}

func (p *Parser) currentOp() (bp, lex.Token, assoc) {
	tok := p.peek()
	switch tok.Kind {
	case lex.Equals2:
		return 10, tok, left
	case lex.BangEqual:
		return 10, tok, left
	default:
		return 0, tok, left
	}
}

func (p *Parser) parseMatchesExpr(expr syn.Expr) *syn.MatchesExpr {
	p.bump(lex.Matches)
	pat := p.parsePat()
	if pat == nil {
		return nil
	}

	return &syn.MatchesExpr{Expr: expr, Pat: pat}
}

func (p *Parser) parseCallExpr(f syn.Expr) *syn.CallExpr {
	p.bump(lex.ParenL)
	args := []syn.Expr{}
	for !p.at(lex.EOF) && !p.at(lex.ParenR) {
		arg := p.parseExpr()
		if arg == nil {
			return nil
		}

		args = append(args, arg)
	}

	p.expect(lex.ParenR)

	return &syn.CallExpr{Fn: f, Args: args}
}

func (p *Parser) parseAtomExpr() syn.Expr {
	tok := p.peek()
	switch tok.Kind {
	case lex.BraceL:
		return p.parseObjectExpr()
	case lex.Query, lex.Mutation:
		return p.parseQueryExpr()
	case lex.Int, lex.Float, lex.String, lex.BlockString, lex.True, lex.False, lex.Null:
		return p.parseLiteralExpr()
	case lex.Name:
		p.bump(lex.Name)
		return &syn.NameExpr{Name: tok.Value}
	default:
		p.errors = append(p.errors, mkError(tok, "expected expression, found `%s`", tok.Kind.String()))
		return nil
	}
}

func (p *Parser) parseObjectExpr() *syn.ObjectExpr {
	p.bump(lex.BraceL)
	fields := orderedmap.New[string, syn.Expr]()
	for !p.at(lex.EOF) && !p.at(lex.BraceR) {
		name, ok := p.expect(lex.Name)
		if !ok {
			return nil
		}

		var expr syn.Expr = &syn.NameExpr{Name: name.Value}
		if p.eat_(lex.Colon) {
			expr = p.parseExpr()
			if expr == nil {
				return nil
			}
		}

		fields.Set(name.Value, expr)
	}

	p.expect(lex.BraceR)

	return &syn.ObjectExpr{Fields: fields}
}

func (p *Parser) parseLiteralExpr() *syn.LiteralExpr {
	return &syn.LiteralExpr{Value: p.parseLiteral()}
}

func (p *Parser) parseLiteral() any {
	if s, ok := p.eat(lex.Int); ok {
		return must(strconv.Atoi(s.Value))
	} else if s, ok := p.eat(lex.Float); ok {
		return must(strconv.ParseFloat(s.Value, 64))
	} else if s, ok := p.eat(lex.String); ok {
		return s.Value
	} else if s, ok := p.eat(lex.BlockString); ok {
		return s.Value
	} else if p.eat_(lex.True) {
		return true
	} else if p.eat_(lex.False) {
		return false
	} else if p.eat_(lex.Null) {
		return nil
	} else {
		panic("unreachable, token types were checked by caller")
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

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

func isNil(i any) bool {
	if i == nil {
		return true
	}

	val := reflect.ValueOf(i)
	return val.Kind() == reflect.Ptr && val.IsNil()
}
