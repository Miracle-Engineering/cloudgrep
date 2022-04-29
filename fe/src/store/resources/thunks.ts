import { createAsyncThunk } from '@reduxjs/toolkit';
import ResourceService from 'services/ResourceService';

import { setResources } from './slice';

const getResources = createAsyncThunk('resources/getResources', async (_, thunkAPI) => {
	try {
		const response = await ResourceService.getResources();
		console.log(response);
		thunkAPI.dispatch(setResources(response.data));
		return response.data;
	} catch (error: any) {
		return thunkAPI.rejectWithValue({ status: error.response?.status, error: error.message });
	}
});

export { getResources };
