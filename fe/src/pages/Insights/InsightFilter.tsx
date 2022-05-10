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

import { accordionStyles, labelClasses, overrideSummaryClasses } from './style';

const InsightFilter: FC = () => {
	const { tags, tagResource } = useAppSelector(state => state.tags);
	const { t } = useTranslation();

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
			}}>
			<Accordion>
				<AccordionSummary
					sx={{ backgroundColor: 'rgb(226, 229, 237)', borderRadius: '4px', minHeight: '14px !important' }}
					expandIcon={<ExpandMoreIcon />}
					aria-controls="panel1a-content"
					id="panel1a-header"
					classes={overrideSummaryClasses}>
					<Typography sx={accordionStyles.accordionHeader}>{t('TAGS')}</Typography>
				</AccordionSummary>
				<AccordionDetails>
					<Typography>
						<FormGroup>
							{tags.map((tag: Tag) => (
								<FormControlLabel
									classes={labelClasses}
									key={tag.Key}
									control={<Checkbox size={'small'} defaultChecked />}
									label={tag.Key}
								/>
							))}
						</FormGroup>
					</Typography>
				</AccordionDetails>
			</Accordion>
			<Accordion>
				<AccordionSummary
					sx={{ backgroundColor: 'rgb(226, 229, 237)', borderRadius: '4px', minHeight: '14px !important' }}
					expandIcon={<ExpandMoreIcon />}
					aria-controls="panel2a-content"
					id="panel2a-header"
					classes={overrideSummaryClasses}>
					<Typography sx={accordionStyles.accordionHeader}>{t('REGIONS')}</Typography>
				</AccordionSummary>
				<AccordionDetails>
					<Typography>
						<FormGroup>
							{regions &&
								Array.from(regions).map((region: string) => (
									<FormControlLabel
										classes={labelClasses}
										key={region}
										control={<Checkbox size={'small'} defaultChecked />}
										label={region}
									/>
								))}
						</FormGroup>
					</Typography>
				</AccordionDetails>
			</Accordion>
			<Accordion>
				<AccordionSummary
					sx={{ backgroundColor: 'rgb(226, 229, 237)', borderRadius: '4px', minHeight: '14px !important' }}
					expandIcon={<ExpandMoreIcon />}
					aria-controls="panel2a-content"
					id="panel2a-header"
					classes={overrideSummaryClasses}>
					<Typography sx={accordionStyles.accordionHeader}>{t('TYPES')}</Typography>
				</AccordionSummary>
				<AccordionDetails>
					<Typography>
						<FormGroup sx={accordionStyles.accordionDetails}>
							{types &&
								Array.from(types).map((type: string) => (
									<FormControlLabel
										classes={labelClasses}
										key={type}
										control={<Checkbox size={'small'} defaultChecked />}
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
