import { getResourcesPath } from 'constants/paths';
import { Resources } from 'models/Resource';
import ApiClient, { Response } from 'utils/apiClient/ApiClient';

import { getResourcesFilters } from './helpers';
import { FilterResourcesProps } from './types';
class ResourceService {
	static async getResources(): Promise<Response<Resources>> {
		return ApiClient.get<string | undefined, Resources>(getResourcesPath());
	}

	static async getFilteredResources(filterData: FilterResourcesProps): Promise<Response<Resources>> {
		const requestData = getResourcesFilters(filterData);
		return ApiClient.post<string | undefined, Resources>(getResourcesPath(), requestData);
	}
}

export default ResourceService;
