import { useState, useEffect } from 'react';
import { apiClient } from '../api/client';
import { Calendar, Clock, ChevronRight, Trophy } from 'lucide-react';

interface WorkoutSession {
	id: number;
	started_at: string;
	ended_at: string;
}

export const HistoryPage = () => {
	const [history, setHistory] = useState<WorkoutSession[]>([]);
	const [loading, setLoading] = useState(true);

	useEffect(() => {
		apiClient.get<WorkoutSession[]>('/workouts/history')
			.then(res => setHistory(res.data || []))
			.finally(() => setLoading(false));
	}, []);

	const formatDate = (dateStr: string) => {
		return new Intl.DateTimeFormat('ru-RU', { day: 'numeric', month: 'long' }).format(new Date(dateStr));
	};

	const calculateDuration = (start: string, end: string) => {
		const diff = new Date(end).getTime() - new Date(start).getTime();
		return Math.floor(diff / 1000 / 60); // в минутах
	};

	return (
		<div className="p-4 animate-in fade-in duration-500">
			<header className="mb-6 mt-4">
				<h1 className="text-3xl font-bold text-white tracking-tight">История</h1>
			</header>

			{history.length === 0 && !loading ? (
				<div className="flex flex-col items-center justify-center mt-20 text-slate-500">
					<Calendar size={48} className="mb-4 opacity-20" />
					<p>Тренировок пока нет. Время начать!</p>
				</div>
			) : (
				<div className="space-y-4">
					{history.map(session => (
						<div key={session.id} className="bg-slate-800/50 border border-slate-700/50 rounded-2xl p-4 flex items-center justify-between group active:scale-[0.98] transition-all">
							<div className="flex items-center gap-4">
								<div className="w-12 h-12 rounded-full bg-blue-500/10 flex items-center justify-center text-blue-500">
									<Trophy size={24} />
								</div>
								<div>
									<h3 className="text-white font-bold">{formatDate(session.started_at)}</h3>
									<div className="flex items-center gap-2 text-slate-500 text-xs">
										<Clock size={12} />
										<span>{calculateDuration(session.started_at, session.ended_at)} мин</span>
									</div>
								</div>
							</div>
							<ChevronRight className="text-slate-600 group-hover:text-slate-400 transition-colors" />
						</div>
					))}
				</div>
			)}
		</div>
	);
};