package filtersql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"vitess.io/vitess/go/vt/sqlparser"
)

type Config struct {
	Allow Allow
	Debug bool
}

type Allow struct {
	Comparisons Comparisons
	Ors         bool
	Ands        bool
}

type (
	Comparisons []ILeft
	ILeft       interface {
		iLeft()
	}
	Right interface {
		iRight()
		valid(any) bool
		nodeType() string
	}
	Rights []Right

	ComparisonOperators []IComparisonOperator
	IComparisonOperator interface {
		iComparisonOperator()
		ToString() string
		Rights() Rights
	}
)

type Column struct {
	Qualifier           string
	Name                string
	ComparisonOperators ComparisonOperators
}

func (Column) iLeft() {}

// EqualsOperator
type (
	IEqualsOperator interface {
		iEqualsOperator()
	}
	EqualsOperatorRights []IEqualsOperatorRight
	IEqualsOperatorRight interface {
		iEqualsOperatorRight()
		Right() Right
	}
	EqualsOperator struct {
		RightsAccessor EqualsOperatorRights
	}
)

func (EqualsOperator) ToString() string     { return "=" }
func (EqualsOperator) iComparisonOperator() {}
func (eo EqualsOperator) Rights() Rights {
	return lo.Map(eo.RightsAccessor, func(item IEqualsOperatorRight, index int) Right {
		return item.Right()
	})
}

// NotEqualsOperator
type (
	INotEqualsOperator interface {
		iNotEqualsOperator()
	}
	NotEqualsOperatorRights []INotEqualsOperatorRight
	INotEqualsOperatorRight interface {
		iNotEqualsOperatorRight()
		Right() Right
	}
	NotEqualsOperator struct {
		RightsAccessor NotEqualsOperatorRights
	}
)

func (NotEqualsOperator) ToString() string     { return "!=" }
func (NotEqualsOperator) iComparisonOperator() {}
func (eo NotEqualsOperator) Rights() Rights {
	return lo.Map(eo.RightsAccessor, func(item INotEqualsOperatorRight, index int) Right {
		return item.Right()
	})
}

// GreaterThanOperator
type (
	IGreaterThanOperator interface {
		iGreaterThanOperator()
	}
	GreaterThanOperatorRights []IGreaterThanOperatorRight
	IGreaterThanOperatorRight interface {
		iGreaterThanOperatorRight()
		Right() Right
	}
	GreaterThanOperator struct {
		RightsAccessor GreaterThanOperatorRights
	}
)

func (GreaterThanOperator) ToString() string     { return ">" }
func (GreaterThanOperator) iComparisonOperator() {}
func (eo GreaterThanOperator) Rights() Rights {
	return lo.Map(eo.RightsAccessor, func(item IGreaterThanOperatorRight, index int) Right {
		return item.Right()
	})
}

// LessThanOperator
type (
	ILessThanOperator interface {
		iLessThanOperator()
	}
	LessThanOperatorRights []ILessThanOperatorRight
	ILessThanOperatorRight interface {
		iLessThanOperatorRight()
		Right() Right
	}
	LessThanOperator struct {
		RightsAccessor LessThanOperatorRights
	}
)

func (LessThanOperator) ToString() string     { return "<" }
func (LessThanOperator) iComparisonOperator() {}
func (eo LessThanOperator) Rights() Rights {
	return lo.Map(eo.RightsAccessor, func(item ILessThanOperatorRight, index int) Right {
		return item.Right()
	})
}

// GreaterThanOrEqualOperator
type (
	IGreaterThanOrEqualOperator interface {
		iGreaterThanOrEqualOperator()
	}
	GreaterThanOrEqualOperatorRights []IGreaterThanOrEqualOperatorRight
	IGreaterThanOrEqualOperatorRight interface {
		iGreaterThanOrEqualOperatorRight()
		Right() Right
	}
	GreaterThanOrEqualOperator struct {
		RightsAccessor GreaterThanOrEqualOperatorRights
	}
)

func (GreaterThanOrEqualOperator) ToString() string     { return ">=" }
func (GreaterThanOrEqualOperator) iComparisonOperator() {}
func (eo GreaterThanOrEqualOperator) Rights() Rights {
	return lo.Map(eo.RightsAccessor, func(item IGreaterThanOrEqualOperatorRight, index int) Right {
		return item.Right()
	})
}

// LessThanOrEqualOperator
type (
	ILessThanOrEqualOperator interface {
		iLessThanOrEqualOperator()
	}
	LessThanOrEqualOperatorRights []ILessThanOrEqualOperatorRight
	ILessThanOrEqualOperatorRight interface {
		iLessThanOrEqualOperatorRight()
		Right() Right
	}
	LessThanOrEqualOperator struct {
		RightsAccessor LessThanOrEqualOperatorRights
	}
)

func (LessThanOrEqualOperator) ToString() string     { return "<=" }
func (LessThanOrEqualOperator) iComparisonOperator() {}
func (eo LessThanOrEqualOperator) Rights() Rights {
	return lo.Map(eo.RightsAccessor, func(item ILessThanOrEqualOperatorRight, index int) Right {
		return item.Right()
	})
}

// LiteralValueType
type (
	ILiteralValueType interface {
		iLiteralValueType()
		valid(any) bool
	}
	LiteralValue struct {
		ValueType ILiteralValueType
	}
)

func (LiteralValue) iEqualsOperatorRight()             {}
func (LiteralValue) iNotEqualsOperatorRight()          {}
func (LiteralValue) iGreaterThanOperatorRight()        {}
func (LiteralValue) iLessThanOperatorRight()           {}
func (LiteralValue) iGreaterThanOrEqualOperatorRight() {}
func (LiteralValue) iLessThanOrEqualOperatorRight()    {}
func (LiteralValue) iRight()                           {}
func (lv LiteralValue) Right() Right                   { return lv }
func (lv LiteralValue) valid(e any) bool               { return lv.ValueType.valid(e) }
func (LiteralValue) nodeType() string                  { return "*sqlparser.Literal" }

// StringValue
type StringValue struct {
	ValidationFunc func(string) bool
}

func (StringValue) iLiteralValueType() {}
func (StringValue) iRight()            {}
func (sv StringValue) Right() Right    { return sv }
func (sv StringValue) valid(s any) bool {
	parent, ok := s.(*sqlparser.Literal)
	return ok && parent.Type == sqlparser.StrVal && sv.ValidationFunc(parent.Val)
}
func (StringValue) nodeType() string { return LiteralValue{}.nodeType() }

// IntegerValue
type IntegerValue struct {
	ValidationFunc func(int) bool
}

func (IntegerValue) iLiteralValueType() {}
func (IntegerValue) iRight()            {}
func (iv IntegerValue) Right() Right    { return iv }
func (iv IntegerValue) valid(s any) bool {
	parent, ok := s.(*sqlparser.Literal)
	return ok && parent.Type == sqlparser.IntVal && iv.ValidationFunc(cast.ToInt(parent.Val))
}
func (IntegerValue) nodeType() string { return LiteralValue{}.nodeType() }

func (config Config) Parse(filter string) (string, error) {
	if strings.Trim(filter, " ") == "" {
		return "", nil
	}

	sqlBlob := fmt.Sprintf("SELECT * FROM `not_a_table` WHERE %s", filter)
	sql, _, err := sqlparser.SplitStatement(sqlBlob)
	if err != nil {
		return "", err
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
