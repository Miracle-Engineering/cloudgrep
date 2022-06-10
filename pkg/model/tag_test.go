package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTagsFind(t *testing.T) {
	t1 := Tag{
		Key:   "cluster",
		Value: "dev-cluster",
	}
	t2 := Tag{
		Key:   "env",
		Value: "dev",
	}
	tags := Tags{t1, t2}
	require.Equal(t, "cluster", tags.Find("cluster").Key)
	AssertEqualsTag(t, &t2, tags.Find("env"))
	require.Nil(t, tags.Find("none"))
}

func TestTagsDelete(t *testing.T) {
	t1 := Tag{
		Key:   "cluster",
		Value: "dev-cluster",
	}
	t2 := Tag{
		Key:   "env",
		Value: "dev",
	}
	tags := Tags{t1, t2}
	AssertEqualsTags(t, Tags{t1}, tags.Delete("env"))
	AssertEqualsTags(t, Tags{t1, t2}, tags.Delete("unknown"))
}
