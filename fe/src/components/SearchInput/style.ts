import { BACKGROUND_COLOR } from 'constants/colors';
import { CSSProperties } from 'react';

export const searchStyle: CSSProperties = {
	padding: '2px 4px',
	display: 'flex',
	alignItems: 'center',
	border: '1px solid #CECDCD',
	borderRadius: '4px',
	boxShadow: 'none',
	backgroundColor: BACKGROUND_COLOR,
};

export const searchText: CSSProperties = {
	fontFamily: 'Montserrat',
	fontStyle: 'normal',
	fontWeight: 400,
	fontSize: '12px',
	lineHeight: '15px',
};
