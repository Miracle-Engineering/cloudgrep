import Box from '@mui/material/Box';
import CircularProgress from '@mui/material/CircularProgress';
import Paper from '@mui/material/Paper';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Typography from '@mui/material/Typography';
import { BORDER_COLOR } from 'constants/colors';
import { DEBOUNCE_PERIOD, PAGE_LENGTH } from 'constants/globals';
import debounce from 'debounce';
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

const InsightTable: FC = () => {
	const { resources, count } = useAppSelector(state => state.resources);
	const { filterTags } = useAppSelector(state => state.tags);
	const { t } = useTranslation();
	const dispatch = useAppDispatch();
	const [isInfiniteScroll, setIsInfiniteScroll] = useState<boolean>(false);
	const [hasNext, setHasNext] = useState<boolean>(true);
	const { currentPage, next } = usePagination(PAGE_LENGTH, count);

	useEffect(() => {
		if (resources && isInfiniteScroll) {
			setIsInfiniteScroll(false);
			next();
		}
	}, [resources]);

	const handleInfiniteScroll = async (e: React.MouseEvent<HTMLInputElement>): Promise<void> => {
		if (isScrolledForInfiniteScroll(e)) {
			setIsInfiniteScroll(true);
			const response = await ResourceService.getFilteredResources(
				filterTags,
				currentPage * PAGE_LENGTH,
				PAGE_LENGTH
			);
			if (response?.data?.resources && response.data.resources.length > 0) {
				setHasNext(true);
				dispatch(
					getFilteredResourcesNextPage({
						resources: response.data.resources,
						limit: PAGE_LENGTH,
						offset: currentPage * PAGE_LENGTH,
					})
				);
			} else {
				setHasNext(false);
				setIsInfiniteScroll(false);
			}
		}
	};

	const onContainerScroll = async (e: React.MouseEvent<HTMLInputElement>): Promise<void> => {
		if (!isInfiniteScroll && hasNext) {
			await handleInfiniteScroll(e);
			e.persist();
		}
	};

	const handleClick = (resource: Resource) => {
		dispatch(setCurrentResource(resource));
		dispatch(toggleMenuVisible());
	};

	const debouncedContainerScroll = debounce(onContainerScroll, DEBOUNCE_PERIOD);

	return (
		<Box
			sx={{
				width: '80%',
				backgroundColor: '#F9F7F6',
				paddingLeft: '28px',
				paddingRight: '44px',
			}}>
			<Typography sx={{ ...tableStyles.bodyRow, display: 'flex', margin: '8px' }}>{`${count} ${t(
				'COUNT_RESOURCES'
			)}`}</Typography>
			<TableContainer
				sx={{ maxHeight: 'calc(100vH - 34px)' }}
				component={Paper}
				onScroll={async (e: React.MouseEvent<HTMLInputElement>): Promise<void> => {
					if (!isInfiniteScroll) {
						await debouncedContainerScroll(e);
					}
				}}>
				<Table stickyHeader sx={{ minWidth: 650, overflowY: 'scroll' }} size="small" aria-label="a dense table">
					<TableHead>
						<TableRow sx={{ height: '46px', borderBottom: `1px solid ${BORDER_COLOR}` }}>
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
									height: '66px',
									// '&:last-child td, &:last-child th': { border: 0 },
									'&:hover': tableStyles.hoverStyle,
								}}>
								<TableCell sx={tableStyles.bodyRow} component="th" scope="row">
									{row.type}
								</TableCell>
								<TableCell sx={tableStyles.bodyRow} align="left">
									{row.id}
								</TableCell>
								<TableCell sx={tableStyles.bodyRow} align="left">
									{row.region}
								</TableCell>
							</TableRow>
						))}
					</TableBody>
				</Table>
				{(!resources?.length || (isInfiniteScroll && hasNext)) && (
					<Box
						sx={{ display: 'flex', justifyContent: 'center', height: '100px', alignItems: 'center' }}
						mt={1}>
						{<CircularProgress />}
					</Box>
				)}
				{!isInfiniteScroll && !hasNext && (
					<Box sx={{ ...tableStyles.bodyRow, display: 'flex', justifyContent: 'center' }} my={1}>
						{t('NO_MORE_RESULTS')}
					</Box>
				)}
			</TableContainer>
		</Box>
	);
};

export default InsightTable;
