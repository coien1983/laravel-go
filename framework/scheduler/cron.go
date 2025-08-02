package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CronExpression Cron 表达式结构
type CronExpression struct {
	Second     []int `json:"second"`
	Minute     []int `json:"minute"`
	Hour       []int `json:"hour"`
	DayOfMonth []int `json:"day_of_month"`
	Month      []int `json:"month"`
	DayOfWeek  []int `json:"day_of_week"`
	Year       []int `json:"year"`
}

// ParseSchedule 解析调度表达式
func ParseSchedule(schedule string) (time.Time, error) {
	// 支持多种格式
	if strings.HasPrefix(schedule, "@") {
		return parseSpecialSchedule(schedule)
	}

	// 标准 Cron 表达式
	if strings.Count(schedule, " ") >= 5 {
		return parseCronExpression(schedule)
	}

	// 简单时间格式
	return parseSimpleSchedule(schedule)
}

// parseCronExpression 解析标准 Cron 表达式
func parseCronExpression(expression string) (time.Time, error) {
	parts := strings.Fields(expression)
	if len(parts) < 5 || len(parts) > 7 {
		return time.Time{}, ErrInvalidCronExpression
	}

	cron := &CronExpression{}

	// 解析各个字段
	if err := parseField(parts[0], &cron.Second, 0, 59); err != nil {
		return time.Time{}, fmt.Errorf("invalid second field: %w", err)
	}

	if err := parseField(parts[1], &cron.Minute, 0, 59); err != nil {
		return time.Time{}, fmt.Errorf("invalid minute field: %w", err)
	}

	if err := parseField(parts[2], &cron.Hour, 0, 23); err != nil {
		return time.Time{}, fmt.Errorf("invalid hour field: %w", err)
	}

	if err := parseField(parts[3], &cron.DayOfMonth, 1, 31); err != nil {
		return time.Time{}, fmt.Errorf("invalid day of month field: %w", err)
	}

	if err := parseField(parts[4], &cron.Month, 1, 12); err != nil {
		return time.Time{}, fmt.Errorf("invalid month field: %w", err)
	}

	if len(parts) > 5 {
		if err := parseField(parts[5], &cron.DayOfWeek, 0, 6); err != nil {
			return time.Time{}, fmt.Errorf("invalid day of week field: %w", err)
		}
	}

	if len(parts) > 6 {
		if err := parseField(parts[6], &cron.Year, 1970, 2099); err != nil {
			return time.Time{}, fmt.Errorf("invalid year field: %w", err)
		}
	}

	// 计算下次运行时间
	return cron.NextRun(time.Now())
}

// parseField 解析单个字段
func parseField(field string, values *[]int, min, max int) error {
	if field == "*" {
		// 所有值
		for i := min; i <= max; i++ {
			*values = append(*values, i)
		}
		return nil
	}

	if strings.Contains(field, ",") {
		// 多个值
		parts := strings.Split(field, ",")
		for _, part := range parts {
			if err := parseSingleValue(strings.TrimSpace(part), values, min, max); err != nil {
				return err
			}
		}
		return nil
	}

	if strings.Contains(field, "-") {
		// 范围
		return parseRange(field, values, min, max)
	}

	if strings.Contains(field, "/") {
		// 步长
		return parseStep(field, values, min, max)
	}

	// 单个值
	return parseSingleValue(field, values, min, max)
}

// parseSingleValue 解析单个值
func parseSingleValue(value string, values *[]int, min, max int) error {
	val, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	if val < min || val > max {
		return fmt.Errorf("value %d out of range [%d, %d]", val, min, max)
	}

	*values = append(*values, val)
	return nil
}

// parseRange 解析范围
func parseRange(rangeStr string, values *[]int, min, max int) error {
	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid range format: %s", rangeStr)
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	if start < min || end > max || start > end {
		return fmt.Errorf("invalid range [%d, %d]", start, end)
	}

	for i := start; i <= end; i++ {
		*values = append(*values, i)
	}

	return nil
}

// parseStep 解析步长
func parseStep(stepStr string, values *[]int, min, max int) error {
	parts := strings.Split(stepStr, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid step format: %s", stepStr)
	}

	step, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	if step <= 0 {
		return fmt.Errorf("step must be positive: %d", step)
	}

	// 解析范围部分
	var rangeValues []int
	if parts[0] == "*" {
		// 从最小值开始
		for i := min; i <= max; i += step {
			rangeValues = append(rangeValues, i)
		}
	} else {
		// 解析指定范围
		if err := parseField(parts[0], &rangeValues, min, max); err != nil {
			return err
		}

		// 应用步长
		var steppedValues []int
		for i := 0; i < len(rangeValues); i += step {
			steppedValues = append(steppedValues, rangeValues[i])
		}
		rangeValues = steppedValues
	}

	*values = append(*values, rangeValues...)
	return nil
}

// NextRun 计算下次运行时间
func (c *CronExpression) NextRun(from time.Time) (time.Time, error) {
	now := from

	// 尝试最多 1000 次来找到下次运行时间
	for i := 0; i < 1000; i++ {
		// 检查年份
		if len(c.Year) > 0 && !contains(c.Year, now.Year()) {
			now = time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, now.Location())
			continue
		}

		// 检查月份
		if len(c.Month) > 0 && !contains(c.Month, int(now.Month())) {
			nextMonth := now.AddDate(0, 1, 0)
			now = time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, now.Location())
			continue
		}

		// 检查日期
		if len(c.DayOfMonth) > 0 && !contains(c.DayOfMonth, now.Day()) {
			nextDay := now.AddDate(0, 0, 1)
			now = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, now.Location())
			continue
		}

		// 检查星期
		if len(c.DayOfWeek) > 0 && !contains(c.DayOfWeek, int(now.Weekday())) {
			nextDay := now.AddDate(0, 0, 1)
			now = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, now.Location())
			continue
		}

		// 检查小时
		if len(c.Hour) > 0 && !contains(c.Hour, now.Hour()) {
			nextHour := now.Add(time.Hour)
			now = time.Date(nextHour.Year(), nextHour.Month(), nextHour.Day(), nextHour.Hour(), 0, 0, 0, now.Location())
			continue
		}

		// 检查分钟
		if len(c.Minute) > 0 && !contains(c.Minute, now.Minute()) {
			nextMinute := now.Add(time.Minute)
			now = time.Date(nextMinute.Year(), nextMinute.Month(), nextMinute.Day(), nextMinute.Hour(), nextMinute.Minute(), 0, 0, now.Location())
			continue
		}

		// 检查秒
		if len(c.Second) > 0 && !contains(c.Second, now.Second()) {
			nextSecond := now.Add(time.Second)
			now = time.Date(nextSecond.Year(), nextSecond.Month(), nextSecond.Day(), nextSecond.Hour(), nextSecond.Minute(), nextSecond.Second(), 0, now.Location())
			continue
		}

		// 找到匹配的时间
		return now, nil
	}

	return time.Time{}, fmt.Errorf("could not find next run time within 1000 iterations")
}

// contains 检查切片是否包含指定值
func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// parseSpecialSchedule 解析特殊调度表达式
func parseSpecialSchedule(schedule string) (time.Time, error) {
	now := time.Now()

	switch schedule {
	case "@yearly", "@annually":
		return time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, now.Location()), nil
	case "@monthly":
		return time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()), nil
	case "@weekly":
		daysUntilNextWeek := 7 - int(now.Weekday())
		return now.AddDate(0, 0, daysUntilNextWeek), nil
	case "@daily", "@midnight":
		return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location()), nil
	case "@hourly":
		return now.Add(time.Hour), nil
	default:
		return time.Time{}, fmt.Errorf("unknown special schedule: %s", schedule)
	}
}

// parseSimpleSchedule 解析简单时间格式
func parseSimpleSchedule(schedule string) (time.Time, error) {
	now := time.Now()

	// 尝试解析时间格式
	layouts := []string{
		"15:04",
		"15:04:05",
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
		"01-02 15:04",
		"01-02 15:04:05",
	}

	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, schedule, now.Location()); err == nil {
			// 如果时间已经过去，设置为明天
			if t.Before(now) {
				if layout == "15:04" || layout == "15:04:05" {
					t = time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), t.Second(), 0, now.Location())
					if t.Before(now) {
						t = t.AddDate(0, 0, 1)
					}
				} else if strings.HasPrefix(layout, "01-02") {
					t = time.Date(now.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, now.Location())
					if t.Before(now) {
						t = t.AddDate(1, 0, 0)
					}
				}
			}
			return t, nil
		}
	}

	return time.Time{}, ErrInvalidTimeFormat
}
