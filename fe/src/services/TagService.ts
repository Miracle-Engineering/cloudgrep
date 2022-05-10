import { getTagsPath } from 'constants/paths';
import ApiClient, { Response } from 'utils/apiClient/ApiClient';

import { TagResource } from '../models/TagResource';

class TagService {
	static async getTags(): Promise<Response<TagResource>> {
		return ApiClient.get<string | undefined, TagResource>(getTagsPath());
	}
}

export default TagService;
