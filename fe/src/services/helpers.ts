import { AND_OPERATOR, OR_OPERATOR, PAGE_LENGTH } from 'constants/globals';
import { Tag } from 'models/Tag';

export const getArrayOfObjects = (data: Tag[]) => {
	return data.map((tag: Tag) => {
		return {
			[tag.key]: tag.value,
		};
	});
};

export const getResourcesFilters = (data: Tag[], offset = 0, limit = PAGE_LENGTH) => {
	const filter: {
		[key: string]: Array<Object>;
	} = {};

	const uniqueTags = new Set(data.map((tag: Tag) => tag.key));

	uniqueTags.forEach((key: string) => {
		const currentTags = data.filter((tag: Tag) => tag.key === key);
		const currentFilters = getArrayOfObjects(currentTags);
		filter[AND_OPERATOR] = [...(filter[AND_OPERATOR] || []), { [OR_OPERATOR]: [...currentFilters] }];
	});

	return { filter, limit: limit, offset: offset };
};
