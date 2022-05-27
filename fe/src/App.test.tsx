import { render } from '@testing-library/react';
import React from 'react';
import { Provider } from 'react-redux';
import { store } from 'store/store';

import App from './App';

test('renders cloud grep react link', () => {
	const { getByText } = render(
		<Provider store={store}>
			<App />
		</Provider>
	);

	expect(getByText(/Contact/i)).toBeInTheDocument();
});
