package model

type CriticalityEntry struct {
	Name             string  `json:"name"`
	Project          string  `json:"project"`
	RiskScore        float64 `json:"risk_score"`
	CriticalityTier  string  `json:"criticality_tier"`
	DownstreamImpact int     `json:"downstream_impact"`
	Betweenness      float64 `json:"betweenness"`
	KCore            int     `json:"kcore"`
	PageRank         float64 `json:"pagerank"`
	SemanticLayer    string  `json:"semantic_layer"`
}

type CommunityEntry struct {
	CommunityID             int      `json:"community_id"`
	Label                   string   `json:"label"`
	NodeCount               int      `json:"node_count"`
	TopNodes                []string `json:"top_nodes"`
	CrossCommunityEdgeCount int      `json:"cross_community_edge_count"`
}

type ProvenanceStep struct {
	NodeName      string `json:"node_name"`
	ColumnName    string `json:"column_name"`
	SemanticLayer string `json:"semantic_layer"`
	TransformType string `json:"transform_type"`
}

type ProvenanceResult struct {
	Column      string             `json:"column"`
	Table       string             `json:"table"`
	Path        []ProvenanceStep   `json:"path"`  // longest single path (backward compat)
	Paths       [][]ProvenanceStep `json:"paths"` // all upstream source paths
	SourceLayer string             `json:"source_layer"`
}

type TierBreakdown struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Standard int `json:"standard"`
}

type BlastRadiusResult struct {
	Table         string        `json:"table"`
	Project       string        `json:"project"`
	TotalImpact   int           `json:"total_impact"`
	TierBreakdown TierBreakdown `json:"tier_breakdown"`
}

type CrossDomainEdge struct {
	SourceName           string `json:"source_name"`
	TargetName           string `json:"target_name"`
	SourceCommunityLabel string `json:"source_community_label"`
	TargetCommunityLabel string `json:"target_community_label"`
}

type ProjectHealthEntry struct {
	Project             string  `json:"project"`
	NodeCount           int     `json:"node_count"`
	AvgCriticalityScore float64 `json:"avg_criticality_score"`
	TopBottleneck       string  `json:"top_bottleneck"`
	DeadCodeCount       int     `json:"dead_code_count"`
}

type PlatformHealthReport struct {
	TotalNodes                int                  `json:"total_nodes"`
	TotalEdges                int                  `json:"total_edges"`
	NodesByLayer              map[string]int       `json:"nodes_by_layer"`
	AvgDownstreamImpactPSG    float64              `json:"avg_downstream_impact_psg"`
	IsolatedNodeCount         int                  `json:"isolated_node_count"`
	CrossCommunityEdgeDensity float64              `json:"cross_community_edge_density"`
	DeadCodeCount             int                  `json:"dead_code_count"`
	LineageCoveragePct        float64              `json:"lineage_coverage_pct"`
	PerProject                []ProjectHealthEntry `json:"per_project"`
}