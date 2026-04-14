import { useState, useEffect } from 'react';
import { CheckCircle2, Circle, Plus, Loader2, CheckSquare } from 'lucide-react';
import { apiClient } from '../api/client';
import { useTelegram } from '../hooks/useTelegram';

interface Habit {
	id: number;
	user_id: number;
	name: string;
	completed: boolean;
}

export const HabitsPage = () => {
	const { triggerHaptic } = useTelegram();
	const [habits, setHabits] = useState<Habit[]>([]);
	const [isLoading, setIsLoading] = useState(true);
	const [newHabitName, setNewHabitName] = useState('');
	const [isCreating, setIsCreating] = useState(false);

	const fetchHabits = async () => {
		try {
			const { data } = await apiClient.get<Habit[]>('/habits/daily');
			setHabits(data || []);
		} catch (error) {
			console.error('Ошибка загрузки привычек:', error);
		} finally {
			setIsLoading(false);
		}
	};

	useEffect(() => {
		fetchHabits();
	}, []);

	const toggleHabit = async (habitId: number, currentStatus: boolean) => {
		const newStatus = !currentStatus;
		triggerHaptic(newStatus ? 'medium' : 'light');

		setHabits(prev => prev.map(h => h.id === habitId ? { ...h, completed: newStatus } : h));

		try {
			await apiClient.post('/habits/toggle', {
				habit_id: habitId,
				completed: newStatus
			});
		} catch (error) {
			console.error('Ошибка переключения:', error);
			setHabits(prev => prev.map(h => h.id === habitId ? { ...h, completed: currentStatus } : h));
		}
	};

	const handleCreate = async (e: React.FormEvent) => {
		e.preventDefault();
		if (!newHabitName.trim()) return;

		try {
			setIsCreating(true);
			triggerHaptic('light');
			const { data } = await apiClient.post<Habit>('/habits', { name: newHabitName });

			setHabits(prev => [...prev, { ...data, completed: false }]);
			setNewHabitName('');
		} catch (error) {
			console.error('Ошибка создания:', error);
		} finally {
			setIsCreating(false);
		}
	};

	const todayStr = new Intl.DateTimeFormat('ru-RU', {
		weekday: 'long', day: 'numeric', month: 'long'
	}).format(new Date());

	return (
		<div className="p-4 flex flex-col h-full animate-in fade-in duration-500">
			<header className="mb-6 mt-4">
				<h1 className="text-3xl font-bold text-white tracking-tight">Привычки</h1>
				<p className="text-slate-400 mt-1 capitalize">{todayStr}</p>
			</header>

			<form onSubmit={handleCreate} className="mb-8 relative">
				<input
					type="text"
					placeholder="Например: Выпить 2л воды..."
					value={newHabitName}
					onChange={(e) => setNewHabitName(e.target.value)}
					className="w-full bg-slate-800 text-white rounded-2xl py-4 pl-4 pr-12 outline-none focus:ring-2 focus:ring-blue-500/50 transition-all placeholder:text-slate-500 border border-slate-700/50"
				/>
				<button
					type="submit"
					disabled={!newHabitName.trim() || isCreating}
					className="absolute right-2 top-1/2 -translate-y-1/2 w-10 h-10 bg-blue-600 hover:bg-blue-500 disabled:bg-slate-700 text-white rounded-xl flex items-center justify-center transition-colors"
				>
					{isCreating ? <Loader2 size={20} className="animate-spin" /> : <Plus size={20} />}
				</button>
			</form>

			<div className="flex-1 overflow-y-auto pb-20">
				{isLoading ? (
					<div className="flex justify-center mt-10">
						<Loader2 className="animate-spin text-blue-500" size={32} />
					</div>
				) : habits.length === 0 ? (
					<div className="text-center text-slate-500 mt-10">
						<CheckSquare size={48} className="mx-auto mb-4 opacity-20" />
						<p>У вас пока нет привычек.</p>
						<p className="text-sm mt-1">Добавьте первую выше 👆</p>
					</div>
				) : (
					<div className="space-y-3">
						{habits.map((habit) => (
							<div
								key={habit.id}
								onClick={() => toggleHabit(habit.id, habit.completed)}
								className={`flex items-center p-4 rounded-2xl cursor-pointer transition-all duration-300 border ${habit.completed
									? 'bg-blue-500/10 border-blue-500/20'
									: 'bg-slate-800/50 border-slate-700/50 hover:bg-slate-800'
									}`}
							>
								<div className="mr-4">
									{habit.completed ? (
										<CheckCircle2 size={28} className="text-blue-500" />
									) : (
										<Circle size={28} className="text-slate-500" />
									)}
								</div>
								<span className={`text-lg font-medium transition-all duration-300 ${habit.completed ? 'text-blue-400 line-through opacity-70' : 'text-slate-200'
									}`}>
									{habit.name}
								</span>
							</div>
						))}
					</div>
				)}
			</div>
		</div>
	);
};