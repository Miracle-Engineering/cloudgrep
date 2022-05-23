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
import React, { ChangeEvent, FC, useState } from 'react';

import { accordionStyles, filterStyles, labelClasses, overrideSummaryClasses } from './style';
import { AccordionFilterProps } from './types';

const AccordionFilter: FC<AccordionFilterProps> = props => {
	const { label, hasSearch, id, field, handleChange } = props;
	const [_, setSearchTerm] = useState('');

	const handleSearchTerm = (e: ChangeEvent<HTMLInputElement>): void => {
		setSearchTerm(e.target.value);
	};

	return (
		<Box>
			<Accordion sx={{ '&:hover': filterStyles.filterHover }}>
				<AccordionSummary
					sx={filterStyles.filterHeader}
					expandIcon={<ExpandMoreIcon sx={{ color: 'white' }} />}
					aria-controls={`${id}-content`}
					id={`${id}-header`}
					classes={overrideSummaryClasses}>
					<Typography sx={accordionStyles.accordionHeader}>{label}</Typography>
				</AccordionSummary>
				<AccordionDetails>
					{hasSearch && <SearchInput onChange={handleSearchTerm} />}
					<Typography>
						<FormGroup>
							{field?.values &&
								field?.values.map((item: ValueType) => (
									<FormControlLabel
										classes={labelClasses}
										key={item.value}
										control={
											<Checkbox
												size={'small'}
												defaultChecked
												onChange={e => handleChange(e, field, item)}
											/>
										}
										label={item.value}
									/>
								))}
						</FormGroup>
					</Typography>
				</AccordionDetails>
			</Accordion>
		</Box>
	);
};

export default AccordionFilter;
