import 'utils/localisation/index';

import Backdrop from '@mui/material/Backdrop';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Snackbar from '@mui/material/Snackbar';
import Typography from '@mui/material/Typography';
import Alert from 'components/Alert/Alert';
import { DARK_BLUE } from 'constants/colors';
import { AUTO_HIDE_DURATION, ENGINE_STATUS_INTERVAL, PAGE_LENGTH, PAGE_START } from 'constants/globals';
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
			// throw { error: 'test error message for Demo purposes' };
			await RefreshService.refresh();
			await handleStatus();
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
		} catch (err: any) {
			setErrorMessage(err.error);
		}
		dispatch(getFields());
		dispatch(getFilteredResources({ data: filterTags, offset: PAGE_START, limit: PAGE_LENGTH }));
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
			dispatch(getFields());
			dispatch(getFilteredResources({ data: filterTags, offset: PAGE_START, limit: PAGE_LENGTH }));
			setEngineStatus(undefined);
		}
	}, [engineStatus, filterTags, dispatch]);

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
			<Backdrop sx={{ color: '#fff', zIndex: theme => theme.zIndex.drawer + 1 }} open={open}>
				<CircularProgress color="inherit" />
			</Backdrop>
			<Snackbar mr={2} open={!!errorMessage} autoHideDuration={AUTO_HIDE_DURATION} onClose={handleCloseBanner}>
				<Alert onClose={handleCloseBanner} severity="error" sx={{ width: '100%' }}>
					{errorMessage}
				</Alert>
			</Snackbar>
			<Snackbar mr={2} open={!!engineStatus} autoHideDuration={AUTO_HIDE_DURATION} onClose={handleCloseBanner}>
				<Alert onClose={handleCloseBanner} severity="success" sx={{ width: '100%' }}>
					{t('REFRESH_SUCCESS')}
				</Alert>
			</Snackbar>
		</Box>
	);
};

export default Header;
