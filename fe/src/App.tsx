import './App.css';
import 'utils/localisation/index';

import ErrorBoundary from 'components/ErrorHandling/ErrorBoundary';
import Insights from 'pages/Insights';
import React from 'react';

function App() {
	return (
		<div className="App">
			<ErrorBoundary>
				<span>
					<span>CloudGrep</span>
				</span>
				<Insights />
			</ErrorBoundary>
		</div>
	);
}

export default App;
