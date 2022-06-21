import Box from '@mui/material/Box';
import Checkbox from '@mui/material/Checkbox';
import FormControlLabel from '@mui/material/FormControlLabel';
import Tooltip from '@mui/material/Tooltip';
import Typography from '@mui/material/Typography';
import React, { FC } from 'react';
import useHover from 'utils/hooks/useHover';

import { accordionItem, labelClasses, onlyAndAll } from './style';
import { AccordionItemProps } from './types';

const AccordionItem: FC<AccordionItemProps> = props => {
	const { field, handleChange, isChecked, item, handleOnly, singleItem, handleAll } = props;
	const [containerRef, isHovered] = useHover<HTMLDivElement>();
	const isOnlyUsed = singleItem && singleItem === item.value;

	const handleClick = () => {
		isOnlyUsed ? handleAll() : handleOnly(item);
	};

	return (
		<Box
			key={item.value}
			ref={containerRef}
			sx={{
				display: 'flex',
				alignItems: 'center',
				justifyContent: 'space-between',
				maxWidth: '100%',
			}}>
			<Tooltip enterDelay={1000} title={item.value}>
				<FormControlLabel
					sx={{
						overflow: 'hidden',
						textOverflow: 'ellipsis',
						whiteSpace: 'nowrap',
					}}
					classes={labelClasses}
					control={
						<Checkbox size={'small'} checked={isChecked} onChange={e => handleChange(e, field, item)} />
					}
					label={item.value}
				/>
			</Tooltip>
			<Box
				sx={{ display: 'flex', cursor: isHovered ? 'pointer' : 'unset' }}
				onClick={() => isHovered && handleClick()}>
				{isHovered && (
					<Typography mr={2} sx={onlyAndAll}>
						{isOnlyUsed ? 'All' : 'Only'}
					</Typography>
				)}
				<Typography sx={accordionItem}>{item.count}</Typography>
			</Box>
		</Box>
	);
};

export default AccordionItem;
