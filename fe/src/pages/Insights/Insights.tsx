import './Insights.css';

import Box from '@mui/material/Box';
import SearchInput from 'components/SearchInput/SearchInput';
import { BACKGROUND_COLOR } from 'constants/colors';
import React, { FC } from 'react';
import { useAppSelector } from 'store/hooks';

import InsightFilter from './InsightFilter';
import InsightTable from './InsightTable';
import SideMenu from './SideMenu';

const Insights: FC = () => {
	const { sideMenuVisible } = useAppSelector(state => state.resources);

	const handleChange = () => {
		// TODO
	};

	return (
		<>
			<Box
				sx={{
					backgroundColor: BACKGROUND_COLOR,
					display: 'none', // todo replace with flex
					justifyContent: 'end',
					paddingRight: '44px',
					flexDirection: 'row-reverse',
				}}
				p={2}>
				<SearchInput
					width={'400px'}
					height={'32px'} // todo if visible substract from page height below
					onChange={handleChange}
					rest={{ flexDirection: 'row-reverse' }}
				/>
			</Box>
			<Box sx={{ display: 'flex', backgroundColor: BACKGROUND_COLOR }}>
				<InsightFilter />
				<InsightTable />
				{sideMenuVisible ? <SideMenu /> : <></>}
			</Box>
		</>
	);
};

export default Insights;
