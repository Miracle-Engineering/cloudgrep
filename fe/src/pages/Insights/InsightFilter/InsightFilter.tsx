import Box from '@mui/material/Box';
import AccordionFilter from 'components/AccordionFilter';
import { Field, ValueType } from 'models/Field';
import { Tag } from 'models/Tag';
import React, { ChangeEvent, FC, useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useAppDispatch, useAppSelector } from 'store/hooks';
import { getFilteredResources, getResources } from 'store/resources/thunks';

const InsightFilter: FC = () => {
	const { fields, tagResource } = useAppSelector(state => state.tags);
	const { t } = useTranslation();
	const [searchTerm, setSearchTerm] = useState('');
	const [searchTypeTerm, setSearchTypeTerm] = useState('');
	const [filterTags, setFilterTags] = useState<Tag[]>([]);
	const dispatch = useAppDispatch();

	const regions = useMemo((): Set<string> => {
		return new Set(tagResource?.Resources?.map(resource => resource.Region) || ['']);
	}, [tagResource.Resources?.length]);

	const types = useMemo((): Set<string> => {
		return new Set(
			tagResource?.Resources?.filter(resource =>
				resource.Type.toUpperCase().includes(searchTypeTerm.toUpperCase())
			)?.map(resource => resource.Type) || ['']
		);
	}, [tagResource.Resources?.length, searchTypeTerm]);

	useEffect(() => {
		if (filterTags?.length) {
			dispatch(getFilteredResources(filterTags));
		} else {
			dispatch(getResources());
		}
	}, [filterTags]);

	const handleSearchTags = (e: ChangeEvent<HTMLInputElement>): void => {
		setSearchTerm(e.target.value);
	};

	const handleSearchTypes = (e: ChangeEvent<HTMLInputElement>): void => {
		setSearchTypeTerm(e.target.value);
	};

	const handleChange = (event: React.ChangeEvent<HTMLInputElement>, field: Field, item: ValueType) => {
		const tag = { key: field.name, value: item.value };
		const existingTag = filterTags?.some(item => item.key === tag.key && item.value === tag.value);
		if (event.target.checked && !existingTag) {
			setFilterTags([...filterTags, tag]);
		} else if (!event.target.checked && existingTag && filterTags) {
			setFilterTags(filterTags.filter(item => item.key !== tag.key && item.value !== tag.value));
		}
	};

	return (
		<Box
			sx={{
				width: '15%',
				height: '100%',
				backgroundColor: '#F9F7F6',
				overflowY: 'scroll',
			}}>
			{fields
				.filter(field => field.name.toUpperCase().includes(searchTerm.toUpperCase()))
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
		</Box>
	);
};

export default InsightFilter;
