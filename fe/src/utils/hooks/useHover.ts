import { useCallback, useRef, useState } from 'react';

const useHover = <T extends HTMLElement>(): [(node?: T | null) => void, boolean] => {
	const [hovered, setHovered] = useState(false);

	const handleMouseOver = useCallback(() => setHovered(true), []);
	const handleMouseLeave = useCallback(() => setHovered(false), []);

	const ref = useRef<T | null>();

	const callbackRef = useCallback<(node?: null | T) => void>(
		node => {
			if (ref.current) {
				ref.current.removeEventListener('mouseover', handleMouseOver);
				ref.current.removeEventListener('mouseleave', handleMouseLeave);
			}

			ref.current = node;

			if (ref.current) {
				ref.current.addEventListener('mouseover', handleMouseOver);
				ref.current.addEventListener('mouseleave', handleMouseLeave);
			}
		},
		[handleMouseOver, handleMouseLeave]
	);

	return [callbackRef, hovered];
};

export default useHover;
