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
	fontFamily: 'NotoSans,Lucida Grande,Lucida Sans Unicode,sans-serif',
};

export const sideMenuRightText = {
	color: 'rgba(28, 43, 52, 0.68)',
	fontSize: '12px',
	lineHeight: '19.2px',
	fontFamily: 'NotoSans,Lucida Grande,Lucida Sans Unicode,sans-serif',
};

export const boxStyle = {
	display: 'flex',
	flexDirection: 'column',
	border: '1px solid rgb(183, 188, 203)',
	minWidth: '152px',
	alignItems: 'baseline',
};

export const boxFirstLine = {
	color: 'rgba(28, 43, 52, 0.98)',
	fontWeight: 700,
	fontSize: '12px',
	fontFamily: 'NotoSans,Lucida Grande,Lucida Sans Unicode,sans-serif',
};

export const boxSecondLine = {
	color: 'rgba(28, 43, 52, 0.98)',
	fontWeight: 400,
	fontSize: '13px',
	fontFamily: 'NotoSans,Lucida Grande,Lucida Sans Unicode,sans-serif',
};
