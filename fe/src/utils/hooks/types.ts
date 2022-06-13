export interface PaginationType {
	next: () => void;
	prev: () => void;
	jump: (page: number) => void;
	currentPage: number;
	maxPage: number;
}
