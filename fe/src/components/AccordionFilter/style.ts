import { BLACK, BLACK_ROW, DARK_BLUE, GREY_ONLY_ITEM, WHITE } from 'constants/colors';
import { CSSProperties } from 'react';

export const overrideSummaryClasses = {
	content: 'summary_content',
	root: 'summary_root',
};

const accordionHeader: CSSProperties = {
	fontStyle: 'normal',
	fontWeight: '400',
	fontSize: '13px',
	lineHeight: '140%',
	marginBottom: '4px',
	fontFamily: 'Montserrat',
	overflow: 'hidden',
	whiteSpace: 'nowrap',
	textOverflow: 'ellipsis',
	color: DARK_BLUE,
};

const details: CSSProperties = {
	overflowY: 'auto',
	maxHeight: '210px', // 5 elements visible before scroll and overflow
};

export const accordionStyles = { accordionHeader, details };

export const labelClasses = {
	label: 'label_label',
};

export const filterHeader = {
	backgroundColor: `${WHITE}`,
	color: `${BLACK}`,
	opacity: '0.8',
	borderRadius: '4px',
	minHeight: '42px !important',
};

export const filterHover = {
	width: '500px',
	position: 'absolute',
	left: '0px ',
	zIndex: '10',
	overflow: 'visible',
};

export const filterStyles = { filterHeader, filterHover };

export const accordionItem: CSSProperties = {
	fontFamily: 'Montserrat',
	fontStyle: 'normal',
	fontWeight: 500,
	fontSize: '12px',
	lineHeight: '140%',
	color: BLACK_ROW,
};

export const onlyAndAll: CSSProperties = {
	fontFamily: 'Montserrat',
	fontStyle: 'normal',
	fontWeight: 500,
	fontSize: '12px',
	lineHeight: '140%',
	color: GREY_ONLY_ITEM,
	cursor: 'pointer',
};
