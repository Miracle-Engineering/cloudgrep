package config

func (c GetTagsAPI) Has() bool {
	return c.Call != ""
}
