package model

type TableNodeData struct {
	Project  string `json:"project"`
	Database string `json:"database"`
	Schema   string `json:"schema"`
	Name     string `json:"name"`
	Version  string `json:"version"`

	// raw graph metrics
	DownstreamImpact int     `json:"downstream_impact"`
	Betweenness      float64 `json:"betweenness"`
	PageRank         float64 `json:"pagerank"`
	KCore            int     `json:"kcore"`
	Community        int     `json:"community"`
	Leiden           int     `json:"leiden"`
	Component        int     `json:"component"`
	Degree           float64 `json:"degree"`

	// derived fields
	RiskScore       float64 `json:"risk_score"`
	CriticalityTier string  `json:"criticality_tier"`
	DomainLabel     string  `json:"domain_label"`
}

// Used as Node.Data when Type == column
type ColumnNodeData struct {
	Project  string `json:"project"`
	Database string `json:"database"`
	Schema   string `json:"schema"`
	Table    string `json:"table"`
	Name     string `json:"name"`
	Version  string `json:"version"`

	// raw graph metrics (no community, component — not on CSG nodes)
	DownstreamImpact int     `json:"downstream_impact"`
	Betweenness      float64 `json:"betweenness"`
	PageRank         float64 `json:"pagerank"`
	KCore            int     `json:"kcore"`
	Leiden           int     `json:"leiden"`
	Degree           float64 `json:"degree"`

	// derived fields
	RiskScore       float64 `json:"risk_score"`
	CriticalityTier string  `json:"criticality_tier"`
	DomainLabel     string  `json:"domain_label"`
}
