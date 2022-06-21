import { BLACK_ROW, BORDER_COLOR, GREY_SIDE_MENU, GREY_TAB_ITEM } from 'constants/colors';
import { CSSProperties } from 'react';

export const sideMenuStyle: CSSProperties = {
	height: '100%',
	maxHeight: '100%',
	overflow: 'scroll',
	width: '50%',
	border: '1px solid black',
	position: 'fixed',
	right: '0px',
	top: '0px',
	backgroundColor: 'white',
};

export const sideMenuLeftText = {
	color: 'rgb(78, 145, 209)',
	fontSize: '12px',
	lineHeight: '19.2px',
	fontFamily: 'Montserrat',
};

export const sideMenuRightText = {
	color: 'rgba(28, 43, 52, 0.68)',
	fontSize: '12px',
	lineHeight: '19.2px',
	fontFamily: 'Montserrat',
};

export const boxStyle = {
	display: 'flex',
	flexDirection: 'column',
	border: `1px solid ${BORDER_COLOR}`,
	minWidth: '152px',
	alignItems: 'baseline',
};

export const boxFirstLine = {
	color: BLACK_ROW,
	fontWeight: 600,
	fontSize: '16px',
	fontFamily: 'Montserrat',
	lineHeight: '24px',
};

export const boxSecondLine = {
	color: GREY_SIDE_MENU,
	fontWeight: 500,
	fontSize: '12px',
	lineHeight: '18px',
	fontFamily: 'Montserrat',
};

export const sideMenuHeader = {
	fontFamily: 'Montserrat',
	fontStyle: 'normal',
	fontWeight: 500,
	fontSize: '24px',
	lineHeight: '32px',
}

export const tabStyle = {
	fontFamily: 'Montserrat',
	fontStyle: 'normal',
	fontWeight: 600,
	fontSize: '16px',
	lineHeight: '24px',
	color: GREY_TAB_ITEM,
	textTransform: 'none',
}
