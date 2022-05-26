import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import Accordion from '@mui/material/Accordion';
import AccordionDetails from '@mui/material/AccordionDetails';
import AccordionSummary from '@mui/material/AccordionSummary';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import AccordionFilter from 'components/AccordionFilter';
import { Field, ValueType } from 'models/Field';
import { Tag } from 'models/Tag';
import React, { FC, useEffect, useMemo, useState } from 'react';
import { useAppDispatch, useAppSelector } from 'store/hooks';
import { getFilteredResources, getResources } from 'store/resources/thunks';

import { accordionStyles, filterStyles, overrideSummaryClasses } from '../style';

const InsightFilter: FC = () => {
	const { fields } = useAppSelector(state => state.tags);
	const [filterTags, setFilterTags] = useState<Tag[]>([]);
	const dispatch = useAppDispatch();

	const groups = useMemo(() => {
		return new Set(fields.map(field => field.group));
	}, [fields]);

	useEffect(() => {
		if (filterTags?.length) {
			dispatch(getFilteredResources(filterTags));
		} else {
			dispatch(getResources());
		}
	}, [filterTags]);

	const handleChange = (event: React.ChangeEvent<HTMLInputElement>, field: Field, item: ValueType) => {
		const tag = { key: field.name, value: item.value };
		const existingTag = filterTags?.some(filterTag => filterTag.key === tag.key && filterTag.value === tag.value);
		if (event.target.checked && !existingTag) {
			setFilterTags([...filterTags, tag]);
		} else if (!event.target.checked && existingTag && filterTags) {
			setFilterTags(filterTags.filter(filterTag => filterTag.key !== tag.key && filterTag.value !== tag.value));
		}
	};

	return (
		<Box
			sx={{
				width: '20%',
				height: '100%',
				backgroundColor: '#F9F7F6',
				overflowY: 'scroll',
			}}>
			{Array.from(groups).map((group: string) => (
				<>
					<Accordion key={group}>
						<AccordionSummary
							sx={filterStyles.filterHeader}
							expandIcon={<ExpandMoreIcon sx={{ color: 'white' }} />}
							aria-controls={`${group}-content`}
							id={`${group}-header`}
							key={group}
							classes={overrideSummaryClasses}>
							<Typography sx={accordionStyles.accordionHeader}>{group}</Typography>
						</AccordionSummary>
						<AccordionDetails key={group} sx={{ padding: '0px' }}>
							{fields
								.filter((field: Field) => field.group === group)
								.map((field: Field, index: number) => (
									<AccordionFilter
										key={field.name + index}
										field={field}
										hasSearch={true}
										label={field.name}
										id={field.name + index}
										handleChange={handleChange}
									/>
								))}
						</AccordionDetails>
					</Accordion>
				</>
			))}
		</Box>
	);
};

export default InsightFilter;
