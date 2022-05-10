import 'utils/localisation/index';

import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import React from 'react';
import { useTranslation } from 'react-i18next';

function Header() {
	const { t } = useTranslation();

	return (
		<Box sx={{ height: '64px', width: '100%', display: 'flex', alignItems: 'center', border: '1px solid #EAEBF0' }}>
			<Typography ml={3} sx={{ color: '#2B3A67', textTransform: 'uppercase', cursor: 'pointer' }}>
				{t('APP_HEADER')}
			</Typography>
		</Box>
	);
}

export default Header;
