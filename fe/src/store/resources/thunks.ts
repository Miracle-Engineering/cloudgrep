import { createAsyncThunk } from '@reduxjs/toolkit';
import { Resource } from 'models/Resource';
import ResourceService from 'services/ResourceService';

import { setFilterTags, setPaging } from '../tags/slice';
import { addResources, setResources } from './slice';
import { FilterResourcesApiParams, ResourcesNextPageParams } from './types';

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
	async (nextPageParams: ResourcesNextPageParams, thunkAPI) => {
		try {
			thunkAPI.dispatch(addResources(nextPageParams.resources));
			thunkAPI.dispatch(setPaging({ limit: nextPageParams.limit, offset: nextPageParams.offset }));
			return nextPageParams.resources;
		} catch (error: any) {
			return thunkAPI.rejectWithValue({ status: error.response?.status, error: error.message });
		}
	}
);

export { getFilteredResources, getFilteredResourcesNextPage, getResources };
