package syn

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/movio/gqlt/gqlparser/lexer"
	"github.com/movio/gqlt/internal/lex"
	"github.com/stretchr/testify/require"
)

func TestDump(t *testing.T) {
	res := Dump(SchemaDefinition{
		Directives: []*Directive{
			{
				Name:      "foo",
				Arguments: []*Argument{{Name: lexer.Token{Value: "bar"}}},
			},
			{Arguments: []*Argument{}},
		},
	})

	expected := `<SchemaDefinition>
  Directives: [Directive]
  - <Directive>
      Name: "foo"
      Arguments: [Argument]
      - <Argument>
          Name: "bar"
  - <Directive>`

	fmt.Println(diff.LineDiff(expected, res))
	require.Equal(t, expected, res)

	require.True(t, shouldSkip(reflect.ValueOf(lexer.Token{})))
	require.True(t, shouldSkip(reflect.ValueOf(lex.Token{})))
}
