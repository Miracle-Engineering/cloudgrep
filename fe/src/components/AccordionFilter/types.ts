import { Field, ValueType } from 'models/Field';

export interface AccordionFilterProps {
	id: string;
	label: string;
	field: Field;
	hasSearch: boolean;
	handleChange: (event: React.ChangeEvent<HTMLInputElement>, field: Field, item: ValueType) => void;
}
