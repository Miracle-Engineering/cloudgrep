package model_test

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/testingutil"
	"github.com/stretchr/testify/require"
)

func TestTagsFind(t *testing.T) {
	t1 := model.Tag{
		Key:   "cluster",
		Value: "dev-cluster",
	}
	t2 := model.Tag{
		Key:   "env",
		Value: "dev",
	}
	tags := model.Tags{t1, t2}
	require.Equal(t, "cluster", tags.Find("cluster").Key)
	testingutil.AssertEqualsTag(t, &t2, tags.Find("env"))
	require.Nil(t, tags.Find("none"))
}

func TestTagsDelete(t *testing.T) {
	t1 := model.Tag{
		Key:   "cluster",
		Value: "dev-cluster",
	}
	t2 := model.Tag{
		Key:   "env",
		Value: "dev",
	}
	tags := model.Tags{t1, t2}
	testingutil.AssertEqualsTags(t, model.Tags{t1}, tags.Delete("env"))
	testingutil.AssertEqualsTags(t, model.Tags{t1, t2}, tags.Delete("unknown"))
}
