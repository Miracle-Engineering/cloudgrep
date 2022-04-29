import Box from '@mui/material/Box';
import React, { FC, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from 'store/hooks';
import { getResources } from 'store/resources/thunks';
import { getTags } from 'store/tags/thunks';

import InsightFilter from './InsightFilter';
import InsightTable from './InsightTable';

const Insights: FC = () => {
	const dispatch = useAppDispatch();
	const { tags } = useAppSelector(state => state.tags);
	const { resources } = useAppSelector(state => state.resources);

	useEffect(() => {
		if (!tags?.length) {
			dispatch(getTags());
		}
	}, [tags?.length]);

	useEffect(() => {
		if (!resources?.length) {
			dispatch(getResources());
		}
	}, [resources?.length]);

	return (
		<Box sx={{ display: 'flex' }}>
			<InsightFilter />
			<InsightTable />
		</Box>
	);
};

export default Insights;
