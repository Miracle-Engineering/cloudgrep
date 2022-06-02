export interface PropertyItemProps {
	keyItem: string;
	value: string;
}

export interface PropertyItemListProps {
	data: Array<Object>;
	renderObjects: (data: Object) => React.ReactNode;
}
