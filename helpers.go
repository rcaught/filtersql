package filtersql

// Equals Operator

func EqualsOperatorStringValueAny() IComparisonOperator {
	return EqualsOperator{
		RightsAccessor: EqualsOperatorRights{
			coLiteralStringAny(),
		},
	}
}

func EqualsOperatorStringValue(fun func(string) bool) IComparisonOperator {
	return EqualsOperator{
		RightsAccessor: EqualsOperatorRights{
			coLiteralStringValidationFunction(fun),
		},
	}
}

func EqualsOperatorIntegerValueAny() IComparisonOperator {
	return EqualsOperator{
		RightsAccessor: EqualsOperatorRights{
			coLiteralIntegerAny(),
		},
	}
}

func EqualsOperatorIntegerValue(fun func(int) bool) IComparisonOperator {
	return EqualsOperator{
		RightsAccessor: EqualsOperatorRights{
			coLiteralIntegerValidationFunction(fun),
		},
	}
}

// Not Equals Operator

func NotEqualsOperatorStringValueAny() IComparisonOperator {
	return NotEqualsOperator{
		RightsAccessor: NotEqualsOperatorRights{
			coLiteralStringAny(),
		},
	}
}

func NotEqualsOperatorStringValue(fun func(string) bool) IComparisonOperator {
	return NotEqualsOperator{
		RightsAccessor: NotEqualsOperatorRights{
			coLiteralStringValidationFunction(fun),
		},
	}
}

func NotEqualsOperatorIntegerValueAny() IComparisonOperator {
	return NotEqualsOperator{
		RightsAccessor: NotEqualsOperatorRights{
			coLiteralIntegerAny(),
		},
	}
}

func NotEqualsOperatorIntegerValue(fun func(int) bool) IComparisonOperator {
	return NotEqualsOperator{
		RightsAccessor: NotEqualsOperatorRights{
			coLiteralIntegerValidationFunction(fun),
		},
	}
}

// In Operator

func InOperatorStringsValueAny() IComparisonOperator {
	return InOperator{
		RightsAccessor: InOperatorRights{
			coTupleStringsAny(),
		},
	}
}

func InOperatorStringsValue(fun func([]string) bool) IComparisonOperator {
	return InOperator{
		RightsAccessor: InOperatorRights{
			coTupleStringsValidationFunction(fun),
		},
	}
}

// Helpers

func coLiteralStringAny() LiteralValue {
	return LiteralValue{
		ValueType: StringValue{
			ValidationFunc: func(s string) bool { return true },
		},
	}
}

func coLiteralStringValidationFunction(fun func(string) bool) LiteralValue {
	return LiteralValue{
		ValueType: StringValue{
			ValidationFunc: fun,
		},
	}
}

func coLiteralIntegerAny() LiteralValue {
	return LiteralValue{
		ValueType: IntegerValue{
			ValidationFunc: func(s int) bool { return true },
		},
	}
}

func coLiteralIntegerValidationFunction(fun func(int) bool) LiteralValue {
	return LiteralValue{
		ValueType: IntegerValue{
			ValidationFunc: fun,
		},
	}
}

func coTupleStringsAny() TupleValue {
	return TupleValue{
		StringValues{
			ValidationFunc: func(vals []string) bool { return true },
		},
	}
}

func coTupleStringsValidationFunction(fun func([]string) bool) TupleValue {
	return TupleValue{
		StringValues{
			ValidationFunc: fun,
		},
	}
}

