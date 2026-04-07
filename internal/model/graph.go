package model

type Node struct {
	ID    string      `json:"id"`
	Type  string      `json:"type"`
	Label string      `json:"label"`
	Data  interface{} `json:"data"`
	Flags NodeFlags   `json:"flags"`
}

type NodeFlags struct {
	Expandable bool `json:"expandable"`
	Expanded   bool `json:"expanded"`
}

type Edge struct {
	ID       string                 `json:"id"`
	Type     EdgeType               `json:"type"`
	From     string                 `json:"from"`
	To       string                 `json:"to"`
	Metadata map[string]interface{} `json:"metadata"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}