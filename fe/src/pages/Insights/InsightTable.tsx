import Box from '@mui/material/Box';
import CircularProgress from '@mui/material/CircularProgress';
import Paper from '@mui/material/Paper';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import TableSortLabel from '@mui/material/TableSortLabel';
import Typography from '@mui/material/Typography';
import { visuallyHidden } from '@mui/utils';
import { BORDER_COLOR } from 'constants/colors';
import { DEBOUNCE_PERIOD, PAGE_LENGTH, PAGE_START } from 'constants/globals';
import debounce from 'debounce';
import { Resource } from 'models/Resource';
import React, { FC, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import ResourceService from 'services/ResourceService';
import { useAppDispatch, useAppSelector } from 'store/hooks';
import { setCurrentResource, toggleMenuVisible } from 'store/resources/slice';
import { getFilteredResources, getFilteredResourcesNextPage } from 'store/resources/thunks';
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
	const { currentPage, next, setCurrentPage } = usePagination(PAGE_LENGTH, count);
	const [order, setOrder] = useState<'asc' | 'desc' | undefined>();
	const [orderBy, setOrderBy] = useState<string | undefined>();

	useEffect(() => {
		if (resources && isInfiniteScroll) {
			setIsInfiniteScroll(false);
			next();
		}
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [resources, isInfiniteScroll]);

	useEffect(() => {
		setOrderBy(undefined);
		setCurrentPage(1);
		setHasNext(true);
	}, [filterTags, setCurrentPage]);

	const handleInfiniteScroll = async (e: React.MouseEvent<HTMLInputElement>): Promise<void> => {
		if (isScrolledForInfiniteScroll(e)) {
			setIsInfiniteScroll(true);
			const response = await ResourceService.getFilteredResources({
				data: filterTags,
				offset: currentPage * PAGE_LENGTH,
				limit: PAGE_LENGTH,
				order,
				orderBy,
			});
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

	const handleRequestSort = (_event: React.MouseEvent<unknown>, property: string) => {
		const isAsc = orderBy === property && order === 'asc';
		const newOrder = isAsc ? 'desc' : 'asc';
		setOrder(newOrder);
		setOrderBy(property);

		dispatch(
			getFilteredResources({
				data: filterTags,
				offset: PAGE_START,
				limit: PAGE_LENGTH,
				order: newOrder,
				orderBy: property,
			})
		);
	};

	const createSortHandler = (property: string) => (event: React.MouseEvent<unknown>) => {
		handleRequestSort(event, property);
	};

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
							<TableCell
								sortDirection={orderBy === t('TYPE') ? order : false}
								sx={tableStyles.headerStyle}>
								<TableSortLabel
									active={orderBy === t('TYPE')}
									direction={orderBy === t('TYPE') ? order : 'asc'}
									onClick={createSortHandler(t('TYPE'))}>
									{t('TYPE')}
									{orderBy === t('TYPE') ? (
										<Box component="span" sx={visuallyHidden}>
											{order === 'desc' ? 'sorted descending' : 'sorted ascending'}
										</Box>
									) : null}
								</TableSortLabel>
							</TableCell>
							<TableCell
								sortDirection={orderBy === t('ID') ? order : false}
								align="left"
								sx={tableStyles.headerStyle}>
								<TableSortLabel
									active={orderBy === t('ID')}
									direction={orderBy === t('ID') ? order : 'asc'}
									onClick={createSortHandler(t('ID'))}>
									{t('ID')}
									{orderBy === t('ID') ? (
										<Box component="span" sx={visuallyHidden}>
											{order === 'desc' ? 'sorted descending' : 'sorted ascending'}
										</Box>
									) : null}
								</TableSortLabel>
							</TableCell>
							<TableCell
								sortDirection={orderBy === t('REGION') ? order : false}
								align="left"
								sx={tableStyles.headerStyle}>
								<TableSortLabel
									active={orderBy === t('REGION')}
									direction={orderBy === t('REGION') ? order : 'asc'}
									onClick={createSortHandler(t('REGION'))}>
									{t('REGION')}
									{orderBy === t('REGION') ? (
										<Box component="span" sx={visuallyHidden}>
											{order === 'desc' ? 'sorted descending' : 'sorted ascending'}
										</Box>
									) : null}
								</TableSortLabel>
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
									'&:hover': tableStyles.hoverStyle,
								}}>
								<TableCell sx={tableStyles.bodyRow} component="th" scope="row">
									{row.type}
								</TableCell>
								<TableCell sx={tableStyles.bodyRow} align="left">
									{row.displayId || row.id}
								</TableCell>
								<TableCell sx={tableStyles.regionRow} align="left">
									{row.region}
								</TableCell>
							</TableRow>
						))}
					</TableBody>
				</Table>
				{(!resources || (isInfiniteScroll && hasNext)) && (
					<Box
						sx={{ display: 'flex', justifyContent: 'center', height: '100px', alignItems: 'center' }}
						mt={1}>
						{<CircularProgress color="primary" />}
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
