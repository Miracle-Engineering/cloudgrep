import { AND_OPERATOR, OR_OPERATOR } from 'constants/globals';
import { Tag } from 'models/Tag';

import { FilterResourcesProps } from './types';

export const getArrayOfObjects = (data: Tag[]) => {
	return data.map((tag: Tag) => {
		return {
			[tag.key]: tag.value,
		};
	});
};

export const getResourcesFilters = (filterData: FilterResourcesProps) => {
	const { data, offset, limit, order, orderBy } = filterData;

	const filter: {
		[key: string]: Array<Object>;
	} = {};

	const uniqueTags = new Set(data.map((tag: Tag) => tag.key));

	uniqueTags.forEach((key: string) => {
		const currentTags = data.filter((tag: Tag) => tag.key === key);
		const currentFilters = getArrayOfObjects(currentTags);
		filter[AND_OPERATOR] = [...(filter[AND_OPERATOR] || []), { [OR_OPERATOR]: [...currentFilters] }];
	});

	const sortByProperty = order === 'desc' ? '-' + orderBy : orderBy;

	return { filter, limit: limit, offset: offset, sort: orderBy ? [sortByProperty?.toLowerCase()] : undefined };
};
