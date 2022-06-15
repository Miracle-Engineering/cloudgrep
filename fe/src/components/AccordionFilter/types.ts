import { Field, ValueType } from 'models/Field';

export interface AccordionFilterProps {
	id: string;
	label: string;
	field: Field;
	hasSearch: boolean;
	handleChange: (event: React.ChangeEvent<HTMLInputElement>, field: Field, item: ValueType) => void;
	checkedByDefault: boolean;
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
