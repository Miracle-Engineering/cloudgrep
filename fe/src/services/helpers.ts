import { Tag } from 'models/Tag';

const PAGE_LIMIT = 100;
const OR_OPERATOR = '$or';

export const getArrayOfObjects = (data: Tag[]) => {
	return data.map((tag: Tag) => {
		return {
			[tag.key]: tag.value,
		};
	});
};

export const getResourcesRequestData = (data: Tag[], offset = 0) => {
	const filter: {
		[key: string]: string;
	} = {};

	data.forEach((tag: Tag) => {
		filter[tag.key] = tag.value;
	});

	return { filter, limit: PAGE_LIMIT, offset: offset };
};

export const getResourcesFilters = (data: Tag[], offset = 0) => {
	const filter: {
		[key: string]: string | Array<Object>;
	} = {};

	const tagsToRemove: string[] = [];

	data.forEach((tag: Tag) => {
		if (tag.key in filter) {
			const existingValue = filter[tag.key];
			filter[OR_OPERATOR] = [
				...(filter[OR_OPERATOR] || []),
				...(!tagsToRemove.includes(tag.key) ? [{ [tag.key]: existingValue }] : []),
				{ [tag.key]: tag.value },
			];
			tagsToRemove.push(tag.key);
		} else {
			filter[tag.key] = tag.value;
		}
	});

	tagsToRemove.forEach(tagKey => {
		delete filter[tagKey];
	});

	return { filter, limit: PAGE_LIMIT, offset: offset };
};
