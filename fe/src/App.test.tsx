import { render } from '@testing-library/react';
import React from 'react';
import { Provider } from 'react-redux';

import App from './App';
import { store } from './app/store';

test('renders cloud grep react link', () => {
	const { getByText } = render(
		<Provider store={store}>
			<App />
		</Provider>
	);

	expect(getByText(/CloudGrep/i)).toBeInTheDocument();
});
