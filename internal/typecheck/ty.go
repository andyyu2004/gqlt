package typecheck

import (
	"fmt"
	"strings"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type Ty interface {
	isTy()
	fmt.Stringer
}

type Bool struct{}

func (Bool) String() string {
	return "bool"
}

func (Bool) isTy() {}

type Number struct{}

func (Number) String() string {
	return "number"
}

func (Number) isTy() {}

type String struct{}

func (String) String() string {
	return "string"
}

func (String) isTy() {}

type List struct {
	Elem Ty
}

func (l List) String() string {
	return fmt.Sprintf("%s[]", l.Elem)
}

func (List) isTy() {}

type Tuple struct {
	Elems []Ty
}

func (t Tuple) String() string {
	s := new(strings.Builder)
	_, _ = s.WriteString("[")
	for i, ty := range t.Elems {
		if i > 0 {
			_, _ = s.WriteString(", ")
		}
		_, _ = s.WriteString(ty.String())
	}
	_, _ = s.WriteString("]")
	return s.String()
}

func (Tuple) isTy() {}

type Object struct {
	Fields *orderedmap.OrderedMap[string, Ty]
}

func (o Object) String() string {
	s := new(strings.Builder)
	_, _ = s.WriteString("{ ")

	first := true
	for entry := o.Fields.Oldest(); entry != nil; entry = entry.Next() {
		if !first {
			_, _ = s.WriteString(", ")
		}
		first = false
		_, _ = fmt.Fprintf(s, "%s: %s", entry.Key, entry.Value)
	}
	_, _ = s.WriteString(" }")
	return s.String()
}

func (Object) isTy() {}

type Any struct{}

func (Any) String() string {
	return "any"
}

func (Any) isTy() {}

type errTy struct{}

func (errTy) String() string {
	return "error"
}

func (errTy) isTy() {}

func compat(a, b Ty) bool {
	// shortcut if either is an error to prevent reporting unnecessary errors
	if _, ok := b.(errTy); ok {
		return true
	}

	if _, ok := b.(Any); ok {
		return true
	}

	switch a := a.(type) {
	case errTy, Any:
		return true
	case Bool, Number, String:
		return a == b
	case List:
		if t, ok := b.(Tuple); ok {
			// empty tuple is compatible with any list
			return len(t.Elems) == 0
		}

		b, ok := b.(List)
		return ok && compat(a.Elem, b.Elem)
	case Tuple:
		if _, ok := b.(List); ok {
			return len(a.Elems) == 0
		}

		b, ok := b.(Tuple)
		return ok && len(a.Elems) == len(b.Elems) && allCompat(a.Elems, b.Elems)
	case Object:
		b, ok := b.(Object)
		if !ok {
			return false
		}

		if a.Fields.Len() != b.Fields.Len() {
			return false
		}

		for entry := a.Fields.Oldest(); entry != nil; entry = entry.Next() {
			v, ok := b.Fields.Get(entry.Key)
			if !ok {
				return false
			}

			if !compat(entry.Value, v) {
				return false
			}
		}

		return true

	default:
		panic(fmt.Sprintf("missing case identical %T = %T", a, b))
	}
}

func allCompat(a, b []Ty) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !compat(a[i], b[i]) {
			return false
		}
	}

	return true
}

func isAny(ty Ty) bool {
	_, ok := ty.(Any)
	return ok
}

func isErr(ty Ty) bool {
	_, ok := ty.(errTy)
	return ok
}
