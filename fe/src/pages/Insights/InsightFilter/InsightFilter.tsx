import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import Accordion from '@mui/material/Accordion';
import AccordionDetails from '@mui/material/AccordionDetails';
import AccordionSummary from '@mui/material/AccordionSummary';
import Box from '@mui/material/Box';
import CircularProgress from '@mui/material/CircularProgress';
import Typography from '@mui/material/Typography';
import AccordionFilter from 'components/AccordionFilter';
import { DARK_BLUE } from 'constants/colors';
import { CHECKED_BY_DEFAULT, PAGE_LENGTH, PAGE_START } from 'constants/globals';
import { Field, FieldGroup, ValueType } from 'models/Field';
import { Tag } from 'models/Tag';
import React, { FC, useEffect, useRef, useState } from 'react';
import { useAppDispatch, useAppSelector } from 'store/hooks';
import { getFilteredResources } from 'store/resources/thunks';
import { setFilterTags } from 'store/tags/slice';

import { accordionStyles, filterStyles, overrideSummaryClasses } from '../style';
import { capitalize } from './helper';

const InsightFilter: FC = () => {
	const { fields } = useAppSelector(state => state.tags);
	const [filterTags, setTags] = useState<Tag[]>([]);
	const dispatch = useAppDispatch();
	const isFirstRun = useRef(true);

	useEffect(() => {
		if (filterTags?.length) {
			if (isFirstRun.current) {
				isFirstRun.current = false;
				return;
			}

			dispatch(getFilteredResources({ data: filterTags, offset: PAGE_START, limit: PAGE_LENGTH }));
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
			setTags(tags);
			dispatch(setFilterTags(tags));
		}
	}, [fields, filterTags?.length, dispatch]);

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
							expandIcon={<ExpandMoreIcon sx={{ color: DARK_BLUE }} />}
							aria-controls={`${field.name}-content`}
							id={`${field.name}-header`}
							key={field.name}
							classes={overrideSummaryClasses}>
							<Typography sx={accordionStyles.accordionHeader}>{capitalize(field.name)}</Typography>
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
				{!fields?.length && (
					<Box
						sx={{ display: 'flex', justifyContent: 'center', height: '100px', alignItems: 'center' }}
						mt={1}>
						{<CircularProgress />}
					</Box>
				)}
			</Box>
		</Box>
	);
};

export default InsightFilter;
