package datastore

import (
	"context"
	"testing"

	"github.com/a8m/rql"
	"github.com/juandiegopalomino/cloudgrep/pkg/model"
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
	ri.fieldColumns.addExplicitFields("core", "type", "region", "id")
	assert.True(t,
		ri.fieldColumns.addDynamicFields("tags", "aws:ec2:fleet-id", "team-name", "cluster", "env"),
	)

	testCases := []testCase{
		{`{
  "limit":5,
  "filter":{
    "core.type":"ec2.Instance",
    "tags.aws:ec2:fleet-id":"fleet-bafee5d7-215d-addb-2632-290ab09da4e7"
  },
  "sort":[
    "-tags.aws:ec2:fleet-id"
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
    "core.type": "s3.Bucket"
  },
  "sort": ["core.region"]
}`, `{
  "filter":{
    "type": "s3.Bucket"
  },
  "sort": ["region"]
}`},
		{
			`{
  "filter":{
    "core.type":"ec2.Volume",
    "$or": [
      { "tags.team-name": "marketplace" },
      { "tags.team-name": "shipping" }
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
    "tags.unknown-field":"ec2.Volume"
  }
}`,
			`{
  "filter":{
    "tags.unknown-field":"ec2.Volume"
  }
}`,
		},
		//only filter and sort are updated
		{
			`{
  "filter2":{
    "tags.aws:ec2:fleet-id":"fleet-bafee5d7-215d-addb-2632-290ab09da4e7"
  }
}`,
			`{
  "filter2":{
    "tags.aws:ec2:fleet-id":"fleet-bafee5d7-215d-addb-2632-290ab09da4e7"
  }
}`,
		},
		//test multiple ORs
		{
			`{
  "filter":{
    "core.type":"ec2.Volume",
    "$or": [
      { "tags.team-name": "marketplace" },
      { "tags.team-name": "shipping" }
    ],
	"$and": [
		{ "$or": [
			{ "tags.cluster": "dev" },
			{ "tags.cluster": "prod" }
		] },
		{ "$or": [
			{ "tags.env": "staging" },
			{ "tags.env": "prod" }
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
			[]interface{}{model.FieldMissing, "us-east-1", "ec2.instance"},
			"col_17 is ? AND region = ? AND type = ?",
			[]interface{}{nil, "us-east-1", "ec2.instance"},
		},
		{
			"col_17 = ? AND region = ? AND type = ?",
			[]interface{}{model.FieldMissing, "us-east-1", model.FieldMissing},
			"col_17 is ? AND region = ? AND type is ?",
			[]interface{}{nil, "us-east-1", nil},
		},
		{
			"(col_17 = ? OR col_17 = ?) AND type = ?",
			[]interface{}{model.FieldMissing, model.FieldMissing, model.FieldMissing},
			"(col_17 is ? OR col_17 is ?) AND type is ?",
			[]interface{}{nil, nil, nil},
		},
		{
			"(col_17 = ? OR col_17 = ?) AND type = ?",
			[]interface{}{model.FieldMissing, model.FieldPresent, model.FieldMissing},
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
