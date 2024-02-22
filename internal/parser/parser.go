package parser

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/movio/gqlt/internal/lex"
	"github.com/movio/gqlt/memosa/lib"
	"github.com/movio/gqlt/syn"

	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/gqlparser/gqlerror"
	"github.com/movio/gqlt/gqlparser/parser"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

// Implementation notes:
// The error handling invariant is that if you return nil, you must have already emitted an error.
type Parser struct {
	steps  int
	lexer  lex.Lexer
	stmts  []syn.Stmt
	errors ast.Errors
}

func NewFromPath(path string) (*Parser, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lexer, err := lex.New(&ast.Source{Name: path, Input: string(bytes)})
	return &Parser{lexer: lexer}, err
}

func New(src *ast.Source) *Parser {
	lexer, err := lex.New(src)
	if err != nil {
		return &Parser{lexer: lexer, errors: err.(ast.Errors)}
	}
	return &Parser{lexer: lexer}
}

var stmtFst = []lex.TokenKind{lex.Let, lex.Assert, lex.Set, lex.Fragment}

func (p *Parser) Parse() (syn.File, error) {
	for !p.at(lex.EOF) {
		stmt := p.parseStmt()
		if !lib.IsNil(stmt) {
			p.eat(lex.Semi)
			p.stmts = append(p.stmts, stmt)
		} else {
			// error recovery by skipping tokens until the next start of statement by searching for a semicolon
			// or searching for the start of the next statement
			for !p.at(lex.EOF) && !p.at(stmtFst...) && !p.eat_(lex.Semi) {
				p.lexer.Next()
			}
		}
	}

	assert(p.at(lex.EOF))

	var err error
	if len(p.errors) > 0 {
		err = p.errors
	}

	pos := ast.Position{Src: p.lexer.Peek().Src}
	if len(p.stmts) > 0 {
		pos.Start = p.stmts[0].Pos().Start
		pos.End = p.stmts[len(p.stmts)-1].Pos().End
	}
	return syn.File{Position: pos, Stmts: p.stmts}, err
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

func (p *Parser) at(kinds ...lex.TokenKind) bool {
	tok := p.peek()
	for _, kind := range kinds {
		return tok.Kind == kind
	}
	return false
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

func (p *Parser) bump(kind lex.TokenKind) lex.Token {
	tok, ok := p.eat(kind)
	assert(ok)
	return tok
}

func (p *Parser) expect(kind lex.TokenKind) (lex.Token, bool) {
	if tok, ok := p.eat(kind); ok {
		return tok, true
	}

	tok := p.lexer.Peek()
	p.error(tok, "expected '%s', found '%s'", kind.String(), tok.String())
	return lex.Token{}, false
}

func (p *Parser) parseStmt() syn.Stmt {
	tok := p.peek()

	switch tok.Kind {
	case lex.Let:
		return p.parseLetStmt()
	case lex.Assert:
		return p.parseAssertStmt()
	case lex.Set:
		return p.parseSetStmt()
	case lex.Use:
		return p.parseUseStmt()
	case lex.Fragment:
		return p.parseFragment()
	default:
		expr := p.parseExpr()
		if lib.IsNil(expr) {
			return nil
		}

		return &syn.ExprStmt{Expr: expr}
	}
}

func (p *Parser) parseFragment() *syn.FragmentStmt {
	parser := parser.New(&p.lexer)
	startPos := p.lexer.Peek()
	fragment := parser.ParseFragmentDefinition()
	if err := parser.Err(); err != nil {
		// we can make this conversion as this is the only error type returned by parser
		// (excluding lexer errors which we have already handled)
		err := err.(*gqlerror.Error)
		p.errors = append(p.errors, ast.Error{
			Position: startPos.Pos(),
			Msg:      err.Message,
		})
		return nil
	}

	endTok := p.lexer.Peek()
	rawFragment := strings.TrimRight(endTok.Src.Input[startPos.Start:endTok.Start], "\n")

	return &syn.FragmentStmt{
		Position:    startPos.Merge(endTok),
		RawFragment: rawFragment,
		Definition:  fragment,
	}
}

func (p *Parser) parseUseStmt() *syn.UseStmt {
	useKw := p.bump(lex.Use)
	path, ok := p.expect(lex.String)
	if !ok {
		return nil
	}

	return &syn.UseStmt{UseKw: useKw, Path: path}
}

func (p *Parser) parseSetStmt() *syn.SetStmt {
	setKw := p.bump(lex.Set)
	key, ok := p.expect(lex.Name)
	if !ok {
		return nil
	}

	// optional equals
	p.eat(lex.Equals)

	expr := p.parseExpr()
	if lib.IsNil(expr) {
		return nil
	}

	return &syn.SetStmt{Position: setKw.Merge(expr), SetKw: setKw, Variable: key, Expr: expr}
}

func (p *Parser) parseAssertStmt() *syn.AssertStmt {
	assert := p.bump(lex.Assert)
	expr := p.parseExpr()
	if lib.IsNil(expr) {
		return nil
	}

	return &syn.AssertStmt{Position: assert.Merge(expr), AssertKw: assert, Expr: expr}
}

func (p *Parser) parseLetStmt() *syn.LetStmt {
	let := p.bump(lex.Let)
	pat := p.parsePat(patOpts{})
	if lib.IsNil(pat) {
		return nil
	}

	equals, ok := p.expect(lex.Equals)
	if !ok {
		return nil
	}

	expr := p.parseExpr()
	if lib.IsNil(expr) {
		return nil
	}

	return &syn.LetStmt{Position: let.Merge(expr), LetKw: let, Pat: pat, Equals: equals, Expr: expr}
}

type patOpts struct {
	// to allow `...<pat>` pattern
	allowSpread bool
	// to allow `...` without a trailing pattern (normal grammar is `...<pat>`)
	allowImplicitWildcard bool
}

func (p *Parser) error(tok ast.HasPosition, msg string, args ...any) {
	err := ast.Error{Position: tok.Pos(), Msg: fmt.Sprintf(msg, args...)}
	p.errors = append(p.errors, err)
}

func (p *Parser) parsePat(opts patOpts) syn.Pat {
	tok := p.lexer.Peek()
	switch tok.Kind {
	case lex.Underscore:
		p.bump(lex.Underscore)
		return &syn.WildcardPat{Position: tok.Pos()}
	case lex.Spread:
		p.bump(lex.Spread)
		if !opts.allowSpread {
			p.error(tok, "spread pattern not allowed here")
			return nil
		}

		pat := p.parsePat(patOpts{allowImplicitWildcard: true})
		if lib.IsNil(pat) {
			return nil
		}

		return &syn.RestPat{Position: tok.Pos().Merge(pat), Pat: pat}

	case lex.Name:
		p.bump(lex.Name)
		return &syn.NamePat{Name: tok}
	case lex.BraceL:
		return p.parseObjectPat()
	case lex.BracketL:
		return p.parseListPat()
	case lex.ParenL:
		p.bump(lex.ParenL)
		pat := p.parsePat(patOpts{})
		if lib.IsNil(pat) {
			return nil
		}
		p.expect(lex.ParenR)
		return pat
	case lex.Int, lex.Float, lex.String, lex.BlockString, lex.True, lex.False, lex.Null:
		return p.parseLiteralPat()
	default:
		if opts.allowImplicitWildcard {
			return &syn.WildcardPat{Position: tok.Pos()}
		}
		p.error(tok, "expected pattern, found `%s`", tok.String())
		return nil
	}
}

func (p *Parser) parseLiteralPat() *syn.LiteralPat {
	token := p.peek()
	return &syn.LiteralPat{Position: token.Position, Token: token, Value: p.parseLiteral()}
}

func (p *Parser) parseListPat() *syn.ListPat {
	start := p.bump(lex.BracketL)
	pats := []syn.Pat{}
	for !p.at(lex.EOF) && !p.at(lex.BracketR) {
		pat := p.parsePat(patOpts{allowSpread: true})
		if lib.IsNil(pat) {
			return nil
		}

		pats = append(pats, pat)

		if !p.eat_(lex.Comma) {
			break
		}
	}

	end, _ := p.expect(lex.BracketR)

	for i, pat := range pats {
		rest, ok := pat.(*syn.RestPat)
		if ok && i != len(pats)-1 {
			p.error(rest, "rest pattern must be last")
			return nil
		}

	}

	return &syn.ListPat{Position: start.Merge(end), Pats: pats}
}

func (p *Parser) parseObjectPat() *syn.ObjectPat {
	start := p.bump(lex.BraceL)
	fields := orderedmap.New[lex.Token, syn.Pat]()
	for !p.at(lex.EOF) && !p.at(lex.BraceR) {
		var name lex.Token
		var pat syn.Pat
		var ok bool
		if p.eat_(lex.Spread) {
			name, ok = p.expect(lex.Name)
			if !ok {
				return nil
			}
			pat = &syn.RestPat{Position: name.Pos(), Pat: &syn.NamePat{Name: name}}
		} else {
			if p.at(lex.String) {
				name = p.bump(lex.String)
			} else {
				name, ok = p.expect(lex.Name)
				if !ok {
					return nil
				}
			}

			pat = &syn.NamePat{Name: name}
			if p.eat_(lex.Colon) {
				pat = p.parsePat(patOpts{})
				if lib.IsNil(pat) {
					return nil
				}
			}

		}

		fields.Set(name, pat)

		if !p.eat_(lex.Comma) {
			break
		}
	}

	end, _ := p.expect(lex.BraceR)

	if fields.Len() > 0 {
		for entry := fields.Newest().Prev(); entry != nil; entry = entry.Prev() {
			if _, ok := entry.Value.(*syn.RestPat); ok {
				p.error(entry.Value, "rest pattern must be last")
				return nil
			}
		}
	}

	return &syn.ObjectPat{Position: start.Merge(end), Fields: fields}
}

func (p *Parser) parseExpr() syn.Expr {
	return p.parseExprBP(1)
}

// pratt parser binding power
type bp uint8

type assoc bool

const (
	left  assoc = false
	right assoc = true
)

func (p *Parser) parseExprBP(minBp bp) syn.Expr {
	var lhs syn.Expr
	if tok, bp := p.prefixOp(); tok != nil {
		p.bump(tok.Kind)
		expr := p.parseExprBP(bp)
		if lib.IsNil(expr) {
			return nil
		}

		lhs = &syn.UnaryExpr{Position: tok.Merge(expr), Op: *tok, Expr: expr}
	} else {
		lhs = p.parseAtomExpr()
	}

	if lib.IsNil(lhs) {
		return nil
	}

	for {
		if tok, bp := p.postfixOp(); tok != nil {
			if bp < minBp {
				break
			}

			switch tok.Kind {
			case lex.BracketL:
				lhs = p.parseIndexExpr(lhs)
			case lex.ParenL:
				lhs = p.parseCallExpr(lhs)
			default:
				panic("unreachable")
			}

			continue
		}

		bp, token, assoc := p.infixOp()
		if bp < minBp {
			break
		}

		if token.Kind == lex.Matches {
			assert(assoc == left)
			lhs = p.parseMatchesExpr(lhs)
			continue
		}

		if token.Kind == lex.Dot {
			assert(assoc == left)
			lhs = p.parseFieldExpr(lhs)
			continue
		}

		p.bump(token.Kind)

		if assoc == left {
			bp++
		}

		if lib.IsNil(lhs) {
			return nil
		}

		rhs := p.parseExprBP(bp)
		if lib.IsNil(rhs) {
			return nil
		}

		lhs = &syn.BinaryExpr{Left: lhs, Op: token, Right: rhs}
	}

	return lhs
}

func (p *Parser) prefixOp() (*lex.Token, bp) {
	tok := p.peek()
	switch tok.Kind {
	// not is the same as ! but with lower precedence
	// it's useful for writing expressions such as `assert not x matches y`
	// opposed to the clunky `assert !(x matches y)`
	case lex.Not:
		return &tok, 90
	case lex.Minus, lex.Bang:
		return &tok, 140
	default:
		return nil, 0
	}
}

func (p *Parser) infixOp() (bp, lex.Token, assoc) {
	tok := p.peek()
	switch tok.Kind {
	case lex.Matches:
		return 100, tok, left
	case lex.Equals2, lex.BangEqual:
		return 110, tok, left
	case lex.EqualsTilde, lex.BangTilde, lex.AngleL, lex.AngleR, lex.AngleLEqual, lex.AngleREqual:
		return 115, tok, left
	case lex.Plus, lex.Minus:
		return 120, tok, left
	case lex.Star, lex.Slash:
		return 130, tok, left
	case lex.Dot:
		return 150, tok, left
	default:
		return 0, tok, left
	}
}

func (p *Parser) postfixOp() (*lex.Token, bp) {
	tok := p.peek()
	switch tok.Kind {
	case lex.ParenL:
		return &tok, 150
	case lex.BracketL:
		return &tok, 150
	default:
		return nil, 0
	}
}

func (p *Parser) parseFieldExpr(expr syn.Expr) *syn.FieldExpr {
	dot := p.bump(lex.Dot)
	field, ok := p.expect(lex.Name)
	if !ok {
		return nil
	}

	return &syn.FieldExpr{Expr: expr, Dot: dot, Field: field}
}

func (p *Parser) parseMatchesExpr(expr syn.Expr) *syn.MatchesExpr {
	matches := p.bump(lex.Matches)
	pat := p.parsePat(patOpts{})
	if lib.IsNil(pat) {
		return nil
	}

	return &syn.MatchesExpr{Position: expr.Pos().Merge(pat), Expr: expr, MatchesKw: matches, Pat: pat}
}

func (p *Parser) parseIndexExpr(expr syn.Expr) *syn.IndexExpr {
	p.bump(lex.BracketL)
	index := p.parseExpr()
	if lib.IsNil(index) {
		return nil
	}

	end, _ := p.expect(lex.BracketR)

	return &syn.IndexExpr{Position: expr.Pos().Merge(end), Expr: expr, Index: index}
}

func (p *Parser) parseCallExpr(f syn.Expr) *syn.CallExpr {
	p.bump(lex.ParenL)
	args := []syn.Expr{}
	for !p.at(lex.EOF) && !p.at(lex.ParenR) {
		arg := p.parseExpr()
		if lib.IsNil(arg) {
			return nil
		}

		args = append(args, arg)

		if !p.eat_(lex.Comma) {
			break
		}
	}

	end, _ := p.expect(lex.ParenR)

	return &syn.CallExpr{Position: f.Pos().Merge(end), Fn: f, Args: args}
}

func (p *Parser) parseAtomExpr() syn.Expr {
	tok := p.peek()
	switch tok.Kind {
	case lex.BraceL:
		return p.parseObjectExpr()
	case lex.BracketL:
		return p.parseListExpr()
	case lex.ParenL:
		p.bump(lex.ParenL)
		expr := p.parseExpr()
		if lib.IsNil(expr) {
			return nil
		}
		p.expect(lex.ParenR)
		return expr
	case lex.Query, lex.Mutation:
		return p.parseQueryExpr()
	case lex.Int, lex.Float, lex.String, lex.BlockString, lex.True, lex.False, lex.Null:
		return p.parseLiteralExpr()
	case lex.Try:
		tryKw := p.bump(lex.Try)
		expr := p.parseAtomExpr()
		if lib.IsNil(expr) {
			return nil
		}
		return &syn.TryExpr{TryKw: tryKw, Expr: expr}
	case lex.Name:
		p.bump(lex.Name)
		return &syn.NameExpr{Name: tok}
	default:
		p.error(tok, "expected expression, found '%s'", tok.String())
		return nil
	}
}

func (p *Parser) parseListExpr() *syn.ListExpr {
	start := p.bump(lex.BracketL)
	exprs := []syn.Expr{}
	commas := []lex.Token{}
	for !p.at(lex.EOF) && !p.at(lex.BracketR) {
		expr := p.parseExpr()
		if lib.IsNil(expr) {
			return nil
		}

		exprs = append(exprs, expr)

		comma, ok := p.eat(lex.Comma)
		if !ok {
			break
		}
		commas = append(commas, comma)
	}

	end, _ := p.expect(lex.BracketR)

	return &syn.ListExpr{OpenBracket: start, Position: start.Merge(end), Exprs: exprs, Commas: commas, CloseBracket: end}
}

func (p *Parser) parseObjectExpr() *syn.ObjectExpr {
	start := p.bump(lex.BraceL)
	fields := orderedmap.New[lex.Token, syn.Expr]()
	commas := []lex.Token{}
	var base syn.Expr
	for !p.at(lex.EOF) && !p.at(lex.BraceR) {
		if p.eat_(lex.Spread) {
			base = p.parseExpr()
			if p.eat_(lex.Comma) {
				p.error(p.peek(), "cannot have comma after the base object")
			}
			continue
		}

		var name lex.Token
		if p.at(lex.String) {
			name = p.bump(lex.String)
		} else {
			var ok bool
			name, ok = p.expect(lex.Name)
			if !ok {
				return nil
			}
		}

		var expr syn.Expr = &syn.NameExpr{Name: name}
		if p.eat_(lex.Colon) {
			expr = p.parseExpr()
			if lib.IsNil(expr) {
				return nil
			}
		}

		fields.Set(name, expr)

		comma, ok := p.eat(lex.Comma)
		if !ok {
			break
		}
		commas = append(commas, comma)
	}

	end, _ := p.expect(lex.BraceR)

	return &syn.ObjectExpr{
		Position:   start.Merge(end),
		OpenBrace:  start,
		Fields:     fields,
		Commas:     commas,
		Base:       base,
		CloseBrace: end,
	}
}

func (p *Parser) parseLiteralExpr() *syn.LiteralExpr {
	tok := p.peek()
	return &syn.LiteralExpr{Token: tok, Value: p.parseLiteral()}
}

func (p *Parser) parseLiteral() any {
	if s, ok := p.eat(lex.Int); ok {
		// we only deal with float64s as values for simplicity
		return float64(must(strconv.Atoi(s.Value)))
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

func (p *Parser) parseQueryExpr() *syn.QueryExpr {
	parser := parser.New(&p.lexer)
	startPos := p.lexer.Peek()
	operation := parser.ParseOperationDefinition()
	if err := parser.Err(); err != nil {
		// we can make this conversion as this is the only error type returned by parser
		// (excluding lexer errors which we have already handled)
		err := err.(*gqlerror.Error)
		p.errors = append(p.errors, ast.Error{
			Position: startPos.Pos(),
			Msg:      err.Message,
		})
		return nil
	}

	endTok := p.lexer.Peek()
	query := strings.TrimRight(endTok.Src.Input[startPos.Start:endTok.Start], "\n")

	return &syn.QueryExpr{Position: startPos.Merge(endTok), Query: query, Operation: operation}
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
