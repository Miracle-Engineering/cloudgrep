import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { FieldGroup } from 'models/Field';
import { Tag } from 'models/Tag';
import { TagResource } from 'models/TagResource';

import { TagState } from './types';

const initialState: TagState = {
	tagResource: { Tags: [], Resources: [] },
	fields: [],
	filterTags: [],
	limit: 0,
	offset: 0,
};

const tagsSlice = createSlice({
	name: 'tags',
	initialState,
	reducers: {
		setFilterTags: (state, action: PayloadAction<Tag[]>) => {
			state.filterTags = action.payload;
		},
		setTagResource: (state, action: PayloadAction<TagResource>) => {
			state.tagResource = action.payload;
		},
		setFields: (state, action: PayloadAction<FieldGroup[]>) => {
			state.fields = action.payload;
		},
	},
});

export const { setFields, setFilterTags, setTagResource } = tagsSlice.actions;

export default tagsSlice;
