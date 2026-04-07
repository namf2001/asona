package model

// ─── BRD Parser Output ──────────────────────────────────────────────────────

// FormattedBRDData is the structured output produced by the BRD Parser LLM step.
type FormattedBRDData struct {
	TargetTable     string          `json:"target_table"`
	BusinessContext string          `json:"business_context"`
	SourceTables    []BRDSourceTable `json:"source_tables"`
	OutputColumns   []BRDColumn      `json:"output_columns"`
	Filters         []BRDFilter      `json:"filters"`
	SpecialLogic    []BRDSpecialLogic `json:"special_logic"`
	Metadata        BRDMetadata      `json:"metadata"`
}

type BRDSourceTable struct {
	TableName      string `json:"table_name"`
	Alias          string `json:"alias"`
	JoinType       string `json:"join_type"`
	JoinCondition  string `json:"join_condition"`
	WhereCondition string `json:"where_condition"`
	Description    string `json:"description"`
	ExtractionNote string `json:"extraction_note"`
}

type BRDColumn struct {
	Name                 string `json:"name"`
	Source               string `json:"source"`
	VietnameseName       string `json:"vietnamese_name"`
	EnglishName          string `json:"english_name"`
	DataType             string `json:"data_type"`
	CalculationType      string `json:"calculation_type"`
	SQLExpression        string `json:"sql_expression"`
	BusinessLogic        string `json:"business_logic"`
	VietnameseDescription string `json:"vietnamese_description"`
	DefaultValue         string `json:"default_value"`
	IsMandatory          bool   `json:"is_mandatory"`
	LOVFile              string `json:"lov_file"`
	ExtractionNote       string `json:"extraction_note"`
}

type BRDFilter struct {
	Condition      string `json:"condition"`
	Description    string `json:"description"`
	ExtractionNote string `json:"extraction_note"`
}

type BRDSpecialLogic struct {
	LogicType      string `json:"logic_type"`
	FieldName      string `json:"field_name"`
	Description    string `json:"description"`
	SQLFormula     string `json:"sql_formula"`
	Example        string `json:"example"`
	ExtractionNote string `json:"extraction_note"`
}

type BRDMetadata struct {
	SourceSystem          string `json:"source_system"`
	BRDVersion            string `json:"brd_version"`
	LastModified          string `json:"last_modified"`
	ExtractionConfidence  string `json:"extraction_confidence"`
}

// ─── Excel Parser Output ─────────────────────────────────────────────────────

// ExcelParseResult is returned by the Excel parser step.
type ExcelParseResult struct {
	FileName   string        `json:"file_name"`
	RawContent string        `json:"raw_content"`
	Metadata   ExcelMetadata `json:"metadata"`
}

type ExcelMetadata struct {
	TotalSheets int      `json:"total_sheets"`
	SheetNames  []string `json:"sheet_names"`
	TotalRows   int      `json:"total_rows"`
	ExtractedAt string   `json:"extracted_at"`
}

// ─── Vector Store ─────────────────────────────────────────────────────────────

// BRDCaseDocument is stored in the Postgres BRDCase node.
type BRDCaseDocument struct {
	ID              string   `json:"id"`
	Content         string   `json:"content"`
	BrdID           string   `json:"brd_id"`
	BrdVersion      string   `json:"brd_version"`
	TargetTable     string   `json:"target_table"`
	SourceSystem    string   `json:"source_system"`
	BusinessContext string   `json:"business_context"`
	TableNames      []string `json:"table_names"`
	ColumnCount     int      `json:"column_count"`
	SQLQuery        string   `json:"sql_query"`
	SQLComplexity   string   `json:"sql_complexity"`
	TransformTypes  []string `json:"transform_types"`
	IsVerified      bool     `json:"is_verified"`
	UsageCount      int      `json:"usage_count"`
	SuccessRate     float64  `json:"success_rate"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
}

// SimilarCaseResult is what the orchestrator uses as context for the SQL Generator.
type SimilarCaseResult struct {
	BrdSummary      string   `json:"brd_summary"`
	SQLExample      string   `json:"sql_example"`
	Similarity      float64  `json:"similarity"`
	TableNames      []string `json:"table_names"`
	BusinessContext string   `json:"business_context"`
}

// ─── Agent Steps ─────────────────────────────────────────────────────────────

// AgentStep tracks one step of the pipeline for observable output.
type AgentStep struct {
	Tool      string      `json:"tool"`
	Input     interface{} `json:"input"`
	Output    interface{} `json:"output"`
	Timestamp string      `json:"timestamp"`
}

// ─── SQL Validation ──────────────────────────────────────────────────────────

// SQLValidationResult is the output of the SQL Validator step.
type SQLValidationResult struct {
	IsValid     bool     `json:"is_valid"`
	Issues      []string `json:"issues"`
	Suggestions []string `json:"suggestions"`
	Dialect     string   `json:"dialect"`
}

// ─── Final Response ──────────────────────────────────────────────────────────

// SQLGenerationResult is the final response returned to the client.
type SQLGenerationResult struct {
	SQL           string              `json:"sql"`
	Explanation   string              `json:"explanation"`
	Confidence    float64             `json:"confidence"`
	Steps         []AgentStep         `json:"steps"`
	SimilarCases  []SimilarCaseResult `json:"similar_cases"`
	FormattedData FormattedBRDData    `json:"formatted_data"`
	Validation    SQLValidationResult `json:"validation"`
}
