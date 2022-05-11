import 'utils/localisation/index';

import CloseIcon from '@mui/icons-material/Close';
import Box from '@mui/material/Box';
import Tab from '@mui/material/Tab';
import Tabs from '@mui/material/Tabs';
import Typography from '@mui/material/Typography';
import { Property } from 'models/Resource';
import { Tag } from 'models/Tag';
import React, { useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { useAppDispatch, useAppSelector } from 'store/hooks';
import { toggleMenuVisible } from 'store/resources/slice';

import { sideMenuStyle } from './style';
import TabPanel from './TabPanel';

const SideMenu = () => {
	const { t } = useTranslation();
	const menuRef = useRef<HTMLElement>(null);
	const dispatch = useAppDispatch();
	const { currentResource, sideMenuVisible } = useAppSelector(state => state.resources);
	const [value, setValue] = React.useState(0);

	const handleChange = (event: React.SyntheticEvent, newValue: number) => {
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

	return (
		<>
			{currentResource ? (
				<Box ref={menuRef} sx={sideMenuStyle}>
					<Box sx={{ display: 'flex', justifyContent: 'end' }}>
						<CloseIcon onClick={handleClose} sx={{ margin: '20px', cursor: 'pointer' }} />
					</Box>
					<Box sx={{ display: 'flex', flexDirection: 'column' }}>
						<Box ml={2} sx={{ display: 'flex' }}>
							<Typography> {`${t('ID')} : ${currentResource.id}`} </Typography>
						</Box>
						<Box ml={2} sx={{ display: 'flex' }}>
							<Typography> {`${t('REGION')} : ${currentResource.region}`} </Typography>
						</Box>
					</Box>
					<Box sx={{ width: '100%' }}>
						<Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
							<Tabs value={value} onChange={handleChange} aria-label="basic tabs example">
								<Tab label={t('TAGS')} {...commonTabProps(0)} />
								<Tab label={t('PROPERTIES')} {...commonTabProps(1)} />
							</Tabs>
						</Box>
						<TabPanel value={value} index={0}>
							<Box>
								{currentResource.tags ? (
									currentResource.tags.map((item: Tag, index: number) => (
										<Box key={index} sx={{ display: 'flex' }}>
											<Typography> {`${item.key} : ${item.value}`} </Typography>
										</Box>
									))
								) : (
									<></>
								)}
							</Box>
						</TabPanel>
						<TabPanel value={value} index={1}>
							<Box>
								{currentResource.properties ? (
									currentResource.properties.map((item: Property, index: number) => (
										<Box key={index} sx={{ display: 'flex' }}>
											<Typography> {`${item.name} : ${item.value}`} </Typography>
										</Box>
									))
								) : (
									<></>
								)}
							</Box>
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
