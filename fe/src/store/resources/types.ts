import { Resource } from 'models/Resource';

export interface ResourceState {
	resources: Resource[];
	currentResource?: Resource;
	sideMenuVisible: boolean;
}
