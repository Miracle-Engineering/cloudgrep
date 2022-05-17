import './index.css';

import React from 'react';
import { createRoot } from 'react-dom/client';
import ReactGA from 'react-ga';
import { Provider } from 'react-redux';
import { store } from 'store/store';
import { initAmplitude } from 'utils/amplitude/amplitude';

// Add Amplitude integration
if (process.env.REACT_APP_ENABLE_AMPLITUDE === 'true' && process.env.NODE_ENV !== 'development') {
	initAmplitude();
}

// Add Google Analytics integration
if (process.env.REACT_APP_ENABLE_GA === 'true' && process.env.NODE_ENV !== 'development') {
	ReactGA.initialize('UA-228613342');
	ReactGA.pageview(window.location.pathname + window.location.search);
}

import App from './App';

const container = document.getElementById('root');
const root = createRoot(container!);

root.render(
	<React.StrictMode>
		<Provider store={store}>
			<App />
		</Provider>
	</React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
// reportWebVitals();
