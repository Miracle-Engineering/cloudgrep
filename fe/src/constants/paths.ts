import { MockTag } from 'models/Tag';

export const GET_TAGS_PATH = 'tags';
export const GET_RESOURCES_PATH = 'resources';

export const getTagsPath = (): string => `${process.env.REACT_APP_MOCK_API_URL}${GET_TAGS_PATH}`;
export const getResourcesPath = (): string => `${process.env.REACT_APP_API_URL}${GET_RESOURCES_PATH}`;
export const getFilteredResourcesPath = (data: MockTag[]): string => {
	let tagsPath = '';

	data.forEach((tag: MockTag, index: number) => {
		tagsPath += `${index === 0 ? '?' : '&'}tags[${tag.Key}]=${tag.Value}`;
	});

	return `${process.env.REACT_APP_API_URL}${GET_RESOURCES_PATH}${tagsPath}`;
};
