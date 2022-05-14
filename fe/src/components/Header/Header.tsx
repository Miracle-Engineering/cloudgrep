import 'utils/localisation/index';

import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import React from 'react';
import { useTranslation } from 'react-i18next';

import { headerStyle } from './style';

function Header() {
	const { t } = useTranslation();

	return (
		<Box sx={headerStyle}>
			<Typography ml={3} sx={{ color: '#2B3A67', textTransform: 'uppercase', cursor: 'pointer' }}>
				{t('APP_HEADER')}
			</Typography>
		</Box>
	);
}

export default Header;
