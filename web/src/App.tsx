import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Layout } from './components/Layout';
import { WorkoutPage } from './pages/WorkoutPage';
import { HabitsPage } from './pages/HabitsPage';
import { HistoryPage } from './pages/HistoryPage';
import { ProfilePage } from './pages/ProfilePage';

function App() {
	return (
		<BrowserRouter>
			<Routes>
				<Route path="/" element={<Layout />}>
					<Route index element={<WorkoutPage />} />

					<Route path="history" element={<HistoryPage />} />
					<Route path="habits" element={<HabitsPage />} />
					<Route path="profile" element={<ProfilePage />} />
				</Route>
			</Routes>
		</BrowserRouter>
	);
}

export default App;