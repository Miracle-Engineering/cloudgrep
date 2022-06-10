import 'utils/localisation/index';

import Backdrop from '@mui/material/Backdrop';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Typography from '@mui/material/Typography';
import { DARK_BLUE } from 'constants/colors';
import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useAppDispatch, useAppSelector } from 'store/hooks';
import { getFilteredResources } from 'store/resources/thunks';
import { getFields } from 'store/tags/thunks';

import { headerStyle, menuItems } from './style';

function Header() {
	const { t } = useTranslation();
	const dispatch = useAppDispatch();
	const { filterTags, limit, offset } = useAppSelector(state => state.tags);
	const [open, setOpen] = useState(false);
	const { resources } = useAppSelector(state => state.resources);

	const handleClick = () => {
		setOpen(true);
		dispatch(getFields());
		dispatch(getFilteredResources({ data: filterTags, offset: offset, limit: limit }));
	};

	const handleClose = () => {
		setOpen(false);
	};

	useEffect(() => {
		if (resources) {
			setOpen(false);
		}
	}, [resources]);

	return (
		<Box sx={headerStyle}>
			<Box>
				<img
					style={{ marginLeft: '24px', height: '40px', cursor: 'pointer' }}
					src={`${process.env.REACT_APP_PATH_PREFIX}/logo.png`}
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
				<Button
					onClick={handleClick}
					sx={{ color: '#697391', borderColor: '#677290', marginRight: '44px' }}
					variant="outlined">
					{t('REFRESH')}
				</Button>
			</Box>
			<Backdrop
				sx={{ color: '#fff', zIndex: theme => theme.zIndex.drawer + 1 }}
				open={open}
				onClick={handleClose}>
				<CircularProgress color="inherit" />
			</Backdrop>
		</Box>
	);
}

export default Header;
