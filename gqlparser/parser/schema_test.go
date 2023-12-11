package parser

import (
	"testing"

	"github.com/andyyu2004/gqlt/gqlparser/gqlerror"
	"github.com/andyyu2004/gqlt/syn"
	"github.com/stretchr/testify/assert"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/gqlparser/parser/testrunner"
)

func TestSchemaDocument(t *testing.T) {
	testrunner.Test(t, "schema_test.yml", func(_ *testing.T, input string) testrunner.Spec {
		doc, err := ParseSchema(&ast.Source{Input: input, Name: "spec"})
		if err != nil {
			return testrunner.Spec{
				Error: err.(*gqlerror.Error),
				AST:   syn.Dump(doc),
			}
		}
		return testrunner.Spec{
			AST: syn.Dump(doc),
		}
	})
}

func TestTypePosition(t *testing.T) {
	t.Run("type line number with no bang", func(t *testing.T) {
		schema, parseErr := ParseSchema(&ast.Source{
			Input: `type query {
						me: User
					}
			`,
		})
		assert.Nil(t, parseErr)
		assert.Equal(t, 2, schema.Definitions.ForName("query").Fields.ForName("me").Type.Position.Line)
	})
	t.Run("type line number with bang", func(t *testing.T) {
		schema, parseErr := ParseSchema(&ast.Source{
			Input: `type query {
						me: User!
					}
			`,
		})
		assert.Nil(t, parseErr)
		assert.Equal(t, 2, schema.Definitions.ForName("query").Fields.ForName("me").Type.Position.Line)
	})
	t.Run("type line number with comments", func(t *testing.T) {
		schema, parseErr := ParseSchema(&ast.Source{
			Input: `type query {
            # comment
						me: User
					}
			`,
		})
		assert.Nil(t, parseErr)
		assert.Equal(t, 3, schema.Definitions.ForName("query").Fields.ForName("me").Type.Position.Line)
	})
}
