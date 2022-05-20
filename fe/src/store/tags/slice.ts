import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { MockTag } from 'models/Tag';
import { TagResource } from 'models/TagResource';

import { TagState } from './types';

const initialState: TagState = {
	tagResource: { Tags: [], Resources: [] },
	tags: [],
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
	},
});

export const { setTags, setTagResource } = tagsSlice.actions;

export default tagsSlice;
