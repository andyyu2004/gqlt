package gqlt

import (
	"fmt"

	"github.com/andyyu2004/gqlt/syn"
)

type binder interface {
	bind(name string, val any)
}

type dummyBinder struct{}

var _ binder = dummyBinder{}

func (dummyBinder) bind(string, any) {}

func bindPat(binder binder, pat syn.Pat, val any) error {
	switch pat := pat.(type) {
	case *syn.WildcardPat:
		return nil

	case *syn.NamePat:
		binder.bind(pat.Name, val)
		return nil

	case *syn.ListPat:
		vals, ok := val.([]any)
		if !ok {
			return fmt.Errorf("cannot bind %T value `%v` to list pattern", val, val)
		}
		return bindListPat(binder, pat, vals)

	case *syn.ObjectPat:
		vals, ok := val.(map[string]any)
		if !ok {
			return fmt.Errorf("cannot bind %T value `%v` to object pattern", val, val)
		}

		return bindObjectPat(binder, pat, vals)

	case *syn.LiteralPat:
		if pat.Value != val {
			return fmt.Errorf("literal pattern does not match value: %v != %v", pat.Value, val)
		}
		return nil

	case *syn.RestPat:
		panic("rest pattern should have special handling in list and object cases")

	default:
		panic(fmt.Sprintf("missing pattern bind case: %T", pat))
	}
}

func bindListPat(binder binder, pat *syn.ListPat, values []any) error {
	for i, pat := range pat.Pats {
		rest, ok := pat.(*syn.RestPat)
		if ok {
			if err := bindPat(binder, rest.Pat, values[i:]); err != nil {
				return err
			}
			return nil

		}

		if i > len(values)-1 {
			return bindPat(binder, pat, nil)
		}

		if err := bindPat(binder, pat, values[i]); err != nil {
			return err
		}
	}

	return nil
}

func bindObjectPat(binder binder, pat *syn.ObjectPat, values map[string]any) error {
	for entry := pat.Fields.Oldest(); entry != nil; entry = entry.Next() {
		name := entry.Key
		pat := entry.Value

		val, ok := values[name]
		if !ok {
			return fmt.Errorf("object missing field specified in pattern %s", name)
		}

		if err := bindPat(binder, pat, val); err != nil {
			return err
		}
	}

	return nil
}
