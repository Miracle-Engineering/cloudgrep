import { Tag } from 'models/Tag';

export interface FilterResourcesProps {
	data: Tag[];
	offset: number;
	limit: number;
	order?: string;
	orderBy?: string;
}
