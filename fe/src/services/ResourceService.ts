import { getResourcesPath } from 'constants/paths';
import { Resources } from 'models/Resource';
import { Tag } from 'models/Tag';
import ApiClient, { Response } from 'utils/apiClient/ApiClient';

import { getResourcesFilters } from './helpers';
class ResourceService {
	static async getResources(): Promise<Response<Resources>> {
		return ApiClient.get<string | undefined, Resources>(getResourcesPath());
	}

	static async getFilteredResources(data: Tag[], offset: number, limit: number): Promise<Response<Resources>> {
		const requestData = getResourcesFilters(data, offset, limit);
		return ApiClient.post<string | undefined, Resources>(getResourcesPath(), requestData);
	}
}

export default ResourceService;
