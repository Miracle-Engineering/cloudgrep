import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import Accordion from '@mui/material/Accordion';
import AccordionDetails from '@mui/material/AccordionDetails';
import AccordionSummary from '@mui/material/AccordionSummary';
import Box from '@mui/material/Box';
import Checkbox from '@mui/material/Checkbox';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormGroup from '@mui/material/FormGroup';
import Typography from '@mui/material/Typography';
import React, { FC } from 'react';

const InsightFilter: FC = () => {
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
					<Typography>Source</Typography>
				</AccordionSummary>
				<AccordionDetails>
					<Typography>
						<FormGroup>
							<FormControlLabel control={<Checkbox defaultChecked />} label="Source 1" />
							<FormControlLabel control={<Checkbox />} label="Source 2" />
						</FormGroup>
					</Typography>
				</AccordionDetails>
			</Accordion>
			<Accordion>
				<AccordionSummary expandIcon={<ExpandMoreIcon />} aria-controls="panel2a-content" id="panel2a-header">
					<Typography>Host</Typography>
				</AccordionSummary>
				<AccordionDetails>
					<Typography>
						<FormGroup>
							<FormControlLabel control={<Checkbox defaultChecked />} label="Host 1" />
							<FormControlLabel control={<Checkbox />} label="Host 2" />
						</FormGroup>
					</Typography>
				</AccordionDetails>
			</Accordion>
			<Accordion>
				<AccordionSummary expandIcon={<ExpandMoreIcon />} aria-controls="panel3a-content" id="panel3a-header">
					<Typography>Service</Typography>
				</AccordionSummary>
				<AccordionDetails>
					<Typography>
						<FormGroup>
							<FormControlLabel control={<Checkbox defaultChecked />} label="Service 1" />
							<FormControlLabel control={<Checkbox />} label="Service 2" />
						</FormGroup>
					</Typography>
				</AccordionDetails>
			</Accordion>
			<Accordion>
				<AccordionSummary expandIcon={<ExpandMoreIcon />} aria-controls="panel3a-content" id="panel3a-header">
					<Typography>Status</Typography>
				</AccordionSummary>
				<AccordionDetails>
					<Typography>
						<FormGroup>
							<FormControlLabel control={<Checkbox defaultChecked />} label="Error" />
							<FormControlLabel control={<Checkbox />} label="Warn" />
							<FormControlLabel control={<Checkbox />} label="Info" />
						</FormGroup>
					</Typography>
				</AccordionDetails>
			</Accordion>
		</Box>
	);
};

export default InsightFilter;
