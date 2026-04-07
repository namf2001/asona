package model

type NodeType string
type NodeLable string
type EdgeType string
type LineageDirection string

const (
	NodeTypeTable  NodeType = "TABLE"
	NodeTypeColumn NodeType = "COLUMN"

	SourceSemanticGroup    NodeLable = "SSG"
	MutualSemanticGroup    NodeLable = "MSG"
	InterfaceSemanticGroup NodeLable = "ISG"
	ProductSemanticGroup   NodeLable = "PSG"
	ColumnSemanticGroup    NodeLable = "CSG"

	EdgeTypeDerive    EdgeType = "DERIVE"
	EdgeTypeTransform EdgeType = "TRANSFORM"
	EdgeTypeImpact    EdgeType = "IMPACT"
	EdgeTypeHas       EdgeType = "HAS"

	LineageUpStream   LineageDirection = "UPSTREAM"
	LineageDownStream LineageDirection = "DOWNSTREAM"
	LineageFull       LineageDirection = "FULL"
)
