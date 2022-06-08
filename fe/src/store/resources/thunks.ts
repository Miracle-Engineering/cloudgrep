import { createAsyncThunk } from '@reduxjs/toolkit';
import { Resource } from 'models/Resource';
import ResourceService from 'services/ResourceService';

import { setFilterTags } from '../tags/slice';
import { addResources, setResources } from './slice';
import { FilterResourcesApiParams } from './types';

const getResources = createAsyncThunk('resources/getResources', async (_, thunkAPI) => {
	try {
		const response = await ResourceService.getResources();
		thunkAPI.dispatch(setResources(response.data));
		return response.data;
	} catch (error: any) {
		return thunkAPI.rejectWithValue({ status: error.response?.status, error: error.message });
	}
});

const getFilteredResources = createAsyncThunk(
	'resources/getFilteredResources',
	async (apiParams: FilterResourcesApiParams, thunkAPI) => {
		const { data, limit, offset } = apiParams;
		try {
			const response = await ResourceService.getFilteredResources(data, offset, limit);
			thunkAPI.dispatch(setResources(response.data));
			thunkAPI.dispatch(setFilterTags(data));
			return response.data;
		} catch (error: any) {
			return thunkAPI.rejectWithValue({ status: error.response?.status, error: error.message });
		}
	}
);

const getFilteredResourcesNextPage = createAsyncThunk(
	'resources/getFilteredResourcesNextPage',
	async (resources: Resource[], thunkAPI) => {
		try {
			thunkAPI.dispatch(addResources(resources));
			return resources;
		} catch (error: any) {
			return thunkAPI.rejectWithValue({ status: error.response?.status, error: error.message });
		}
	}
);

export { getFilteredResources, getFilteredResourcesNextPage, getResources };
