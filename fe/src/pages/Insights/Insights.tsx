import './Insights.css';

import Box from '@mui/material/Box';
import SearchInput from 'components/SearchInput/SearchInput';
import { BACKGROUND_COLOR } from 'constants/colors';
import React, { FC, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from 'store/hooks';
import { getResources } from 'store/resources/thunks';
import { getFields } from 'store/tags/thunks';

import InsightFilter from './InsightFilter';
import InsightTable from './InsightTable';
import SideMenu from './SideMenu';

const Insights: FC = () => {
	const dispatch = useAppDispatch();
	const { fields } = useAppSelector(state => state.tags);
	const { resources, sideMenuVisible } = useAppSelector(state => state.resources);

	useEffect(() => {
		if (!fields?.length) {
			dispatch(getFields());
		}
	}, []);

	useEffect(() => {
		if (!resources?.length) {
			dispatch(getResources());
		}
	}, []);

	const handleChange = () => {
		// TODO
	};

	return (
		<>
			<Box
				sx={{
					backgroundColor: BACKGROUND_COLOR,
					display: 'flex',
					justifyContent: 'end',
					paddingRight: '44px',
				}}
				p={2}>
				<SearchInput width={'400px'} height={'32px'} onChange={handleChange} />
			</Box>
			<Box sx={{ display: 'flex', height: 'calc(100% - 136px)' }}>
				<InsightFilter />
				<InsightTable />
				{sideMenuVisible ? <SideMenu /> : <></>}
			</Box>
		</>
	);
};

export default Insights;
