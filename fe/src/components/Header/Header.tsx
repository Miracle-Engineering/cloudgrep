import 'utils/localisation/index';

import RefreshIcon from '@mui/icons-material/Refresh';
import Backdrop from '@mui/material/Backdrop';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Link from '@mui/material/Link';
import Snackbar from '@mui/material/Snackbar';
import Alert from 'components/Alert/Alert';
import { BORDER_COLOR, DARK_BLUE } from 'constants/colors';
import { AUTO_HIDE_DURATION, ENGINE_STATUS_INTERVAL, GITHUB, PAGE_LENGTH, PAGE_START, SLACK } from 'constants/globals';
import { EngineStatus, EngineStatusEnum } from 'models/EngineStatus';
import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import RefreshService from 'services/RefreshService';
import { useAppDispatch, useAppSelector } from 'store/hooks';
import { getFilteredResources } from 'store/resources/thunks';
import { getFields } from 'store/tags/thunks';

import { headerStyle, menuItems } from './style';

const Header = () => {
	const { t } = useTranslation();
	const dispatch = useAppDispatch();
	const { filterTags } = useAppSelector(state => state.tags);
	const [open, setOpen] = useState(false);
	const { resources } = useAppSelector(state => state.resources);
	const [errorMessage, setErrorMessage] = useState('');
	const [engineStatus, setEngineStatus] = useState<EngineStatus | undefined>();

	const handleStatus = async () => {
		const response = await RefreshService.getStatus();
		if (response.data.status === EngineStatusEnum.FAILED) {
			setErrorMessage(response.data.errorMessage);
		} else if (response.data.status === EngineStatusEnum.SUCCESS) {
			setEngineStatus(response.data);
		} else if (response.data.status === EngineStatusEnum.FETCHING) {
			setTimeout(handleStatus, ENGINE_STATUS_INTERVAL);
		}
	};

	const handleClick = async () => {
		setOpen(true);
		try {
			await RefreshService.refresh();
			await handleStatus();
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
		} catch (err: any) {
			setErrorMessage(err.error);
		}
	};

	const handleCloseBanner = (_event: React.SyntheticEvent | Event, reason?: string) => {
		if (reason === 'clickaway') {
			return;
		}

		setErrorMessage('');
		setEngineStatus(undefined);
	};

	useEffect(() => {
		if (resources) {
			setOpen(false);
		}
	}, [resources]);

	useEffect(() => {
		if (engineStatus?.status === EngineStatusEnum.SUCCESS) {
			dispatch(getFilteredResources({ data: filterTags, offset: PAGE_START, limit: PAGE_LENGTH }));
			setTimeout(() => setEngineStatus(undefined), AUTO_HIDE_DURATION);
		}
	}, [engineStatus, filterTags, dispatch]);

	return (
		<Box sx={headerStyle}>
			<Box sx={{ display: 'flex' }}>
				<Box>
					<img
						style={{ marginLeft: '24px', height: '28px', cursor: 'pointer' }}
						src={`${process.env.REACT_APP_PATH_PREFIX}/logo.png`}
					/>
				</Box>
				<Box sx={{ display: 'flex', marginLeft: '203.25px', alignItems: 'center' }}>
					<Link ml={4} sx={menuItems} href={SLACK} underline="none" target="_blank" rel="noopener">
						{t('SLACK')}
					</Link>
					<Link ml={4} sx={menuItems} href={GITHUB} underline="none" target="_blank" rel="noopener">
						{t('GITHUB')}
					</Link>
				</Box>
			</Box>
			<Box
				onClick={handleClick}
				sx={{
					height: '100%',
					display: 'flex',
					borderLeft: `1px solid ${BORDER_COLOR}`,
					width: '163px',
					justifyContent: 'center',
					cursor: 'pointer',
				}}>
				<Button sx={{ color: DARK_BLUE, textTransform: 'none' }} startIcon={<RefreshIcon />}>
					{t('REFRESH')}
				</Button>
			</Box>
			<Backdrop sx={{ color: '#fff', zIndex: theme => theme.zIndex.drawer + 1 }} open={open}>
				<CircularProgress color="inherit" />
			</Backdrop>
			<Snackbar
				sx={{ marginRight: '24px' }}
				open={!!errorMessage}
				autoHideDuration={AUTO_HIDE_DURATION}
				onClose={handleCloseBanner}>
				<Alert onClose={handleCloseBanner} severity="error" sx={{ width: '100%' }}>
					{errorMessage}
				</Alert>
			</Snackbar>
			<Snackbar
				sx={{ marginRight: '24px' }}
				open={!!engineStatus}
				autoHideDuration={AUTO_HIDE_DURATION}
				onClose={handleCloseBanner}>
				<Alert onClose={handleCloseBanner} severity="success" sx={{ width: '100%' }}>
					{t('REFRESH_SUCCESS')}
				</Alert>
			</Snackbar>
		</Box>
	);
};

export default Header;
