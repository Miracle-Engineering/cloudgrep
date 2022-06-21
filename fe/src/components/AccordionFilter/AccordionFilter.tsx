import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import Accordion from '@mui/material/Accordion';
import AccordionDetails from '@mui/material/AccordionDetails';
import AccordionSummary from '@mui/material/AccordionSummary';
import Box from '@mui/material/Box';
import FormGroup from '@mui/material/FormGroup';
import Typography from '@mui/material/Typography';
import SearchInput from 'components/SearchInput/SearchInput';
import { SEARCH_ELEMENTS_NUMBER } from 'constants/globals';
import { PAGE_LENGTH, PAGE_START } from 'constants/globals';
import { Field, ValueType } from 'models/Field';
import { Tag } from 'models/Tag';
import React, { ChangeEvent, FC, useEffect, useMemo, useRef, useState } from 'react';
import { useAppDispatch, useAppSelector } from 'store/hooks';
import { getFilteredResources } from 'store/resources/thunks';
import useHover from 'utils/hooks/useHover';

import AccordionItem from './AccordionItem';
import { accordionStyles, filterStyles, overrideSummaryClasses } from './style';
import { AccordionFilterProps } from './types';

const AccordionFilter: FC<AccordionFilterProps> = props => {
	const { label, hasSearch, id, field } = props;
	const [searchTerm, setSearchTerm] = useState('');
	const [applyHover, setApplyHover] = useState(false);
	const [boxHeight, setBoxHeight] = useState('unset');
	const [containerRef, isHovered] = useHover<HTMLDivElement>();
	const accordionRef = useRef<HTMLElement>();
	const [expanded, setExpanded] = useState(false);
	const [singleItem, setSingleItem] = useState('');
	const initialTags = useMemo(
		() =>
			field?.values?.map((item: ValueType) => {
				return { key: field.name, value: item.value };
			}),
		[field]
	);
	const [tags, setTags] = useState<Tag[]>(initialTags);
	const { filterTags } = useAppSelector(state => state.tags);
	const dispatch = useAppDispatch();
	const isFirstRun = useRef(true);
	const [isChanged, setIsChanged] = useState(false);

	useEffect(() => {
		if (expanded && isHovered && accordionRef?.current?.clientHeight) {
			setBoxHeight(`${accordionRef?.current.clientHeight + 4}px`);
			setApplyHover(true);
		} else {
			setBoxHeight('unset');
			setApplyHover(false);
		}
	}, [isHovered, expanded, accordionRef]);

	useEffect(() => {
		if (isFirstRun.current) {
			isFirstRun.current = false;
			return;
		}

		if (!isChanged) {
			return;
		}

		const allOtherTags = filterTags.filter(tag => tag.key !== field.name);
		const newTags = allOtherTags.concat(tags);
		dispatch(getFilteredResources({ data: newTags, offset: PAGE_START, limit: PAGE_LENGTH }));
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [tags, isChanged]);

	const handleExpand = () => {
		setExpanded(!expanded);
	};

	const handleSearchTerm = (e: ChangeEvent<HTMLInputElement>): void => {
		setSearchTerm(e.target.value);
	};

	const handleOnly = (item: ValueType) => {
		setTags([{ key: field.name, value: item.value }]);
		setSingleItem(item.value);
		setIsChanged(true);
	};

	const handleAll = () => {
		setTags(initialTags);
		setSingleItem('');
		setIsChanged(true);
	};

	const handleChange = (event: React.ChangeEvent<HTMLInputElement>, field: Field, item: ValueType) => {
		const tag = { key: field.name, value: item.value };
		const existingTag = tags?.some(filterTag => filterTag.key === tag.key && filterTag.value === tag.value);

		if (event.target.checked && !existingTag) {
			setTags([...tags, tag]);
		} else if (!event.target.checked && existingTag && tags?.length > 0) {
			const newFilters = tags.filter(filterTag => !(filterTag.key === tag.key && filterTag.value === tag.value));
			setTags(newFilters);
		}
		setIsChanged(true);
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
							<SearchInput onChange={handleSearchTerm} height={32} />
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
												isChecked={tags.some(
													tag => tag.key === field.name && tag.value === item.value
												)}
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
