import SearchIcon from '@mui/icons-material/Search';
import IconButton from '@mui/material/IconButton';
import InputBase from '@mui/material/InputBase';
import Paper from '@mui/material/Paper';
import React, { FC } from 'react';
import { useTranslation } from 'react-i18next';

import { searchStyle, searchText } from './style';
import { SearchProps } from './types';

const SearchInput: FC<SearchProps> = props => {
	const { onChange, width, height, rest } = props;
	const { t } = useTranslation();

	return (
		<Paper
			component="form"
			sx={{
				...searchStyle,
				width: width,
				height: height || 24,
				maxWidth: width,
				...rest,
			}}>
			<IconButton type="submit" sx={{ p: '6px', paddingLeft: '16px' }} aria-label={t('SEARCH')}>
				<SearchIcon width="20px" height="20px" fontSize="small" />
			</IconButton>
			<InputBase
				sx={{
					...searchText,
					ml: 1,
					flex: 1,
				}}
				placeholder={t('SEARCH')}
				inputProps={{ 'aria-label': t('SEARCH_TERM') }}
				onChange={onChange}
			/>
		</Paper>
	);
};

export default SearchInput;
