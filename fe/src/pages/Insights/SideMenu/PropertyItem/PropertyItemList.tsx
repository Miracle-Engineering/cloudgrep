import AddBoxIcon from '@mui/icons-material/AddBox';
import IndeterminateCheckBoxIcon from '@mui/icons-material/IndeterminateCheckBox';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import { TEXT_COLOR } from 'constants/colors';
import React, { FC, useState } from 'react';

import { sideMenuLeftText } from '../style';
import { PropertyItemListProps } from './types';

const PropertyItemList: FC<PropertyItemListProps> = props => {
	const { data, renderObjects, keyItem } = props;
	const [expanded, setExpanded] = useState(false);

	const handleClick = () => {
		setExpanded(!expanded);
	};

	return (
		<>
			{expanded ? (
				<>
					<Box sx={{ display: 'flex' }}>
						<Typography mr={2} {...sideMenuLeftText} sx={{ display: 'flex' }}>{`${keyItem} `}</Typography>
						<IndeterminateCheckBoxIcon
							color={'primary'}
							fontSize="small"
							onClick={handleClick}
							sx={{ cursor: 'pointer', display: 'flex' }}
						/>
						<Typography
							color={TEXT_COLOR}
							sx={{ opacity: '0.9', display: 'flex', lineHeight: '1' }}>{`[`}</Typography>
					</Box>
					<Box>
						<Box ml={2}>{data.map(item => renderObjects(item))}</Box>
						<Typography color={TEXT_COLOR} sx={{ opacity: '0.9', display: 'flex' }}>{`]`}</Typography>
					</Box>
				</>
			) : (
				<Box sx={{ display: 'flex' }}>
					<Typography mr={2} {...sideMenuLeftText} sx={{ display: 'flex' }}>{`${keyItem} `}</Typography>
					<AddBoxIcon
						color={'primary'}
						fontSize="small"
						onClick={handleClick}
						sx={{ cursor: 'pointer', display: 'flex' }}
					/>
					<Typography
						color={TEXT_COLOR}
						sx={{ opacity: '0.9', display: 'flex', lineHeight: '1' }}>{`[ ... ]`}</Typography>
				</Box>
			)}
		</>
	);
};

export default PropertyItemList;
