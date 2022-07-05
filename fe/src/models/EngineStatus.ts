export enum EngineStatusEnum {
	SUCCESS = 'success',
	FETCHING = 'fetching',
	FAILED = 'failed',
}

export interface EngineStatus {
	eventType: string;
	error: string;
	status: EngineStatusEnum;
	createdAt: string;
	updatedAt: string;
	runId: string;
	providerName: string;
}
