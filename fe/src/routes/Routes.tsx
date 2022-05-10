import Home from 'pages/Home';
import Insights from 'pages/Insights';
import React, { FC } from 'react';
import { Navigate, Route, Routes as SwitchRoutes } from 'react-router-dom';

import { getHomePage } from './helpers';
// ROUTES
import { HOME, INSIGHTS } from './routePaths';

const Routes: FC = () => {
	return (
		<SwitchRoutes>
			<Route path={HOME} element={<Home />} />
			<Route path={INSIGHTS} element={<Insights />} />
			<Route path="*" element={<Navigate to={getHomePage()} />} />
		</SwitchRoutes>
	);
};

export default Routes;
