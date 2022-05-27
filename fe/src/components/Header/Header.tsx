import 'utils/localisation/index';

import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import { DARK_BLUE } from 'constants/colors';
import React from 'react';
import { useTranslation } from 'react-i18next';

import { headerStyle, menuItems } from './style';

function Header() {
	const { t } = useTranslation();

	return (
		<Box sx={headerStyle}>
			<Box>
				<img
					style={{ marginLeft: '24px', width: '160px', height: '40px', cursor: 'pointer' }}
					src="./logo.png"
				/>
			</Box>
			<Box sx={{ display: 'flex' }}>
				<Typography sx={{ ...menuItems, color: DARK_BLUE }}>{t('HOME')}</Typography>
				<Typography ml={4} sx={menuItems}>
					{t('SLACK')}
				</Typography>
				<Typography ml={4} sx={menuItems}>
					{t('GITHUB')}
				</Typography>
				<Typography ml={4} sx={menuItems}>
					{t('CONTACT')}
				</Typography>
			</Box>
			<Box>
				<Button sx={{ color: '#697391', borderColor: '#677290', marginRight: '44px' }} variant="outlined">
					{t('REFRESH')}
				</Button>
			</Box>
		</Box>
	);
}

export default Header;
