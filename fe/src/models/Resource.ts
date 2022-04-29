import { Tag } from './Tag';

export interface Property {
	Name: string;
	Value: string;
}

export interface Resource {
	Type: string;
	Id: string;
	Region: string;
	Properties?: Property[];
	Tags?: Tag[];
}
