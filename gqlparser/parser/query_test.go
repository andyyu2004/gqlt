package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/gqlparser/gqlerror"
	"github.com/andyyu2004/gqlt/gqlparser/parser/testrunner"
	"github.com/andyyu2004/gqlt/syn"
)

func TestQueryDocument(t *testing.T) {
	testrunner.Test(t, "query_test.yml", func(_ *testing.T, input string) testrunner.Spec {
		doc, err := ParseQuery(&ast.Source{Input: input, Name: "spec"})
		if err != nil {
			gqlErr := err.(*gqlerror.Error)
			return testrunner.Spec{
				Error: gqlErr,
				AST:   syn.Dump(doc),
			}
		}
		return testrunner.Spec{
			AST: syn.Dump(doc),
		}
	})
}

func TestQueryPosition(t *testing.T) {
	t.Run("query line number with comments", func(t *testing.T) {
		query, err := ParseQuery(&ast.Source{
			Input: `
	# comment 1
query SomeOperation {
	# comment 2
	myAction {
		id
	}
}
      `,
		})
		assert.Nil(t, err)
		assert.Equal(t, 3, query.Operations.ForName("SomeOperation").Position.Line)
		assert.Equal(t, 5, query.Operations.ForName("SomeOperation").SelectionSet[0].GetPosition().Line)
	})
}
