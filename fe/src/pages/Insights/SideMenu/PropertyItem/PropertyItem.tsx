import AddBoxIcon from '@mui/icons-material/AddBox';
import IndeterminateCheckBoxIcon from '@mui/icons-material/IndeterminateCheckBox';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import { TEXT_COLOR } from 'constants/colors';
import React, { FC, useState } from 'react';

import { sideMenuLeftText, sideMenuRightText } from '../style';
import { PropertyItemProps } from './types';

const PropertyItem: FC<PropertyItemProps> = props => {
	const { keyItem, value } = props;
	const [expanded, setExpanded] = useState(false);

	const renderObjects = (data: Object): React.ReactNode =>
		Object.entries(data).map(([key, value]) => <PropertyItem key={`${key}${value}`} keyItem={key} value={value} />);

	const renderArrayOrObjects = (data: Object | Array<Object>): React.ReactNode => {
		if (Array.isArray(data)) {
			return (
				<Box sx={{ lineHeight: '1' }}>
					<Typography color={TEXT_COLOR} sx={{ opacity: '0.9', display: 'flex' }}>{`[`}</Typography>
					{data.map(item => renderObjects(item))}
					<Typography color={TEXT_COLOR} sx={{ opacity: '0.9', display: 'flex' }}>{`]`}</Typography>
				</Box>
			);
		} else {
			return renderObjects(data);
		}
	};

	const handleClick = () => {
		setExpanded(!expanded);
	};

	return (
		<>
			{typeof value === 'object' && !Array.isArray(value) && value !== null ? (
				<Box sx={{ display: expanded ? 'block' : 'flex' }}>
					<Box sx={{ display: 'flex' }}>
						{expanded ? (
							<IndeterminateCheckBoxIcon
								color={'primary'}
								fontSize="small"
								onClick={handleClick}
								sx={{ cursor: 'pointer', display: 'flex' }}
							/>
						) : (
							<AddBoxIcon
								color={'primary'}
								fontSize="small"
								onClick={handleClick}
								sx={{ cursor: 'pointer', display: 'flex' }}
							/>
						)}
						<Typography mr={2} {...sideMenuLeftText} sx={{ display: 'flex' }}>{`${keyItem} `}</Typography>
					</Box>
					{expanded ? (
						<>
							<Typography color={TEXT_COLOR} sx={{ opacity: '0.9', display: 'flex' }}>{`{`}</Typography>
							<Box sx={{ display: 'block' }}>
								<Box ml={2}>{renderArrayOrObjects(value)}</Box>
								<Typography
									color={TEXT_COLOR}
									sx={{ opacity: '0.9', display: 'flex' }}>{`}`}</Typography>
							</Box>
						</>
					) : (
						<Typography color={TEXT_COLOR} sx={{ opacity: '0.9', display: 'flex', lineHeight: '1' }}>
							{` { ... } `}
						</Typography>
					)}
				</Box>
			) : (
				<Box sx={{ display: 'flex' }}>
					<Typography mr={2} {...sideMenuLeftText}>{`${keyItem} `}</Typography>
					<Typography {...sideMenuRightText}>
						{Array.isArray(value) ? renderArrayOrObjects(value) : ` ${value}`}
						{/* {` ${value}`} */}
					</Typography>
				</Box>
			)}
		</>
	);
};

export default PropertyItem;
