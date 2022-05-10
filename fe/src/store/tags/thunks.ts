import { createAsyncThunk } from '@reduxjs/toolkit';
import TagService from 'services/TagService';

import { setTagResource } from './slice';

const getTags = createAsyncThunk('tags/getTags', async (_, thunkAPI) => {
	try {
		const response = await TagService.getTags();
		thunkAPI.dispatch(setTagResource(response.data));
		return response.data;
	} catch (error: any) {
		return thunkAPI.rejectWithValue({ status: error.response?.status, error: error.message });
	}
});

export { getTags };
