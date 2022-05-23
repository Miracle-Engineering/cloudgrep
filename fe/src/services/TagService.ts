import { getFieldsPath } from 'constants/paths';
import ApiClient, { Response } from 'utils/apiClient/ApiClient';

import { Field } from '../models/Field';

class TagService {
	static async getFields(): Promise<Response<Field[]>> {
		return ApiClient.get<string | undefined, Field[]>(getFieldsPath());
	}
}

export default TagService;
