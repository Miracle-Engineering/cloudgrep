import AddBoxIcon from '@mui/icons-material/AddBox';
import IndeterminateCheckBoxIcon from '@mui/icons-material/IndeterminateCheckBox';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import { TEXT_COLOR } from 'constants/colors';
import React, { FC, useState } from 'react';

import { PropertyItemListProps } from './types';

const PropertyItemList: FC<PropertyItemListProps> = props => {
	const { data, renderObjects } = props;
	const [expanded, setExpanded] = useState(false);

	const handleClick = () => {
		setExpanded(!expanded);
	};

	return (
		<>
			<Box sx={{ lineHeight: '1' }}>
				{expanded ? (
					<>
						<IndeterminateCheckBoxIcon
							color={'primary'}
							fontSize="small"
							onClick={handleClick}
							sx={{ cursor: 'pointer', display: 'flex' }}
						/>
						<Typography color={TEXT_COLOR} sx={{ opacity: '0.9', display: 'flex' }}>{`[`}</Typography>
						{data.map(item => renderObjects(item))}
						<Typography color={TEXT_COLOR} sx={{ opacity: '0.9', display: 'flex' }}>{`]`}</Typography>
					</>
				) : (
					<>
						<AddBoxIcon
							color={'primary'}
							fontSize="small"
							onClick={handleClick}
							sx={{ cursor: 'pointer', display: 'flex' }}
						/>
						<Typography color={TEXT_COLOR} sx={{ opacity: '0.9', display: 'flex' }}>{`[ ... ]`}</Typography>
					</>
				)}
			</Box>
		</>
	);
};

export default PropertyItemList;
