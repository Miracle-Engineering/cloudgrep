import './App.css';
import 'utils/localisation/index';

import ErrorBoundary from 'components/ErrorHandling/ErrorBoundary';
import React from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import Routes from 'routes/Routes';

import Header from './components/Header';

function App() {
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
}

export default App;
