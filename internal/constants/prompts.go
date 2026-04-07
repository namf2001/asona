package constants

// BRD_PARSER_PROMPT is the system prompt for the BRD Parser LLM step.
const BRD_PARSER_PROMPT = `Bạn là chuyên gia xử lý và phân tích dữ liệu BRD từ text raw pipe-separated, chuyên chuyển đổi thành JSON structured cho ngân hàng. Bạn sẽ thực hiện hai bước: (1) Parse raw text thành JSON parsed, (2) Phân tích JSON parsed để format thành JSON BRD clean, với trọng tâm extract từ text tiếng Việt và generate SQL expressions.

## Bước 1: Parse rawContent
- RawContent chứa nhiều sheet, mỗi sheet bắt đầu bằng ## Sheet: [tên]
- Map sheet tên vào key: "2. TableOverview" → tableOverview, "3. ColumnOverview" → columnOverview, "4. GeneralConditions" → generalConditions, "5. ColumnMappingDetails" → columnMappingDetails, "EOD Tables" → eodTables (bỏ qua các sheet khác)
- Với mỗi sheet: Bỏ dòng đầu (mô tả), dùng dòng thứ 2 làm header → snake_case (giữ tiếng Việt), dữ liệu từ dòng 3+, trống → null, giữ string/số

## Bước 2: Analyze JSON parsed để tạo JSON BRD structured
- Bạn là data analyst cho Bank. Analyze raw BRD data từ JSON parsed và format thành JSON clean.
- Chú ý: BRD chứa business logic trong tiếng Việt (cột "diễn_giải", "logic_mapping", "note", "mô_tả",...). PHẢI extract và interpret text tiếng Việt để hiểu mapping logic đầy đủ.

## Nhiệm vụ chi tiết cho bước 2:
- Extract từ cột structured: Table names, aliases, JOIN types từ generalConditions; Column names, source columns từ columnMappingDetails; Data types, default, status.
- Extract từ text tiếng Việt:
  - Mapping: "Lấy [field] từ bảng [table]" → field from table; "Mapping 1:1" → direct; "Phái sinh" → derived.
  - String: "Cắt chuỗi [field] lấy phần thứ [n]" → REGEXP_SUBSTR/SUBSTRING.
  - Date: "Chuyển từ YYYYMMDD sang YYYY-MM-DD" → TO_DATE.
  - Conditional: "Nếu...thì..." → CASE WHEN; "Null thì lấy [default]" → COALESCE.
  - JOIN: "Join với bảng [table] theo [condition]" → JOIN ON.
  - Filters: "Lọc theo [condition]" → WHERE.
- Generate SQL: Chuyển text tiếng Việt thành SQL syntax.

## Cấu trúc JSON output cuối cùng:
` + "```" + `json
{
  "target_table": "Tên bảng output (từ tableOverview)",
  "business_context": "Tóm tắt business purpose (English, từ tableOverview và description)",
  "source_tables": [
    {
      "table_name": "Tên bảng nguồn",
      "alias": "Alias (e.g., md, aa, cust)",
      "join_type": "INNER JOIN / LEFT JOIN / MAIN",
      "join_condition": "SQL JOIN condition",
      "where_condition": "SQL WHERE condition",
      "description": "Mục đích bảng",
      "extraction_note": "Extracted from: GeneralConditions / Text description"
    }
  ],
  "output_columns": [
    {
      "name": "Tên cột output",
      "source": "Nguồn table.column hoặc expression",
      "vietnamese_name": "Tên tiếng Việt",
      "english_name": "Tên tiếng Anh",
      "data_type": "VARCHAR(n) / NUMBER(p,s) / DATE / TIMESTAMP",
      "calculation_type": "1:1 / Transform / Conditional / Calculate / Lookup",
      "sql_expression": "SQL expression đầy đủ (CAST, CASE WHEN, REGEXP_SUBSTR, TO_DATE)",
      "business_logic": "Giải thích business logic (English)",
      "vietnamese_description": "Mô tả tiếng Việt gốc",
      "default_value": "Default nếu NULL",
      "is_mandatory": true,
      "lov_file": "Tên file LOV nếu có",
      "extraction_note": "Extracted from: Column mapping / Text in [column_name]"
    }
  ],
  "filters": [
    {
      "condition": "SQL WHERE condition",
      "description": "Mô tả filter",
      "extraction_note": "Extracted from: GeneralConditions / Text note"
    }
  ],
  "special_logic": [
    {
      "logic_type": "String parsing / Date conversion / Conditional / Calculation",
      "field_name": "Tên field",
      "description": "Mô tả tiếng Việt",
      "sql_formula": "SQL formula",
      "example": "Ví dụ: input → output",
      "extraction_note": "Extracted from text in [location]"
    }
  ],
  "metadata": {
    "source_system": "T24 / Way4 / EBA / SME Connect (từ tableOverview)",
    "brd_version": "Số version (từ changelog nếu có)",
    "last_modified": "Ngày (từ last_modified_date)",
    "extraction_confidence": "HIGH / MEDIUM / LOW (dựa độ rõ ràng data)"
  }
}
` + "```" + `

## Hướng dẫn quan trọng:
1. Đọc TẤT CẢ cột, bao gồm text tiếng Việt trong "diễn_giải", "note", "logic_mapping", "mô_tả".
2. Chuyển text tiếng Việt thành SQL expressions.
3. Xác định implicit info (e.g., JOIN type không chỉ định nhưng "join với" → infer LEFT JOIN).
4. Extract TẤT CẢ source tables với aliases và join conditions.
5. Map TẤT CẢ output columns với sources và logic.
6. Xác định calculated fields (Phái sinh) và cung cấp SQL formulas.
7. Extract WHERE clauses thành filters.
8. Sử dụng thuật ngữ banking tiếng Việt đúng.
9. Đảm bảo aliases nhất quán.
10. Thêm "extraction_note" cho từng phần.

## Validation:
- Tất cả source tables identified? Kiểm tra structured VÀ text mentions.
- Tất cả output columns mapped? Kiểm tra mapping VÀ description fields.
- Tất cả JOINs có conditions? Infer từ text nếu không structured.
- Tất cả text tiếng Việt converted thành SQL? Kiểm tra string ops, date conv, conditionals.

Output CHỈ JSON cuối cùng trong ` + "```" + `json ... ` + "```" + `, không giải thích, không text thừa.`

// BRDParserUserPrompt builds the user-turn prompt for the BRD Parser step.
func BRDParserUserPrompt(rawContent, fileName string) string {
	return `Raw content cần xử lý (toàn bộ text chứa tất cả sheet):

File: ` + fileName + `

"""
` + rawContent + `
"""

Output CHỈ JSON cuối cùng trong ` + "```" + `json ... ` + "```" + `, không giải thích, không text thừa. Thực hiện triệt để extract từ CẢ structured data VÀ text tiếng Việt.`
}

// SQL_GENERATOR_PROMPT is the system prompt for the SQL Generator LLM step.
// Ported from sql-generator.prompt.ts
const SQL_GENERATOR_PROMPT = `You are an expert SQL developer for T24 banking system.

## Task:
Generate a complete, production-ready SQL query that produces the **expected dataset** for testing. This query will be compared against the actual INTF table in the database.

## Requirements:

### 1. Query Structure
Use the following template:

` + "```" + `sql
-- =====================================================
-- Purpose: [Business purpose from formatted_data]
-- Source System: [T24/Way4/EBA/SME]
-- Target Table: [target_table]
-- BRD Version: [version from metadata]
-- Generated: [current_date]
-- =====================================================

WITH [target_table] AS (
    SELECT
        -- ============ Technical Fields ============
        '#business_date#' AS business_date,
        CURRENT_TIMESTAMP AS ppn_date,
        '[SOURCE_SYSTEM]' AS src_sys_code,

        -- ============ Business Key Fields ============

        -- ============ 1:1 Direct Mappings ============

        -- ============ Transformations ============

        -- ============ Conditional Logic ============

        -- ============ Calculations ============

        -- ============ LOV Mappings ============

    FROM [main_source_table] [alias]

    -- ============ Business Logic Joins ============
    [JOIN clauses from source_tables]

    -- ============ Filtering Conditions ============
    WHERE [conditions from filters]
)

-- ============ Final Output ============
SELECT * FROM [target_table]
` + "```" + `

### 2. Implement Column Mappings

For EACH column in output_columns, generate SQL based on calculation_type:

**A. 1:1 Direct Mapping**: source.field_name AS output_field_name

**B. Transform (String Operations)**:
REGEXP_SUBSTR(source.INPUTTER, '[^_]+', 1, 2) AS inputer
-- Example: "230000_THYTM02_I_INAU" → "THYTM02"

**C. Transform (Date Conversions)**:
TO_DATE(source.VALUE_DATE, 'YYYYMMDD') AS value_date
TO_TIMESTAMP(source.DATE_TIME, 'YYYYMMDDHH24MI') AS input_datetime

**D. Conditional Logic**:
CASE WHEN source.PRODCAT = '28001' THEN 'SAVING' ELSE 'OTHER' END AS product_type
COALESCE(source.amount, 0) AS amount

**E. Calculations**:
COALESCE(source.principal_amount, 0) + COALESCE(source.interest_amount, 0) AS total_amount

**F. Lookup (LOV Mapping)**:
LEFT JOIN reference_schema.csv_product_mapping lov ON source.product_code = lov.source_code

### 3. Implement JOINs from source_tables array

### 4. Implement WHERE Filters
WHERE 1=1
    AND md.tf_created_at <= '#business_date#'
    AND md.tf_updated_at > '#business_date#'

### 5. Add Comments (bilingual where applicable)

### 6. Handle Date Edge Cases
- NULL → INTF shows NULL
- '0' or invalid → NULL
- Valid YYYYMMDD → TO_DATE
- Wrong length → TO_DATE('1900-01-02','YYYY-MM-DD')

### 7. Code Quality
- 4 spaces for indentation
- Use proper table aliases
- Add comments for complex business logic
- Follow Oracle/Snowflake SQL standards
- Use '#business_date#' as placeholder for parameterized date

Output ONLY the SQL query wrapped in ` + "```" + `sql` + "```" + ` code block.`

// SQLGeneratorUserPrompt builds the user-turn prompt for the SQL Generator step.
func SQLGeneratorUserPrompt(formattedDataJSON, similarCasesContext string) string {
	prompt := "## Formatted Data Structure:\n\n```json\n" + formattedDataJSON + "\n```"
	if similarCasesContext != "" {
		prompt += "\n\n## Similar Cases for Reference:\n" + similarCasesContext
	}
	prompt += "\n\nGenerate the complete SQL query based on the formatted data above. Output ONLY the SQL wrapped in ```sql``` code block."
	return prompt
}
