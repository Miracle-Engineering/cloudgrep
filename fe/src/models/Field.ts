export interface ValueType {
	value: string;
	count: number;
}

export interface Field {
	name: string;
	group: string;
	count: number;
	values: ValueType[];
}
