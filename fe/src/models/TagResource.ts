import { Resource } from './Resource';
import { Tag } from './Tag';

export interface TagResource {
	Tags: Tag[];
	Resources: Resource[];
}
