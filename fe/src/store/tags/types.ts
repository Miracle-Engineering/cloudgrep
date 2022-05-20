import { MockTag } from 'models/Tag';
import { TagResource } from 'models/TagResource';

export interface TagState {
	tagResource: TagResource;
	tags: MockTag[];
}

export interface ErrorType {
	response?: { status: string };
	message: string;
}
