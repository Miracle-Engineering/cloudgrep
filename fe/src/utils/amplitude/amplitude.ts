import amplitude from 'amplitude-js';

export const initAmplitude = () => {
    amplitude.getInstance().init(process.env.REACT_APP_AMPLITUDE_API_KEY || '');
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