package filtersql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/samber/lo"
	"vitess.io/vitess/go/vt/sqlparser"
)

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

	err = config.validate(where)
	if err != nil {
		return "", err
	}

	result := sqlparser.String(where)
	result = strings.Replace(result, " where ", "", 1)
	return result, nil
}

func (config Config) validate(filter *sqlparser.Where) error {
	fun := func(node sqlparser.SQLNode) (bool, error) {
		switch node := node.(type) {
		case *sqlparser.Where:
			return true, nil
		case (*sqlparser.Literal):
			return true, nil
		case (sqlparser.ValTuple):
			return true, nil
		case *sqlparser.AndExpr:
			if config.Allow.Ands {
				return true, nil
			} else {
				return config.walkError("unsupported and: %s", node)
			}
		case *sqlparser.OrExpr:
			if config.Allow.Ors {
				return true, nil
			} else {
				return config.walkError("unsupported or: %s", node)
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
		case *sqlparser.ComparisonExpr:
			switch lhs := node.Left.(type) {
			case *sqlparser.ColName:
				if columnConfig, found := lo.Find(config.allowedLeftColumns(), func(col Column) bool {
					return lhs.Qualifier.Name.String() == col.Qualifier && lhs.Name.EqualString(col.Name)
				}); found {
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
