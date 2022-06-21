import { GREY, WHITE } from 'constants/colors';
import { CSSProperties } from 'react';

export const headerStyle: CSSProperties = {
	height: '64px',
	width: '100%',
	display: 'flex',
	alignItems: 'center',
	border: '1px solid #EAEBF0',
	backgroundColor: WHITE,
	justifyContent: 'space-between',
};

export const menuItems: CSSProperties = {
	color: GREY,
	lineHeight: '18px',
	letterSpacing: '0em',
	cursor: 'pointer',
	fontSize: '14px',
	fontWeight: '600',
	fontFamily: 'Montserrat',
};
