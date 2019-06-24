package checker

import (
	"gopkg.in/antonmedv/expr.v2/ast"
	"gopkg.in/antonmedv/expr.v2/internal/conf"
	"gopkg.in/antonmedv/expr.v2/parser"
)

type operatorPatcher struct {
	ast.BaseVisitor
	ops   map[string][]string
	types conf.TypesTable
}

func (p *operatorPatcher) Node(node *ast.Node) {
	binaryNode, ok := (*node).(*ast.BinaryNode)
	if !ok {
		return
	}

	fns, ok := p.ops[binaryNode.Operator]
	if !ok {
		return
	}

	leftType := binaryNode.Left.GetType()
	rightType := binaryNode.Right.GetType()
	for _, fn := range fns {
		fnType := p.types[fn]

		firstArgType := fnType.Type.In(0)
		secondArgType := fnType.Type.In(1)

		if leftType == firstArgType && rightType == secondArgType {
			*node = &ast.FunctionNode{
				Name:      fn,
				Arguments: []ast.Node{binaryNode.Left, binaryNode.Right},
			}
		}
	}
}

func patchOperators(tree *parser.Tree, config *conf.Config) {
	if len(config.Operators) == 0 {
		return
	}
	patcher := &operatorPatcher{ops: config.Operators, types: config.Types}
	ast.Walk(&tree.Node, patcher)
}
