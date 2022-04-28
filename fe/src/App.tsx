import './App.css';

import Insights from 'pages/Insights';
import React from 'react';

import logo from './logo.svg';

function App() {
	return (
		<div className="App">
			<header className="App-header">
				<img src={logo} className="App-logo" alt="logo" />
				<span>
					<span>CloudGrep</span>
				</span>
			</header>
			<Insights />
		</div>
	);
}

export default App;
