import Box from '@mui/material/Box';
import Paper from '@mui/material/Paper';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import { Resource } from 'models/Resource';
import React, { FC } from 'react';
import { useTranslation } from 'react-i18next';
import { useAppDispatch, useAppSelector } from 'store/hooks';
import { setCurrentResource, toggleMenuVisible } from 'store/resources/slice';

import { tableStyles } from './style';

const InsightTable: FC = () => {
	const { resources } = useAppSelector(state => state.resources);
	const { t } = useTranslation();
	const dispatch = useAppDispatch();

	const handleClick = (resource: Resource) => {
		dispatch(setCurrentResource(resource));
		dispatch(toggleMenuVisible());
	};

	return (
		<Box
			sx={{
				width: '80%',
				height: '100%',
				backgroundColor: '#F8F9FD',
			}}>
			<TableContainer component={Paper}>
				<Table sx={{ minWidth: 650 }} size="small" aria-label="a dense table">
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
						{resources.map((row: Resource, index: number) => (
							<TableRow
								onClick={() => handleClick(row)}
								key={row.id + row.type + index}
								sx={{
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
			</TableContainer>
		</Box>
	);
};

export default InsightTable;
