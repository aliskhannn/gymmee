import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Layout } from './components/Layout';
import { WorkoutPage } from './pages/WorkoutPage';
import { HabitsPage } from './pages/HabitsPage';

const DummyPage = ({ title }: { title: string }) => (
	<div className="p-4 flex items-center justify-center h-full text-slate-400">
		<h2 className="text-xl">{title} (В разработке)</h2>
	</div>
);

function App() {
	return (
		<BrowserRouter>
			<Routes>
				<Route path="/" element={<Layout />}>
					<Route index element={<WorkoutPage />} />

					<Route path="history" element={<DummyPage title="История тренировок" />} />
					<Route path="habits" element={<HabitsPage />} />
					<Route path="profile" element={<DummyPage title="Профиль и настройки" />} />
				</Route>
			</Routes>
		</BrowserRouter>
	);
}

export default App;