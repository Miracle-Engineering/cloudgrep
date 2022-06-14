import { Resource } from 'models/Resource';
import { Tag } from 'models/Tag';

export interface ResourceState {
	resources: Resource[];
	currentResource?: Resource;
	sideMenuVisible: boolean;
}

export interface FilterResourcesApiParams {
	data: Tag[];
	offset: number;
	limit: number;
}

export interface ResourcesNextPageParams {
	resources: Resource[];
	offset: number;
	limit: number;
}
