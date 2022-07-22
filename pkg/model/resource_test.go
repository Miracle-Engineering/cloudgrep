package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResourceFindById(t *testing.T) {
	r1 := Resource{
		Id: "i-123", Region: "us-east-1", Type: "test.Instance",
	}
	r2 := Resource{
		Id: "i-124", Region: "us-east-1", Type: "test.Instance",
	}
	resources := Resources{&r1, &r2}
	assert.Equal(t, "i-123", resources.FindById("i-123").Id)
	assert.Equal(t, "i-124", resources.FindById("i-124").Id)
	assert.Nil(t, resources.FindById("i-123-not-found"))
}

func TestResourceIds(t *testing.T) {
	r1 := Resource{
		Id: "i-123", Region: "us-east-1", Type: "test.Instance",
	}
	r2 := Resource{
		Id: "i-124",
	}
	resources := Resources{&r1, &r2}
	assert.ElementsMatch(t, []ResourceId{"i-123", "i-124"}, resources.Ids())
}

func TestClean(t *testing.T) {
	r1 := Resource{
		Id: "i-123", Region: "us-east-1", Type: "test.Instance", UpdatedAt: time.Now(),
	}
	rs := Resources{&r1}
	assert.False(t, r1.UpdatedAt.IsZero())
	assert.True(t, rs.Clean()[0].UpdatedAt.IsZero())
}

func TestResource_EffectiveDisplayId(t *testing.T) {
	r1 := Resource{
		Id: "foo", DisplayId: "bar",
	}
	r2 := Resource{
		Id: "foo",
	}

	assert.Equal(t, "bar", r1.EffectiveDisplayId())
	assert.Equal(t, "foo", r2.EffectiveDisplayId())
}
