import { BLACK, WHITE } from 'constants/colors';
import { CSSProperties } from 'react';

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

export const accordionStyles = { accordionHeader };

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
	width: '400px',
	position: 'absolute',
	left: '0px ',
	zIndex: '10',
	backgroundColor: 'white',
	overflow: 'visible',
	top: '4px',
};

export const filterStyles = { filterHeader, filterHover };
