import { MockResource } from './Resource';
import { MockTag } from './Tag';

export interface TagResource {
	Tags: MockTag[];
	Resources: MockResource[];
}
