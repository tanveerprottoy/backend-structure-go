package sqlext

import "strconv"

func BuildInsertQuery(tableName string, columns []string, clause string) string {
	q := "INSERT INTO " + tableName
	cols := " ("
	vals := " VALUES ("
	for i, v := range columns {
		cols += v + ", "
		vals += "$" + strconv.Itoa(i+1) + ", "
	}
	// remove trailing space & comma
	// there are two chars to be removed
	cols = cols[:len(cols)-2]
	vals = vals[:len(vals)-2]
	// add closing parentheses
	cols += ")"
	vals += ")"
	return q + cols + vals + " " + clause
}

func BuildSelectQuery(tableName string, projections []string, whereClauseCols []string, clause string, args ...string) string {
	q := "SELECT "
	p := ""
	w := ""
	op := "AND"
	cl := ""
	// set clause & operator
	if clause != "" {
		cl = clause
	}
	if len(args) > 0 {
		// operator is in the first arg
		op = args[0]
	}
	if len(projections) > 0 {
		for _, v := range projections {
			p += v + ", "
		}
		// remove trailing space & comma
		// there are two chars to be removed
		p = p[:len(p)-2]
	} else {
		p += "*"
	}
	if len(whereClauseCols) > 0 {
		opLen := len(op)
		w += " WHERE "
		for i, v := range whereClauseCols {
			w += v + " = $" + strconv.Itoa(i+1) + " " + op + " "
		}
		// remove trailing space & operator
		// there are two spaces and operator length to be removed
		w = w[:len(w)-opLen-2]
	}
	return q + p + " FROM " + tableName + w + " " + cl
}

func BuildWhereClause(tableName string, whereClauseCols []string, clause string, args ...string) string {
	w := " WHERE "
	op := "AND"
	cl := ""
	// set clause & operator
	if clause != "" {
		cl = clause
	}
	if len(args) > 0 {
		// operator is in the first arg
		op = args[0]
	}
	for i, v := range whereClauseCols {
		w += v + " = $" + strconv.Itoa(i+1) + " " + op + " "
	}
	// remove trailing space & operator
	w = w[:len(w)-len(op)-2]
	return w + " " + cl
}

func BuildUpdateQuery(tableName string, columns []string, whereClauseCols []string, clause string, args ...any) string {
	q := "UPDATE " + tableName
	c := " SET "
	w := " WHERE "
	op := "AND"
	if len(args) > 1 {
		op = args[0].(string)
	}
	// col length
	lenCols := len(columns)
	if lenCols > 0 {
		for i, v := range columns {
			c += v + " = $" + strconv.Itoa(i+1) + ", "
		}
		// remove trailing space & comma
		// there are two chars to be removed
		c = c[:len(c)-2]
	}
	for i, v := range whereClauseCols {
		w += v + " = $" + strconv.Itoa(i+1+lenCols) + " " + op + " "
	}
	// remove trailing space & comma
	w = w[:len(w)-len(op)-2]
	if clause == "" {
		return q + c + w
	}
	return q + c + w + " " + clause
}

func BuildDeleteQuery(tableName string, whereClauseCols []string, clause string) string {
	q := "DELETE FROM " + tableName
	w := " WHERE "
	for i, v := range whereClauseCols {
		w += v + " = $" + strconv.Itoa(i+1) + ", "
	}
	// remove trailing space & comma
	// there are two chars to be removed
	w = w[:len(w)-2]
	return q + w + " " + clause
}

func BuildQueryPlaceholder(count int) string {
	q := "("
	for i := 0; i < count; i++ {
		q += "$" + strconv.Itoa(i+1) + ", "
	}
	// remove trailing space & comma
	q = q[:len(q)-2]
	q += ")"
	return q
}

func BuildSetColumnsQuery(columns []string) string {
	q := "SET "
	// col length
	lenCols := len(columns)
	if lenCols > 0 {
		for i, v := range columns {
			q += v + " = $" + strconv.Itoa(i+1) + ", "
		}
		// remove trailing space & comma
		// there are two chars to be removed
		q = q[:len(q)-2]
	}
	return q
}
