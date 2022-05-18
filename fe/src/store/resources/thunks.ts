import { createAsyncThunk } from '@reduxjs/toolkit';
import { MockTag, TagType } from 'models/Tag';
import ResourceService from 'services/ResourceService';

import { setResources } from './slice';

const getResources = createAsyncThunk('resources/getResources', async (_, thunkAPI) => {
	try {
		const response = await ResourceService.getResources();
		thunkAPI.dispatch(setResources(response.data));
		return response.data;
	} catch (error: any) {
		return thunkAPI.rejectWithValue({ status: error.response?.status, error: error.message });
	}
});

const getFilteredResources = createAsyncThunk('resources/getFilteredResources', async (data: MockTag, thunkAPI) => {
	try {
		const response = await ResourceService.getFilteredResources(data);
		thunkAPI.dispatch(setResources(response.data));
		return response.data;
	} catch (error: any) {
		return thunkAPI.rejectWithValue({ status: error.response?.status, error: error.message });
	}
});

export { getFilteredResources, getResources };
