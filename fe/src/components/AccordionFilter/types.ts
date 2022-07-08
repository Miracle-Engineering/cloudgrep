import { Field, FieldGroup, ValueType } from 'models/Field';

export interface AccordionFilterProps {
	id: string;
	label: string;
	field: Field;
	fieldGroup: FieldGroup;
	hasSearch: boolean;
}
export interface AccordionItemProps {
	field: Field;
	handleChange: (
		event: React.ChangeEvent<HTMLInputElement>,
		fieldGroup: FieldGroup,
		field: Field,
		item: ValueType
	) => void;
	isChecked: boolean;
	item: ValueType;
	handleOnly: (item: ValueType) => void;
	handleAll: () => void;
	singleItem: string;
	fieldGroup: FieldGroup;
}
