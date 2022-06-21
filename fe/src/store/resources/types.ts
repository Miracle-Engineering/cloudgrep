import { Resource } from 'models/Resource';
import { Tag } from 'models/Tag';

export interface ResourceState {
	resources: Resource[];
	currentResource?: Resource;
	sideMenuVisible: boolean;
	count: number;
}

export interface FilterResourcesApiParams {
	data: Tag[];
	limit: number;
	offset: number;
}

export interface ResourcesNextPageParams {
	resources: Resource[];
	limit: number;
	offset: number;
}
