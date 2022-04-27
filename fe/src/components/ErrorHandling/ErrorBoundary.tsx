import ErrorScreen from 'pages/ErrorScreen/ErrorScreen';
import React, { ReactNode } from 'react';

type ErrorObject = {
	hasError: boolean;
};

class ErrorBoundary extends React.Component<{ children: ReactNode }> {
	state = { hasError: false };

	// parameters: error
	static getDerivedStateFromError(): ErrorObject {
		// Update state so the next render will show the fallback UI.
		return { hasError: true };
	}

	// parameters: error, errorInfo
	componentDidCatch(error: Error, errorInfo: React.ErrorInfo): void {
		// Log the error
		if (error) {
			// here errors should be logged instead
			// eslint-disable-next-line no-console
			console.log(error, errorInfo);
		}
	}

	render(): React.ReactNode {
		if (this.state.hasError) {
			return <ErrorScreen />;
		}

		return this.props.children;
	}
}

export default ErrorBoundary;
