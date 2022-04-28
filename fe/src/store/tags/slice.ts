import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Tag } from 'models/Tag';

import { TagState } from './types';

const initialState: TagState = {
	tags: [],
};

const tagsSlice = createSlice({
	name: 'tags',
	initialState,
	reducers: {
		setTags: (state, action: PayloadAction<Tag[]>) => {
			state.tags = action.payload;
		},
	},
});

export const { setTags } = tagsSlice.actions;

export default tagsSlice;
