import { getFieldsPath } from 'constants/paths';
import ApiClient, { Response } from 'utils/apiClient/ApiClient';

import { FieldGroup } from '../models/Field';

class TagService {
	static async getFields(): Promise<Response<FieldGroup[]>> {
		return ApiClient.get<string | undefined, FieldGroup[]>(getFieldsPath());
	}
}

export default TagService;
