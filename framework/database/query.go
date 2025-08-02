package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"laravel-go/framework/errors"
)

// QueryBuilder ORM 查询构建器
type QueryBuilder struct {
	connection Connection
	table      string
	selects    []string
	wheres     []WhereCondition
	orders     []OrderBy
	groupBy    []string
	having     []WhereCondition
	joins      []Join
	limit      int
	offset     int
	distinct   bool
	lock       string
	ctx        context.Context
	withTrashed bool // 是否包含软删除的记录
}

// WhereCondition WHERE 条件
type WhereCondition struct {
	Column    string
	Operator  string
	Value     interface{}
	Logical   string // AND, OR
	Raw       bool   // 是否为原始 SQL
	RawSQL    string // 原始 SQL 语句
	RawArgs   []interface{}
}

// OrderBy 排序条件
type OrderBy struct {
	Column string
	Direction string // ASC, DESC
}

// Join 连接条件
type Join struct {
	Type      string // INNER, LEFT, RIGHT, FULL
	Table     string
	Condition string
	Args      []interface{}
}

// NewQueryBuilder 创建新的查询构建器
func NewQueryBuilder(connection Connection) *QueryBuilder {
	return &QueryBuilder{
		connection: connection,
		selects:    make([]string, 0),
		wheres:     make([]WhereCondition, 0),
		orders:     make([]OrderBy, 0),
		groupBy:    make([]string, 0),
		having:     make([]WhereCondition, 0),
		joins:      make([]Join, 0),
		ctx:        context.Background(),
	}
}

// Table 设置表名
func (qb *QueryBuilder) Table(table string) *QueryBuilder {
	qb.table = table
	return qb
}

// Select 设置查询字段
func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	if len(columns) == 0 {
		qb.selects = []string{"*"}
	} else {
		qb.selects = append(qb.selects, columns...)
	}
	return qb
}

// Distinct 设置去重
func (qb *QueryBuilder) Distinct() *QueryBuilder {
	qb.distinct = true
	return qb
}

// Where 添加 WHERE 条件
func (qb *QueryBuilder) Where(column string, operator string, value interface{}) *QueryBuilder {
	qb.wheres = append(qb.wheres, WhereCondition{
		Column:   column,
		Operator: operator,
		Value:    value,
		Logical:  "AND",
	})
	return qb
}

// WhereEq 等于条件
func (qb *QueryBuilder) WhereEq(column string, value interface{}) *QueryBuilder {
	return qb.Where(column, "=", value)
}

// WhereNe 不等于条件
func (qb *QueryBuilder) WhereNe(column string, value interface{}) *QueryBuilder {
	return qb.Where(column, "!=", value)
}

// WhereGt 大于条件
func (qb *QueryBuilder) WhereGt(column string, value interface{}) *QueryBuilder {
	return qb.Where(column, ">", value)
}

// WhereGte 大于等于条件
func (qb *QueryBuilder) WhereGte(column string, value interface{}) *QueryBuilder {
	return qb.Where(column, ">=", value)
}

// WhereLt 小于条件
func (qb *QueryBuilder) WhereLt(column string, value interface{}) *QueryBuilder {
	return qb.Where(column, "<", value)
}

// WhereLte 小于等于条件
func (qb *QueryBuilder) WhereLte(column string, value interface{}) *QueryBuilder {
	return qb.Where(column, "<=", value)
}

// WhereLike LIKE 条件
func (qb *QueryBuilder) WhereLike(column string, value string) *QueryBuilder {
	return qb.Where(column, "LIKE", "%"+value+"%")
}

// WhereNotLike NOT LIKE 条件
func (qb *QueryBuilder) WhereNotLike(column string, value string) *QueryBuilder {
	return qb.Where(column, "NOT LIKE", "%"+value+"%")
}

// WhereIn IN 条件
func (qb *QueryBuilder) WhereIn(column string, values []interface{}) *QueryBuilder {
	if len(values) == 0 {
		return qb.WhereRaw("1 = 0") // 空数组返回空结果
	}
	
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
	}
	
	return qb.WhereRaw(fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ",")), values...)
}

// WhereNotIn NOT IN 条件
func (qb *QueryBuilder) WhereNotIn(column string, values []interface{}) *QueryBuilder {
	if len(values) == 0 {
		return qb
	}
	
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
	}
	
	return qb.WhereRaw(fmt.Sprintf("%s NOT IN (%s)", column, strings.Join(placeholders, ",")), values...)
}

// WhereNull NULL 条件
func (qb *QueryBuilder) WhereNull(column string) *QueryBuilder {
	return qb.Where(column, "IS NULL", nil)
}

// WhereNotNull NOT NULL 条件
func (qb *QueryBuilder) WhereNotNull(column string) *QueryBuilder {
	return qb.Where(column, "IS NOT NULL", nil)
}

// WhereBetween BETWEEN 条件
func (qb *QueryBuilder) WhereBetween(column string, min, max interface{}) *QueryBuilder {
	return qb.WhereRaw(fmt.Sprintf("%s BETWEEN ? AND ?", column), min, max)
}

// WhereNotBetween NOT BETWEEN 条件
func (qb *QueryBuilder) WhereNotBetween(column string, min, max interface{}) *QueryBuilder {
	return qb.WhereRaw(fmt.Sprintf("%s NOT BETWEEN ? AND ?", column), min, max)
}

// WhereRaw 原始 WHERE 条件
func (qb *QueryBuilder) WhereRaw(sql string, args ...interface{}) *QueryBuilder {
	qb.wheres = append(qb.wheres, WhereCondition{
		Raw:     true,
		RawSQL:  sql,
		RawArgs: args,
		Logical: "AND",
	})
	return qb
}

// OrWhere OR WHERE 条件
func (qb *QueryBuilder) OrWhere(column string, operator string, value interface{}) *QueryBuilder {
	qb.wheres = append(qb.wheres, WhereCondition{
		Column:   column,
		Operator: operator,
		Value:    value,
		Logical:  "OR",
	})
	return qb
}

// OrWhereRaw OR 原始 WHERE 条件
func (qb *QueryBuilder) OrWhereRaw(sql string, args ...interface{}) *QueryBuilder {
	qb.wheres = append(qb.wheres, WhereCondition{
		Raw:     true,
		RawSQL:  sql,
		RawArgs: args,
		Logical: "OR",
	})
	return qb
}

// Join 添加连接
func (qb *QueryBuilder) Join(table, condition string, args ...interface{}) *QueryBuilder {
	qb.joins = append(qb.joins, Join{
		Type:      "INNER",
		Table:     table,
		Condition: condition,
		Args:      args,
	})
	return qb
}

// LeftJoin 左连接
func (qb *QueryBuilder) LeftJoin(table, condition string, args ...interface{}) *QueryBuilder {
	qb.joins = append(qb.joins, Join{
		Type:      "LEFT",
		Table:     table,
		Condition: condition,
		Args:      args,
	})
	return qb
}

// RightJoin 右连接
func (qb *QueryBuilder) RightJoin(table, condition string, args ...interface{}) *QueryBuilder {
	qb.joins = append(qb.joins, Join{
		Type:      "RIGHT",
		Table:     table,
		Condition: condition,
		Args:      args,
	})
	return qb
}

// OrderBy 排序
func (qb *QueryBuilder) OrderBy(column, direction string) *QueryBuilder {
	qb.orders = append(qb.orders, OrderBy{
		Column:    column,
		Direction: strings.ToUpper(direction),
	})
	return qb
}

// OrderByAsc 升序排序
func (qb *QueryBuilder) OrderByAsc(column string) *QueryBuilder {
	return qb.OrderBy(column, "ASC")
}

// OrderByDesc 降序排序
func (qb *QueryBuilder) OrderByDesc(column string) *QueryBuilder {
	return qb.OrderBy(column, "DESC")
}

// GroupBy 分组
func (qb *QueryBuilder) GroupBy(columns ...string) *QueryBuilder {
	qb.groupBy = append(qb.groupBy, columns...)
	return qb
}

// Having HAVING 条件
func (qb *QueryBuilder) Having(column string, operator string, value interface{}) *QueryBuilder {
	qb.having = append(qb.having, WhereCondition{
		Column:   column,
		Operator: operator,
		Value:    value,
		Logical:  "AND",
	})
	return qb
}

// HavingRaw 原始 HAVING 条件
func (qb *QueryBuilder) HavingRaw(sql string, args ...interface{}) *QueryBuilder {
	qb.having = append(qb.having, WhereCondition{
		Raw:     true,
		RawSQL:  sql,
		RawArgs: args,
		Logical: "AND",
	})
	return qb
}

// Limit 限制结果数量
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

// Offset 偏移量
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

// Skip 跳过记录数
func (qb *QueryBuilder) Skip(offset int) *QueryBuilder {
	return qb.Offset(offset)
}

// Take 获取记录数
func (qb *QueryBuilder) Take(limit int) *QueryBuilder {
	return qb.Limit(limit)
}

// ForUpdate 锁定更新
func (qb *QueryBuilder) ForUpdate() *QueryBuilder {
	qb.lock = "FOR UPDATE"
	return qb
}

// SharedLock 共享锁
func (qb *QueryBuilder) SharedLock() *QueryBuilder {
	qb.lock = "LOCK IN SHARE MODE"
	return qb
}

// Context 设置上下文
func (qb *QueryBuilder) Context(ctx context.Context) *QueryBuilder {
	qb.ctx = ctx
	return qb
}

// WithTrashed 包含软删除的记录
func (qb *QueryBuilder) WithTrashed() *QueryBuilder {
	qb.withTrashed = true
	return qb
}

// Get 执行查询并返回结果
func (qb *QueryBuilder) Get() ([]map[string]interface{}, error) {
	query, args := qb.buildSelectQuery()
	
	rows, err := qb.connection.QueryContext(qb.ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query")
	}
	defer rows.Close()
	
	return qb.scanRows(rows)
}

// First 获取第一条记录
func (qb *QueryBuilder) First() (map[string]interface{}, error) {
	results, err := qb.Limit(1).Get()
	if err != nil {
		return nil, err
	}
	
	if len(results) == 0 {
		return nil, sql.ErrNoRows
	}
	
	return results[0], nil
}

// Find 根据主键查找
func (qb *QueryBuilder) Find(id interface{}) (map[string]interface{}, error) {
	return qb.WhereEq("id", id).First()
}

// Count 统计记录数
func (qb *QueryBuilder) Count() (int64, error) {
	// 保存原始查询
	originalSelects := qb.selects
	originalDistinct := qb.distinct
	
	// 修改为 COUNT 查询
	qb.selects = []string{"COUNT(*)"}
	qb.distinct = false
	
	query, args := qb.buildSelectQuery()
	
	var count int64
	err := qb.connection.QueryRowContext(qb.ctx, query, args...).Scan(&count)
	
	// 恢复原始查询
	qb.selects = originalSelects
	qb.distinct = originalDistinct
	
	if err != nil {
		return 0, errors.Wrap(err, "failed to count records")
	}
	
	return count, nil
}

// Exists 检查是否存在记录
func (qb *QueryBuilder) Exists() (bool, error) {
	count, err := qb.Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// DoesntExist 检查是否不存在记录
func (qb *QueryBuilder) DoesntExist() (bool, error) {
	exists, err := qb.Exists()
	if err != nil {
		return false, err
	}
	return !exists, nil
}

// Sum 求和
func (qb *QueryBuilder) Sum(column string) (float64, error) {
	return qb.aggregate("SUM", column)
}

// Avg 平均值
func (qb *QueryBuilder) Avg(column string) (float64, error) {
	return qb.aggregate("AVG", column)
}

// Min 最小值
func (qb *QueryBuilder) Min(column string) (float64, error) {
	return qb.aggregate("MIN", column)
}

// Max 最大值
func (qb *QueryBuilder) Max(column string) (float64, error) {
	return qb.aggregate("MAX", column)
}

// aggregate 聚合函数
func (qb *QueryBuilder) aggregate(function, column string) (float64, error) {
	// 保存原始查询
	originalSelects := qb.selects
	originalDistinct := qb.distinct
	
	// 修改为聚合查询
	qb.selects = []string{fmt.Sprintf("%s(%s)", function, column)}
	qb.distinct = false
	
	query, args := qb.buildSelectQuery()
	
	var result sql.NullFloat64
	err := qb.connection.QueryRowContext(qb.ctx, query, args...).Scan(&result)
	
	// 恢复原始查询
	qb.selects = originalSelects
	qb.distinct = originalDistinct
	
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("failed to execute %s", function))
	}
	
	if !result.Valid {
		return 0, nil
	}
	
	return result.Float64, nil
}

// Paginate 分页查询
func (qb *QueryBuilder) Paginate(page, perPage int) (map[string]interface{}, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15
	}
	
	offset := (page - 1) * perPage
	
	// 获取总数
	total, err := qb.Count()
	if err != nil {
		return nil, err
	}
	
	// 获取数据
	data, err := qb.Offset(offset).Limit(perPage).Get()
	if err != nil {
		return nil, err
	}
	
	// 计算分页信息
	lastPage := int((total + int64(perPage) - 1) / int64(perPage))
	
	return map[string]interface{}{
		"data":        data,
		"total":       total,
		"per_page":    perPage,
		"current_page": page,
		"last_page":   lastPage,
		"from":        offset + 1,
		"to":          offset + len(data),
	}, nil
}

// buildSelectQuery 构建 SELECT 查询
func (qb *QueryBuilder) buildSelectQuery() (string, []interface{}) {
	var args []interface{}
	
	// SELECT 子句
	selectClause := "SELECT "
	if qb.distinct {
		selectClause += "DISTINCT "
	}
	
	if len(qb.selects) == 0 {
		selectClause += "*"
	} else {
		selectClause += strings.Join(qb.selects, ", ")
	}
	
	// FROM 子句
	fromClause := fmt.Sprintf(" FROM %s", qb.table)
	
	// JOIN 子句
	joinClause := ""
	for _, join := range qb.joins {
		joinClause += fmt.Sprintf(" %s JOIN %s ON %s", join.Type, join.Table, join.Condition)
		args = append(args, join.Args...)
	}
	
	// WHERE 子句
	whereClause, whereArgs := qb.buildWhereClause()
	args = append(args, whereArgs...)
	
	// GROUP BY 子句
	groupByClause := ""
	if len(qb.groupBy) > 0 {
		groupByClause = fmt.Sprintf(" GROUP BY %s", strings.Join(qb.groupBy, ", "))
	}
	
	// HAVING 子句
	havingClause, havingArgs := qb.buildHavingClause()
	args = append(args, havingArgs...)
	
	// ORDER BY 子句
	orderByClause := ""
	if len(qb.orders) > 0 {
		orderParts := make([]string, len(qb.orders))
		for i, order := range qb.orders {
			orderParts[i] = fmt.Sprintf("%s %s", order.Column, order.Direction)
		}
		orderByClause = fmt.Sprintf(" ORDER BY %s", strings.Join(orderParts, ", "))
	}
	
	// LIMIT 和 OFFSET 子句
	limitClause := ""
	if qb.limit > 0 {
		limitClause = fmt.Sprintf(" LIMIT %d", qb.limit)
	}
	if qb.offset > 0 {
		limitClause += fmt.Sprintf(" OFFSET %d", qb.offset)
	}
	
	// LOCK 子句
	lockClause := ""
	if qb.lock != "" {
		lockClause = " " + qb.lock
	}
	
	query := selectClause + fromClause + joinClause + whereClause + groupByClause + havingClause + orderByClause + limitClause + lockClause
	
	return query, args
}

// buildWhereClause 构建 WHERE 子句
func (qb *QueryBuilder) buildWhereClause() (string, []interface{}) {
	var args []interface{}
	var conditions []string
	
	// 添加软删除条件（除非明确要求包含软删除的记录）
	if !qb.withTrashed {
		conditions = append(conditions, "deleted_at IS NULL")
	}
	
	// 添加用户定义的WHERE条件
	for _, where := range qb.wheres {
		if len(conditions) > 0 {
			conditions = append(conditions, where.Logical)
		}
		
		if where.Raw {
			conditions = append(conditions, where.RawSQL)
			args = append(args, where.RawArgs...)
		} else {
			if where.Value == nil {
				conditions = append(conditions, fmt.Sprintf("%s %s", where.Column, where.Operator))
			} else {
				conditions = append(conditions, fmt.Sprintf("%s %s ?", where.Column, where.Operator))
				args = append(args, where.Value)
			}
		}
	}
	
	if len(conditions) == 0 {
		return "", nil
	}
	
	return " WHERE " + strings.Join(conditions, " "), args
}

// buildHavingClause 构建 HAVING 子句
func (qb *QueryBuilder) buildHavingClause() (string, []interface{}) {
	if len(qb.having) == 0 {
		return "", nil
	}
	
	var args []interface{}
	var conditions []string
	
	for i, having := range qb.having {
		if i > 0 {
			conditions = append(conditions, having.Logical)
		}
		
		if having.Raw {
			conditions = append(conditions, having.RawSQL)
			args = append(args, having.RawArgs...)
		} else {
			if having.Value == nil {
				conditions = append(conditions, fmt.Sprintf("%s %s", having.Column, having.Operator))
			} else {
				conditions = append(conditions, fmt.Sprintf("%s %s ?", having.Column, having.Operator))
				args = append(args, having.Value)
			}
		}
	}
	
	return " HAVING " + strings.Join(conditions, " "), args
}

// scanRows 扫描查询结果
func (qb *QueryBuilder) scanRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get columns")
	}
	
	var results []map[string]interface{}
	
	for rows.Next() {
		// 创建值的切片
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		
		// 扫描行
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		
		// 构建结果映射
		result := make(map[string]interface{})
		for i, column := range columns {
			val := values[i]
			
			// 处理特殊类型
			switch v := val.(type) {
			case []byte:
				result[column] = string(v)
			case time.Time:
				result[column] = v
			default:
				result[column] = v
			}
		}
		
		results = append(results, result)
	}
	
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating rows")
	}
	
	return results, nil
}

// ToSQL 生成 SQL 语句（用于调试）
func (qb *QueryBuilder) ToSQL() (string, []interface{}) {
	return qb.buildSelectQuery()
} 