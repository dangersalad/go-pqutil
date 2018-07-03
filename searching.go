package pqutil

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ParseAmountQuery parses a set of queries for a key and returns an
// SQL string and an array of arguments.
//
// The SQL generated uses ? placeholders. It is meant to be used in
// conjunction with github.com/Masterminds/squirrel
func ParseAmountQuery(key string, qs []string) (interface{}, []interface{}, error) {

	var sqls []string
	var args []interface{}

	for _, q := range qs {
		parts := strings.Split(q, "_")

		opFlag := "eq"
		value := ""

		// if only an amount, set value to first part
		if len(parts) == 1 {
			value = parts[0]
		} else if len(parts) > 2 { // if too many parts, error out
			return nil, nil, errors.Errorf("invalid query for %s: %s", key, q)
		} else { // assign first part to opFlag, second to value
			value = parts[1]
			opFlag = parts[0]
		}

		// if value ended up blank, we have a problem
		if value == "" {
			return nil, nil, errors.Errorf("invalid query for %s: %s", key, q)
		}

		op := ""

		switch opFlag {
		case "gt":
			op = ">"
		case "lt":
			op = "<"
		case "gte":
			op = ">="
		case "lte":
			op = "<="
		default:
			op = "="
		}

		sqlStr := fmt.Sprintf("%s %s ?", key, op)
		sqls = append(sqls, sqlStr)

		if dateVal, err := time.Parse(time.RFC3339, value); err == nil {
			args = append(args, dateVal)
		} else if dateVal, err := time.Parse(time.RFC3339Nano, value); err == nil {
			args = append(args, dateVal)
		} else if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			args = append(args, intVal)
		} else {
			args = append(args, value)
		}
	}

	return fmt.Sprintf(`(%s)`, strings.Join(sqls, " AND ")), args, nil

}

var sortFilter = regexp.MustCompile(`[^a-z]`)

// ParseSort will take a string containing a column and direction and
// return them or an error if the string is invalid
func ParseSort(sort string) (column, order string, err error) {
	sort = strings.TrimSpace(sort)
	parts := strings.Split(sort, ":")
	if len(parts) == 0 || len(parts) > 2 {
		return "", "", errors.Errorf("invalid sort: %s", sort)
	}
	column = parts[0]
	o := "desc"
	if len(parts) == 2 {
		o = sortFilter.ReplaceAllString(strings.ToLower(parts[1]), "")
	}

	switch o {
	case "desc", "d", "down", "descending", "hightolow", "":
		order = "DESC"
	case "asc", "a", "up", "ascending", "lowtohigh":
		order = "ASC"
	default:
		return "", "", errors.Errorf("invalid sort order: %s", o)
	}

	return

}
