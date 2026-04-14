import WebApp from '@twa-dev/sdk';

function App() {
	const user = WebApp.initDataUnsafe?.user;

	return (
		<div className="p-4 flex flex-col items-center justify-center min-h-screen">
			<h1 className="text-3xl font-bold text-blue-500 mb-2">Gymmee</h1>
			<p className="text-slate-400">
				Привет, {user?.first_name || 'Атлет'}!
			</p>
		</div>
	)
}

export default App;