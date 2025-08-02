package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Model 基础模型结构体
type Model struct {
	ID        int64      `db:"id"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// Hook 钩子接口
type Hook interface{}

// BeforeSave 保存前钩子
type BeforeSave interface {
	BeforeSave(conn Connection) error
}

// AfterSave 保存后钩子
type AfterSave interface {
	AfterSave(conn Connection) error
}

// BeforeDelete 删除前钩子
type BeforeDelete interface {
	BeforeDelete(conn Connection) error
}

// AfterDelete 删除后钩子
type AfterDelete interface {
	AfterDelete(conn Connection) error
}

// TableName 返回表名（可重写）
func (m *Model) TableName() string {
	// 默认表名为结构体名小写+s
	typeName := reflect.Indirect(reflect.ValueOf(m)).Type().Name()
	return strings.ToLower(typeName) + "s"
}

// PrimaryKey 返回主键字段名
func (m *Model) PrimaryKey() string {
	return "id"
}

// callHook 调用钩子方法
func callHook(model interface{}, hookName string, conn Connection) error {
	switch hookName {
	case "BeforeSave":
		if h, ok := model.(BeforeSave); ok {
			return h.BeforeSave(conn)
		}
	case "AfterSave":
		if h, ok := model.(AfterSave); ok {
			return h.AfterSave(conn)
		}
	case "BeforeDelete":
		if h, ok := model.(BeforeDelete); ok {
			return h.BeforeDelete(conn)
		}
	case "AfterDelete":
		if h, ok := model.(AfterDelete); ok {
			return h.AfterDelete(conn)
		}
	}
	return nil
}

// Find 根据主键查找
func (m *Model) Find(conn Connection, id interface{}, dest interface{}) error {
	table := getTableName(dest)
	pk := getPrimaryKey(dest)
	qb := NewQueryBuilder(conn).Table(table).WhereEq(pk, id).Limit(1)
	row, err := qb.First()
	if err != nil {
		return err
	}
	return mapToStruct(row, dest)
}

// First 获取第一条记录
func (m *Model) First(conn Connection, dest interface{}) error {
	table := getTableName(dest)
	qb := NewQueryBuilder(conn).Table(table).Limit(1)
	row, err := qb.First()
	if err != nil {
		return err
	}
	return mapToStruct(row, dest)
}

// All 获取所有记录
func (m *Model) All(conn Connection, destSlice interface{}) error {
	table := getTableName(destSlice)
	qb := NewQueryBuilder(conn).Table(table)
	rows, err := qb.Get()
	if err != nil {
		return err
	}
	return mapToSlice(rows, destSlice)
}

// Where 条件查询
func (m *Model) Where(conn Connection, destSlice interface{}, column string, operator string, value interface{}) error {
	table := getTableName(destSlice)
	qb := NewQueryBuilder(conn).Table(table).Where(column, operator, value)
	rows, err := qb.Get()
	if err != nil {
		return err
	}
	return mapToSlice(rows, destSlice)
}

// Count 统计记录数
func (m *Model) Count(conn Connection, dest interface{}) (int64, error) {
	table := getTableName(dest)
	return NewQueryBuilder(conn).Table(table).Count()
}

// Exists 检查记录是否存在
func (m *Model) Exists(conn Connection, dest interface{}, column string, value interface{}) (bool, error) {
	table := getTableName(dest)
	return NewQueryBuilder(conn).Table(table).WhereEq(column, value).Exists()
}

// Save 保存记录（插入或更新）
func (m *Model) Save(conn Connection, model interface{}) error {
	// 调用 BeforeSave 钩子
	if err := callHook(model, "BeforeSave", conn); err != nil {
		return err
	}

	table := getTableName(model)
	pk := getPrimaryKey(model)

	// 获取模型值
	modelVal := reflect.ValueOf(model).Elem()

	// 处理嵌套字段路径
	var pkField reflect.Value
	modelField := modelVal.FieldByName("Model")
	if modelField.IsValid() {
		// 有嵌入的Model结构体
		pkField = modelField.FieldByName("ID")
	} else {
		// 直接字段
		pkField = modelVal.FieldByName(pk)
	}

	// 检查主键字段是否有效
	if !pkField.IsValid() {
		// 调试：打印所有字段名
		var fieldNames []string
		for i := 0; i < modelVal.NumField(); i++ {
			fieldNames = append(fieldNames, modelVal.Type().Field(i).Name)
		}
		return fmt.Errorf("primary key field '%s' not found. Available fields: %v", pk, fieldNames)
	}

	// 自动时间戳
	now := time.Now()
	var createdAtField, updatedAtField reflect.Value

	if modelField.IsValid() {
		// 有嵌入的Model结构体
		createdAtField = modelField.FieldByName("CreatedAt")
		updatedAtField = modelField.FieldByName("UpdatedAt")
	} else {
		// 直接字段
		createdAtField = modelVal.FieldByName("CreatedAt")
		updatedAtField = modelVal.FieldByName("UpdatedAt")
	}

	// 检查主键值
	var pkValue int64
	if pkField.Kind() == reflect.Int64 {
		pkValue = pkField.Int()
	} else if pkField.Kind() == reflect.Int {
		pkValue = int64(pkField.Int())
	} else {
		return errors.New("primary key field must be int or int64")
	}

	if pkValue == 0 {
		// 插入新记录
		if createdAtField.IsValid() && createdAtField.IsNil() {
			createdAtField.Set(reflect.ValueOf(&now))
		}
		if updatedAtField.IsValid() && updatedAtField.IsNil() {
			updatedAtField.Set(reflect.ValueOf(&now))
		}

		// 构建插入SQL
		data := structToMap(model)
		columns := make([]string, 0, len(data))
		values := make([]interface{}, 0, len(data))
		placeholders := make([]string, 0, len(data))

		for col, val := range data {
			if col != pk { // 跳过主键字段
				columns = append(columns, col)
				values = append(values, val)
				placeholders = append(placeholders, "?")
			}
		}

		sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			table, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

		result, err := conn.Exec(sqlStr, values...)
		if err != nil {
			return err
		}

		// 设置自增ID
		if id, err := result.LastInsertId(); err == nil {
			if pkField.Kind() == reflect.Int64 {
				pkField.SetInt(id)
			} else if pkField.Kind() == reflect.Int {
				pkField.SetInt(id)
			}
		}
	} else {
		// 更新记录
		if updatedAtField.IsValid() {
			updatedAtField.Set(reflect.ValueOf(&now))
		}

		// 构建更新SQL
		data := structToMap(model)
		sets := make([]string, 0, len(data))
		values := make([]interface{}, 0, len(data))

		for col, val := range data {
			if col != pk { // 跳过主键
				sets = append(sets, fmt.Sprintf("%s = ?", col))
				values = append(values, val)
			}
		}
		values = append(values, pkValue) // WHERE条件值

		sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?",
			table, strings.Join(sets, ", "), pk)

		_, err := conn.Exec(sqlStr, values...)
		if err != nil {
			return err
		}
	}

	// 调用 AfterSave 钩子
	return callHook(model, "AfterSave", conn)
}

// Delete 删除记录（软删除）
func (m *Model) Delete(conn Connection, model interface{}) error {
	// 调用 BeforeDelete 钩子
	if err := callHook(model, "BeforeDelete", conn); err != nil {
		return err
	}

	table := getTableName(model)
	pk := getPrimaryKey(model)

	// 获取模型值
	modelVal := reflect.ValueOf(model).Elem()

	// 处理嵌套字段路径
	var pkField, deletedAtField reflect.Value
	modelField := modelVal.FieldByName("Model")
	if modelField.IsValid() {
		// 有嵌入的Model结构体
		pkField = modelField.FieldByName("ID")
		deletedAtField = modelField.FieldByName("DeletedAt")
	} else {
		// 直接字段
		pkField = modelVal.FieldByName(pk)
		deletedAtField = modelVal.FieldByName("DeletedAt")
	}

	// 检查主键字段是否有效
	if !pkField.IsValid() {
		return errors.New("primary key field not found")
	}

	// 获取主键值
	var pkValue int64
	if pkField.Kind() == reflect.Int64 {
		pkValue = pkField.Int()
	} else if pkField.Kind() == reflect.Int {
		pkValue = int64(pkField.Int())
	} else {
		return errors.New("primary key field must be int or int64")
	}

	if deletedAtField.IsValid() {
		// 软删除：设置 deleted_at 字段
		now := time.Now()
		deletedAtField.Set(reflect.ValueOf(&now))

		sqlStr := fmt.Sprintf("UPDATE %s SET deleted_at = ? WHERE %s = ?", table, pk)
		_, err := conn.Exec(sqlStr, &now, pkValue)
		if err != nil {
			return err
		}
	} else {
		// 硬删除
		sqlStr := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", table, pk)
		_, err := conn.Exec(sqlStr, pkValue)
		if err != nil {
			return err
		}
	}

	// 调用 AfterDelete 钩子
	return callHook(model, "AfterDelete", conn)
}

// HasOne 一对一关联
func (m *Model) HasOne(conn Connection, dest interface{}, foreignKey, localKey string) error {
	modelVal := reflect.ValueOf(m).Elem()

	// 处理嵌套字段路径
	var localField reflect.Value
	modelField := modelVal.FieldByName("Model")
	if modelField.IsValid() && localKey == "ID" {
		// 有嵌入的Model结构体，且查找的是ID字段
		localField = modelField.FieldByName("ID")
	} else {
		// 直接字段
		localField = modelVal.FieldByName(localKey)
	}

	if !localField.IsValid() {
		return errors.New("local key field not found")
	}

	localVal := localField.Interface()
	table := getTableName(dest)

	qb := NewQueryBuilder(conn).Table(table).WhereEq(foreignKey, localVal).Limit(1)
	row, err := qb.First()
	if err != nil {
		return err
	}
	return mapToStruct(row, dest)
}

// HasMany 一对多关联
func (m *Model) HasMany(conn Connection, destSlice interface{}, foreignKey, localKey string) error {
	modelVal := reflect.ValueOf(m).Elem()

	// 处理嵌套字段路径
	var localField reflect.Value
	modelField := modelVal.FieldByName("Model")
	if modelField.IsValid() && localKey == "ID" {
		// 有嵌入的Model结构体，且查找的是ID字段
		localField = modelField.FieldByName("ID")
	} else {
		// 直接字段
		localField = modelVal.FieldByName(localKey)
	}

	if !localField.IsValid() {
		return errors.New("local key field not found")
	}

	localVal := localField.Interface()
	table := getTableName(destSlice)

	qb := NewQueryBuilder(conn).Table(table).WhereEq(foreignKey, localVal)
	rows, err := qb.Get()
	if err != nil {
		return err
	}
	return mapToSlice(rows, destSlice)
}

// BelongsTo 反向一对一关联
// 注意：由于Go的嵌入机制限制，此方法在嵌入Model的结构体上调用时可能无法正确获取字段
// 建议使用 BelongsToModel 方法替代
func (m *Model) BelongsTo(conn Connection, dest interface{}, ownerKey, foreignKey string) error {
	// 获取调用者的结构体（如Post）
	// 由于m是嵌入的Model，我们需要获取包含它的结构体
	modelVal := reflect.ValueOf(m).Elem()

	// 如果当前结构体就是Model，我们需要找到包含它的父结构体
	if modelVal.Type().Name() == "Model" {
		// 如果调用者是Model类型，我们需要通过其他方式获取父结构体
		// 这种情况下，我们需要修改调用方式
		return fmt.Errorf("BelongsTo called on Model directly. This method should be called on a concrete model instance that embeds Model")
	}

	// 查找外键字段
	var foreignField reflect.Value

	// 首先尝试直接查找字段（如Post结构体中的UserID字段）
	foreignField = modelVal.FieldByName(foreignKey)

	// 如果没找到，检查是否有嵌入的Model结构体
	if !foreignField.IsValid() {
		modelField := modelVal.FieldByName("Model")
		if modelField.IsValid() && foreignKey == "ID" {
			foreignField = modelField.FieldByName("ID")
		}
	}

	if !foreignField.IsValid() {
		// 调试：打印所有字段名
		var fieldNames []string
		for i := 0; i < modelVal.NumField(); i++ {
			fieldNames = append(fieldNames, modelVal.Type().Field(i).Name)
		}
		return fmt.Errorf("foreign key field '%s' not found. Available fields: %v", foreignKey, fieldNames)
	}

	foreignVal := foreignField.Interface()
	table := getTableName(dest)

	qb := NewQueryBuilder(conn).Table(table).WhereEq(ownerKey, foreignVal).Limit(1)
	row, err := qb.First()
	if err != nil {
		return err
	}
	return mapToStruct(row, dest)
}

// BelongsToModel 反向一对一关联（接收具体模型实例）
// 这是BelongsTo方法的推荐替代方案，可以正确处理嵌入Model的结构体
// 参数说明：
// - conn: 数据库连接
// - model: 包含外键的具体模型实例（如Post）
// - dest: 目标模型实例（如User）
// - ownerKey: 目标模型的主键字段名（如"ID"）
// - foreignKey: 当前模型的外键字段名（如"UserID"）
func (m *Model) BelongsToModel(conn Connection, model interface{}, dest interface{}, ownerKey, foreignKey string) error {
	// 获取具体模型的结构体（如Post）
	modelVal := reflect.ValueOf(model).Elem()

	// 查找外键字段
	var foreignField reflect.Value

	// 首先尝试直接查找字段（如Post结构体中的UserID字段）
	foreignField = modelVal.FieldByName(foreignKey)

	// 如果没找到，检查是否有嵌入的Model结构体
	if !foreignField.IsValid() {
		modelField := modelVal.FieldByName("Model")
		if modelField.IsValid() && foreignKey == "ID" {
			foreignField = modelField.FieldByName("ID")
		}
	}

	if !foreignField.IsValid() {
		// 调试：打印所有字段名
		var fieldNames []string
		for i := 0; i < modelVal.NumField(); i++ {
			fieldNames = append(fieldNames, modelVal.Type().Field(i).Name)
		}
		return fmt.Errorf("foreign key field '%s' not found. Available fields: %v", foreignKey, fieldNames)
	}

	foreignVal := foreignField.Interface()
	table := getTableName(dest)

	qb := NewQueryBuilder(conn).Table(table).WhereEq(ownerKey, foreignVal).Limit(1)
	row, err := qb.First()
	if err != nil {
		return err
	}
	return mapToStruct(row, dest)
}

// ManyToMany 多对多关联
func (m *Model) ManyToMany(conn Connection, destSlice interface{}, pivotTable, localKey, foreignPivotKey, relatedPivotKey, relatedKey string) error {
	modelVal := reflect.ValueOf(m).Elem()

	// 处理嵌套字段路径
	var localField reflect.Value
	modelField := modelVal.FieldByName("Model")
	if modelField.IsValid() && localKey == "ID" {
		// 有嵌入的Model结构体，且查找的是ID字段
		localField = modelField.FieldByName("ID")
	} else {
		// 直接字段
		localField = modelVal.FieldByName(localKey)
	}

	if !localField.IsValid() {
		return errors.New("local key field not found")
	}

	localVal := localField.Interface()
	table := getTableName(destSlice)

	sqlStr := fmt.Sprintf(
		"SELECT t.* FROM %s t JOIN %s p ON t.%s = p.%s WHERE p.%s = ?",
		table, pivotTable, relatedKey, relatedPivotKey, foreignPivotKey,
	)

	rows, err := conn.QueryContext(context.Background(), sqlStr, localVal)
	if err != nil {
		return err
	}
	defer rows.Close()

	return scanRowsToSlice(rows, destSlice)
}

// 辅助函数

// getTableName 获取表名
func getTableName(model interface{}) string {
	if t, ok := model.(interface{ TableName() string }); ok {
		return t.TableName()
	}

	// 处理切片类型
	modelVal := reflect.ValueOf(model)
	if modelVal.Kind() == reflect.Ptr {
		modelVal = modelVal.Elem()
	}

	// 如果是切片，获取元素类型
	if modelVal.Kind() == reflect.Slice {
		elementType := modelVal.Type().Elem()
		// 创建元素类型的实例来调用TableName方法
		element := reflect.New(elementType)
		if t, ok := element.Interface().(interface{ TableName() string }); ok {
			return t.TableName()
		}
		// 默认表名为元素类型名小写+s
		return strings.ToLower(elementType.Name()) + "s"
	}

	// 默认表名为结构体名小写+s
	typeName := modelVal.Type().Name()
	return strings.ToLower(typeName) + "s"
}

// getPrimaryKey 获取主键字段名
func getPrimaryKey(model interface{}) string {
	if p, ok := model.(interface{ PrimaryKey() string }); ok {
		return p.PrimaryKey()
	}
	return "id"
}

// getPrimaryKeyField 获取主键字段的完整路径
func getPrimaryKeyField(model interface{}) string {
	modelVal := reflect.ValueOf(model).Elem()
	modelType := modelVal.Type()

	// 检查是否有嵌入的Model结构体
	for i := 0; i < modelVal.NumField(); i++ {
		fieldType := modelType.Field(i)

		// 如果是嵌入的Model结构体
		if fieldType.Anonymous && fieldType.Type.Name() == "Model" {
			return "Model.ID"
		}
	}

	// 如果没有嵌入Model，直接返回主键名
	return getPrimaryKey(model)
}

// mapToStruct 将map映射到结构体
func mapToStruct(data map[string]interface{}, dest interface{}) error {
	if data == nil {
		return errors.New("data is nil")
	}

	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr {
		return errors.New("dest must be a pointer")
	}

	destVal = destVal.Elem()
	destType := destVal.Type()

	for i := 0; i < destVal.NumField(); i++ {
		field := destVal.Field(i)
		fieldType := destType.Field(i)

		// 处理嵌入的Model结构体
		if fieldType.Anonymous && fieldType.Type.Name() == "Model" {
			// 递归处理Model结构体的字段
			if err := mapToStruct(data, field.Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		// 获取db标签
		dbTag := fieldType.Tag.Get("db")
		if dbTag == "" {
			dbTag = strings.ToLower(fieldType.Name)
		}

		if value, exists := data[dbTag]; exists && value != nil {
			// 设置字段值
			fieldVal := reflect.ValueOf(value)
			if fieldVal.Type().ConvertibleTo(field.Type()) {
				field.Set(fieldVal.Convert(field.Type()))
			}
		}
	}

	return nil
}

// mapToSlice 将map切片映射到结构体切片
func mapToSlice(data []map[string]interface{}, destSlice interface{}) error {
	if data == nil {
		return errors.New("data is nil")
	}

	destVal := reflect.ValueOf(destSlice)
	if destVal.Kind() != reflect.Ptr || destVal.Elem().Kind() != reflect.Slice {
		return errors.New("destSlice must be a pointer to slice")
	}

	sliceVal := destVal.Elem()
	elementType := sliceVal.Type().Elem()

	for _, item := range data {
		element := reflect.New(elementType)
		if err := mapToStruct(item, element.Interface()); err != nil {
			return err
		}
		sliceVal.Set(reflect.Append(sliceVal, element.Elem()))
	}

	return nil
}

// structToMap 将结构体转换为map
func structToMap(model interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	modelVal := reflect.ValueOf(model).Elem()
	modelType := modelVal.Type()

	for i := 0; i < modelVal.NumField(); i++ {
		field := modelVal.Field(i)
		fieldType := modelType.Field(i)

		// 跳过未导出字段
		if !field.CanInterface() {
			continue
		}

		// 处理嵌入的Model结构体
		if fieldType.Anonymous && fieldType.Type.Name() == "Model" {
			// 递归处理Model结构体的字段
			modelMap := structToMap(field.Addr().Interface())
			for k, v := range modelMap {
				result[k] = v
			}
			continue
		}

		// 获取db标签
		dbTag := fieldType.Tag.Get("db")
		if dbTag == "" {
			dbTag = strings.ToLower(fieldType.Name)
		}

		// 跳过db标签为"-"的字段
		if dbTag == "-" {
			continue
		}

		// 处理指针类型
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			result[dbTag] = field.Elem().Interface()
		} else if field.Kind() != reflect.Ptr {
			result[dbTag] = field.Interface()
		}

		// 特殊处理：如果字段值不为零，也要包含在结果中
		if field.Kind() == reflect.Ptr && field.IsNil() {
			// 跳过nil指针
		} else if !isZero(field) {
			result[dbTag] = field.Interface()
		}
	}

	return result
}

// isZero 检查值是否为零值
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	}
	return false
}

// scanRowsToSlice 扫描查询结果到切片
func scanRowsToSlice(rows *sql.Rows, destSlice interface{}) error {
	destVal := reflect.ValueOf(destSlice)
	if destVal.Kind() != reflect.Ptr || destVal.Elem().Kind() != reflect.Slice {
		return errors.New("destSlice must be a pointer to slice")
	}

	sliceVal := destVal.Elem()
	elementType := sliceVal.Type().Elem()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	for rows.Next() {
		// 创建临时map存储行数据
		rowData := make(map[string]interface{})
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return err
		}

		for i, col := range columns {
			rowData[col] = values[i]
		}

		// 创建结构体实例并映射数据
		element := reflect.New(elementType)
		if err := mapToStruct(rowData, element.Interface()); err != nil {
			return err
		}
		sliceVal.Set(reflect.Append(sliceVal, element.Elem()))
	}

	return rows.Err()
}
