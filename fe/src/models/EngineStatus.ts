export enum EngineStatusEnum {
	SUCCESS = 'success',
	FETCHING = 'fetching',
	FAILED = 'failed',
}

export interface EngineStatus {
	resourceType: string;
	errorMessage: string;
	status: EngineStatusEnum;
	fetchedAt: string;
}
