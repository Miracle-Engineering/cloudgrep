import React from 'react';

const SCROLLED_BOTTOM_VALUE = 0.9;
const INFINITE_SCROLL_VALUE = 0.75;

export const isScrolledToBottom = (element: React.UIEvent<HTMLElement>): boolean => {
	const target = element.target as HTMLElement;
	const height = target.scrollHeight - target.clientHeight;
	const scrolled = target.scrollTop / height;
	// Used as scroll offset check. Max scrolled value is 1.
	return scrolled > SCROLLED_BOTTOM_VALUE;
};

export const isScrolledForInfiniteScroll = (element: React.UIEvent<HTMLElement>): boolean => {
	const target = element.target as HTMLElement;
	const height = target.scrollHeight - target.clientHeight;
	const scrolled = target.scrollTop / height;
	// Used as scroll offset check. Max scrolled value is 1.
	return scrolled > INFINITE_SCROLL_VALUE;
};
