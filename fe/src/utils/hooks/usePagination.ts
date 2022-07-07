import { useState } from 'react';
import { PaginationType } from 'utils/hooks/types';

function usePagination(itemsPerPage: number, maxData: number): PaginationType {
	const [currentPage, setCurrentPage] = useState(1);
	const maxPage = Math.ceil(maxData / itemsPerPage);

	function next(): void {
		setCurrentPage(nextPage => Math.min(nextPage + 1, maxPage));
	}

	function prev(): void {
		setCurrentPage(previousPage => Math.max(previousPage - 1, 1));
	}

	function jump(page: number): void {
		const pageNumber = Math.max(1, page);
		setCurrentPage(Math.min(pageNumber, maxPage));
	}

	return { next, prev, jump, currentPage, maxPage, setCurrentPage };
}

export default usePagination;
