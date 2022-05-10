import { Tag } from './Tag';

export interface Property {
	name: string;
	value: string;
}

export interface Resource {
	type: string;
	id: string;
	region: string;
	properties?: Property[];
	tags?: Tag[];
}

export interface MockResource {
	Type: string;
	Id: string;
	Region: string;
}
