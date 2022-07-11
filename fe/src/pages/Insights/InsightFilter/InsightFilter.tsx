import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import FilterListOffIcon from '@mui/icons-material/FilterListOff';
import Accordion from '@mui/material/Accordion';
import AccordionDetails from '@mui/material/AccordionDetails';
import AccordionSummary from '@mui/material/AccordionSummary';
import Box from '@mui/material/Box';
import CircularProgress from '@mui/material/CircularProgress';
import Tooltip from '@mui/material/Tooltip';
import Typography from '@mui/material/Typography';
import AccordionFilter from 'components/AccordionFilter';
import { DARK_BLUE } from 'constants/colors';
import { CHECKED_BY_DEFAULT, PAGE_LENGTH, PAGE_START } from 'constants/globals';
import { t } from 'i18next';
import { Field, FieldGroup, ValueType } from 'models/Field';
import { Tag } from 'models/Tag';
import React, { FC, useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useAppDispatch, useAppSelector } from 'store/hooks';
import { getFilteredResources } from 'store/resources/thunks';
import { setFilterTags } from 'store/tags/slice';

import { accordionStyles, filterStyles, overrideSummaryClasses } from '../style';
import { capitalize } from './helper';
import { filterHeader, filterIcon } from './style';

const InsightFilter: FC = () => {
	const { fields, filterTags } = useAppSelector(state => state.tags);
	const [tags, setTags] = useState<Tag[]>([]);
	const dispatch = useAppDispatch();
	const isFirstRun = useRef(true);
	const { t } = useTranslation();

	useEffect(() => {
		if (tags?.length) {
			if (isFirstRun.current) {
				isFirstRun.current = false;
				return;
			}

			dispatch(getFilteredResources({ data: tags, offset: PAGE_START, limit: PAGE_LENGTH }));
		}
	}, [dispatch, tags]);

	useEffect(() => {
		if (fields?.length && !tags?.length && CHECKED_BY_DEFAULT) {
			const allTags = fields.flatMap(field =>
				field.fields.flatMap((fieldItem: Field) =>
					fieldItem.values.flatMap((valueItem: ValueType) => {
						return { key: fieldItem.name, value: valueItem.value, group: field.name };
					})
				)
			);
			setTags(allTags);
			dispatch(setFilterTags(allTags));
		}
	}, [fields, tags?.length, dispatch]);

	const handleResetFilters = (e: React.MouseEvent, fieldName: string) => {
		e.stopPropagation();
		e.preventDefault();

		const resetTags = fields
			.filter(field => field.name === fieldName)
			.flatMap(field =>
				field.fields.flatMap((fieldItem: Field) =>
					fieldItem.values.flatMap((valueItem: ValueType) => {
						return { key: fieldItem.name, value: valueItem.value, group: field.name };
					})
				)
			);

		const restOfTags = filterTags.filter(tag => tag.group !== fieldName);
		const newTags = restOfTags.concat(resetTags);
		dispatch(getFilteredResources({ data: newTags, offset: PAGE_START, limit: PAGE_LENGTH }));
	};

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
							<Box sx={filterHeader}>
								<Typography sx={accordionStyles.accordionHeader}>{capitalize(field.name)}</Typography>
								<Box
									onClick={(e: React.MouseEvent) => {
										handleResetFilters(e, field.name);
									}}>
									<Tooltip enterDelay={1000} title={t('RESET_FILTERS')}>
										<FilterListOffIcon fontSize={'medium'} sx={filterIcon} />
									</Tooltip>
								</Box>
							</Box>
						</AccordionSummary>
						<AccordionDetails sx={{ padding: '0px' }}>
							{field.fields.map((fieldItem: Field, index: number) => (
								<AccordionFilter
									key={fieldItem.name + index}
									field={fieldItem}
									hasSearch={true}
									label={fieldItem.name}
									id={fieldItem.name + index}
									fieldGroup={field}
								/>
							))}
						</AccordionDetails>
					</Accordion>
				))}
				{!fields?.length && (
					<Box
						sx={{ display: 'flex', justifyContent: 'center', height: '100px', alignItems: 'center' }}
						mt={1}>
						{<CircularProgress color="primary" />}
					</Box>
				)}
			</Box>
		</Box>
	);
};

export default InsightFilter;
