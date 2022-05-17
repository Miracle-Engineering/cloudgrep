import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import Accordion from '@mui/material/Accordion';
import AccordionDetails from '@mui/material/AccordionDetails';
import AccordionSummary from '@mui/material/AccordionSummary';
import Box from '@mui/material/Box';
import Checkbox from '@mui/material/Checkbox';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormGroup from '@mui/material/FormGroup';
import Typography from '@mui/material/Typography';
import SearchInput from 'components/SearchInput/SearchInput';
import { MockTag } from 'models/Tag';
import React, { ChangeEvent, FC, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useAppSelector } from 'store/hooks';

import { accordionStyles, filterHeader, labelClasses, overrideSummaryClasses } from './style';

const InsightFilter: FC = () => {
	const { tags, tagResource } = useAppSelector(state => state.tags);
	const { t } = useTranslation();
	const [searchTerm, setSearchTerm] = useState('');
	const [searchTypeTerm, setSearchTypeTerm] = useState('');

	const regions = useMemo((): Set<string> => {
		return new Set(tagResource?.Resources?.map(resource => resource.Region) || ['']);
	}, [tagResource.Resources?.length]);

	const types = useMemo((): Set<string> => {
		return new Set(
			tagResource?.Resources?.filter(resource =>
				resource.Type.toUpperCase().includes(searchTypeTerm.toUpperCase())
			)?.map(resource => resource.Type) || ['']
		);
	}, [tagResource.Resources?.length, searchTypeTerm]);

	const handleSearchTags = (e: ChangeEvent<HTMLInputElement>): void => {
		setSearchTerm(e.target.value);
	};

	const handleSearchTypes = (e: ChangeEvent<HTMLInputElement>): void => {
		setSearchTypeTerm(e.target.value);
	};

	return (
		<Box
			sx={{
				width: '15%',
				height: '100%',
				backgroundColor: '#F9F7F6',
			}}>
			<Accordion>
				<AccordionSummary
					sx={filterHeader}
					expandIcon={<ExpandMoreIcon sx={{ color: 'white' }} />}
					aria-controls="panel1a-content"
					id="panel1a-header"
					classes={overrideSummaryClasses}>
					<Typography sx={accordionStyles.accordionHeader}>{t('TAGS')}</Typography>
				</AccordionSummary>
				<AccordionDetails>
					<SearchInput onChange={handleSearchTags} />
					<Typography mt={1}>
						<FormGroup>
							{tags
								.filter(tag => tag.Key.toUpperCase().includes(searchTerm.toUpperCase()))
								.map((tag: MockTag) => (
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
					sx={filterHeader}
					expandIcon={<ExpandMoreIcon sx={{ color: 'white' }} />}
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
					sx={filterHeader}
					expandIcon={<ExpandMoreIcon sx={{ color: 'white' }} />}
					aria-controls="panel2a-content"
					id="panel2a-header"
					classes={overrideSummaryClasses}>
					<Typography sx={accordionStyles.accordionHeader}>{t('TYPES')}</Typography>
				</AccordionSummary>
				<AccordionDetails>
					<SearchInput onChange={handleSearchTypes} />
					<Typography mt={1}>
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
			{tags
				.filter(tag => tag.Key.toUpperCase().includes(searchTerm.toUpperCase()))
				.map((tag: MockTag, index: number) => (
					<Accordion key={`${tag.Key}${index}`}>
						<AccordionSummary
							sx={filterHeader}
							expandIcon={<ExpandMoreIcon sx={{ color: 'white' }} />}
							aria-controls={`accordion${tag.Key}${index}`}
							id={`accordion${tag.Key}${index}`}
							classes={overrideSummaryClasses}>
							<Typography sx={accordionStyles.accordionHeader}>{tag.Key}</Typography>
						</AccordionSummary>
						<AccordionDetails>
							<SearchInput onChange={handleSearchTags} />
							<Typography mt={1}>
								<FormGroup>
									<FormControlLabel
										classes={labelClasses}
										key={tag.Value}
										control={<Checkbox size={'small'} defaultChecked />}
										label={tag.Value}
									/>
								</FormGroup>
							</Typography>
						</AccordionDetails>
					</Accordion>
				))}
		</Box>
	);
};

export default InsightFilter;
