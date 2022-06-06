package datastore

import (
	"context"
	"testing"

	"github.com/a8m/rql"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

type testCase struct {
	input, output interface{}
}

func TestUpdateQueryFields(t *testing.T) {

	ri, err := newResourceIndexer(context.Background(), zaptest.NewLogger(t), nil)
	//for this test we don't have a DB, we can still proceed
	assert.Error(t, err, "no DB provided")
	ri.fieldColumns.addExplicitFields("type", "region", "id")
	assert.True(t,
		ri.fieldColumns.addDynamicFields("aws:ec2:fleet-id", "team-name", "cluster", "env"),
	)

	testCases := []testCase{
		{`{
  "limit":5,
  "filter":{
    "type":"ec2.Instance",
    "aws:ec2:fleet-id":"fleet-bafee5d7-215d-addb-2632-290ab09da4e7"
  },
  "sort":[
    "-aws:ec2:fleet-id"
  ]
}`, `{
  "limit":5,
  "filter":{
    "type":"ec2.Instance",
    "col_1":"fleet-bafee5d7-215d-addb-2632-290ab09da4e7"
  },
  "sort":[
    "-col_1"
  ]
}`},
		{`{
  "filter":{
    "type": "s3.Bucket"
  },
  "sort": ["region"]
}`, `{
  "filter":{
    "type": "s3.Bucket"
  },
  "sort": ["region"]
}`},
		{
			`{
  "filter":{
    "type":"ec2.Volume",
    "$or": [
      { "team-name": "marketplace" },
      { "team-name": "shipping" }
    ]
  }
}`,
			`{
  "filter":{
    "type":"ec2.Volume",
    "$or": [
      { "col_2": "marketplace" },
      { "col_2": "shipping" }
    ]
  }
}`,
		},
		{
			`{
  "filter":{
    "unknown-field":"ec2.Volume"
  }
}`,
			`{
  "filter":{
    "unknown-field":"ec2.Volume"
  }
}`,
		},
		//only filter and sort are updated
		{
			`{
  "filter2":{
    "aws:ec2:fleet-id":"fleet-bafee5d7-215d-addb-2632-290ab09da4e7"
  }
}`,
			`{
  "filter2":{
    "aws:ec2:fleet-id":"fleet-bafee5d7-215d-addb-2632-290ab09da4e7"
  }
}`,
		},
		//test multiple ORs
		{
			`{
  "filter":{
    "type":"ec2.Volume",
    "$or": [
      { "team-name": "marketplace" },
      { "team-name": "shipping" }
    ],
	"$and": [
		{ "$or": [
			{ "cluster": "dev" },
			{ "cluster": "prod" }
		] },
		{ "$or": [
			{ "env": "staging" },
			{ "env": "prod" }
		] }
	]
  }
}`,
			`{
  "filter":{
    "type":"ec2.Volume",
    "$or": [
      { "col_2": "marketplace" },
      { "col_2": "shipping" }
    ],
	"$and": [
		{ "$or": [
			{ "col_3": "dev" },
			{ "col_3": "prod" }
		] },
		{ "$or": [
			{ "col_4": "staging" },
			{ "col_4": "prod" }
		] }
	]
  }
}`,
		},
	}

	for _, tc := range testCases {
		actual, err := ri.updateQueryFields([]byte(tc.input.(string)))
		assert.NoError(t, err)
		assert.JSONEq(t, tc.output.(string), string(actual))
	}
}
func TestReplaceNullValues(t *testing.T) {
	testCases := []struct {
		InFilterExp   string
		InFilterArgs  []interface{}
		OutFilterExp  string
		OutFilterArgs []interface{}
	}{
		{
			"col_17 = ? AND region = ? AND type = ?",
			[]interface{}{model.NullValue, "us-east-1", "ec2.instance"},
			"col_17 is ? AND region = ? AND type = ?",
			[]interface{}{nil, "us-east-1", "ec2.instance"},
		},
		{
			"col_17 = ? AND region = ? AND type = ?",
			[]interface{}{model.NullValue, "us-east-1", model.NullValue},
			"col_17 is ? AND region = ? AND type is ?",
			[]interface{}{nil, "us-east-1", nil},
		},
		{
			"(col_17 = ? OR col_17 = ?) AND type = ?",
			[]interface{}{model.NullValue, model.NullValue, model.NullValue},
			"(col_17 is ? OR col_17 is ?) AND type is ?",
			[]interface{}{nil, nil, nil},
		},
		{
			"(col_17 = ? OR col_17 = ?) AND type = ?",
			[]interface{}{model.NullValue, model.NotNullValue, model.NullValue},
			"(col_17 is ? OR col_17 is not ?) AND type is ?",
			[]interface{}{nil, nil, nil},
		},
	}
	for _, tc := range testCases {
		inParams := rql.Params{
			FilterExp:  tc.InFilterExp,
			FilterArgs: tc.InFilterArgs,
		}
		outParams := replaceNullValues(&inParams)
		assert.Equal(t, tc.OutFilterExp, outParams.FilterExp)
		assert.EqualValues(t, tc.OutFilterArgs, outParams.FilterArgs)
	}
}
