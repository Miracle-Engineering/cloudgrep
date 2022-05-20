import amplitude from 'amplitude-js';

export const initAmplitude = () => {
	amplitude.getInstance().init('2b0167b9ea1dacf8f0dae96326abd879');
	amplitude.getInstance().logEvent('PAGE LOAD');
};

export const setAmplitudeUserDevice = (installationToken: string) => {
	amplitude.getInstance().setDeviceId(installationToken);
};

export const setAmplitudeUserId = (userId: string | null) => {
	amplitude.getInstance().setUserId(userId);
};

export const setAmplitudeUserProperties = (properties: any) => {
	amplitude.getInstance().setUserProperties(properties);
};

export const sendAmplitudeData = (eventType: string, eventProperties: any) => {
	amplitude.getInstance().logEvent(eventType, eventProperties);
};
