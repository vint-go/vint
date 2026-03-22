package no_duplicate_code

import (
	"go/ast"
	"go/token"
)

// AST node type constants, matching dupl's golang package for compatibility.
const (
	nodeBad = iota
	nodeFile
	nodeArrayType
	nodeAssignStmt
	nodeBasicLit
	nodeBinaryExpr
	nodeBlockStmt
	nodeBranchStmt
	nodeCallExpr
	nodeCaseClause
	nodeChanType
	nodeCommClause
	nodeCompositeLit
	nodeDeclStmt
	nodeDeferStmt
	nodeEllipsis
	nodeEmptyStmt
	nodeExprStmt
	nodeField
	nodeFieldList
	nodeForStmt
	nodeFuncDecl
	nodeFuncLit
	nodeFuncType
	nodeGenDecl
	nodeGoStmt
	nodeIdent
	nodeIfStmt
	nodeIncDecStmt
	nodeIndexExpr
	nodeInterfaceType
	nodeKeyValueExpr
	nodeLabeledStmt
	nodeMapType
	nodeParenExpr
	nodeRangeStmt
	nodeReturnStmt
	nodeSelectStmt
	nodeSelectorExpr
	nodeSendStmt
	nodeSliceExpr
	nodeStarExpr
	nodeStructType
	nodeSwitchStmt
	nodeTypeAssertExpr
	nodeTypeSpec
	nodeTypeSwitchStmt
	nodeUnaryExpr
	nodeValueSpec

	nodeSentinel = -1
)

// tokenData holds a serialized AST as flat, GC-friendly parallel arrays.
// No pointers, no interfaces — just value types that the GC can skip.
type tokenData struct {
	types []int32 // node type (used by suffix tree)
	pos   []int32 // byte offset of node start in source
	end   []int32 // byte offset of node end in source
	owns  []int32 // number of descendant tokens (for syntax unit detection)
}

func (td *tokenData) len() int { return len(td.types) }

// fileSpan records which slice of the global token sequence belongs to which file.
type fileSpan struct {
	start int    // first index in global tokenData
	end   int    // one past last index
	file  string // source filename
}

// serializeAST converts an already-parsed ast.File into flat token arrays.
// This avoids re-parsing from disk (which golang.Parse does) and avoids
// allocating *syntax.Node objects (which causes GC pressure).
func serializeAST(file *ast.File, fset *token.FileSet) tokenData {
	s := serializer{tokenFile: fset.File(file.Pos())}
	s.types = make([]int32, 0, 256)
	s.pos = make([]int32, 0, 256)
	s.end = make([]int32, 0, 256)
	s.owns = make([]int32, 0, 256)
	s.serialize(file)
	return tokenData{
		types: s.types,
		pos:   s.pos,
		end:   s.end,
		owns:  s.owns,
	}
}

type serializer struct {
	tokenFile *token.File // cached per-file, avoids fset.Position() string allocs
	types     []int32
	pos       []int32
	end       []int32
	owns      []int32
}

// serialize recursively walks an AST node and appends to the flat arrays.
// Returns the number of tokens added (including descendants) so the caller
// can compute the "owns" count.
func (s *serializer) serialize(node ast.Node) int {
	idx := len(s.types)

	st, en := node.Pos(), node.End()
	startOff := int32(s.tokenFile.Offset(st))
	endOff := int32(s.tokenFile.Offset(en))

	// Append placeholder — owns filled in after children.
	s.types = append(s.types, 0)
	s.pos = append(s.pos, startOff)
	s.end = append(s.end, endOff)
	s.owns = append(s.owns, 0)

	var childCount int

	switch n := node.(type) {
	case *ast.File:
		s.types[idx] = nodeFile
		for _, decl := range n.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
				continue // skip imports, same as dupl
			}
			childCount += s.serialize(decl)
		}

	case *ast.ArrayType:
		s.types[idx] = nodeArrayType
		if n.Len != nil {
			childCount += s.serialize(n.Len)
		}
		childCount += s.serialize(n.Elt)

	case *ast.AssignStmt:
		s.types[idx] = nodeAssignStmt
		for _, e := range n.Rhs {
			childCount += s.serialize(e)
		}
		for _, e := range n.Lhs {
			childCount += s.serialize(e)
		}

	case *ast.BasicLit:
		s.types[idx] = nodeBasicLit

	case *ast.BinaryExpr:
		s.types[idx] = nodeBinaryExpr
		childCount += s.serialize(n.X)
		childCount += s.serialize(n.Y)

	case *ast.BlockStmt:
		s.types[idx] = nodeBlockStmt
		for _, stmt := range n.List {
			childCount += s.serialize(stmt)
		}

	case *ast.BranchStmt:
		s.types[idx] = nodeBranchStmt
		if n.Label != nil {
			childCount += s.serialize(n.Label)
		}

	case *ast.CallExpr:
		s.types[idx] = nodeCallExpr
		childCount += s.serialize(n.Fun)
		for _, arg := range n.Args {
			childCount += s.serialize(arg)
		}

	case *ast.CaseClause:
		s.types[idx] = nodeCaseClause
		for _, e := range n.List {
			childCount += s.serialize(e)
		}
		for _, stmt := range n.Body {
			childCount += s.serialize(stmt)
		}

	case *ast.ChanType:
		s.types[idx] = nodeChanType
		childCount += s.serialize(n.Value)

	case *ast.CommClause:
		s.types[idx] = nodeCommClause
		if n.Comm != nil {
			childCount += s.serialize(n.Comm)
		}
		for _, stmt := range n.Body {
			childCount += s.serialize(stmt)
		}

	case *ast.CompositeLit:
		s.types[idx] = nodeCompositeLit
		if n.Type != nil {
			childCount += s.serialize(n.Type)
		}
		for i, e := range n.Elts {
			if i > 10000 { // same guard as dupl's maxChildrenSerial
				break
			}
			childCount += s.serialize(e)
		}

	case *ast.DeclStmt:
		s.types[idx] = nodeDeclStmt
		childCount += s.serialize(n.Decl)

	case *ast.DeferStmt:
		s.types[idx] = nodeDeferStmt
		childCount += s.serialize(n.Call)

	case *ast.Ellipsis:
		s.types[idx] = nodeEllipsis
		if n.Elt != nil {
			childCount += s.serialize(n.Elt)
		}

	case *ast.EmptyStmt:
		s.types[idx] = nodeEmptyStmt

	case *ast.ExprStmt:
		s.types[idx] = nodeExprStmt
		childCount += s.serialize(n.X)

	case *ast.Field:
		s.types[idx] = nodeField
		for _, name := range n.Names {
			childCount += s.serialize(name)
		}
		childCount += s.serialize(n.Type)

	case *ast.FieldList:
		s.types[idx] = nodeFieldList
		for _, field := range n.List {
			childCount += s.serialize(field)
		}

	case *ast.ForStmt:
		s.types[idx] = nodeForStmt
		if n.Init != nil {
			childCount += s.serialize(n.Init)
		}
		if n.Cond != nil {
			childCount += s.serialize(n.Cond)
		}
		if n.Post != nil {
			childCount += s.serialize(n.Post)
		}
		childCount += s.serialize(n.Body)

	case *ast.FuncDecl:
		s.types[idx] = nodeFuncDecl
		if n.Recv != nil {
			childCount += s.serialize(n.Recv)
		}
		childCount += s.serialize(n.Name)
		childCount += s.serialize(n.Type)
		if n.Body != nil {
			childCount += s.serialize(n.Body)
		}

	case *ast.FuncLit:
		s.types[idx] = nodeFuncLit
		childCount += s.serialize(n.Type)
		childCount += s.serialize(n.Body)

	case *ast.FuncType:
		s.types[idx] = nodeFuncType
		childCount += s.serialize(n.Params)
		if n.Results != nil {
			childCount += s.serialize(n.Results)
		}

	case *ast.GenDecl:
		s.types[idx] = nodeGenDecl
		for _, spec := range n.Specs {
			childCount += s.serialize(spec)
		}

	case *ast.GoStmt:
		s.types[idx] = nodeGoStmt
		childCount += s.serialize(n.Call)

	case *ast.Ident:
		s.types[idx] = nodeIdent

	case *ast.IfStmt:
		s.types[idx] = nodeIfStmt
		if n.Init != nil {
			childCount += s.serialize(n.Init)
		}
		childCount += s.serialize(n.Cond)
		childCount += s.serialize(n.Body)
		if n.Else != nil {
			childCount += s.serialize(n.Else)
		}

	case *ast.IncDecStmt:
		s.types[idx] = nodeIncDecStmt
		childCount += s.serialize(n.X)

	case *ast.IndexExpr:
		s.types[idx] = nodeIndexExpr
		childCount += s.serialize(n.X)
		childCount += s.serialize(n.Index)

	case *ast.InterfaceType:
		s.types[idx] = nodeInterfaceType
		childCount += s.serialize(n.Methods)

	case *ast.KeyValueExpr:
		s.types[idx] = nodeKeyValueExpr
		childCount += s.serialize(n.Key)
		childCount += s.serialize(n.Value)

	case *ast.LabeledStmt:
		s.types[idx] = nodeLabeledStmt
		childCount += s.serialize(n.Label)
		childCount += s.serialize(n.Stmt)

	case *ast.MapType:
		s.types[idx] = nodeMapType
		childCount += s.serialize(n.Key)
		childCount += s.serialize(n.Value)

	case *ast.ParenExpr:
		s.types[idx] = nodeParenExpr
		childCount += s.serialize(n.X)

	case *ast.RangeStmt:
		s.types[idx] = nodeRangeStmt
		if n.Key != nil {
			childCount += s.serialize(n.Key)
		}
		if n.Value != nil {
			childCount += s.serialize(n.Value)
		}
		childCount += s.serialize(n.X)
		childCount += s.serialize(n.Body)

	case *ast.ReturnStmt:
		s.types[idx] = nodeReturnStmt
		for _, e := range n.Results {
			childCount += s.serialize(e)
		}

	case *ast.SelectStmt:
		s.types[idx] = nodeSelectStmt
		childCount += s.serialize(n.Body)

	case *ast.SelectorExpr:
		s.types[idx] = nodeSelectorExpr
		childCount += s.serialize(n.X)
		childCount += s.serialize(n.Sel)

	case *ast.SendStmt:
		s.types[idx] = nodeSendStmt
		childCount += s.serialize(n.Chan)
		childCount += s.serialize(n.Value)

	case *ast.SliceExpr:
		s.types[idx] = nodeSliceExpr
		childCount += s.serialize(n.X)
		if n.Low != nil {
			childCount += s.serialize(n.Low)
		}
		if n.High != nil {
			childCount += s.serialize(n.High)
		}
		if n.Max != nil {
			childCount += s.serialize(n.Max)
		}

	case *ast.StarExpr:
		s.types[idx] = nodeStarExpr
		childCount += s.serialize(n.X)

	case *ast.StructType:
		s.types[idx] = nodeStructType
		childCount += s.serialize(n.Fields)

	case *ast.SwitchStmt:
		s.types[idx] = nodeSwitchStmt
		if n.Init != nil {
			childCount += s.serialize(n.Init)
		}
		if n.Tag != nil {
			childCount += s.serialize(n.Tag)
		}
		childCount += s.serialize(n.Body)

	case *ast.TypeAssertExpr:
		s.types[idx] = nodeTypeAssertExpr
		childCount += s.serialize(n.X)
		if n.Type != nil {
			childCount += s.serialize(n.Type)
		}

	case *ast.TypeSpec:
		s.types[idx] = nodeTypeSpec
		childCount += s.serialize(n.Name)
		childCount += s.serialize(n.Type)

	case *ast.TypeSwitchStmt:
		s.types[idx] = nodeTypeSwitchStmt
		if n.Init != nil {
			childCount += s.serialize(n.Init)
		}
		childCount += s.serialize(n.Assign)
		childCount += s.serialize(n.Body)

	case *ast.UnaryExpr:
		s.types[idx] = nodeUnaryExpr
		childCount += s.serialize(n.X)

	case *ast.ValueSpec:
		s.types[idx] = nodeValueSpec
		for _, name := range n.Names {
			childCount += s.serialize(name)
		}
		if n.Type != nil {
			childCount += s.serialize(n.Type)
		}
		for _, val := range n.Values {
			childCount += s.serialize(val)
		}

	default:
		s.types[idx] = nodeBad
	}

	s.owns[idx] = int32(childCount)
	return childCount + 1
}
