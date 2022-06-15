import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import Accordion from '@mui/material/Accordion';
import AccordionDetails from '@mui/material/AccordionDetails';
import AccordionSummary from '@mui/material/AccordionSummary';
import Box from '@mui/material/Box';
import FormGroup from '@mui/material/FormGroup';
import Typography from '@mui/material/Typography';
import SearchInput from 'components/SearchInput/SearchInput';
import { SEARCH_ELEMENTS_NUMBER } from 'constants/globals';
import { ValueType } from 'models/Field';
import React, { ChangeEvent, FC, useEffect, useRef, useState } from 'react';
import useHover from 'utils/hooks/useHover';

import AccordionItem from './AccordionItem';
import { accordionStyles, filterStyles, overrideSummaryClasses } from './style';
import { AccordionFilterProps } from './types';

const AccordionFilter: FC<AccordionFilterProps> = props => {
	const { label, hasSearch, id, field, handleChange, checkedByDefault } = props;
	const [searchTerm, setSearchTerm] = useState('');
	const [applyHover, setApplyHover] = useState(false);
	const [boxHeight, setBoxHeight] = useState('unset');
	const [containerRef, isHovered] = useHover<HTMLDivElement>();
	const accordionRef = useRef<HTMLElement>();
	const [expanded, setExpanded] = useState(false);
	const [allIncluded, setAllIncluded] = useState(false);
	const [singleItem, setSingleItem] = useState('');

	const handleExpand = () => {
		setExpanded(!expanded);
	};

	const handleSearchTerm = (e: ChangeEvent<HTMLInputElement>): void => {
		setSearchTerm(e.target.value);
	};

	useEffect(() => {
		if (expanded && isHovered && accordionRef?.current?.clientHeight) {
			setBoxHeight(`${accordionRef?.current.clientHeight + 4}px`);
			setApplyHover(true);
		} else {
			setBoxHeight('unset');
			setApplyHover(false);
		}
	}, [isHovered, accordionRef?.current?.clientHeight, expanded, accordionRef]);

	const handleOnly = (item: ValueType) => {
		setSingleItem(item.value);
		setAllIncluded(false);
	};

	const handleAll = () => {
		setAllIncluded(true);
		setSingleItem('');
	};

	return (
		<Box ref={accordionRef} key={id} sx={{ position: 'relative', height: boxHeight }}>
			<Box sx={expanded && applyHover ? { ...filterStyles.filterHover } : undefined}>
				<Accordion expanded={expanded} onChange={handleExpand}>
					<AccordionSummary
						sx={filterStyles.filterHeader}
						expandIcon={<ExpandMoreIcon />}
						aria-controls={`${id}-content`}
						id={`${id}-header`}
						classes={overrideSummaryClasses}>
						<Typography sx={accordionStyles.accordionHeader}>{label}</Typography>
					</AccordionSummary>
					<AccordionDetails ref={containerRef} sx={accordionStyles.details}>
						{hasSearch && field?.values?.length > SEARCH_ELEMENTS_NUMBER && (
							<SearchInput onChange={handleSearchTerm} />
						)}
						<Box>
							<FormGroup>
								{field?.values &&
									field?.values
										.filter(item => item.value?.toLowerCase()?.includes(searchTerm?.toLowerCase()))
										.map((item: ValueType) => (
											<AccordionItem
												key={item.value}
												field={field}
												item={item}
												handleChange={handleChange}
												isChecked={
													singleItem
														? allIncluded || item.value === singleItem
														: checkedByDefault
												}
												handleOnly={handleOnly}
												handleAll={handleAll}
												singleItem={singleItem}
											/>
										))}
							</FormGroup>
						</Box>
					</AccordionDetails>
				</Accordion>
			</Box>
		</Box>
	);
};

export default AccordionFilter;
