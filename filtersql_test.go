package filtersql_test

import (
	"testing"
	"time"

	fs "github.com/rcaught/filtersql"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func commonConfig() fs.Config {
	return fs.Config{
		Debug: true,
		Allow: fs.Allow{
			Ands:           true,
			Ors:            true,
			Nots:           true,
			GroupingParens: fs.UNLIMITED,
			Comparisons: fs.Comparisons{
				fs.Column{
					Qualifier: "",
					Name:      "a",
					ComparisonOperators: fs.ComparisonOperators{
						fs.EqualsOperator{
							fs.EqualsOperatorRights{
								fs.LiteralValue{
									fs.StringValue{
										ValidationFunc: func(val string) bool { return val == "test" },
									},
								},
							},
						},
						fs.NotEqualsOperator{
							fs.NotEqualsOperatorRights{
								fs.LiteralValue{
									fs.StringValue{
										ValidationFunc: func(val string) bool { return val == "test" },
									},
								},
							},
						},
						fs.InOperator{
							fs.InOperatorRights{
								fs.TupleValue{
									fs.StringValues{
										ValidationFunc: func(vals []string) bool { return lo.Every([]string{"test", "test2"}, vals) },
									},
								},
							},
						},
						fs.NotInOperator{
							fs.NotInOperatorRights{
								fs.TupleValue{
									fs.StringValues{
										ValidationFunc: func(vals []string) bool { return lo.Every([]string{"test", "test2"}, vals) },
									},
								},
							},
						},
					},
				},
				fs.Column{
					Qualifier: "",
					Name:      "b",
					ComparisonOperators: fs.ComparisonOperators{
						fs.EqualsOperator{
							fs.EqualsOperatorRights{
								fs.LiteralValue{
									fs.IntegerValue{
										ValidationFunc: func(val int) bool { return val == 2 },
									},
								},
							},
						},
						fs.GreaterThanOperator{
							fs.GreaterThanOperatorRights{
								fs.LiteralValue{
									fs.IntegerValue{
										ValidationFunc: func(val int) bool { return val == 2 },
									},
								},
							},
						},
						fs.LessThanOperator{
							fs.LessThanOperatorRights{
								fs.LiteralValue{
									fs.IntegerValue{
										ValidationFunc: func(val int) bool { return val == 2 },
									},
								},
							},
						},
						fs.GreaterThanOrEqualOperator{
							fs.GreaterThanOrEqualOperatorRights{
								fs.LiteralValue{
									fs.IntegerValue{
										ValidationFunc: func(val int) bool { return val == 2 },
									},
								},
							},
						},
						fs.LessThanOrEqualOperator{
							fs.LessThanOrEqualOperatorRights{
								fs.LiteralValue{
									fs.IntegerValue{
										ValidationFunc: func(val int) bool { return val == 2 },
									},
								},
							},
						},
						fs.InOperatorIntegersValue(func(i []int) bool { return lo.EveryBy(i, func(item int) bool { return item%2 == 0 }) }),
					},
				},
				fs.Column{
					Qualifier: "something",
					Name:      "d",
					ComparisonOperators: fs.ComparisonOperators{
						fs.EqualsOperator{
							fs.EqualsOperatorRights{
								fs.LiteralValue{
									fs.StringValue{
										ValidationFunc: func(val string) bool { return val == "test" },
									},
								},
							},
						},
					},
				},
				fs.Column{
					Qualifier: "",
					Name:      "e",
					ComparisonOperators: fs.ComparisonOperators{
						fs.EqualsOperator{
							fs.EqualsOperatorRights{
								fs.LiteralValue{
									fs.StringValue{
										ValidationFunc: func(val string) bool { return true },
									},
								},
							},
						},
					},
				},
				fs.Column{
					Qualifier: "",
					Name:      "t",
					BetweenOperator: fs.BetweenOperator{
						fs.BetweenOperatorFroms{
							fs.LiteralValue{
								fs.StringValue{
									ValidationFunc: func(val string) bool { return cast.ToTime(val) != time.Time{} },
								},
							},
						},
						fs.BetweenOperatorTos{
							fs.LiteralValue{
								fs.StringValue{
									ValidationFunc: func(val string) bool { return cast.ToTime(val) != time.Time{} },
								},
							},
						},
					},
				},
			},
		},
	}
}

func TestFilterSQLParseEmpty(t *testing.T) {
	config := commonConfig()

	query := ""
	parsedQuery, err := config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "", parsedQuery)
}

func TestFilterSQLParseMultipleQueries(t *testing.T) {
	config := commonConfig()

	query := "a = 'test'; select * from passwords;"
	parsedQuery, err := config.Parse(query)
	assert.EqualError(t, err, "unsupported syntax")
	assert.Equal(t, "", parsedQuery)

	query = "a = 'test' AND (SELECT * FROM passwords)"
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported syntax: (select * from passwords)")
	assert.Equal(t, "", parsedQuery)
}

func TestFilterSQLParseUnsupportedExpression(t *testing.T) {
	config := commonConfig()

	query := "GROUP BY a"
	parsedQuery, err := config.Parse(query)
	assert.EqualError(t, err, "unsupported syntax")
	assert.Equal(t, "", parsedQuery)
}

func TestFilterSQLParseUnsupportedQualifier(t *testing.T) {
	config := commonConfig()

	query := "passwords.a = 'test"
	parsedQuery, err := config.Parse(query)
	assert.EqualError(t, err, "unsupported syntax")
	assert.Equal(t, "", parsedQuery)
}

func TestFilterSQLParseAnd(t *testing.T) {
	config := commonConfig()

	query := "a = 'test' AND b = 2"
	parsedQuery, err := config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "a = 'test' and b = 2", parsedQuery)

	config.Allow.Ands = false
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported and: a = 'test' and b = 2")
	assert.Equal(t, "", parsedQuery)
}

func TestFilterSQLParseOr(t *testing.T) {
	config := commonConfig()

	query := "a = 'test' OR b = 2"
	parsedQuery, err := config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "a = 'test' or b = 2", parsedQuery)

	config.Allow.Ors = false
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported or: a = 'test' or b = 2")
	assert.Equal(t, "", parsedQuery)
}

func TestFilterSQLParseNot(t *testing.T) {
	config := commonConfig()

	query := "NOT (a = 'test' OR b = 2)"
	parsedQuery, err := config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "not (a = 'test' or b = 2)", parsedQuery)

	config.Allow.Nots = false
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported not: not (a = 'test' or b = 2)")
	assert.Equal(t, "", parsedQuery)
}

func TestFilterSQLParseComparisonsColumn(t *testing.T) {
	config := commonConfig()

	query := "a = 'test'"
	parsedQuery, err := config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "a = 'test'", parsedQuery)

	query = "something.d = 'test'"
	parsedQuery, err = config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "something.d = 'test'", parsedQuery)

	query = "c = 'test'"
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported comparison: c = 'test'")
	assert.Equal(t, "", parsedQuery)

	query = "'a' = 'a'"
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported comparison: 'a' = 'a'")
	assert.Equal(t, "", parsedQuery)

	query = "a = a"
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported or invalid RHS: a = a")
	assert.Equal(t, "", parsedQuery)
}

func TestFilterSQLParseComparisonsLiteralString(t *testing.T) {
	config := commonConfig()

	query := "a = 'test'"
	parsedQuery, err := config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "a = 'test'", parsedQuery)

	query = "a = 'fail'"
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported or invalid RHS: a = 'fail'")
	assert.Equal(t, "", parsedQuery)
}

func TestFilterSQLParseComparisonsLiteralInteger(t *testing.T) {
	config := commonConfig()

	query := "b = 2"
	parsedQuery, err := config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "b = 2", parsedQuery)

	query = "b = 3"
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported or invalid RHS: b = 3")
	assert.Equal(t, "", parsedQuery)
}

// Operators

func TestFilterSQLParseUnsupportedOperator(t *testing.T) {
	config := commonConfig()

	query := "a > 'test'"
	parsedQuery, err := config.Parse(query)
	assert.EqualError(t, err, "unsupported operator: a > 'test'")
	assert.Equal(t, "", parsedQuery)
}

func TestFilterSQLParseEqualsOperator(t *testing.T) {
	config := commonConfig()

	query := "a = 'test'"
	parsedQuery, err := config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "a = 'test'", parsedQuery)
}

func TestFilterSQLParseNotEqualsOperator(t *testing.T) {
	config := commonConfig()

	query := "a != 'test'"
	parsedQuery, err := config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "a != 'test'", parsedQuery)
}

func TestFilterSQLParseGreaterThanOperator(t *testing.T) {
	config := commonConfig()

	query := "b > 2"
	parsedQuery, err := config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "b > 2", parsedQuery)
}

func TestFilterSQLParseLessThanOperator(t *testing.T) {
	config := commonConfig()

	query := "b < 2"
	parsedQuery, err := config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "b < 2", parsedQuery)
}

func TestFilterSQLParseGreaterThanOrEqualOperator(t *testing.T) {
	config := commonConfig()

	query := "b >= 2"
	parsedQuery, err := config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "b >= 2", parsedQuery)
}

func TestFilterSQLParseLessThanOrEqualOperator(t *testing.T) {
	config := commonConfig()

	query := "b <= 2"
	parsedQuery, err := config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "b <= 2", parsedQuery)
}

func TestFilterSQLParseInOperatorStringValues(t *testing.T) {
	config := commonConfig()

	query := "a IN ('test', 'test2')"
	parsedQuery, err := config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "a in ('test', 'test2')", parsedQuery)

	query = "a IN ('test', 'test3')"
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported or invalid RHS: a in ('test', 'test3')")
	assert.Equal(t, "", parsedQuery)

	query = "a IN ('test', 2)"
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported or invalid RHS: a in ('test', 2)")
	assert.Equal(t, "", parsedQuery)
}

func TestFilterSQLParseInOperatorIntegerValues(t *testing.T) {
	config := commonConfig()

	query := "b IN (2, 4)"
	parsedQuery, err := config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "b in (2, 4)", parsedQuery)

	query = "b IN (2, 3)"
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported or invalid RHS: b in (2, 3)")
	assert.Equal(t, "", parsedQuery)

	query = "b IN ('test', 2)"
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported or invalid RHS: b in ('test', 2)")
	assert.Equal(t, "", parsedQuery)
}

func TestFilterSQLParseNotInOperatorStringValues(t *testing.T) {
	config := commonConfig()

	query := "a NOT IN ('test', 'test2')"
	parsedQuery, err := config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "a not in ('test', 'test2')", parsedQuery)

	query = "a NOT IN ('test', 'test3')"
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported or invalid RHS: a not in ('test', 'test3')")
	assert.Equal(t, "", parsedQuery)

	query = "a NOT IN ('test', 2)"
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported or invalid RHS: a not in ('test', 2)")
	assert.Equal(t, "", parsedQuery)
}

func TestFilterSQLParseNotInOperatorIntegerValues(t *testing.T) {
	config := commonConfig()

	query := "a NOT IN (2, 4)"
	parsedQuery, err := config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "a not in (2, 4)", parsedQuery)

	query = "a NOT IN ('test', 'test3')"
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported or invalid RHS: a not in ('test', 'test3')")
	assert.Equal(t, "", parsedQuery)

	query = "a NOT IN ('test', 2)"
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported or invalid RHS: a not in ('test', 2)")
	assert.Equal(t, "", parsedQuery)
}

func TestFilterSQLParseBetweenOperator(t *testing.T) {
	config := commonConfig()

	query := "t BETWEEN '2023-05-14 00:00:00.000000000' AND '2023-05-14 03:00:00.000000000'"
	parsedQuery, err := config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "t between '2023-05-14 00:00:00.000000000' and '2023-05-14 03:00:00.000000000'", parsedQuery)

	query = "a BETWEEN '2023-05-14 00:00:00.000000000' AND '2023-05-14 03:00:00.000000000'"
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported operator: a between '2023-05-14 00:00:00.000000000' and '2023-05-14 03:00:00.000000000'")
	assert.Equal(t, "", parsedQuery)

	query = "t BETWEEN 'cheese' AND '2023-05-14 03:00:00.000000000'"
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported or invalid RHS: t between 'cheese' and '2023-05-14 03:00:00.000000000'")
	assert.Equal(t, "", parsedQuery)
}

func TestFilterSQLParseGroupingParens(t *testing.T) {
	config := commonConfig()

	query := "a = 'test' AND (a = 'test' OR (b = 2 OR (a != 'test' AND (a = 'test' OR b = 2))))"
	config.Allow.GroupingParens = 0
	parsedQuery, err := config.Parse(query)
	assert.EqualError(t, err, "unsupported parens: a = 'test' AND (a = 'test' OR (b = 2 OR (a != 'test' AND (a = 'test' OR b = 2))))")
	assert.Equal(t, "", parsedQuery)

	config.Allow.GroupingParens = 20
	parsedQuery, err = config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "a = 'test' and (a = 'test' or (b = 2 or a != 'test' and (a = 'test' or b = 2)))", parsedQuery)

	config.Allow.GroupingParens = 0
	query = "e = '(test)' AND e = ')' AND e = '('"
	parsedQuery, err = config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "e = '(test)' and e = ')' and e = '('", parsedQuery)
}

