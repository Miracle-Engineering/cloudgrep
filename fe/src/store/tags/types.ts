import { FieldGroup } from 'models/Field';
import { Tag } from 'models/Tag';
import { TagResource } from 'models/TagResource';

export interface TagState {
	tagResource: TagResource;
	filterTags: Tag[];
	fields: FieldGroup[];
	limit: number;
	offset: number;
}

export interface ErrorType {
	response?: { status: string };
	message: string;
}

