import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import Accordion from '@mui/material/Accordion';
import AccordionDetails from '@mui/material/AccordionDetails';
import AccordionSummary from '@mui/material/AccordionSummary';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import AccordionFilter from 'components/AccordionFilter';
import { CHECKED_BY_DEFAULT, PAGE_LENGTH, PAGE_START } from 'constants/globals';
import { Field, FieldGroup, ValueType } from 'models/Field';
import { Tag } from 'models/Tag';
import React, { FC, useEffect, useState } from 'react';
import { useAppDispatch, useAppSelector } from 'store/hooks';
import { getFilteredResources, getResources } from 'store/resources/thunks';

import { accordionStyles, filterStyles, overrideSummaryClasses } from '../style';

const InsightFilter: FC = () => {
	const { fields } = useAppSelector(state => state.tags);
	const [filterTags, setFilterTags] = useState<Tag[]>([]);
	const dispatch = useAppDispatch();

	useEffect(() => {
		if (filterTags?.length) {
			dispatch(getFilteredResources({ data: filterTags, offset: PAGE_START, limit: PAGE_LENGTH }));
		} else if (filterTags?.length === 0) {
			dispatch(getResources());
		}
	}, [dispatch, filterTags]);

	useEffect(() => {
		if (fields?.length && !filterTags?.length && CHECKED_BY_DEFAULT) {
			const tags = fields.flatMap(field =>
				field.fields.flatMap((fieldItem: Field) =>
					fieldItem.values.flatMap((valueItem: ValueType) => {
						return { key: fieldItem.name, value: valueItem.value };
					})
				)
			);
			setFilterTags(tags);
		}
	}, [fields, filterTags?.length]);

	return (
		<Box
			sx={{
				backgroundColor: '#F9F7F6',
				width: '20%',
				overflowX: 'visible',
			}}>
			<Box>
				{fields.map((field: FieldGroup) => (
					<Accordion key={field.name}>
						<AccordionSummary
							sx={filterStyles.filterHeader}
							expandIcon={<ExpandMoreIcon sx={{ color: 'white' }} />}
							aria-controls={`${field.name}-content`}
							id={`${field.name}-header`}
							key={field.name}
							classes={overrideSummaryClasses}>
							<Typography sx={accordionStyles.accordionHeader}>{field.name.toUpperCase()}</Typography>
						</AccordionSummary>
						<AccordionDetails sx={{ padding: '0px' }}>
							{field.fields.map((fieldItem: Field, index: number) => (
								<AccordionFilter
									key={fieldItem.name + index}
									field={fieldItem}
									hasSearch={true}
									label={fieldItem.name}
									id={fieldItem.name + index}
								/>
							))}
						</AccordionDetails>
					</Accordion>
				))}
			</Box>
		</Box>
	);
};

export default InsightFilter;
