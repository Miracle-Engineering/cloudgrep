import axios, { AxiosInstance } from 'axios';

const configureAPI = (): AxiosInstance => {
	const clientAPI: AxiosInstance = axios.create({
		responseType: 'json',
		headers: {
			'Content-Type': 'application/json',
		},
	});

	return clientAPI;
};

export default configureAPI;
