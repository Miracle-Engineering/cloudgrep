package testingutil

import (
	"testing"

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
