package filtersql

import (
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"vitess.io/vitess/go/vt/sqlparser"
)

type Config struct {
	Allow Allow
	Debug bool
}

type Allow struct {
	Comparisons    Comparisons
	Ors            int
	Ands           int
	Nots           bool
	GroupingParens int
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

	From interface {
		iFrom()
		valid(any) bool
		nodeType() string
	}
	Froms []From

	To interface {
		iTo()
		valid(any) bool
		nodeType() string
	}
	Tos []To

	ComparisonOperators []IComparisonOperator
	IComparisonOperator interface {
		iComparisonOperator()
		ToString() string
		Rights() Rights
	}

	IBetweenOperator interface {
		iBetweenOperator()
		ToString() string
		Froms() Froms
		Tos() Tos
	}
)

type Column struct {
	Qualifier           string
	Name                string
	ComparisonOperators ComparisonOperators
	BetweenOperator     IBetweenOperator
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

// InOperator
type (
	IInOperator interface {
		iInOperator()
	}
	InOperatorRights []IInOperatorRight
	IInOperatorRight interface {
		iInOperatorRight()
		Right() Right
	}
	InOperator struct {
		RightsAccessor InOperatorRights
	}
)

func (InOperator) ToString() string     { return "in" }
func (InOperator) iComparisonOperator() {}
func (eo InOperator) Rights() Rights {
	return lo.Map(eo.RightsAccessor, func(item IInOperatorRight, index int) Right {
		return item.Right()
	})
}

// NotInOperator
type (
	INotInOperator interface {
		iNotInOperator()
	}
	NotInOperatorRights []INotInOperatorRight
	INotInOperatorRight interface {
		iNotInOperatorRight()
		Right() Right
	}
	NotInOperator struct {
		RightsAccessor NotInOperatorRights
	}
)

func (NotInOperator) ToString() string     { return "not in" }
func (NotInOperator) iComparisonOperator() {}
func (eo NotInOperator) Rights() Rights {
	return lo.Map(eo.RightsAccessor, func(item INotInOperatorRight, index int) Right {
		return item.Right()
	})
}

// BetweenOperator
type (
	BetweenOperatorFroms []IBetweenOperatorFrom
	BetweenOperatorTos   []IBetweenOperatorTo
	IBetweenOperatorFrom interface {
		iBetweenOperatorFrom()
		From() From
	}
	IBetweenOperatorTo interface {
		iBetweenOperatorTo()
		To() To
	}
	BetweenOperator struct {
		FromsAccessor BetweenOperatorFroms
		TosAccessor   BetweenOperatorTos
	}
)

func (BetweenOperator) ToString() string  { return "between" }
func (BetweenOperator) iBetweenOperator() {}
func (bo BetweenOperator) Froms() Froms {
	return lo.Map(bo.FromsAccessor, func(item IBetweenOperatorFrom, index int) From {
		return item.From()
	})
}
func (bo BetweenOperator) Tos() Tos {
	return lo.Map(bo.TosAccessor, func(item IBetweenOperatorTo, index int) To {
		return item.To()
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
func (LiteralValue) iBetweenOperatorFrom()             {}
func (LiteralValue) iBetweenOperatorTo()               {}
func (LiteralValue) iRight()                           {}
func (lv LiteralValue) Right() Right                   { return lv }
func (LiteralValue) iFrom()                            {}
func (lv LiteralValue) From() From                     { return lv }
func (LiteralValue) iTo()                              {}
func (lv LiteralValue) To() To                         { return lv }
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

// StringValues
type StringValues struct {
	ValidationFunc func([]string) bool
}

func (StringValues) iTupleValueType() {}
func (StringValues) iRight()          {}
func (sv StringValues) Right() Right  { return sv }
func (sv StringValues) valid(s any) bool {
	parent, ok := s.(sqlparser.ValTuple)

	if ok {
		values := lo.FilterMap(parent, func(item sqlparser.Expr, index int) (string, bool) {
			value, ok := item.(*sqlparser.Literal)

			if ok && value.Type == sqlparser.StrVal {
				return value.Val, true
			} else {
				return "", false
			}
		})

		if len(values) == len(parent) && sv.ValidationFunc(values) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
func (StringValues) nodeType() string { return TupleValue{}.nodeType() }

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

// IntegerValues
type IntegerValues struct {
	ValidationFunc func([]int) bool
}

func (IntegerValues) iTupleValueType() {}
func (IntegerValues) iRight()          {}
func (iv IntegerValues) Right() Right  { return iv }
func (iv IntegerValues) valid(s any) bool {
	parent, ok := s.(sqlparser.ValTuple)

	if ok {
		values := lo.FilterMap(parent, func(item sqlparser.Expr, index int) (int, bool) {
			value, ok := item.(*sqlparser.Literal)

			if ok && value.Type == sqlparser.IntVal {
				return cast.ToInt(value.Val), true
			} else {
				return 0, false
			}
		})

		if len(values) == len(parent) && iv.ValidationFunc(values) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
func (IntegerValues) nodeType() string { return TupleValue{}.nodeType() }

// TupleValueType
type (
	ITupleValueType interface {
		iTupleValueType()
		valid(any) bool
	}
	TupleValue struct {
		ValueType ITupleValueType
	}
)

func (TupleValue) iInOperatorRight()    {}
func (TupleValue) iNotInOperatorRight() {}
func (TupleValue) iRight()              {}
func (tv TupleValue) Right() Right      { return tv }
func (tv TupleValue) valid(e any) bool  { return tv.ValueType.valid(e) }
func (TupleValue) nodeType() string     { return "sqlparser.ValTuple" }
