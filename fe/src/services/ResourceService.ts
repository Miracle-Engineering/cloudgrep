import { getFilteredResourcesPath, getResourcesPath } from 'constants/paths';
import { Resource } from 'models/Resource';
import { MockTag } from 'models/Tag';
import ApiClient, { Response } from 'utils/apiClient/ApiClient';
class ResourceService {
	static async getResources(): Promise<Response<Resource[]>> {
		return ApiClient.get<string | undefined, Resource[]>(getResourcesPath());
	}

	static async getFilteredResources(data: MockTag): Promise<Response<Resource[]>> {
		const path = getFilteredResourcesPath(data);
		return ApiClient.get<string | undefined, Resource[]>(path);
	}
}

export default ResourceService;
