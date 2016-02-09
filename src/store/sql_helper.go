package store

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"github.com/golang/glog"
)

func SqlSelect(tableName string, fields []string) string {
	return FormatSql("s", tableName, fields)
}

func SqlInsert(tableName string, fields []string) string {
	return FormatSql("i", tableName, fields)
}

func SqlUpdate(tableName string, fields []string) string {
	return FormatSql("u", tableName, fields)
}

// Postgres = $#
// Oracle = :
// Default = ?
func FormateToPQuery(query string) string {
	var qrx = regexp.MustCompile(`\?`)
	var pref = "$"

	n := 0
	return qrx.ReplaceAllStringFunc(query, func(string) string {
		n++
		return pref + strconv.Itoa(n)
	})
}

func FormateToPG(query string, fields []string) string {
	var qrx = regexp.MustCompile(`\?`)
	var pref = "?"

	n := 0
	return qrx.ReplaceAllStringFunc(query, func(string) string {

		res := pref + fields[n]
		n++
		return res
	})
}

func FormatSql(_type string, tableName string, arraysField []string) (sql string) {
	switch _type {
	case "s":
		sql = fmt.Sprintf("SELECT %s FROM %s", strings.Join(arraysField, ", "), tableName)
	case "i":
		sql = fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", tableName, strings.Join(arraysField, ", "), strings.Repeat("?,", len(arraysField))[:len(arraysField)*2-1])
	case "u":
		setFields := []string{}
		for _, fieldName := range arraysField {
			setFields = append(setFields, fmt.Sprintf("%s = ?", fieldName))
		}
		sql = fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(setFields, ", "))
	default:
		glog.Warningf("format sql: not supported type '%s'", _type)
	}

	return sql
}
