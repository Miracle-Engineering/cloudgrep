import Home from 'pages/Home';
import Login from 'pages/Login';
import React, { FC } from 'react';
import { Navigate, Route, Routes as SwitchRoutes } from 'react-router-dom';

import { getHomePage } from './helpers';
// ROUTES
import { HOME, LOGIN } from './routePaths';

const Routes: FC = () => {
	return (
		<SwitchRoutes>
			<Route path={LOGIN} element={<Login />} />
			<Route path={HOME} element={<Home />} />
			<Route path="*" element={<Navigate to={getHomePage()} />} />
		</SwitchRoutes>
	);
};

export default Routes;
