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
import React, { FC, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { useAppSelector } from 'store/hooks';

import { accordionStyles, overrideSummaryClasses } from './style';

const InsightFilter: FC = () => {
	const { tags, tagResource } = useAppSelector(state => state.tags);
	const { t } = useTranslation();
	const [expanded, setExpanded] = React.useState<string | false>('tagsPanel');

	const handleChange = (panel: string) => (event: React.SyntheticEvent, newExpanded: boolean) => {
		setExpanded(newExpanded ? panel : false);
	};

	const regions = useMemo((): Set<string> => {
		return new Set(tagResource?.Resources?.map(resource => resource.Region) || ['']);
	}, [tagResource.Resources?.length]);

	const types = useMemo((): Set<string> => {
		return new Set(tagResource?.Resources?.map(resource => resource.Type) || ['']);
	}, [tagResource.Resources?.length]);

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
			<Accordion expanded={expanded === 'tagsPanel'} onChange={handleChange('tagsPanel')}>
				<AccordionSummary
					sx={{ backgroundColor: 'rgb(226, 229, 237)', borderRadius: '4px', minHeight: '14px !important' }}
					expandIcon={<ExpandMoreIcon />}
					aria-controls="panel1a-content"
					id="panel1a-header"
					classes={overrideSummaryClasses}>
					<Typography sx={accordionStyles.accordionHeader}>{t('TAGS')}</Typography>
				</AccordionSummary>
				<AccordionDetails sx={accordionStyles.accordionDetails}>
					<Typography sx={accordionStyles.accordionDetails}>
						<FormGroup sx={accordionStyles.accordionDetails}>
							{tags.map((tag: Tag) => (
								<FormControlLabel
									sx={accordionStyles.accordionDetails}
									key={tag.Key}
									control={<Checkbox defaultChecked />}
									label={tag.Key}
								/>
							))}
						</FormGroup>
					</Typography>
				</AccordionDetails>
			</Accordion>
			<Accordion expanded={expanded === 'regionPanel'} onChange={handleChange('regionPanel')}>
				<AccordionSummary
					sx={{ backgroundColor: 'rgb(226, 229, 237)', borderRadius: '4px', minHeight: '14px !important' }}
					expandIcon={<ExpandMoreIcon />}
					aria-controls="panel2a-content"
					id="panel2a-header"
					classes={overrideSummaryClasses}>
					<Typography sx={accordionStyles.accordionHeader}>{t('REGIONS')}</Typography>
				</AccordionSummary>
				<AccordionDetails sx={accordionStyles.accordionDetails}>
					<Typography sx={accordionStyles.accordionDetails}>
						<FormGroup sx={accordionStyles.accordionDetails}>
							{regions &&
								Array.from(regions).map((region: string) => (
									<FormControlLabel
										sx={accordionStyles.accordionDetails}
										key={region}
										control={<Checkbox defaultChecked />}
										label={region}
									/>
								))}
						</FormGroup>
					</Typography>
				</AccordionDetails>
			</Accordion>
			<Accordion expanded={expanded === 'typePanel'} onChange={handleChange('typePanel')}>
				<AccordionSummary
					sx={{ backgroundColor: 'rgb(226, 229, 237)', borderRadius: '4px', minHeight: '14px !important' }}
					expandIcon={<ExpandMoreIcon />}
					aria-controls="panel2a-content"
					id="panel2a-header"
					classes={overrideSummaryClasses}>
					<Typography sx={accordionStyles.accordionHeader}>{t('TYPES')}</Typography>
				</AccordionSummary>
				<AccordionDetails sx={accordionStyles.accordionDetails}>
					<Typography sx={accordionStyles.accordionDetails}>
						<FormGroup sx={accordionStyles.accordionDetails}>
							{types &&
								Array.from(types).map((type: string) => (
									<FormControlLabel
										sx={accordionStyles.accordionDetails}
										key={type}
										control={<Checkbox defaultChecked />}
										label={type}
									/>
								))}
						</FormGroup>
					</Typography>
				</AccordionDetails>
			</Accordion>
		</Box>
	);
};

export default InsightFilter;
