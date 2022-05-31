import { DARK_BLUE } from 'constants/colors';
import { CSSProperties } from 'react';

const headerStyle: CSSProperties = {
	backgroundColor: DARK_BLUE,
	color: 'white',
	height: 56,
	fontFamily: 'Montserrat',
	fontSize: 32,
};

const hoverStyle: CSSProperties = {
	cursor: 'pointer',
	backgroundColor: '#8FCAF9',
};

export const tableStyles = { hoverStyle, headerStyle };

export const overrideSummaryClasses = {
	content: 'summary_content',
	root: 'summary_root',
};

const accordionHeader: CSSProperties = {
	fontWeight: '400',
	fontSize: '18px',
	lineHeight: '18px',
	marginBottom: '4px',
	fontFamily: 'Montserrat',
	overflow: 'hidden',
	whiteSpace: 'nowrap',
	textOverflow: 'ellipsis',
};

const accordionDetails: CSSProperties = {
	overflow: 'hidden',
};

export const accordionStyles = { accordionHeader, accordionDetails };

export const labelClasses = {
	label: 'label_label',
};

export const filterHeader = {
	backgroundColor: `${DARK_BLUE}`,
	color: '#FFFFFF',
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
