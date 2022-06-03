package testingutil

import (
	"testing"

	"gorm.io/datatypes"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/stretchr/testify/assert"
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
