import './App.css';
import 'utils/localisation/index';

import ErrorBoundary from 'components/ErrorHandling/ErrorBoundary';
import React, { FC, useEffect } from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import Routes from 'routes/Routes';
import { useAppDispatch, useAppSelector } from 'store/hooks';
import { getResources } from 'store/resources/thunks';

import Header from './components/Header';

const App: FC = () => {
	const dispatch = useAppDispatch();
	const { fields } = useAppSelector(state => state.tags);

	useEffect(() => {
		if (!fields?.length) {
			dispatch(getResources());
		}
	}, [dispatch, fields?.length]);

	return (
		<div className="App">
			<Router>
				<ErrorBoundary>
					<Header />
					<Routes />
				</ErrorBoundary>
			</Router>
		</div>
	);
};

export default App;
