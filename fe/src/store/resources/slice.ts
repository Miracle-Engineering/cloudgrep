import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Resource } from 'models/Resource';

import { ResourceState } from './types';

const initialState: ResourceState = {
	resources: [],
};

const resourcesSlice = createSlice({
	name: 'resources',
	initialState,
	reducers: {
		setResources: (state, action: PayloadAction<Resource[]>) => {
			state.resources = action.payload;
		},
	},
});

export const { setResources } = resourcesSlice.actions;

export default resourcesSlice;
