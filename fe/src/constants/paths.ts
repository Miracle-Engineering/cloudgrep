import { Tag } from 'models/Tag';

export const GET_FIELDS_PATH = 'fields';
export const GET_RESOURCES_PATH = 'resources';

export const getFieldsPath = (): string => `${process.env.REACT_APP_MOCK_API_URL}${GET_FIELDS_PATH}`;
export const getResourcesPath = (): string => `${process.env.REACT_APP_API_URL}${GET_RESOURCES_PATH}`;
export const getFilteredResourcesPath = (data: Tag[]): string => {
	let tagsPath = '';

	data.forEach((tag: Tag, index: number) => {
		tagsPath += `${index === 0 ? '?' : '&'}tags[${tag.key}]=${tag.value}`;
	});

	return `${process.env.REACT_APP_API_URL}${GET_RESOURCES_PATH}${tagsPath}`;
};
