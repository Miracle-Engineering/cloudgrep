import Box from '@mui/material/Box';
import CircularProgress from '@mui/material/CircularProgress';
import Paper from '@mui/material/Paper';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import { PAGE_LENGTH } from 'constants/globals';
import { Resource } from 'models/Resource';
import React, { FC, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import ResourceService from 'services/ResourceService';
import { useAppDispatch, useAppSelector } from 'store/hooks';
import { setCurrentResource, toggleMenuVisible } from 'store/resources/slice';
import { getFilteredResourcesNextPage } from 'store/resources/thunks';
import usePagination from 'utils/hooks/usePagination';
import { isScrolledForInfiniteScroll } from 'utils/uiHelper';

import { tableStyles } from './style';

const TOTAL_RECORDS = 10000; // todo update with real value from API response when available

const InsightTable: FC = () => {
	const { resources } = useAppSelector(state => state.resources);
	const { filterTags } = useAppSelector(state => state.tags);
	const { t } = useTranslation();
	const dispatch = useAppDispatch();
	const [isInfiniteScroll, setIsInfiniteScroll] = useState<boolean>(false);
	const [hasNext, setHasNext] = useState<boolean>(true);
	const { currentPage, next } = usePagination(PAGE_LENGTH, TOTAL_RECORDS);

	useEffect(() => {
		if (resources) {
			setIsInfiniteScroll(false);
		}
	}, [resources]);

	const handleClick = (resource: Resource) => {
		dispatch(setCurrentResource(resource));
		dispatch(toggleMenuVisible());
	};

	const handleInfiniteScroll = async (e: React.MouseEvent<HTMLInputElement>): Promise<void> => {
		if (isScrolledForInfiniteScroll(e) && hasNext) {
			setIsInfiniteScroll(true);
			next();
			const response = await ResourceService.getFilteredResources(
				filterTags,
				currentPage * PAGE_LENGTH,
				PAGE_LENGTH
			);
			if (response?.data && response.data.length > 0) {
				setHasNext(true);
				dispatch(getFilteredResourcesNextPage(response.data));
			} else {
				setHasNext(false);
			}
		}
	};

	return (
		<Box
			sx={{
				width: '80%',
				height: '100%',
				backgroundColor: '#F9F7F6',
				paddingLeft: '28px',
				paddingRight: '44px',
			}}>
			<TableContainer
				component={Paper}
				sx={{ height: '100%' }}
				onScroll={async (e: React.MouseEvent<HTMLInputElement>): Promise<void> => {
					if (!isInfiniteScroll) {
						await handleInfiniteScroll(e);
					}
				}}>
				<Table
					sx={{ minWidth: 650, maxHeight: '100%', overflowY: 'scroll' }}
					size="small"
					aria-label="a dense table">
					<TableHead>
						<TableRow>
							<TableCell sx={tableStyles.headerStyle}>{t('TYPE')} </TableCell>
							<TableCell align="left" sx={tableStyles.headerStyle}>
								{t('ID')}
							</TableCell>
							<TableCell align="left" sx={tableStyles.headerStyle}>
								{t('REGION')}
							</TableCell>
						</TableRow>
					</TableHead>
					<TableBody>
						{resources?.map((row: Resource, index: number) => (
							<TableRow
								onClick={() => handleClick(row)}
								key={row.id + row.type + index}
								sx={{
									height: '77px',
									'&:last-child td, &:last-child th': { border: 0 },
									'&:hover': tableStyles.hoverStyle,
								}}>
								<TableCell component="th" scope="row">
									{row.type}
								</TableCell>
								<TableCell align="left">{row.id}</TableCell>
								<TableCell align="left">{row.region}</TableCell>
							</TableRow>
						))}
					</TableBody>
				</Table>
				{isInfiniteScroll && hasNext && (
					<Box sx={{ display: 'flex', justifyContent: 'center' }} mt={1}>
						{<CircularProgress />}
					</Box>
				)}
			</TableContainer>
		</Box>
	);
};

export default InsightTable;
