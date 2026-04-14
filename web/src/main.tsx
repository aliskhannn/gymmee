import React from 'react';
import ReactDOM from 'react-dom/client';
import WebApp from '@twa-dev/sdk';
import App from './App.tsx';
import './index.css';

WebApp.ready();
WebApp.expand();
WebApp.setHeaderColor('#0f172a'); // slate-900

ReactDOM.createRoot(document.getElementById('root')!).render(
	<React.StrictMode>
		<App />
	</React.StrictMode>,
)