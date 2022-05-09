import { MockResource } from './Resource';
import { Tag } from './Tag';

export interface TagResource {
	Tags: Tag[];
	Resources: MockResource[];
}
