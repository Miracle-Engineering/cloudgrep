import './index.css';

import React from 'react';
import { createRoot } from 'react-dom/client';
import { Provider } from 'react-redux';
import { store } from 'store/store';
import { initAmplitude } from 'utils/amplitude/amplitude';

// Add Amplitude integration
if (process.env.REACT_APP_ENABLE_AMPLITUDE === 'true' && process.env.NODE_ENV !== 'development') {
	initAmplitude();
}

// Add Google Analytics integration
if (process.env.REACT_APP_ENABLE_GA === 'true' && process.env.NODE_ENV !== 'development') {
	const gtagScript = document.createElement('script'); // Make a script DOM node

	// Google Analytics (ga) script is added to public folder
	gtagScript.src = `${process.env.PUBLIC_URL?.toString() || ''}${process.env.REACT_APP_JS_PATH_PREFIX}/ga.js`; // Set script source path

	const gaScript = document.createElement('script');
	gaScript.async = true;
	gaScript.src = 'https://www.googletagmanager.com/gtag/js?id=G-J94VVQXGPT';

	if (document.head) {
		document.head.appendChild(gaScript);
		document.head.appendChild(gtagScript);
	}
}

import App from './App';

const container = document.getElementById('root');
const root = createRoot(container!);

root.render(
	<Provider store={store}>
		<App />
	</Provider>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
// reportWebVitals();
