import CloseIcon from '@mui/icons-material/Close';
import IconButton from '@mui/material/IconButton';
import React, { FC } from 'react';

import { Props } from './types';

const BannerAction: FC<Props> = props => {
	const { handleClose } = props;
	return (
		<>
			<IconButton size="small" aria-label="close" color="inherit" onClick={handleClose}>
				<CloseIcon fontSize="small" />
			</IconButton>
		</>
	);
};

export default BannerAction;
