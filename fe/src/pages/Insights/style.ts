import { BACKGROUND_COLOR, BLACK_ROW, DARK_BLUE, WHITE } from 'constants/colors';
import { CSSProperties } from 'react';

const headerStyle: CSSProperties = {
	backgroundColor: WHITE,
	color: DARK_BLUE,
	height: '46px',
	fontFamily: 'Montserrat',
	fontSize: '12px',
	paddingTop: 0,
	paddingBottom: 0,
	fontStyle: 'normal',
	fontWeight: '600',
	lineHeight: '130%',
};

const hoverStyle: CSSProperties = {
	cursor: 'pointer',
	backgroundColor: '#8FCAF9',
};

const bodyRow: CSSProperties = {
	fontFamily: 'Montserrat',
	fontStyle: 'normal',
	fontWeight: 500,
	fontSize: '14px',
	lineHeight: '18px',
	color: BLACK_ROW,
}

export const tableStyles = { hoverStyle, headerStyle, bodyRow };

export const overrideSummaryClasses = {
	content: 'summary_content',
	root: 'summary_root',
};

const accordionHeader: CSSProperties = {
	fontWeight: '600',
	fontSize: '14px',
	lineHeight: '120%',
	marginBottom: '4px',
	fontFamily: 'Montserrat',
	overflow: 'hidden',
	whiteSpace: 'nowrap',
	textOverflow: 'ellipsis',
	color: DARK_BLUE,
};

const accordionDetails: CSSProperties = {
	overflow: 'hidden',
};

export const accordionStyles = { accordionHeader, accordionDetails };

export const labelClasses = {
	label: 'label_label',
};

export const filterHeader = {
	backgroundColor: BACKGROUND_COLOR,
	color: DARK_BLUE,
	borderRadius: '4px',
	minHeight: '42px !important',
};

export const filterHover = {
	width: '400px',
	position: 'absolute',
	left: '0px ',
	zIndex: '10',
	backgroundColor: 'white',
	overflow: 'visible',
};

export const filterStyles = { filterHeader, filterHover };
