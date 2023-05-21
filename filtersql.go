package filtersql

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/samber/lo"
	"vitess.io/vitess/go/vt/sqlparser"
)

const UNLIMITED = -1

func (config Config) Parse(filter string) (string, error) {
	if strings.Trim(filter, " ") == "" {
		return "", nil
	}

	sqlBlob := fmt.Sprintf("SELECT * FROM `not_a_table` WHERE %s", filter)
	sql, remainder, err := sqlparser.SplitStatement(sqlBlob)
	if err != nil {
		return "", err
	}
	if remainder != "" {
		return "", errors.New("unsupported syntax")
	}

	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return "", errors.New("unsupported syntax")
	}

	where := stmt.(*sqlparser.Select).Where

	if err = config.validateGroupingParens(sql); err != nil {
		return "", err
	}

	if err = config.validateAnds(sql); err != nil {
		return "", err
	}

	if err = config.validateOrs(sql); err != nil {
		return "", err
	}

	if err = config.validateNots(sql); err != nil {
		return "", err
	}

	if err = config.validateAST(where); err != nil {
		return "", err
	}

	result := sqlparser.String(where)
	result = strings.Replace(result, " where ", "", 1)
	return result, nil
}

func (config Config) validateAST(filter *sqlparser.Where) error {
	fun := func(node sqlparser.SQLNode) (bool, error) {
		switch node := node.(type) {
		case *sqlparser.Where, *sqlparser.Literal, sqlparser.ValTuple:
			return true, nil
		case *sqlparser.AndExpr:
			max := config.Allow.Ands
			if max == UNLIMITED || max > 0 {
				return true, nil
			} else {
				return config.walkError("unsupported and: %s", node)
			}
		case *sqlparser.OrExpr:
			max := config.Allow.Ors
			if max == UNLIMITED || max > 0 {
				return true, nil
			} else {
				return config.walkError("unsupported or: %s", node)
			}
		case *sqlparser.NotExpr:
			max := config.Allow.Nots
			if max == UNLIMITED || max > 0 {
				return true, nil
			} else {
				return config.walkError("unsupported not: %s", node)
			}
		case *sqlparser.ColName:
			result := lo.ContainsBy(config.allowedLeftColumns(), func(item Column) bool { return node.Name.EqualString(item.Name) })

			if result {
				return true, nil
			} else {
				return config.walkError("unsupported column name: %s", node.Name)
			}
		case sqlparser.IdentifierCI:
			result := lo.ContainsBy(config.allowedLeftColumns(), func(item Column) bool { return node.EqualString(item.Name) })

			if result {
				return true, nil
			} else {
				return config.walkError("unsupported column name: %s", node)
			}
		case sqlparser.IdentifierCS:
			result := lo.ContainsBy(config.allowedLeftColumns(), func(item Column) bool { return node.String() == item.Qualifier })

			if result {
				return true, nil
			} else {
				return config.walkError("unsupported table name: %s", node)
			}
		case sqlparser.TableName:
			result := lo.ContainsBy(config.allowedLeftColumns(), func(item Column) bool { return node.Name.String() == item.Qualifier })

			if result {
				return true, nil
			} else {
				return config.walkError("unsupported table name: %s", node)
			}
		case *sqlparser.BetweenExpr:
			switch lhs := node.Left.(type) {
			case *sqlparser.ColName:
				if columnConfig, found := config.findColumn(lhs); found {
					if columnConfig.BetweenOperator == nil {
						return config.walkError("unsupported operator: %s", node)
					}

					from, fromFound := lo.Find(columnConfig.BetweenOperator.Froms(), func(from From) bool {
						fromNodeType := fmt.Sprintf("%T", node.From)
						return fromNodeType == from.nodeType()
					})

					to, toFound := lo.Find(columnConfig.BetweenOperator.Tos(), func(to To) bool {
						toNodeType := fmt.Sprintf("%T", node.To)
						return toNodeType == to.nodeType()
					})

					if (fromFound && toFound) && (from.valid(node.From) && to.valid(node.To)) {
						return true, nil
					} else {
						return config.walkError("unsupported or invalid RHS: %s", node)
					}
				} else {
					return config.walkError("unsupported operator: %s", node)
				}
			}

			return config.walkError("unsupported between: %s", node)
		case *sqlparser.ComparisonExpr:
			switch lhs := node.Left.(type) {
			case *sqlparser.ColName:
				if columnConfig, found := config.findColumn(lhs); found {
					cop, copFound := lo.Find(columnConfig.ComparisonOperators, func(op IComparisonOperator) bool {
						return node.Operator.ToString() == op.ToString()
					})

					if copFound {
						right, rightFound := lo.Find(cop.Rights(), func(right Right) bool {
							rightNodeType := fmt.Sprintf("%T", node.Right)
							return rightNodeType == right.nodeType()
						})

						if rightFound && right.valid(node.Right) {
							return true, nil
						} else {
							return config.walkError("unsupported or invalid RHS: %s", node)
						}
					} else {
						return config.walkError("unsupported operator: %s", node)
					}
				}
			}

			return config.walkError("unsupported comparison: %s", node)
		default:
			{
				return config.walkError("unsupported syntax: %s", node)
			}
		}
	}

	err := sqlparser.Walk(fun, filter)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (config Config) findColumn(lhs *sqlparser.ColName) (Column, bool) {
	return lo.Find(config.allowedLeftColumns(), func(col Column) bool {
		return lhs.Qualifier.Name.String() == col.Qualifier && lhs.Name.EqualString(col.Name)
	})
}

func (config Config) allowedLeftColumns() []Column {
	return lo.FilterMap(config.Allow.Comparisons, func(item ILeft, index int) (Column, bool) {
		a, b := item.(Column)
		return a, b
	})
}

func (config Config) walkError(message string, node sqlparser.SQLNode) (bool, error) {
	if config.Debug {
		spew.Dump(node)
	}
	return false, fmt.Errorf(message, sqlparser.String(node))
}

func (config Config) validateGroupingParens(query string) error {
	noStringValues := noStringValues(query)

	leftParens := strings.Count(noStringValues, "(")
	rightParents := strings.Count(noStringValues, ")")

	max := config.Allow.GroupingParens

	if leftParens == rightParents && (max == UNLIMITED || leftParens <= max) {
		return nil
	} else {
		return fmt.Errorf("unsupported parens")
	}
}

func (config Config) validateAnds(query string) error {
	noStringValues := noStringValues(query)

	lowercaseQuery := strings.ToLower(noStringValues)
	andsCount := strings.Count(lowercaseQuery, " and ")

	max := config.Allow.Ands

	if max == UNLIMITED || andsCount <= max {
		return nil
	} else {
		return fmt.Errorf("unsupported and")
	}
}

func (config Config) validateOrs(query string) error {
	noStringValues := noStringValues(query)

	lowercaseQuery := strings.ToLower(noStringValues)
	orsCount := strings.Count(lowercaseQuery, " or ")

	max := config.Allow.Ors

	if max == UNLIMITED || orsCount <= max {
		return nil
	} else {
		return fmt.Errorf("unsupported or")
	}
}

func (config Config) validateNots(query string) error {
	noStringValues := noStringValues(query)

	lowercaseQuery := strings.ToLower(noStringValues)
	notsCount := strings.Count(lowercaseQuery, " not ")

	max := config.Allow.Nots

	if max == UNLIMITED || notsCount <= max {
		return nil
	} else {
		return fmt.Errorf("unsupported not")
	}
}

func noStringValues(query string) string {
	stringValsRegex := regexp.MustCompile(`(\"|\').*?(\"|\')`)
	return stringValsRegex.ReplaceAllString(query, "''")
}
