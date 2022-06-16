import { Field, ValueType } from 'models/Field';
import { Tag } from 'models/Tag';

export interface AccordionFilterProps {
	id: string;
	label: string;
	field: Field;
	hasSearch: boolean;
}
export interface AccordionItemProps {
	field: Field;
	handleChange: (event: React.ChangeEvent<HTMLInputElement>, field: Field, item: ValueType) => void;
	isChecked: boolean;
	item: ValueType;
	handleOnly: (item: ValueType) => void;
	handleAll: () => void;
	singleItem: string;
}
