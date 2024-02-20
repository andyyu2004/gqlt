package eval

import (
	"fmt"
	"maps"

	"github.com/movio/gqlt/syn"
)

type binder interface {
	Lookup(name string) (any, bool)
	bind(name string, val any)
}

type noopBinder struct {
	binder
}

var _ binder = noopBinder{}

// override the noopBinder's bind method to do nothing
func (noopBinder) bind(string, any) {}

func bindPat(binder binder, pat syn.Pat, val any) error {
	switch pat := pat.(type) {
	case *syn.WildcardPat:
		return nil

	case *syn.NamePat:
		binder.bind(pat.Name.Value, val)
		return nil

	case *syn.ListPat:
		vals, ok := val.([]any)
		if !ok {
			return errorf(pat, "cannot bind %T value `%v` to list pattern", val, val)
		}
		return bindListPat(binder, pat, vals)

	case *syn.ObjectPat:
		vals, ok := val.(map[string]any)
		if !ok {
			return errorf(pat, "cannot bind %T value `%v` to object pattern", val, val)
		}

		return bindObjectPat(binder, pat, vals)

	case *syn.LiteralPat:
		if pat.Value != val {
			return errorf(pat, "literal pattern does not match value: %v != %v", pat.Value, val)
		}
		return nil

	case *syn.VarPat:
		target, ok := binder.Lookup(pat.Name.Value)
		if !ok {
			return errorf(pat, "variable pattern references undefined variable: '$%s'", pat.Name.Value)
		}

		if val != target {
			return errorf(pat, "variable pattern does not match value: %v != %v", target, val)
		}

		return nil

	case *syn.RestPat:
		panic("rest pattern should have special handling in list and object cases (eval)")

	default:
		panic(fmt.Sprintf("missing pattern bind case: %T", pat))
	}
}

func bindListPat(binder binder, listPat *syn.ListPat, values []any) error {
	for i, pat := range listPat.Pats {
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

func bindObjectPat(binder binder, objPat *syn.ObjectPat, values map[string]any) error {
	// Important to clone, we can't modify in place because the value may be a variable and used multiple times
	values = maps.Clone(values)

	for entry := objPat.Fields.Oldest(); entry != nil; entry = entry.Next() {
		name := entry.Key
		pat := entry.Value

		rest, ok := pat.(*syn.RestPat)
		if ok {
			if err := bindPat(binder, rest.Pat, values); err != nil {
				return err
			}
			return nil
		}

		val, ok := values[name.Value]
		if !ok {
			return errorf(entry.Key, "object missing field specified in pattern `%s`", name.Value)
		}

		if err := bindPat(binder, pat, val); err != nil {
			return err
		}

		// delete the value once it's been bound for the rest pattern
		// it's fine to mutate the map in place as we cloned it above
		delete(values, name.Value)
	}

	return nil
}
