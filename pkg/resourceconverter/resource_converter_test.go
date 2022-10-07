package resourceconverter

import (
	"context"
	"testing"

	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/testingutil"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

type TestEntry struct {
	ID        string
	Attr1     int
	Attr2     string
	Attr3     map[string]interface{}
	WeirdTags []WeirdTags
}

type WeirdTags struct {
	WeirdKey   string
	WeirdValue string
}

func TestReflectionConverter(t *testing.T) {
	ctx := context.Background()

	factory := func() model.Resource {
		return model.Resource{
			Type:   "DummyResource",
			Region: "dummyRegion",
		}
	}

	t.Run("SimpleConversion", func(t *testing.T) {
		entry := TestEntry{
			ID:        "id1",
			Attr1:     1,
			Attr2:     "hi",
			Attr3:     map[string]interface{}{"a": "b", "c": 2},
			WeirdTags: []WeirdTags{{WeirdKey: "key1", WeirdValue: "val1"}, {WeirdKey: "key2", WeirdValue: "val2"}},
		}
		rC := &ReflectionConverter{
			IdField:         "ID",
			DisplayIdField:  "Attr2",
			TagField:        TagField{Name: "WeirdTags", Key: "WeirdKey", Value: "WeirdValue"},
			ResourceFactory: factory,
		}
		resource, err := rC.ToResource(ctx, entry, nil)
		require.NoError(t, err)
		expectedResource := model.Resource{
			Region:    "dummyRegion",
			Id:        "id1",
			DisplayId: "hi",
			Type:      "DummyResource",
			Tags:      model.Tags{{Key: "key1", Value: "val1"}, {Key: "key2", Value: "val2"}},
			RawData:   datatypes.JSON([]byte(`{"ID":"id1","Attr1":1,"Attr2":"hi","Attr3":{"a":"b","c":2},"WeirdTags":[{"WeirdKey":"key1","WeirdValue":"val1"},{"WeirdKey":"key2","WeirdValue":"val2"}]}`)),
		}
		testingutil.AssertEqualsResource(t, expectedResource, resource)
	})

	t.Run("TagsPassedIn", func(t *testing.T) {
		entry := TestEntry{
			ID:        "id1",
			Attr1:     1,
			Attr2:     "hi",
			Attr3:     map[string]interface{}{"a": "b", "c": 2},
			WeirdTags: []WeirdTags{{WeirdKey: "key1", WeirdValue: "val1"}, {WeirdKey: "key2", WeirdValue: "val2"}},
		}
		rC := &ReflectionConverter{
			IdField:         "ID",
			ResourceFactory: factory,
		}
		resource, err := rC.ToResource(ctx, entry, model.Tags{{Key: "key1", Value: "val3"}, {Key: "key2", Value: "val4"}})
		require.NoError(t, err)
		expectedResource := model.Resource{
			Region:  "dummyRegion",
			Id:      "id1",
			Type:    "DummyResource",
			Tags:    model.Tags{{Key: "key1", Value: "val3"}, {Key: "key2", Value: "val4"}},
			RawData: datatypes.JSON([]byte(`{"ID":"id1","Attr1":1,"Attr2":"hi","Attr3":{"a":"b","c":2},"WeirdTags":[{"WeirdKey":"key1","WeirdValue":"val1"},{"WeirdKey":"key2","WeirdValue":"val2"}]}`)),
		}
		testingutil.AssertEqualsResource(t, expectedResource, resource)
	})

	t.Run("MissingTags", func(t *testing.T) {
		entry := TestEntry{
			ID:        "id1",
			Attr1:     1,
			Attr2:     "hi",
			Attr3:     map[string]interface{}{"a": "b", "c": 2},
			WeirdTags: []WeirdTags{{WeirdKey: "key1", WeirdValue: "val1"}, {WeirdKey: "key2", WeirdValue: "val2"}},
		}
		rC := &ReflectionConverter{
			IdField:         "ID",
			ResourceFactory: factory,
			TagField:        TagField{Name: "WeirdTags2", Key: "WeirdKey", Value: "WeirdValue"},
		}
		_, err := rC.ToResource(ctx, entry, nil)
		require.ErrorContains(t, err, "could not find tag field 'WeirdTags2' for type 'DummyResource'")
	})

	t.Run("MissingIdField", func(t *testing.T) {
		entry := TestEntry{
			ID:        "id1",
			Attr1:     1,
			Attr2:     "hi",
			Attr3:     map[string]interface{}{"a": "b", "c": 2},
			WeirdTags: []WeirdTags{{WeirdKey: "key1", WeirdValue: "val1"}, {WeirdKey: "key2", WeirdValue: "val2"}},
		}
		rC := &ReflectionConverter{
			IdField:         "ID2",
			ResourceFactory: factory,
		}
		_, err := rC.ToResource(ctx, entry, model.Tags{{Key: "key1", Value: "val3"}, {Key: "key2", Value: "val4"}})
		require.ErrorContains(t, err, "could not find id field 'ID2' for type 'DummyResource'")
	})

	t.Run("MissingDisplayIdField", func(t *testing.T) {
		entry := TestEntry{
			ID:        "id1",
			Attr1:     1,
			Attr2:     "hi",
			Attr3:     map[string]interface{}{"a": "b", "c": 2},
			WeirdTags: []WeirdTags{{WeirdKey: "key1", WeirdValue: "val1"}, {WeirdKey: "key2", WeirdValue: "val2"}},
		}
		rC := &ReflectionConverter{
			IdField:         "ID",
			DisplayIdField:  "DispID",
			ResourceFactory: factory,
		}
		_, err := rC.ToResource(ctx, entry, model.Tags{{Key: "key1", Value: "val3"}, {Key: "key2", Value: "val4"}})
		require.ErrorContains(t, err, "could not find display id field 'DispID' for type 'DummyResource'")
	})
}
