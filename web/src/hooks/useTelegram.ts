import WebApp from '@twa-dev/sdk';

export function useTelegram() {
	const user = WebApp.initDataUnsafe?.user;

	const triggerHaptic = (style: 'light' | 'medium' | 'heavy' | 'rigid' | 'soft' = 'light') => {
		if (WebApp.HapticFeedback) {
			WebApp.HapticFeedback.impactOccurred(style);
		}
	};

	const showPopup = (message: string) => {
		WebApp.showPopup({ message });
	};

	return {
		tg: WebApp,
		user,
		triggerHaptic,
		showPopup,
	};
}