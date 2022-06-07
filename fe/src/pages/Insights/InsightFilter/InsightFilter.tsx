import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import Accordion from '@mui/material/Accordion';
import AccordionDetails from '@mui/material/AccordionDetails';
import AccordionSummary from '@mui/material/AccordionSummary';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import AccordionFilter from 'components/AccordionFilter';
import { Field, FieldGroup, ValueType } from 'models/Field';
import { Tag } from 'models/Tag';
import React, { FC, useEffect, useState } from 'react';
import { useAppDispatch, useAppSelector } from 'store/hooks';
import { getFilteredResources, getResources } from 'store/resources/thunks';

import { accordionStyles, filterStyles, overrideSummaryClasses } from '../style';

const CHECKED_BY_DEFAULT = true;

const InsightFilter: FC = () => {
	const { fields } = useAppSelector(state => state.tags);
	const { resources } = useAppSelector(state => state.resources);
	const [filterTags, setFilterTags] = useState<Tag[]>([]);
	const dispatch = useAppDispatch();

	useEffect(() => {
		if (filterTags?.length) {
			dispatch(getFilteredResources(filterTags));
		} else if (filterTags?.length === 0 && resources?.length === 0) {
			dispatch(getResources());
		}
	}, [dispatch, filterTags, resources?.length]);

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

	const handleChange = (event: React.ChangeEvent<HTMLInputElement>, field: Field, item: ValueType) => {
		const tag = { key: field.name, value: item.value };
		const existingTag = filterTags?.some(filterTag => filterTag.key === tag.key && filterTag.value === tag.value);

		if (event.target.checked && !existingTag) {
			setFilterTags([...filterTags, tag]);
		} else if (!event.target.checked && existingTag && filterTags?.length > 0) {
			const newFilters = filterTags.filter(
				filterTag => !(filterTag.key === tag.key && filterTag.value === tag.value)
			);
			setFilterTags(newFilters);
		}
	};

	return (
		<Box
			sx={{
				width: '20%',
				height: '100%',
				backgroundColor: '#F9F7F6',
				overflowY: 'auto',
			}}>
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
								handleChange={handleChange}
								checkedByDefault={CHECKED_BY_DEFAULT}
							/>
						))}
					</AccordionDetails>
				</Accordion>
			))}
		</Box>
	);
};

export default InsightFilter;
