import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import Accordion from '@mui/material/Accordion';
import AccordionDetails from '@mui/material/AccordionDetails';
import AccordionSummary from '@mui/material/AccordionSummary';
import Box from '@mui/material/Box';
import Checkbox from '@mui/material/Checkbox';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormGroup from '@mui/material/FormGroup';
import Typography from '@mui/material/Typography';
import { Tag } from 'models/Tag';
import React, { FC } from 'react';
import { useTranslation } from 'react-i18next';
import { useAppSelector } from 'store/hooks';

const InsightFilter: FC = () => {
	const { tags } = useAppSelector(state => state.tags);
	const { t } = useTranslation();

	return (
		<Box
			sx={{
				width: '20%',
				height: '100%',
				'&:hover': {
					backgroundColor: 'primary.main',
					opacity: [0.9, 0.8, 0.7],
				},
			}}>
			<Accordion>
				<AccordionSummary expandIcon={<ExpandMoreIcon />} aria-controls="panel1a-content" id="panel1a-header">
					<Typography>{t('TAGS')}</Typography>
				</AccordionSummary>
				<AccordionDetails>
					<Typography>
						<FormGroup>
							{tags.map((tag: Tag) => (
								<FormControlLabel key={tag.Key} control={<Checkbox defaultChecked />} label={tag.Key} />
							))}
						</FormGroup>
					</Typography>
				</AccordionDetails>
			</Accordion>
		</Box>
	);
};

export default InsightFilter;
