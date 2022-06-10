// Filtering constants
export const CHECKED_BY_DEFAULT = true; // all filters are checked by default on page load
export const SEARCH_ELEMENTS_NUMBER = 3; // minimum number of element per tag for search to appear
export const OR_OPERATOR = '$or';
export const AND_OPERATOR = '$and';

// Pagination constants
export const PAGE_LENGTH = 25;
export const PAGE_START = 0;
export const INFINITE_SCROLL_VALUE = 0.75; // If 75% of table is scrolled then fetch next page of results
export const TOTAL_RECORDS = 10000; // todo update with real value from API response when available