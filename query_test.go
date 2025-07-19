package sqlbuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test NewPaginationParams
func TestNewPaginationParams(t *testing.T) {
	tests := []struct {
		name           string
		page           int
		limit          int
		expectedPage   int
		expectedLimit  int
		expectedOffset int
	}{
		{
			name:           "valid params",
			page:           2,
			limit:          10,
			expectedPage:   2,
			expectedLimit:  10,
			expectedOffset: 10,
		},
		{
			name:           "zero page defaults to 1",
			page:           0,
			limit:          10,
			expectedPage:   1,
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name:           "negative page defaults to 1",
			page:           -1,
			limit:          10,
			expectedPage:   1,
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name:           "zero limit defaults to 10",
			page:           1,
			limit:          0,
			expectedPage:   1,
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name:           "negative limit defaults to 10",
			page:           1,
			limit:          -5,
			expectedPage:   1,
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name:           "both zero defaults",
			page:           0,
			limit:          0,
			expectedPage:   1,
			expectedLimit:  10,
			expectedOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := NewPaginationParams(tt.page, tt.limit)
			assert.Equal(t, tt.expectedPage, params.Page)
			assert.Equal(t, tt.expectedLimit, params.Limit)
			assert.Equal(t, tt.expectedOffset, params.Offset)
		})
	}
}

// Test NewQueryParams
func TestNewQueryParams(t *testing.T) {
	params := NewQueryParams()
	assert.NotNil(t, params)
	assert.NotNil(t, params.Search)
	assert.NotNil(t, params.Filters)
	assert.NotNil(t, params.Sort)
	assert.Equal(t, 0, len(params.Search))
	assert.Equal(t, 0, len(params.Filters))
	assert.Equal(t, 0, len(params.Sort))
	assert.Equal(t, 1, params.Pagination.Page)
	assert.Equal(t, 10, params.Pagination.Limit)
	assert.Equal(t, 0, params.Pagination.Offset)
}

// Test NewAdvancedQueryParams
func TestNewAdvancedQueryParams(t *testing.T) {
	params := NewAdvancedQueryParams()
	assert.NotNil(t, params)
	assert.NotNil(t, params.SearchGroups)
	assert.NotNil(t, params.Filters)
	assert.NotNil(t, params.Sort)
	assert.Equal(t, 0, len(params.SearchGroups))
	assert.Equal(t, 0, len(params.Filters))
	assert.Equal(t, 0, len(params.Sort))
	assert.Equal(t, 1, params.Pagination.Page)
	assert.Equal(t, 10, params.Pagination.Limit)
	assert.Equal(t, 0, params.Pagination.Offset)
}

// Test QueryParams methods
func TestQueryParams_AddSearch(t *testing.T) {
	params := NewQueryParams()
	params.AddSearch("title", OpEqual, "test")
	params.AddSearch("content", OpContains, "example")

	assert.Equal(t, 2, len(params.Search))
	assert.Equal(t, "title", params.Search[0].Field)
	assert.Equal(t, OpEqual, params.Search[0].Operator)
	assert.Equal(t, "test", params.Search[0].Value)
	assert.Equal(t, "content", params.Search[1].Field)
	assert.Equal(t, OpContains, params.Search[1].Operator)
	assert.Equal(t, "example", params.Search[1].Value)
}

func TestQueryParams_AddFilter(t *testing.T) {
	params := NewQueryParams()
	params.AddFilter("status", OpEqual, "active")
	params.AddFilter("created_at", OpGreaterThan, "2024-01-01")

	assert.Equal(t, 2, len(params.Filters))
	assert.Equal(t, "status", params.Filters[0].Field)
	assert.Equal(t, OpEqual, params.Filters[0].Operator)
	assert.Equal(t, "active", params.Filters[0].Value)
	assert.Equal(t, "created_at", params.Filters[1].Field)
	assert.Equal(t, OpGreaterThan, params.Filters[1].Operator)
	assert.Equal(t, "2024-01-01", params.Filters[1].Value)
}

func TestQueryParams_AddSort(t *testing.T) {
	params := NewQueryParams()
	params.AddSort("title", "ASC")
	params.AddSort("created_at", "DESC")

	assert.Equal(t, 2, len(params.Sort))
	assert.Equal(t, "title", params.Sort[0].Field)
	assert.Equal(t, "asc", params.Sort[0].Order)
	assert.Equal(t, "created_at", params.Sort[1].Field)
	assert.Equal(t, "desc", params.Sort[1].Order)
}

func TestQueryParams_SetPagination(t *testing.T) {
	params := NewQueryParams()
	params.SetPagination(3, 25)

	assert.Equal(t, 3, params.Pagination.Page)
	assert.Equal(t, 25, params.Pagination.Limit)
	assert.Equal(t, 50, params.Pagination.Offset)
}

// Test AdvancedQueryParams methods
func TestAdvancedQueryParams_AddSearchGroup(t *testing.T) {
	params := NewAdvancedQueryParams()
	conditions := []SearchCriteria{
		{Field: "title", Operator: OpEqual, Value: "test"},
		{Field: "content", Operator: OpContains, Value: "example"},
	}
	params.AddSearchGroup(LogicAnd, conditions)

	assert.Equal(t, 1, len(params.SearchGroups))
	assert.Equal(t, LogicAnd, params.SearchGroups[0].Operator)
	assert.Equal(t, 2, len(params.SearchGroups[0].Conditions))
	assert.Equal(t, 0, len(params.SearchGroups[0].Groups))
}

func TestAdvancedQueryParams_AddNestedSearchGroup(t *testing.T) {
	params := NewAdvancedQueryParams()

	// Add parent group
	parentConditions := []SearchCriteria{
		{Field: "title", Operator: OpEqual, Value: "test"},
	}
	params.AddSearchGroup(LogicAnd, parentConditions)

	// Add nested group
	nestedConditions := []SearchCriteria{
		{Field: "content", Operator: OpContains, Value: "example"},
	}
	params.AddNestedSearchGroup(0, LogicOr, nestedConditions)

	assert.Equal(t, 1, len(params.SearchGroups))
	assert.Equal(t, 1, len(params.SearchGroups[0].Groups))
	assert.Equal(t, LogicOr, params.SearchGroups[0].Groups[0].Operator)
	assert.Equal(t, 1, len(params.SearchGroups[0].Groups[0].Conditions))
}

func TestAdvancedQueryParams_AddNestedSearchGroup_InvalidIndex(t *testing.T) {
	params := NewAdvancedQueryParams()

	// Try to add nested group to non-existent parent
	nestedConditions := []SearchCriteria{
		{Field: "content", Operator: OpContains, Value: "example"},
	}
	params.AddNestedSearchGroup(10, LogicOr, nestedConditions)

	assert.Equal(t, 0, len(params.SearchGroups))
}

// Test NewSQLBuilder
func TestNewSQLBuilder(t *testing.T) {
	builder := NewSQLBuilder()
	assert.NotNil(t, builder)
	assert.NotNil(t, builder.whereConditions)
	assert.NotNil(t, builder.params)
	assert.Equal(t, 0, len(builder.whereConditions))
	assert.Equal(t, 0, len(builder.params))
	assert.Equal(t, 0, builder.paramIndex)
}

// Test SQLBuilder BuildSearchConditions
func TestSQLBuilder_BuildSearchConditions(t *testing.T) {
	tests := []struct {
		name           string
		search         []SearchCriteria
		expectedSQL    string
		expectedParams []any
	}{
		{
			name:           "empty search",
			search:         []SearchCriteria{},
			expectedSQL:    "",
			expectedParams: []any{},
		},
		{
			name: "single search condition",
			search: []SearchCriteria{
				{Field: "title", Operator: OpEqual, Value: "test"},
			},
			expectedSQL:    "(title = ?)",
			expectedParams: []any{"test"},
		},
		{
			name: "multiple search conditions",
			search: []SearchCriteria{
				{Field: "title", Operator: OpEqual, Value: "test"},
				{Field: "content", Operator: OpContains, Value: "example"},
			},
			expectedSQL:    "(title = ? OR content LIKE ?)",
			expectedParams: []any{"test", "%example%"},
		},
		{
			name: "case insensitive search",
			search: []SearchCriteria{
				{Field: "title", Operator: OpIContains, Value: "Test"},
			},
			expectedSQL:    "(LOWER(title) LIKE LOWER(?))",
			expectedParams: []any{"%Test%"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSQLBuilder()
			result := builder.BuildSearchConditions(tt.search)
			assert.Equal(t, tt.expectedSQL, result)
			assert.Equal(t, tt.expectedParams, builder.GetParams())
		})
	}
}

// Test SQLBuilder BuildFilterConditions
func TestSQLBuilder_BuildFilterConditions(t *testing.T) {
	tests := []struct {
		name           string
		filters        []FilterCriteria
		expectedSQL    string
		expectedParams []any
	}{
		{
			name:           "empty filters",
			filters:        []FilterCriteria{},
			expectedSQL:    "",
			expectedParams: []any{},
		},
		{
			name: "single filter condition",
			filters: []FilterCriteria{
				{Field: "status", Operator: OpEqual, Value: "active"},
			},
			expectedSQL:    "status = ?",
			expectedParams: []any{"active"},
		},
		{
			name: "multiple filter conditions",
			filters: []FilterCriteria{
				{Field: "status", Operator: OpEqual, Value: "active"},
				{Field: "created_at", Operator: OpGreaterThan, Value: "2024-01-01"},
			},
			expectedSQL:    "status = ? AND created_at > ?",
			expectedParams: []any{"active", "2024-01-01"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSQLBuilder()
			result := builder.BuildFilterConditions(tt.filters)
			assert.Equal(t, tt.expectedSQL, result)
			assert.Equal(t, tt.expectedParams, builder.GetParams())
		})
	}
}

// Test SQLBuilder BuildAdvancedSearchConditions
func TestSQLBuilder_BuildAdvancedSearchConditions(t *testing.T) {
	tests := []struct {
		name           string
		groups         []LogicalGroup
		expectedSQL    string
		expectedParams []any
	}{
		{
			name:           "empty groups",
			groups:         []LogicalGroup{},
			expectedSQL:    "",
			expectedParams: []any{},
		},
		{
			name: "single group with AND",
			groups: []LogicalGroup{
				{
					Operator: LogicAnd,
					Conditions: []SearchCriteria{
						{Field: "title", Operator: OpEqual, Value: "test"},
						{Field: "content", Operator: OpContains, Value: "example"},
					},
				},
			},
			expectedSQL:    "(title = ? AND content LIKE ?)",
			expectedParams: []any{"test", "%example%"},
		},
		{
			name: "single group with OR",
			groups: []LogicalGroup{
				{
					Operator: LogicOr,
					Conditions: []SearchCriteria{
						{Field: "title", Operator: OpEqual, Value: "test"},
						{Field: "content", Operator: OpContains, Value: "example"},
					},
				},
			},
			expectedSQL:    "(title = ? OR content LIKE ?)",
			expectedParams: []any{"test", "%example%"},
		},
		{
			name: "multiple groups",
			groups: []LogicalGroup{
				{
					Operator: LogicAnd,
					Conditions: []SearchCriteria{
						{Field: "title", Operator: OpEqual, Value: "test"},
					},
				},
				{
					Operator: LogicOr,
					Conditions: []SearchCriteria{
						{Field: "content", Operator: OpContains, Value: "example"},
					},
				},
			},
			expectedSQL:    "(title = ?) AND (content LIKE ?)",
			expectedParams: []any{"test", "%example%"},
		},
		{
			name: "nested groups",
			groups: []LogicalGroup{
				{
					Operator: LogicAnd,
					Conditions: []SearchCriteria{
						{Field: "title", Operator: OpEqual, Value: "test"},
					},
					Groups: []LogicalGroup{
						{
							Operator: LogicOr,
							Conditions: []SearchCriteria{
								{Field: "content", Operator: OpContains, Value: "example"},
								{Field: "description", Operator: OpContains, Value: "desc"},
							},
						},
					},
				},
			},
			expectedSQL:    "(title = ? AND ((content LIKE ? OR description LIKE ?)))",
			expectedParams: []any{"test", "%example%", "%desc%"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSQLBuilder()
			result := builder.BuildAdvancedSearchConditions(tt.groups)
			assert.Equal(t, tt.expectedSQL, result)
			assert.Equal(t, tt.expectedParams, builder.GetParams())
		})
	}
}

// Test SQLBuilder BuildOrderBy
func TestSQLBuilder_BuildOrderBy(t *testing.T) {
	tests := []struct {
		name                     string
		sort                     []SortCriteria
		expectedSQL              string
		expectedSQLWithoutPrefix string
	}{
		{
			name:                     "empty sort",
			sort:                     []SortCriteria{},
			expectedSQL:              "",
			expectedSQLWithoutPrefix: "",
		},
		{
			name: "single sort ASC",
			sort: []SortCriteria{
				{Field: "title", Order: "asc"},
			},
			expectedSQL:              "ORDER BY title ASC",
			expectedSQLWithoutPrefix: "title ASC",
		},
		{
			name: "single sort DESC",
			sort: []SortCriteria{
				{Field: "title", Order: "desc"},
			},
			expectedSQL:              "ORDER BY title DESC",
			expectedSQLWithoutPrefix: "title DESC",
		},
		{
			name: "multiple sort criteria",
			sort: []SortCriteria{
				{Field: "title", Order: "asc"},
				{Field: "created_at", Order: "desc"},
			},
			expectedSQL:              "ORDER BY title ASC, created_at DESC",
			expectedSQLWithoutPrefix: "title ASC, created_at DESC",
		},
		{
			name: "case insensitive order",
			sort: []SortCriteria{
				{Field: "title", Order: "ASC"},
				{Field: "created_at", Order: "DESC"},
			},
			expectedSQL:              "ORDER BY title ASC, created_at DESC",
			expectedSQLWithoutPrefix: "title ASC, created_at DESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSQLBuilder()

			// Test with prefix (default behavior)
			result := builder.BuildOrderBy(tt.sort)
			assert.Equal(t, tt.expectedSQL, result)

			// Test with prefix explicitly
			result = builder.BuildOrderBy(tt.sort, true)
			assert.Equal(t, tt.expectedSQL, result)

			// Test without prefix
			result = builder.BuildOrderBy(tt.sort, false)
			assert.Equal(t, tt.expectedSQLWithoutPrefix, result)
		})
	}
}

// Test SQLBuilder AddWhereCondition and GetWhereClause
func TestSQLBuilder_WhereClause(t *testing.T) {
	builder := NewSQLBuilder()

	// Test empty where clause
	assert.Equal(t, "", builder.GetWhereClause())
	assert.Equal(t, "", builder.GetWhereClause(true))
	assert.Equal(t, "", builder.GetWhereClause(false))

	// Add conditions
	builder.AddWhereCondition("title = ?")
	builder.AddWhereCondition("status = ?")
	builder.AddWhereCondition("") // Empty condition should be ignored

	// Test with prefix (default behavior)
	expected := "WHERE title = ? AND status = ?"
	assert.Equal(t, expected, builder.GetWhereClause())
	assert.Equal(t, expected, builder.GetWhereClause(true))

	// Test without prefix
	expectedWithoutPrefix := "title = ? AND status = ?"
	assert.Equal(t, expectedWithoutPrefix, builder.GetWhereClause(false))
}

// Test all buildCondition operators
func TestSQLBuilder_buildCondition(t *testing.T) {
	tests := []struct {
		name           string
		field          string
		operator       string
		value          any
		expectedSQL    string
		expectedParams []any
	}{
		{
			name:           "OpEqual",
			field:          "title",
			operator:       OpEqual,
			value:          "test",
			expectedSQL:    "title = ?",
			expectedParams: []any{"test"},
		},
		{
			name:           "OpNotEqual",
			field:          "title",
			operator:       OpNotEqual,
			value:          "test",
			expectedSQL:    "title != ?",
			expectedParams: []any{"test"},
		},
		{
			name:           "OpGreaterThan",
			field:          "count",
			operator:       OpGreaterThan,
			value:          10,
			expectedSQL:    "count > ?",
			expectedParams: []any{10},
		},
		{
			name:           "OpGreaterThanEq",
			field:          "count",
			operator:       OpGreaterThanEq,
			value:          10,
			expectedSQL:    "count >= ?",
			expectedParams: []any{10},
		},
		{
			name:           "OpLessThan",
			field:          "count",
			operator:       OpLessThan,
			value:          10,
			expectedSQL:    "count < ?",
			expectedParams: []any{10},
		},
		{
			name:           "OpLessThanEq",
			field:          "count",
			operator:       OpLessThanEq,
			value:          10,
			expectedSQL:    "count <= ?",
			expectedParams: []any{10},
		},
		{
			name:           "OpContains",
			field:          "title",
			operator:       OpContains,
			value:          "test",
			expectedSQL:    "title LIKE ?",
			expectedParams: []any{"%test%"},
		},
		{
			name:           "OpIContains",
			field:          "title",
			operator:       OpIContains,
			value:          "test",
			expectedSQL:    "LOWER(title) LIKE LOWER(?)",
			expectedParams: []any{"%test%"},
		},
		{
			name:           "OpStartsWith",
			field:          "title",
			operator:       OpStartsWith,
			value:          "test",
			expectedSQL:    "title LIKE ?",
			expectedParams: []any{"test%"},
		},
		{
			name:           "OpIStartsWith",
			field:          "title",
			operator:       OpIStartsWith,
			value:          "test",
			expectedSQL:    "LOWER(title) LIKE LOWER(?)",
			expectedParams: []any{"test%"},
		},
		{
			name:           "OpEndsWith",
			field:          "title",
			operator:       OpEndsWith,
			value:          "test",
			expectedSQL:    "title LIKE ?",
			expectedParams: []any{"%test"},
		},
		{
			name:           "OpIEndsWith",
			field:          "title",
			operator:       OpIEndsWith,
			value:          "test",
			expectedSQL:    "LOWER(title) LIKE LOWER(?)",
			expectedParams: []any{"%test"},
		},
		{
			name:           "OpLike",
			field:          "title",
			operator:       OpLike,
			value:          "%test%",
			expectedSQL:    "title LIKE ?",
			expectedParams: []any{"%test%"},
		},
		{
			name:           "OpILike",
			field:          "title",
			operator:       OpILike,
			value:          "%test%",
			expectedSQL:    "LOWER(title) LIKE LOWER(?)",
			expectedParams: []any{"%test%"},
		},
		{
			name:           "OpFullText",
			field:          "title",
			operator:       OpFullText,
			value:          "test search",
			expectedSQL:    "MATCH(title) AGAINST(? IN NATURAL LANGUAGE MODE)",
			expectedParams: []any{"test search"},
		},
		{
			name:           "OpRegex",
			field:          "title",
			operator:       OpRegex,
			value:          "^test.*",
			expectedSQL:    "title REGEXP ?",
			expectedParams: []any{"^test.*"},
		},
		{
			name:           "OpIRegex",
			field:          "title",
			operator:       OpIRegex,
			value:          "^test.*",
			expectedSQL:    "title REGEXP ?",
			expectedParams: []any{"^test.*"},
		},
		{
			name:           "OpIsNull",
			field:          "title",
			operator:       OpIsNull,
			value:          nil,
			expectedSQL:    "title IS NULL",
			expectedParams: []any{},
		},
		{
			name:           "OpIsNotNull",
			field:          "title",
			operator:       OpIsNotNull,
			value:          nil,
			expectedSQL:    "title IS NOT NULL",
			expectedParams: []any{},
		},
		{
			name:           "OpIn",
			field:          "status",
			operator:       OpIn,
			value:          []any{"active", "inactive"},
			expectedSQL:    "status IN (?, ?)",
			expectedParams: []any{"active", "inactive"},
		},
		{
			name:           "OpNotIn",
			field:          "status",
			operator:       OpNotIn,
			value:          []any{"active", "inactive"},
			expectedSQL:    "status NOT IN (?, ?)",
			expectedParams: []any{"active", "inactive"},
		},
		{
			name:           "OpIn empty array",
			field:          "status",
			operator:       OpIn,
			value:          []any{},
			expectedSQL:    "",
			expectedParams: []any{},
		},
		{
			name:           "OpNotIn empty array",
			field:          "status",
			operator:       OpNotIn,
			value:          []any{},
			expectedSQL:    "",
			expectedParams: []any{},
		},
		{
			name:           "OpIn invalid type",
			field:          "status",
			operator:       OpIn,
			value:          "not_an_array",
			expectedSQL:    "",
			expectedParams: []any{},
		},
		{
			name:           "Unknown operator",
			field:          "title",
			operator:       "unknown_op",
			value:          "test",
			expectedSQL:    "",
			expectedParams: []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSQLBuilder()
			result := builder.buildCondition(tt.field, tt.operator, tt.value)
			assert.Equal(t, tt.expectedSQL, result)
			assert.Equal(t, tt.expectedParams, builder.GetParams())
		})
	}
}

// Test CalculatePaginationMeta
func TestCalculatePaginationMeta(t *testing.T) {
	tests := []struct {
		name            string
		totalRecords    int
		page            int
		limit           int
		expectedItems   int
		expectedTotal   int
		expectedCurrent int
		expectedLimit   int
	}{
		{
			name:            "exact division",
			totalRecords:    20,
			page:            2,
			limit:           10,
			expectedItems:   20,
			expectedTotal:   2,
			expectedCurrent: 2,
			expectedLimit:   10,
		},
		{
			name:            "remainder division",
			totalRecords:    23,
			page:            1,
			limit:           10,
			expectedItems:   23,
			expectedTotal:   3,
			expectedCurrent: 1,
			expectedLimit:   10,
		},
		{
			name:            "zero records",
			totalRecords:    0,
			page:            1,
			limit:           10,
			expectedItems:   0,
			expectedTotal:   1,
			expectedCurrent: 1,
			expectedLimit:   10,
		},
		{
			name:            "single record",
			totalRecords:    1,
			page:            1,
			limit:           10,
			expectedItems:   1,
			expectedTotal:   1,
			expectedCurrent: 1,
			expectedLimit:   10,
		},
		{
			name:            "large numbers",
			totalRecords:    1000,
			page:            10,
			limit:           25,
			expectedItems:   1000,
			expectedTotal:   40,
			expectedCurrent: 10,
			expectedLimit:   25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := CalculatePaginationMeta(tt.totalRecords, tt.page, tt.limit)
			assert.Equal(t, tt.expectedItems, meta.TotalItems)
			assert.Equal(t, tt.expectedTotal, meta.TotalPage)
			assert.Equal(t, tt.expectedCurrent, meta.CurrentPage)
			assert.Equal(t, tt.expectedLimit, meta.PageLimit)
		})
	}
}

// Test helper functions
func TestCreateSearchGroup(t *testing.T) {
	conditions := []SearchCriteria{
		{Field: "title", Operator: OpEqual, Value: "test"},
		{Field: "content", Operator: OpContains, Value: "example"},
	}

	group := CreateSearchGroup(LogicAnd, conditions...)
	assert.Equal(t, LogicAnd, group.Operator)
	assert.Equal(t, 2, len(group.Conditions))
	assert.Equal(t, 0, len(group.Groups))
	assert.Equal(t, "title", group.Conditions[0].Field)
	assert.Equal(t, OpEqual, group.Conditions[0].Operator)
	assert.Equal(t, "test", group.Conditions[0].Value)
}

func TestCreateSearchCondition(t *testing.T) {
	condition := CreateSearchCondition("title", OpEqual, "test")
	assert.Equal(t, "title", condition.Field)
	assert.Equal(t, OpEqual, condition.Operator)
	assert.Equal(t, "test", condition.Value)
}

func TestCreateFilterCondition(t *testing.T) {
	condition := CreateFilterCondition("status", OpEqual, "active")
	assert.Equal(t, "status", condition.Field)
	assert.Equal(t, OpEqual, condition.Operator)
	assert.Equal(t, "active", condition.Value)
}

// Test QueryParams convenience methods
func TestQueryParams_ConvenienceMethods(t *testing.T) {
	params := NewQueryParams()

	// Test empty params
	assert.False(t, params.HasSearch())
	assert.False(t, params.HasFilters())
	assert.False(t, params.HasSort())

	// Add some criteria
	params.AddSearch("title", OpEqual, "test")
	params.AddFilter("status", OpEqual, "active")
	params.AddSort("created_at", "desc")

	assert.True(t, params.HasSearch())
	assert.True(t, params.HasFilters())
	assert.True(t, params.HasSort())

	// Test ApplySearch
	builder := NewSQLBuilder()
	params.ApplySearch(builder)
	assert.Equal(t, 1, len(builder.GetParams()))

	// Test ApplyFilters
	builder = NewSQLBuilder()
	params.ApplyFilters(builder)
	assert.Equal(t, 1, len(builder.GetParams()))

	// Test ApplySort
	builder = NewSQLBuilder()
	orderBy := params.ApplySort(builder, false)
	assert.Equal(t, "created_at DESC", orderBy)

	orderBy = params.ApplySort(builder, true)
	assert.Equal(t, "ORDER BY created_at DESC", orderBy)
}

func TestQueryParams_ConvenienceMethods_Empty(t *testing.T) {
	params := NewQueryParams()
	builder := NewSQLBuilder()

	// Test empty params don't add anything
	params.ApplySearch(builder)
	params.ApplyFilters(builder)
	orderBy := params.ApplySort(builder, false)

	assert.Equal(t, 0, len(builder.GetParams()))
	assert.Equal(t, "", builder.GetWhereClause(false))
	assert.Equal(t, "", orderBy)
}

// Test AdvancedQueryParams convenience methods
func TestAdvancedQueryParams_ConvenienceMethods(t *testing.T) {
	params := NewAdvancedQueryParams()

	// Test empty params
	assert.False(t, params.HasSearchGroups())
	assert.False(t, params.HasFilters())
	assert.False(t, params.HasSort())

	// Add some criteria
	conditions := []SearchCriteria{
		{Field: "title", Operator: OpEqual, Value: "test"},
	}
	params.AddSearchGroup(LogicAnd, conditions)
	params.Filters = append(params.Filters, FilterCriteria{Field: "status", Operator: OpEqual, Value: "active"})
	params.Sort = append(params.Sort, SortCriteria{Field: "created_at", Order: "desc"})

	assert.True(t, params.HasSearchGroups())
	assert.True(t, params.HasFilters())
	assert.True(t, params.HasSort())

	// Test ApplyAdvancedSearch
	builder := NewSQLBuilder()
	params.ApplyAdvancedSearch(builder)
	assert.Equal(t, 1, len(builder.GetParams()))

	// Test ApplyFilters
	builder = NewSQLBuilder()
	params.ApplyFilters(builder)
	assert.Equal(t, 1, len(builder.GetParams()))

	// Test ApplySort
	builder = NewSQLBuilder()
	orderBy := params.ApplySort(builder, false)
	assert.Equal(t, "created_at DESC", orderBy)

	// Test SetPagination
	params.SetPagination(2, 20)
	assert.Equal(t, 2, params.Pagination.Page)
	assert.Equal(t, 20, params.Pagination.Limit)
	assert.Equal(t, 20, params.Pagination.Offset)
}

func TestAdvancedQueryParams_ConvenienceMethods_Empty(t *testing.T) {
	params := NewAdvancedQueryParams()
	builder := NewSQLBuilder()

	// Test empty params don't add anything
	params.ApplyAdvancedSearch(builder)
	params.ApplyFilters(builder)
	orderBy := params.ApplySort(builder, false)

	assert.Equal(t, 0, len(builder.GetParams()))
	assert.Equal(t, "", builder.GetWhereClause(false))
	assert.Equal(t, "", orderBy)
}

// Test complex integration scenarios
func TestSQLBuilder_IntegrationScenarios(t *testing.T) {
	t.Run("complete query building", func(t *testing.T) {
		builder := NewSQLBuilder()

		// Add search conditions
		search := []SearchCriteria{
			{Field: "title", Operator: OpIContains, Value: "test"},
			{Field: "content", Operator: OpIContains, Value: "example"},
		}
		searchCondition := builder.BuildSearchConditions(search)
		builder.AddWhereCondition(searchCondition)

		// Add filter conditions
		filters := []FilterCriteria{
			{Field: "status", Operator: OpEqual, Value: "active"},
			{Field: "created_at", Operator: OpGreaterThan, Value: "2024-01-01"},
		}
		filterCondition := builder.BuildFilterConditions(filters)
		builder.AddWhereCondition(filterCondition)

		// Build sort
		sort := []SortCriteria{
			{Field: "title", Order: "asc"},
			{Field: "created_at", Order: "desc"},
		}
		orderBy := builder.BuildOrderBy(sort)

		expectedWhere := "WHERE (LOWER(title) LIKE LOWER(?) OR LOWER(content) LIKE LOWER(?)) AND status = ? AND created_at > ?"
		expectedOrderBy := "ORDER BY title ASC, created_at DESC"
		expectedParams := []any{"%test%", "%example%", "active", "2024-01-01"}

		assert.Equal(t, expectedWhere, builder.GetWhereClause())
		assert.Equal(t, expectedOrderBy, orderBy)
		assert.Equal(t, expectedParams, builder.GetParams())

		// Test without prefixes
		assert.Equal(t, "(LOWER(title) LIKE LOWER(?) OR LOWER(content) LIKE LOWER(?)) AND status = ? AND created_at > ?", builder.GetWhereClause(false))
		assert.Equal(t, "title ASC, created_at DESC", builder.BuildOrderBy(sort, false))
	})

	t.Run("complete query building with convenience methods", func(t *testing.T) {
		params := NewQueryParams()

		// Add search, filter, and sort criteria
		params.AddSearch("title", OpIContains, "test")
		params.AddSearch("content", OpIContains, "example")
		params.AddFilter("status", OpEqual, "active")
		params.AddFilter("created_at", OpGreaterThan, "2024-01-01")
		params.AddSort("title", "asc")
		params.AddSort("created_at", "desc")

		// Apply using convenience methods
		builder := NewSQLBuilder()
		params.ApplySearch(builder)
		params.ApplyFilters(builder)
		orderBy := params.ApplySort(builder, false)

		expectedWhere := "(LOWER(title) LIKE LOWER(?) OR LOWER(content) LIKE LOWER(?)) AND status = ? AND created_at > ?"
		expectedOrderBy := "title ASC, created_at DESC"
		expectedParams := []any{"%test%", "%example%", "active", "2024-01-01"}

		assert.Equal(t, expectedWhere, builder.GetWhereClause(false))
		assert.Equal(t, expectedOrderBy, orderBy)
		assert.Equal(t, expectedParams, builder.GetParams())
	})

	t.Run("advanced query with nested groups", func(t *testing.T) {
		builder := NewSQLBuilder()

		groups := []LogicalGroup{
			{
				Operator: LogicAnd,
				Conditions: []SearchCriteria{
					{Field: "title", Operator: OpEqual, Value: "test"},
				},
				Groups: []LogicalGroup{
					{
						Operator: LogicOr,
						Conditions: []SearchCriteria{
							{Field: "content", Operator: OpContains, Value: "example"},
							{Field: "description", Operator: OpContains, Value: "desc"},
						},
						Groups: []LogicalGroup{
							{
								Operator: LogicAnd,
								Conditions: []SearchCriteria{
									{Field: "status", Operator: OpEqual, Value: "active"},
									{Field: "published", Operator: OpEqual, Value: true},
								},
							},
						},
					},
				},
			},
		}

		searchCondition := builder.BuildAdvancedSearchConditions(groups)
		expected := "(title = ? AND ((content LIKE ? OR description LIKE ? OR ((status = ? AND published = ?)))))"
		expectedParams := []any{"test", "%example%", "%desc%", "active", true}

		assert.Equal(t, expected, searchCondition)
		assert.Equal(t, expectedParams, builder.GetParams())
	})
}

// Test edge cases and error conditions
func TestSQLBuilder_EdgeCases(t *testing.T) {
	t.Run("nil values handled gracefully", func(t *testing.T) {
		builder := NewSQLBuilder()

		// Test with nil search criteria
		result := builder.BuildSearchConditions(nil)
		assert.Equal(t, "", result)

		// Test with nil filter criteria
		result = builder.BuildFilterConditions(nil)
		assert.Equal(t, "", result)

		// Test with nil sort criteria
		result = builder.BuildOrderBy(nil)
		assert.Equal(t, "", result)

		// Test with nil groups
		result = builder.BuildAdvancedSearchConditions(nil)
		assert.Equal(t, "", result)
	})

	t.Run("empty conditions filtered out", func(t *testing.T) {
		builder := NewSQLBuilder()

		// Mix of valid and invalid conditions
		search := []SearchCriteria{
			{Field: "title", Operator: OpEqual, Value: "test"},
			{Field: "content", Operator: "invalid_op", Value: "example"}, // Invalid operator
		}

		result := builder.BuildSearchConditions(search)
		assert.Equal(t, "(title = ?)", result)
		assert.Equal(t, []any{"test"}, builder.GetParams())
	})

	t.Run("empty groups handled", func(t *testing.T) {
		builder := NewSQLBuilder()

		groups := []LogicalGroup{
			{
				Operator:   LogicAnd,
				Conditions: []SearchCriteria{}, // Empty conditions
				Groups: []LogicalGroup{
					{
						Operator: LogicOr,
						Conditions: []SearchCriteria{
							{Field: "title", Operator: OpEqual, Value: "test"},
						},
					},
				},
			},
		}

		result := builder.BuildAdvancedSearchConditions(groups)
		assert.Equal(t, "(((title = ?)))", result)
		assert.Equal(t, []any{"test"}, builder.GetParams())
	})
}

// Benchmark tests
func BenchmarkSQLBuilder_BuildSearchConditions(b *testing.B) {
	search := []SearchCriteria{
		{Field: "title", Operator: OpIContains, Value: "test"},
		{Field: "content", Operator: OpIContains, Value: "example"},
		{Field: "description", Operator: OpIContains, Value: "desc"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder := NewSQLBuilder()
		builder.BuildSearchConditions(search)
	}
}

func BenchmarkSQLBuilder_BuildAdvancedSearchConditions(b *testing.B) {
	groups := []LogicalGroup{
		{
			Operator: LogicAnd,
			Conditions: []SearchCriteria{
				{Field: "title", Operator: OpEqual, Value: "test"},
			},
			Groups: []LogicalGroup{
				{
					Operator: LogicOr,
					Conditions: []SearchCriteria{
						{Field: "content", Operator: OpContains, Value: "example"},
						{Field: "description", Operator: OpContains, Value: "desc"},
					},
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder := NewSQLBuilder()
		builder.BuildAdvancedSearchConditions(groups)
	}
}
