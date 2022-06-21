import { getEngineStatusPath, getResfreshPath } from 'constants/paths';
import { EngineStatus } from 'models/EngineStatus';
import ApiClient, { Response } from 'utils/apiClient/ApiClient';

class RefreshService {
	static async refresh(): Promise<Response<{}>> {
		return ApiClient.post<string | undefined, {}>(getResfreshPath());
	}

	static async getStatus(): Promise<Response<EngineStatus>> {
		return ApiClient.get<string | undefined, EngineStatus>(getEngineStatusPath());
	}
}

export default RefreshService;
