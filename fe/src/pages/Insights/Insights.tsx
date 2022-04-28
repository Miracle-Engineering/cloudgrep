import React, { FC } from 'react';

import InsightFilter from './InsightFilter';
import InsightTable from './InsightTable';

const Insights: FC = () => {
	return (
		<>
			<InsightFilter />
			<InsightTable />
		</>
	);
};

export default Insights;
