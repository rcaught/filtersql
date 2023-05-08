package filtersql_test

import (
	"testing"

	fs "github.com/rcaught/filtersql"
	"github.com/stretchr/testify/assert"
)

func commonConfig() fs.Config {
	return fs.Config{
		Debug: true,
		Allow: fs.Allow{
			Ands: true,
			Ors:  true,
			Comparisons: fs.Comparisons{
				fs.Column{
					"",
					"a",
					fs.ComparisonOperators{
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
					},
				},
				fs.Column{
					"",
					"b",
					fs.ComparisonOperators{
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

func TestFilterSQLParseComparisonsColumn(t *testing.T) {
	config := commonConfig()

	query := "a = 'test'"
	parsedQuery, err := config.Parse(query)
	assert.NoError(t, err)
	assert.Equal(t, "a = 'test'", parsedQuery)

	query = "c = 'test'"
	parsedQuery, err = config.Parse(query)
	assert.EqualError(t, err, "unsupported comparison: c = 'test'")
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

