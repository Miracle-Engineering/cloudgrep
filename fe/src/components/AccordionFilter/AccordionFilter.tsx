import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import Accordion from '@mui/material/Accordion';
import AccordionDetails from '@mui/material/AccordionDetails';
import AccordionSummary from '@mui/material/AccordionSummary';
import Box from '@mui/material/Box';
import Checkbox from '@mui/material/Checkbox';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormGroup from '@mui/material/FormGroup';
import Typography from '@mui/material/Typography';
import SearchInput from 'components/SearchInput/SearchInput';
import { ValueType } from 'models/Field';
import React, { ChangeEvent, FC, useEffect, useRef, useState } from 'react';
import useHover from 'utils/hooks/useHover';

import { accordionStyles, filterStyles, labelClasses, overrideSummaryClasses } from './style';
import { AccordionFilterProps } from './types';

const SEARCH_ELEMENTS_NUMBER = 3;

const AccordionFilter: FC<AccordionFilterProps> = props => {
	const { label, hasSearch, id, field, handleChange, checkedByDefault } = props;
	const [searchTerm, setSearchTerm] = useState('');
	const [applyHover, setApplyHover] = useState(false);
	const [boxHeight, setBoxHeight] = useState('unset');
	const [containerRef, isHovered] = useHover<HTMLDivElement>();
	const accordionRef = useRef<HTMLElement>();
	const [expanded, setExpanded] = useState(false);

	const handleExpand = () => {
		setExpanded(!expanded);
	};

	const handleSearchTerm = (e: ChangeEvent<HTMLInputElement>): void => {
		setSearchTerm(e.target.value);
	};

	useEffect(() => {
		if (expanded && isHovered && accordionRef?.current?.clientHeight) {
			setBoxHeight(`${accordionRef?.current.clientHeight}px`);
			setApplyHover(true);
		} else {
			setBoxHeight('unset');
			setApplyHover(false);
		}
	}, [isHovered, accordionRef?.current?.clientHeight, expanded, accordionRef]);

	return (
		<Box ref={accordionRef} key={id} sx={{ height: boxHeight, position: 'relative' }}>
			<Accordion expanded={expanded} onChange={handleExpand}>
				{/* <Accordion expanded={expanded} onChange={handleExpand} sx={expanded ? filterStyles.filterHover : undefined}> */}
				{/* sx={expanded && applyHover ? filterStyles.filterHover : undefined}> */}
				<AccordionSummary
					sx={filterStyles.filterHeader}
					expandIcon={<ExpandMoreIcon />}
					aria-controls={`${id}-content`}
					id={`${id}-header`}
					classes={overrideSummaryClasses}>
					<Typography sx={accordionStyles.accordionHeader}>{label}</Typography>
				</AccordionSummary>
				<AccordionDetails ref={containerRef}>
					{hasSearch && field?.values?.length > SEARCH_ELEMENTS_NUMBER && (
						<SearchInput onChange={handleSearchTerm} />
					)}
					<Typography>
						<FormGroup>
							{field?.values &&
								field?.values
									.filter(item => item.value?.toLowerCase()?.includes(searchTerm?.toLowerCase()))
									.map((item: ValueType) => (
										<Box
											key={item.value}
											sx={{
												display: 'flex',
												alignItems: 'center',
												justifyContent: 'space-between',
												maxWidth: '100%',
											}}>
											<FormControlLabel
												sx={{
													overflow: 'hidden',
													textOverflow: 'ellipsis',
													whiteSpace: 'nowrap',
												}}
												classes={labelClasses}
												control={
													<Checkbox
														size={'small'}
														defaultChecked={checkedByDefault}
														onChange={e => handleChange(e, field, item)}
													/>
												}
												label={item.value}
											/>
											<Typography sx={{ fontSize: '13px', fontWeight: 600 }}>
												{item.count}
											</Typography>
										</Box>
									))}
						</FormGroup>
					</Typography>
				</AccordionDetails>
			</Accordion>
		</Box>
	);
};

export default AccordionFilter;
