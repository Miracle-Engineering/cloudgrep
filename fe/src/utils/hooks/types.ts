import { Dispatch, SetStateAction } from 'react';

export interface PaginationType {
	next: () => void;
	prev: () => void;
	setCurrentPage: Dispatch<SetStateAction<number>>;
	jump: (page: number) => void;
	currentPage: number;
	maxPage: number;
}
