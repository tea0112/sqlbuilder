package sqlbuilder

import (
	"fmt"
	"strings"
)

// SearchCriteria represents a single search criterion
type SearchCriteria struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

// FilterCriteria represents a single filter criterion
type FilterCriteria struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

// SortCriteria represents sorting parameters
type SortCriteria struct {
	Field string `json:"field"`
	Order string `json:"order"` // ASC or DESC
}

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page   int `json:"page"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// QueryParams represents comprehensive query parameters
type QueryParams struct {
	Search     []SearchCriteria `json:"search"`
	Filters    []FilterCriteria `json:"filters"`
	Sort       []SortCriteria   `json:"sort"`
	Pagination PaginationParams `json:"pagination"`
}

// LogicalGroup represents a group of conditions with AND/OR logic
type LogicalGroup struct {
	Operator   string           `json:"operator"` // AND or OR
	Conditions []SearchCriteria `json:"conditions"`
	Groups     []LogicalGroup   `json:"groups"` // Nested groups for complex logic
}

// AdvancedQueryParams supports complex AND/OR logic
type AdvancedQueryParams struct {
	SearchGroups []LogicalGroup   `json:"search_groups"`
	Filters      []FilterCriteria `json:"filters"`
	Sort         []SortCriteria   `json:"sort"`
	Pagination   PaginationParams `json:"pagination"`
}

// SearchOperators defines available search operators
const (
	OpEqual         = "eq"
	OpNotEqual      = "ne"
	OpGreaterThan   = "gt"
	OpGreaterThanEq = "gte"
	OpLessThan      = "lt"
	OpLessThanEq    = "lte"
	OpContains      = "contains"
	OpIContains     = "icontains" // Case-insensitive contains
	OpStartsWith    = "starts_with"
	OpIStartsWith   = "istarts_with" // Case-insensitive starts with
	OpEndsWith      = "ends_with"
	OpIEndsWith     = "iends_with" // Case-insensitive ends with
	OpLike          = "like"       // Case-sensitive pattern matching
	OpILike         = "ilike"      // Case-insensitive pattern matching
	OpIn            = "in"
	OpNotIn         = "not_in"
	OpIsNull        = "is_null"
	OpIsNotNull     = "is_not_null"
	OpFullText      = "full_text"
	OpRegex         = "regex"  // Regular expression matching
	OpIRegex        = "iregex" // Case-insensitive regex
)

// Logical operators
const (
	LogicAnd = "AND"
	LogicOr  = "OR"
)

// Sort orders
const (
	SortAsc  = "asc"
	SortDesc = "desc"
)

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	TotalItems  int `json:"total_items"`
	TotalPage   int `json:"total_page"`
	CurrentPage int `json:"current_page"`
	PageLimit   int `json:"page_limit"`
}

// NewPaginationParams creates pagination parameters with defaults
func NewPaginationParams(page, limit int) PaginationParams {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	return PaginationParams{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}

// NewQueryParams creates a new QueryParams instance
func NewQueryParams() *QueryParams {
	return &QueryParams{
		Search:     make([]SearchCriteria, 0),
		Filters:    make([]FilterCriteria, 0),
		Sort:       make([]SortCriteria, 0),
		Pagination: NewPaginationParams(1, 10),
	}
}

// NewAdvancedQueryParams creates a new AdvancedQueryParams instance
func NewAdvancedQueryParams() *AdvancedQueryParams {
	return &AdvancedQueryParams{
		SearchGroups: make([]LogicalGroup, 0),
		Filters:      make([]FilterCriteria, 0),
		Sort:         make([]SortCriteria, 0),
		Pagination:   NewPaginationParams(1, 10),
	}
}

// AddSearch adds a search criterion
func (q *QueryParams) AddSearch(field, operator string, value any) {
	q.Search = append(q.Search, SearchCriteria{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
}

// AddFilter adds a filter criterion
func (q *QueryParams) AddFilter(field, operator string, value any) {
	q.Filters = append(q.Filters, FilterCriteria{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
}

// AddSort adds a sort criterion
func (q *QueryParams) AddSort(field, order string) {
	q.Sort = append(q.Sort, SortCriteria{
		Field: field,
		Order: strings.ToLower(order),
	})
}

// SetPagination sets pagination parameters
func (q *QueryParams) SetPagination(page, limit int) {
	q.Pagination = NewPaginationParams(page, limit)
}

// HasSearch returns true if there are search criteria
func (q *QueryParams) HasSearch() bool {
	return len(q.Search) > 0
}

// HasFilters returns true if there are filter criteria
func (q *QueryParams) HasFilters() bool {
	return len(q.Filters) > 0
}

// HasSort returns true if there are sort criteria
func (q *QueryParams) HasSort() bool {
	return len(q.Sort) > 0
}

// ApplySearch applies search conditions to the SQL builder
func (q *QueryParams) ApplySearch(builder *SQLBuilder) {
	if q.HasSearch() {
		searchConditions := builder.BuildSearchConditions(q.Search)
		builder.AddWhereCondition(searchConditions)
	}
}

// ApplyFilters applies filter conditions to the SQL builder
func (q *QueryParams) ApplyFilters(builder *SQLBuilder) {
	if q.HasFilters() {
		filterConditions := builder.BuildFilterConditions(q.Filters)
		builder.AddWhereCondition(filterConditions)
	}
}

// ApplySort applies sort conditions to the SQL builder and returns the ORDER BY clause
func (q *QueryParams) ApplySort(builder *SQLBuilder, includePrefix ...bool) string {
	if q.HasSort() {
		return builder.BuildOrderBy(q.Sort, includePrefix...)
	}
	return ""
}

// AddSearchGroup adds a logical group of search conditions
func (q *AdvancedQueryParams) AddSearchGroup(operator string, conditions []SearchCriteria) {
	q.SearchGroups = append(q.SearchGroups, LogicalGroup{
		Operator:   operator,
		Conditions: conditions,
		Groups:     make([]LogicalGroup, 0),
	})
}

// AddNestedSearchGroup adds a nested logical group
func (q *AdvancedQueryParams) AddNestedSearchGroup(parentIndex int, operator string, conditions []SearchCriteria) {
	if parentIndex < len(q.SearchGroups) {
		q.SearchGroups[parentIndex].Groups = append(q.SearchGroups[parentIndex].Groups, LogicalGroup{
			Operator:   operator,
			Conditions: conditions,
			Groups:     make([]LogicalGroup, 0),
		})
	}
}

// HasSearchGroups returns true if there are search groups
func (q *AdvancedQueryParams) HasSearchGroups() bool {
	return len(q.SearchGroups) > 0
}

// HasFilters returns true if there are filter criteria
func (q *AdvancedQueryParams) HasFilters() bool {
	return len(q.Filters) > 0
}

// HasSort returns true if there are sort criteria
func (q *AdvancedQueryParams) HasSort() bool {
	return len(q.Sort) > 0
}

// ApplyAdvancedSearch applies advanced search conditions to the SQL builder
func (q *AdvancedQueryParams) ApplyAdvancedSearch(builder *SQLBuilder) {
	if q.HasSearchGroups() {
		searchConditions := builder.BuildAdvancedSearchConditions(q.SearchGroups)
		builder.AddWhereCondition(searchConditions)
	}
}

// ApplyFilters applies filter conditions to the SQL builder
func (q *AdvancedQueryParams) ApplyFilters(builder *SQLBuilder) {
	if q.HasFilters() {
		filterConditions := builder.BuildFilterConditions(q.Filters)
		builder.AddWhereCondition(filterConditions)
	}
}

// ApplySort applies sort conditions to the SQL builder and returns the ORDER BY clause
func (q *AdvancedQueryParams) ApplySort(builder *SQLBuilder, includePrefix ...bool) string {
	if q.HasSort() {
		return builder.BuildOrderBy(q.Sort, includePrefix...)
	}
	return ""
}

// SetPagination sets pagination parameters
func (q *AdvancedQueryParams) SetPagination(page, limit int) {
	q.Pagination = NewPaginationParams(page, limit)
}

// SQLBuilder helps build SQL queries from QueryParams with improved practices
type SQLBuilder struct {
	whereConditions []string
	params          []any
	paramIndex      int
}

// NewSQLBuilder creates a new SQL builder
func NewSQLBuilder() *SQLBuilder {
	return &SQLBuilder{
		whereConditions: make([]string, 0),
		params:          make([]any, 0),
		paramIndex:      0,
	}
}

// BuildSearchConditions builds WHERE conditions for search (OR logic)
func (s *SQLBuilder) BuildSearchConditions(search []SearchCriteria) string {
	if len(search) == 0 {
		return ""
	}

	var conditions []string
	for _, criterion := range search {
		condition := s.buildCondition(criterion.Field, criterion.Operator, criterion.Value)
		if condition != "" {
			conditions = append(conditions, condition)
		}
	}

	if len(conditions) == 0 {
		return ""
	}

	return "(" + strings.Join(conditions, " OR ") + ")"
}

// BuildFilterConditions builds WHERE conditions for filters (AND logic)
func (s *SQLBuilder) BuildFilterConditions(filters []FilterCriteria) string {
	if len(filters) == 0 {
		return ""
	}

	var conditions []string
	for _, filter := range filters {
		condition := s.buildCondition(filter.Field, filter.Operator, filter.Value)
		if condition != "" {
			conditions = append(conditions, condition)
		}
	}

	if len(conditions) == 0 {
		return ""
	}

	return strings.Join(conditions, " AND ")
}

// BuildAdvancedSearchConditions builds complex search conditions with AND/OR logic
func (s *SQLBuilder) BuildAdvancedSearchConditions(groups []LogicalGroup) string {
	if len(groups) == 0 {
		return ""
	}

	var groupConditions []string
	for _, group := range groups {
		groupCondition := s.buildLogicalGroup(group)
		if groupCondition != "" {
			groupConditions = append(groupConditions, groupCondition)
		}
	}

	if len(groupConditions) == 0 {
		return ""
	}

	return strings.Join(groupConditions, " AND ")
}

// buildLogicalGroup builds conditions for a logical group
func (s *SQLBuilder) buildLogicalGroup(group LogicalGroup) string {
	var conditions []string

	// Add direct conditions
	for _, criterion := range group.Conditions {
		condition := s.buildCondition(criterion.Field, criterion.Operator, criterion.Value)
		if condition != "" {
			conditions = append(conditions, condition)
		}
	}

	// Add nested groups
	for _, nestedGroup := range group.Groups {
		nestedCondition := s.buildLogicalGroup(nestedGroup)
		if nestedCondition != "" {
			conditions = append(conditions, "("+nestedCondition+")")
		}
	}

	if len(conditions) == 0 {
		return ""
	}

	operator := " " + strings.ToUpper(group.Operator) + " "
	return "(" + strings.Join(conditions, operator) + ")"
}

// BuildOrderBy builds ORDER BY clause
// If includePrefix is false, returns the order clauses without "ORDER BY" prefix
func (s *SQLBuilder) BuildOrderBy(sort []SortCriteria, includePrefix ...bool) string {
	if len(sort) == 0 {
		return ""
	}

	var orderByClauses []string
	for _, criterion := range sort {
		order := "ASC"
		if strings.ToLower(criterion.Order) == "desc" {
			order = "DESC"
		}
		orderByClauses = append(orderByClauses, fmt.Sprintf("%s %s", criterion.Field, order))
	}

	orderClause := strings.Join(orderByClauses, ", ")

	// Default to including prefix if not specified
	shouldIncludePrefix := true
	if len(includePrefix) > 0 {
		shouldIncludePrefix = includePrefix[0]
	}

	if shouldIncludePrefix {
		return "ORDER BY " + orderClause
	}
	return orderClause
}

// GetParams returns the accumulated parameters
func (s *SQLBuilder) GetParams() []any {
	return s.params
}

// GetWhereClause returns the complete WHERE clause
// If includePrefix is false, returns the conditions without "WHERE" prefix
func (s *SQLBuilder) GetWhereClause(includePrefix ...bool) string {
	if len(s.whereConditions) == 0 {
		return ""
	}

	conditions := strings.Join(s.whereConditions, " AND ")

	// Default to including prefix if not specified
	shouldIncludePrefix := true
	if len(includePrefix) > 0 {
		shouldIncludePrefix = includePrefix[0]
	}

	if shouldIncludePrefix {
		return "WHERE " + conditions
	}
	return conditions
}

// AddWhereCondition adds a condition to the WHERE clause
func (s *SQLBuilder) AddWhereCondition(condition string) {
	if condition != "" {
		s.whereConditions = append(s.whereConditions, condition)
	}
}

// buildCondition builds a single condition and adds parameters
// Uses ILIKE for better search flexibility by default, but supports both LIKE and ILIKE
func (s *SQLBuilder) buildCondition(field, operator string, value any) string {
	switch operator {
	case OpEqual:
		s.params = append(s.params, value)
		return fmt.Sprintf("%s = ?", field)
	case OpNotEqual:
		s.params = append(s.params, value)
		return fmt.Sprintf("%s != ?", field)
	case OpGreaterThan:
		s.params = append(s.params, value)
		return fmt.Sprintf("%s > ?", field)
	case OpGreaterThanEq:
		s.params = append(s.params, value)
		return fmt.Sprintf("%s >= ?", field)
	case OpLessThan:
		s.params = append(s.params, value)
		return fmt.Sprintf("%s < ?", field)
	case OpLessThanEq:
		s.params = append(s.params, value)
		return fmt.Sprintf("%s <= ?", field)
	case OpContains:
		s.params = append(s.params, fmt.Sprintf("%%%v%%", value))
		return fmt.Sprintf("%s LIKE ?", field)
	case OpIContains:
		s.params = append(s.params, fmt.Sprintf("%%%v%%", value))
		return fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", field)
	case OpStartsWith:
		s.params = append(s.params, fmt.Sprintf("%v%%", value))
		return fmt.Sprintf("%s LIKE ?", field)
	case OpIStartsWith:
		s.params = append(s.params, fmt.Sprintf("%v%%", value))
		return fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", field)
	case OpEndsWith:
		s.params = append(s.params, fmt.Sprintf("%%%v", value))
		return fmt.Sprintf("%s LIKE ?", field)
	case OpIEndsWith:
		s.params = append(s.params, fmt.Sprintf("%%%v", value))
		return fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", field)
	case OpLike:
		s.params = append(s.params, value)
		return fmt.Sprintf("%s LIKE ?", field)
	case OpILike:
		s.params = append(s.params, value)
		return fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", field)
	case OpFullText:
		s.params = append(s.params, value)
		return fmt.Sprintf("MATCH(%s) AGAINST(? IN NATURAL LANGUAGE MODE)", field)
	case OpRegex:
		s.params = append(s.params, value)
		return fmt.Sprintf("%s REGEXP ?", field)
	case OpIRegex:
		s.params = append(s.params, value)
		return fmt.Sprintf("%s REGEXP ?", field) // MySQL regex is case-insensitive by default
	case OpIsNull:
		return fmt.Sprintf("%s IS NULL", field)
	case OpIsNotNull:
		return fmt.Sprintf("%s IS NOT NULL", field)
	case OpIn:
		if values, ok := value.([]any); ok && len(values) > 0 {
			placeholders := make([]string, len(values))
			for i, v := range values {
				placeholders[i] = "?"
				s.params = append(s.params, v)
			}
			return fmt.Sprintf("%s IN (%s)", field, strings.Join(placeholders, ", "))
		}
	case OpNotIn:
		if values, ok := value.([]any); ok && len(values) > 0 {
			placeholders := make([]string, len(values))
			for i, v := range values {
				placeholders[i] = "?"
				s.params = append(s.params, v)
			}
			return fmt.Sprintf("%s NOT IN (%s)", field, strings.Join(placeholders, ", "))
		}
	}
	return ""
}

// CalculatePaginationMeta calculates pagination metadata
// Formula explanation: totalPages = (totalRecords + limit - 1) / limit
// This formula ensures we round up to the nearest integer:
// - If totalRecords is exactly divisible by limit, we get the exact number of pages
// - If there's a remainder, we get one extra page to accommodate the remaining records
// Example: 23 records with 10 per page = (23 + 10 - 1) / 10 = 32 / 10 = 3 pages
func CalculatePaginationMeta(totalRecords int, page, limit int) *PaginationMeta {
	totalPages := (totalRecords + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}
	return &PaginationMeta{
		TotalItems:  totalRecords,
		TotalPage:   totalPages,
		CurrentPage: page,
		PageLimit:   limit,
	}
}

// Helper functions for creating search groups

// CreateSearchGroup creates a simple search group
func CreateSearchGroup(operator string, conditions ...SearchCriteria) LogicalGroup {
	return LogicalGroup{
		Operator:   operator,
		Conditions: conditions,
		Groups:     make([]LogicalGroup, 0),
	}
}

// CreateSearchCondition creates a search criterion
func CreateSearchCondition(field, operator string, value any) SearchCriteria {
	return SearchCriteria{
		Field:    field,
		Operator: operator,
		Value:    value,
	}
}

// CreateFilterCondition creates a filter criterion
func CreateFilterCondition(field, operator string, value any) FilterCriteria {
	return FilterCriteria{
		Field:    field,
		Operator: operator,
		Value:    value,
	}
}
