import { INFINITE_SCROLL_VALUE } from 'constants/globals';
import React from 'react';

export const isScrolledForInfiniteScroll = (element: React.UIEvent<HTMLElement>): boolean => {
	const target = element.target as HTMLElement;
	const height = target.scrollHeight - target.clientHeight;
	const scrolled = target.scrollTop / height;
	// Used as scroll offset check. Max scrolled value is 1.
	return scrolled > INFINITE_SCROLL_VALUE;
};
