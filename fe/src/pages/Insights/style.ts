import { AccordionDetailsClasses } from '@mui/material/AccordionDetails/accordionDetailsClasses';
import { AccordionSummaryClasses } from '@mui/material/AccordionSummary/accordionSummaryClasses';
import { CSSProperties } from 'react';

const headerStyle: CSSProperties = {
    fontWeight: 'bold',
    textTransform: 'uppercase',
};

const hoverStyle: CSSProperties = {
    cursor: 'pointer',
    backgroundColor: '#8FCAF9',
};

export const tableStyles = { hoverStyle, headerStyle };

export const overrideSummaryClasses: AccordionSummaryClasses = {
    content: 'summary_content',
    root: '',
    expanded: '',
    focusVisible: '',
    disabled: '',
    gutters: '',
    contentGutters: '',
    expandIconWrapper: '',
};

const accordionHeader: CSSProperties = {
    fontWeight: '400',
    fontSize: '11px',
    textTransform: 'uppercase',
    marginBottom: '4px',
};

const accordionDetails: CSSProperties = {
    fontWeight: '400',
    fontSize: '13px',
    color: 'rgba(28, 43, 52, 0.98)',
    fontFamily: 'NotoSans,Lucida Grande,Lucida Sans Unicode,sans-serif!important',
    overflow: 'hidden',
};

export const accordionStyles = { accordionHeader, accordionDetails };
