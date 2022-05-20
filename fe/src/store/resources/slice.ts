import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Resource } from 'models/Resource';

import { ResourceState } from './types';

const initialState: ResourceState = {
	resources: [],
	currentResource: undefined,
	sideMenuVisible: false,
};

const resourcesSlice = createSlice({
	name: 'resources',
	initialState,
	reducers: {
		setResources: (state, action: PayloadAction<Resource[]>) => {
			state.resources = action.payload;
		},
		setCurrentResource: (state, action: PayloadAction<Resource>) => {
			state.currentResource = action.payload;
		},
		toggleMenuVisible: state => {
			state.sideMenuVisible = !state.sideMenuVisible;
		},
	},
});

export const { setCurrentResource, setResources, toggleMenuVisible } = resourcesSlice.actions;

export default resourcesSlice;
