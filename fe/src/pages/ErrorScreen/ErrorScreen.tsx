import { BACKGROUND_COLOR } from 'constants/colors';
import { ERROR_PAGE_TEST_ID } from 'constants/testIds';
import { ErrorContainer, TextPlaceholder } from 'pages/ErrorScreen/styles';
import React, { FC } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { LOGIN } from 'routes/routePaths';

const ErrorScreen: FC = () => {
	const navigate = useNavigate();
	const { t } = useTranslation();

	return (
		<ErrorContainer data-testid={ERROR_PAGE_TEST_ID}>
			<TextPlaceholder>{t('SOMETHING_WENT_WRONG')}</TextPlaceholder>
			<button style={{ background: BACKGROUND_COLOR }} onClick={(): void => navigate(LOGIN)}>
				{t('TRY_AGAIN')}
			</button>
		</ErrorContainer>
	);
};

export default ErrorScreen;
