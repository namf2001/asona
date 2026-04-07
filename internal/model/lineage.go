package model

type Lineage struct {
	Graph Graph      `json:"graph"`
	Stats Stats      `json:"stats"`
	Meta  *GraphMeta `json:"meta,omitempty"`
}

type Stats struct {
	NodeCount int `json:"node_count"`
	EdgeCount int `json:"edge_count"`
	MaxDepth  int `json:"max_depth"`
}

type GraphMeta struct {
	Direction string `json:"direction"`
	Depth     int    `json:"depth"`
	Partial   bool   `json:"partial"`
}

type Pagination struct {
	HasMore bool `json:"has_more"`
}
