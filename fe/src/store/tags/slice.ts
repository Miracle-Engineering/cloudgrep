import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Field } from 'models/Field';
import { MockTag } from 'models/Tag';
import { TagResource } from 'models/TagResource';

import { TagState } from './types';

const initialState: TagState = {
	tagResource: { Tags: [], Resources: [] },
	tags: [],
	fields: [],
};

const tagsSlice = createSlice({
	name: 'tags',
	initialState,
	reducers: {
		setTags: (state, action: PayloadAction<MockTag[]>) => {
			state.tags = action.payload;
		},
		setTagResource: (state, action: PayloadAction<TagResource>) => {
			state.tagResource = action.payload;
			state.tags = action.payload.Tags;
		},
		setFields: (state, action: PayloadAction<Field[]>) => {
			state.fields = action.payload;
		},
	},
});

export const { setFields, setTags, setTagResource } = tagsSlice.actions;

export default tagsSlice;
