import { Tag } from 'models/Tag';
import { TagResource } from 'models/TagResource';

export interface TagState {
	tagResource: TagResource;
	tags: Tag[];
}

export interface ErrorType {
	response?: { status: string };
	message: string;
}
