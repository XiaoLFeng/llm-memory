package utils

import (
	"fmt"
	"time"
)

// FormatDate 格式化日期为 YYYY-MM-DD 格式
// 参数: t - 要格式化的时间
// 返回: 格式化后的日期字符串
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDateTime 格式化日期时间为 YYYY-MM-DD HH:MM:SS 格式
// 参数: t - 要格式化的时间
// 返回: 格式化后的日期时间字符串
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatRelative 相对时间格式化，返回如"3天前"、"刚刚"等友好格式
// 参数: t - 要格式化的时间
// 返回: 相对时间字符串
func FormatRelative(t time.Time) string {
	now := time.Now()
	duration := now.Sub(t)

	if duration < time.Minute {
		if duration < time.Second {
			return "刚刚"
		}
		seconds := int(duration.Seconds())
		return fmt.Sprintf("%d秒前", seconds)
	}

	if duration < time.Hour {
		minutes := int(duration.Minutes())
		return fmt.Sprintf("%d分钟前", minutes)
	}

	if duration < time.Hour*24 {
		hours := int(duration.Hours())
		return fmt.Sprintf("%d小时前", hours)
	}

	if duration < time.Hour*24*7 {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%d天前", days)
	}

	if duration < time.Hour*24*30 {
		weeks := int(duration.Hours() / 24 / 7)
		return fmt.Sprintf("%d周前", weeks)
	}

	if duration < time.Hour*24*365 {
		months := int(duration.Hours() / 24 / 30)
		return fmt.Sprintf("%d个月前", months)
	}

	years := int(duration.Hours() / 24 / 365)
	return fmt.Sprintf("%d年前", years)
}

// IsToday 判断给定时间是否为今天
// 参数: t - 要判断的时间
// 返回: 如果是今天返回true，否则返回false
func IsToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() &&
		t.Month() == now.Month() &&
		t.Day() == now.Day()
}

// IsPast 判断给定时间是否已经过期
// 参数: t - 要判断的时间
// 返回: 如果时间已过期返回true，否则返回false
func IsPast(t time.Time) bool {
	return t.Before(time.Now())
}

// IsUpcoming 判断给定时间是否在未来N天内
// 参数: t - 要判断的时间
// 参数: days - 未来的天数
// 返回: 如果在未来N天内返回true，否则返回false
func IsUpcoming(t time.Time, days int) bool {
	if t.Before(time.Now()) {
		return false
	}

	future := time.Now().AddDate(0, 0, days)
	return t.Before(future) || t.Equal(future)
}

// FormatISO8601 格式化时间为 ISO8601 标准格式
// 参数: t - 要格式化的时间
// 返回: ISO8601 格式的时间字符串
func FormatISO8601(t time.Time) string {
	return t.Format(time.RFC3339)
}

// FormatTime 格式化时间为 HH:MM:SS 格式
// 参数: t - 要格式化的时间
// 返回: 格式化后的时间字符串
func FormatTime(t time.Time) string {
	return t.Format("15:04:05")
}

// StartOfDay 获取指定时间的当天的开始时间 (00:00:00)
// 参数: t - 基准时间
// 返回: 当天的开始时间
func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// EndOfDay 获取指定时间的当天的结束时间 (23:59:59)
// 参数: t - 基准时间
// 返回: 当天的结束时间
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

// DaysBetween 计算两个时间之间的天数差
// 参数: start - 开始时间
// 参数: end - 结束时间
// 返回: 天数差
func DaysBetween(start, end time.Time) int {
	duration := end.Sub(start)
	return int(duration.Hours() / 24)
}
