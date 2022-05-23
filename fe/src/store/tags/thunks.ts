import { createAsyncThunk } from '@reduxjs/toolkit';
import TagService from 'services/TagService';

import { setFields } from './slice';

const getFields = createAsyncThunk('tags/getFields', async (_, thunkAPI) => {
	try {
		const response = await TagService.getFields();
		thunkAPI.dispatch(setFields(response.data));
		return response.data;
	} catch (error: any) {
		return thunkAPI.rejectWithValue({ status: error.response?.status, error: error.message });
	}
});

export { getFields };
