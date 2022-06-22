import 'utils/localisation/index';

import CloseIcon from '@mui/icons-material/Close';
import DataObjectIcon from '@mui/icons-material/DataObject';
import Box from '@mui/material/Box';
import Tab from '@mui/material/Tab';
import Tabs from '@mui/material/Tabs';
import Typography from '@mui/material/Typography';
import { BACKGROUND_COLOR } from 'constants/colors';
import { Tag } from 'models/Tag';
import React, { useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { useAppDispatch, useAppSelector } from 'store/hooks';
import { toggleMenuVisible } from 'store/resources/slice';

import PropertyItem from './PropertyItem/';
import {
	boxFirstLine,
	boxSecondLine,
	boxStyle,
	sideMenuHeader,
	sideMenuLeftText,
	sideMenuRightText,
	sideMenuStyle,
	tabStyle,
} from './style';
import TabPanel from './TabPanel';

const SideMenu = () => {
	const { t } = useTranslation();
	const menuRef = useRef<HTMLElement>(null);
	const dispatch = useAppDispatch();
	const { currentResource, sideMenuVisible } = useAppSelector(state => state.resources);
	const [value, setValue] = React.useState(0);

	const handleChange = (_event: React.SyntheticEvent, newValue: number) => {
		setValue(newValue);
	};

	const handleEvent = (e: MouseEvent) => {
		if (menuRef.current && !menuRef.current.contains(e.target as HTMLElement) && sideMenuVisible) {
			dispatch(toggleMenuVisible());
		}
	};

	const handleClose = () => {
		dispatch(toggleMenuVisible());
	};

	useEffect(() => {
		window.addEventListener('mouseup', handleEvent);
		return () => {
			window.removeEventListener('mouseup', handleEvent);
		};
	});

	const commonTabProps = (index: number) => {
		return {
			id: `simple-tab-${index}`,
			'aria-controls': `simple-tabpanel-${index}`,
		};
	};

	const renderObjects = (data: Object): React.ReactNode =>
		Object.entries(data).map(([key, value]) => <PropertyItem key={`${key}${value}`} keyItem={key} value={value} />);

	return (
		<>
			{currentResource ? (
				<Box ref={menuRef} sx={sideMenuStyle}>
					<Box
						sx={{
							display: 'flex',
							justifyContent: 'space-between',
							height: '72px',
							backgroundColor: BACKGROUND_COLOR,
						}}
						mb={3}>
						<Box ml={'16px'} sx={{ display: 'flex', alignItems: 'center' }}>
							<DataObjectIcon sx={{ margin: '10px', color: '#AAB0C2' }} />
							<Typography sx={sideMenuHeader}>{t('DETAILS')}</Typography>
						</Box>
						<CloseIcon onClick={handleClose} sx={{ margin: '20px', cursor: 'pointer', color: '#AAB0C2' }} />
					</Box>
					<Box sx={{ display: 'flex', flexDirection: 'row' }}>
						<Box ml={2} p={'12px'} sx={boxStyle}>
							<Typography {...boxFirstLine}>{`${t('ID')}`.toUpperCase()}</Typography>
							<Typography {...boxSecondLine}> {`${currentResource.id}`} </Typography>
						</Box>
						<Box ml={2} p={1} sx={boxStyle}>
							<Typography {...boxFirstLine}> {`${t('REGION')}`} </Typography>
							<Typography {...boxSecondLine}> {`${currentResource.region}`} </Typography>
						</Box>
						<Box ml={2} p={1} sx={boxStyle}>
							<Typography {...boxFirstLine}> {`${t('TYPE')}`} </Typography>
							<Typography {...boxSecondLine}> {`${currentResource.type}`} </Typography>
						</Box>
					</Box>
					<Box sx={{ width: '100%' }} mt={2}>
						<Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
							<Tabs
								sx={{
									color: '#959DB3 !important',
									'&.Mui-selected': {
										color: '#2B3A67 !important',
									},
									textTransform: 'none',
								}}
								value={value}
								onChange={handleChange}
								aria-label="basic tabs example">
								<Tab
									sx={{
										...tabStyle,
										'&.Mui-selected': {
											color: '#2B3A67 !important',
										},
									}}
									label={t('TAGS')}
									{...commonTabProps(0)}
								/>
								<Tab
									sx={{
										...tabStyle,
										'&.Mui-selected': {
											color: '#2B3A67 !important',
										},
									}}
									label={t('PROPERTIES')}
									{...commonTabProps(1)}
								/>
							</Tabs>
						</Box>
						<TabPanel value={value} index={0}>
							<Box>
								{currentResource.tags ? (
									currentResource.tags.map((item: Tag, index: number) => (
										<Box key={index} sx={{ display: 'flex' }}>
											<Typography mr={2} {...sideMenuLeftText}>{`${item.key} `}</Typography>
											<Typography {...sideMenuRightText}>{` ${item.value}`}</Typography>
										</Box>
									))
								) : (
									<></>
								)}
							</Box>
						</TabPanel>
						<TabPanel value={value} index={1}>
							<Box>{currentResource.rawData ? renderObjects(currentResource.rawData) : <></>}</Box>
						</TabPanel>
					</Box>
				</Box>
			) : (
				<></>
			)}
		</>
	);
};

export default SideMenu;
