package testingutil

import (
	"testing"

	"gorm.io/datatypes"

	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssertResourceCount_keyOnly(t *testing.T) {
	t.Run("good", func(t *testing.T) {
		fake := Fake(t)
		in := testingResources()

		AssertResourceCount(fake, in, "", 2)

		assert.False(t, fake.IsFail, "should not fail")
		assert.Zero(t, fake.Logs, "should not log")
	})

	t.Run("bad", func(t *testing.T) {
		fake := Fake(t)
		in := testingResources()

		AssertResourceCount(fake, in, "", 3)

		assert.True(t, fake.IsFail, "should fail")
		assert.Len(t, fake.Logs, 1)
		assert.Contains(t, fake.Logs[0], "expected 3 resource(s) with tag test=")
	})
}

func TestAssertResourceCount_keyValue(t *testing.T) {
	t.Run("good", func(t *testing.T) {
		fake := Fake(t)
		in := testingResources()

		AssertResourceCount(fake, in, "spam-1", 1)

		assert.False(t, fake.IsFail, "should not fail")
		assert.Zero(t, fake.Logs, "should not log")
	})

	t.Run("bad", func(t *testing.T) {
		fake := Fake(t)
		in := testingResources()

		AssertResourceCount(fake, in, "spam-1", 2)

		assert.True(t, fake.IsFail, "should fail")
		assert.Len(t, fake.Logs, 1)
		assert.Contains(t, fake.Logs[0], "expected 2 resource(s) with tag test=spam-1")
	})
}

func testingResources() []model.Resource {
	return []model.Resource{
		{
			Id: "foo",
			Tags: []model.Tag{
				{
					Key:   "test",
					Value: "spam-1",
				},
			},
		},
		{
			Id: "bar",
			Tags: []model.Tag{
				{
					Key:   "test",
					Value: "ham-1",
				},
			},
		},
		{
			Id:   "spam",
			Tags: []model.Tag{},
		},
	}
}

func TestAssertResource(t *testing.T) {
	r1 := model.Resource{
		Id: "i-123", Region: "us-east-1", Type: "test.Instance",
		Tags: []model.Tag{
			{Key: "enabled", Value: "true"},
			{Key: "eks:nodegroup", Value: "staging-default"},
		},
		RawData: datatypes.JSON([]byte(`{"name": "jinzhu", "age": 18, "tags": ["tag1", "tag2"], "orgs": {"orga": "orga"}}`)),
	}
	r2 := model.Resource{
		Id: "i-123", Region: "us-east-1", Type: "test.Instance",
		Tags: []model.Tag{
			{Key: "eks:nodegroup", Value: "staging-default"},
			{Key: "enabled", Value: "true"},
		},
		RawData: datatypes.JSON([]byte(`{"name": "jinzhu", "age": 18, "tags": ["tag1", "tag2"], "orgs": {"orga": "orga"}}`)),
	}
	//r1 and r2 should be equals even though the order of their tags/raw data are different
	AssertEqualsResource(t, r1, r2)
}

func TestAssertEqualTag(t *testing.T) {
	t1 := model.Tag{
		Key:   "cluster",
		Value: "dev-cluster",
	}
	t2 := model.Tag{
		Key:   "cluster",
		Value: "dev-cluster",
	}
	AssertEqualsTag(t, &t1, &t2)
}

func TestAssertResourceFilteredCount_success(t *testing.T) {
	f := ResourceFilter{
		AccountId: "foo",
	}

	in := []model.Resource{
		{
			AccountId: "foo",
			Id:        "spam",
		},
		{
			AccountId: "foo",
			Id:        "ham",
		},
		{
			AccountId: "bar",
			Id:        "a",
		},
	}

	tb := &FakeTB{}

	filtered := AssertResourceFilteredCount(tb, in, 2, f)
	assert.False(t, tb.IsFail)
	assert.ElementsMatch(t, in[0:2], filtered)
	assert.Empty(t, tb.Logs)
}

func TestAssertResourceFilteredCount_fail(t *testing.T) {
	f := ResourceFilter{
		AccountId: "foo",
		Region:    "us",
	}

	in := []model.Resource{
		{
			AccountId: "foo",
			Region:    "us",
			Id:        "spam",
		},
		{
			AccountId: "foo",
			Region:    "eu",
			Id:        "ham",
		},
		{
			AccountId: "bar",
			Region:    "us",
			Id:        "a",
		},
	}

	tb := &FakeTB{}

	filtered := AssertResourceFilteredCount(tb, in, 2, f)
	assert.True(t, tb.IsFail)
	assert.ElementsMatch(t, in[0:1], filtered)
	require.Len(t, tb.Logs, 2)
	assert.Contains(t, tb.Logs[0], "expected 2 resource(s) with filter ResourceFilter{AccountId=foo, Region=us}")
	assert.Equal(t, "filter ResourceFilter{AccountId=foo, Region=us} partial matches: AccountId=2, Region=2", tb.Logs[1])
}
