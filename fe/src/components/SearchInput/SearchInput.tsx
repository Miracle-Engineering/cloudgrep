import SearchIcon from '@mui/icons-material/Search';
import IconButton from '@mui/material/IconButton';
import InputBase from '@mui/material/InputBase';
import Paper from '@mui/material/Paper';
import React, { FC } from 'react';
import { useTranslation } from 'react-i18next';

import { SearchProps } from './types';

const SearchInput: FC<SearchProps> = props => {
	const { onChange } = props;
	const { t } = useTranslation();

	return (
		<Paper component="form" sx={{ p: '2px 4px', display: 'flex', alignItems: 'center', width: '100%', height: 24 }}>
			<InputBase
				sx={{ ml: 1, flex: 1, fontSize: 12, boxShadow: 'none', borderColor: 'rgba(28, 43, 52, 0.98)' }}
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
