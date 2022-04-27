import { LocalStorageKey } from './types';
class LocalStorageHelper {
	static set<T>(key: LocalStorageKey | string, value: T): void {
		window.localStorage.setItem(key, JSON.stringify(value));
	}

	static get<T>(key: LocalStorageKey | string): T | null {
		const dataString = window.localStorage.getItem(key);
		if (dataString) {
			return JSON.parse(dataString) as T;
		}
		return null;
	}

	static remove(key: LocalStorageKey | string): void {
		window.localStorage.removeItem(key);
	}
}

export default LocalStorageHelper;
