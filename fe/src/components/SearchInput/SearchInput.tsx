import SearchIcon from '@mui/icons-material/Search';
import IconButton from '@mui/material/IconButton';
import InputBase from '@mui/material/InputBase';
import Paper from '@mui/material/Paper';
import { BACKGROUND_COLOR } from 'constants/colors';
import React, { FC } from 'react';
import { useTranslation } from 'react-i18next';

import { SearchProps } from './types';

const SearchInput: FC<SearchProps> = props => {
	const { onChange, width, height, rest } = props;
	const { t } = useTranslation();

	return (
		<Paper
			component="form"
			sx={{
				p: '2px 4px',
				display: 'flex',
				alignItems: 'center',
				width: width,
				height: height || 24,
				maxWidth: width,
				border: '1px solid #CECDCD',
				borderRadius: '0px',
				boxShadow: 'none',
				backgroundColor: BACKGROUND_COLOR,
				...rest,
			}}>
			<InputBase
				sx={{
					ml: 1,
					flex: 1,
					fontSize: 12,
				}}
				placeholder={t('SEARCH')}
				inputProps={{ 'aria-label': t('SEARCH_TERM') }}
				onChange={onChange}
			/>
			<IconButton type="submit" sx={{ p: '10px' }} aria-label={t('SEARCH')}>
				<SearchIcon fontSize="small" />
			</IconButton>
		</Paper>
	);
};

export default SearchInput;
