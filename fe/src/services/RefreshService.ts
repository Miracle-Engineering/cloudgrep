import { getResfreshPath } from 'constants/paths';
import ApiClient, { Response } from 'utils/apiClient/ApiClient';

class RefreshService {
	static async refresh(): Promise<Response<{}>> {
		return ApiClient.post<string | undefined, {}>(getResfreshPath());
	}
}

export default RefreshService;
