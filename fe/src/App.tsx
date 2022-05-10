import './App.css';
import 'utils/localisation/index';

import ErrorBoundary from 'components/ErrorHandling/ErrorBoundary';
import Insights from 'pages/Insights';
import React from 'react';

import Header from './components/Header/Header';

function App() {
	return (
		<div className="App">
			<ErrorBoundary>
				<Header />
				<Insights />
			</ErrorBoundary>
		</div>
	);
}

export default App;
