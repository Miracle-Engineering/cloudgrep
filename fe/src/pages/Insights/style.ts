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

export const overrideSummaryClasses = {
    content: 'summary_content',
    root: 'summary_root',
};

const accordionHeader: CSSProperties = {
    fontWeight: '400',
    fontSize: '11px',
    textTransform: 'uppercase',
    marginBottom: '4px',
};

const accordionDetails: CSSProperties = {
    overflow: 'hidden',
};

export const accordionStyles = { accordionHeader, accordionDetails };

export const labelClasses = {
    label: 'label_label',
};
