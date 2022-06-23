import { FieldGroup } from './Field';
import { Tag } from './Tag';

export interface Property {
	name: string;
	value: string;
}

export interface Resource {
	type: string;
	id: string;
	region: string;
	rawData: Object;
	tags?: Tag[];
}

export interface MockResource {
	Type: string;
	Id: string;
	Region: string;
}

export interface Resources {
	count: number;
	resources: Resource[];
	fieldGroups: FieldGroup[];
}
