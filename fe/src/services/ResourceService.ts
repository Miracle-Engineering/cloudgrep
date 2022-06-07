import { getResourcesPath } from 'constants/paths';
import { Resource } from 'models/Resource';
import { Tag } from 'models/Tag';
import ApiClient, { Response } from 'utils/apiClient/ApiClient';

import { getResourcesFilters } from './helpers';
class ResourceService {
	static async getResources(): Promise<Response<Resource[]>> {
		return ApiClient.get<string | undefined, Resource[]>(getResourcesPath());
	}

	static async getFilteredResources(data: Tag[]): Promise<Response<Resource[]>> {
		const requestData = getResourcesFilters(data);
		return ApiClient.post<string | undefined, Resource[]>(getResourcesPath(), requestData);
	}
}

export default ResourceService;
