import { getResourcesPath } from 'constants/paths';
import { Resource } from 'models/Resource';
import ApiClient, { Response } from 'utils/apiClient/ApiClient';

class ResourceService {
	static async getResources(): Promise<Response<Resource[]>> {
		return ApiClient.get<string | undefined, Resource[]>(getResourcesPath());
	}
}

export default ResourceService;
