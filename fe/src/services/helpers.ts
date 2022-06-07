import { Tag } from 'models/Tag';

const PAGE_LIMIT = 100;
const OR_OPERATOR = '$or';
const AND_OPERATOR = '$and';

export const getArrayOfObjects = (data: Tag[]) => {
	return data.map((tag: Tag) => {
		return {
			[tag.key]: tag.value,
		};
	});
};

export const getResourcesFilters = (data: Tag[], offset = 0, limit = PAGE_LIMIT) => {
	const filter: {
		[key: string]: Array<Object>;
	} = {};

	const uniqueTags = new Set(data.map((tag: Tag) => tag.key));

	uniqueTags.forEach((key: string) => {
		const currentTags = data.filter((tag: Tag) => tag.key === key);
		let currentFilters: { [key: string]: string }[] = [];

		currentTags.forEach((tag: Tag) => {
			currentFilters = [...currentFilters, { [tag.key]: tag.value }];
		});

		filter[AND_OPERATOR] = [...(filter[AND_OPERATOR] || []), { [OR_OPERATOR]: currentFilters }];
	});

	return { filter, limit: limit, offset: offset };
};
