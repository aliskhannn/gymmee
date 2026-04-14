import { useState } from 'react';
import { apiClient } from '../api/client';
import { useTelegram } from '../hooks/useTelegram';
import { Dumbbell, Play, Plus } from 'lucide-react';
import { ExerciseCatalog, type Exercise } from '../components/ExerciseCatalog';
import { ActiveExercise } from '../components/ActiveExercise';
import WebApp from '@twa-dev/sdk';

interface WorkoutSession {
	id: number;
	user_id: number;
	started_at: string;
}

export const WorkoutPage = () => {
	const { user, triggerHaptic, showPopup } = useTelegram();
	const [session, setSession] = useState<WorkoutSession | null>(null);
	const [isLoading, setIsLoading] = useState(false);

	const [isCatalogOpen, setIsCatalogOpen] = useState(false);
	const [currentExercise, setCurrentExercise] = useState<Exercise | null>(null);

	const handleStartWorkout = async () => {
		try {
			triggerHaptic('medium');
			setIsLoading(true);
			const { data } = await apiClient.post<WorkoutSession>('/workouts/start', {
				plan_day_id: null
			});
			setSession(data);
		} catch (error) {
			console.error('Ошибка старта тренировки:', error);
		} finally {
			setIsLoading(false);
		}
	};

	const handleExerciseSelect = (exercise: Exercise) => {
		setCurrentExercise(exercise);
		setIsCatalogOpen(false);
	};


	const handleSaveSet = async (weight: number, reps: number) => {
		if (!session || !currentExercise) return;

		try {
			await apiClient.post('/workouts/sets', {
				workout_session_id: session.id,
				exercise_id: currentExercise.id,
				weight: weight,
				reps: reps
			});
			triggerHaptic('medium');
			showPopup(`Записано: ${weight}кг × ${reps}`);
		} catch (error) {
			console.error('Ошибка сохранения подхода:', error);
			showPopup('Ошибка сети. Попробуйте еще раз.');
		}
	};

	const handleFinishWorkout = async () => {
		if (!session) return;

		WebApp.showConfirm("Завершить текущую тренировку?", async (confirmed) => {
			if (confirmed) {
				try {
					setIsLoading(true);
					await apiClient.post('/workouts/finish', {
						workout_session_id: session.id
					});
					triggerHaptic('heavy');
					setSession(null);
					setCurrentExercise(null);
				} catch (error) {
					console.error('Ошибка завершения:', error);
				} finally {
					setIsLoading(false);
				}
			}
		});
	};

	return (
		<div className="p-4 flex flex-col h-full relative">
			<header className="mb-6 mt-4">
				<h1 className="text-3xl font-bold text-white tracking-tight">Тренировка</h1>
				<p className="text-slate-400 mt-1">Готов поработать, {user?.first_name}?</p>
			</header>

			<div className="flex-1 flex flex-col pb-20">
				{!session ? (
					<div className="flex-1 flex items-center justify-center">
						<button
							onClick={handleStartWorkout}
							disabled={isLoading}
							className="flex flex-col items-center justify-center w-48 h-48 rounded-full bg-blue-600 hover:bg-blue-500 active:bg-blue-700 shadow-xl shadow-blue-500/20 transition-all touch-manipulation disabled:opacity-50"
						>
							<Play size={48} className="text-white mb-2 ml-2" fill="currentColor" />
							<span className="text-xl font-bold text-white tracking-wide">
								{isLoading ? 'Загрузка...' : 'НАЧАТЬ'}
							</span>
						</button>
					</div>
				) : (
					<div className="flex flex-col h-full animate-in fade-in duration-500">
						<div className="flex justify-between items-center mb-6 bg-slate-800/50 p-3 rounded-xl border border-slate-700">
							<div className="flex items-center gap-2">
								<div className="w-2 h-2 rounded-full bg-green-500 animate-pulse" />
								<span className="text-slate-300 text-sm font-medium">В процессе</span>
							</div>
							<span className="text-slate-500 text-sm">
								Старт: {new Date(session.started_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
							</span>
						</div>

						{currentExercise ? (
							<div className="flex-1 flex flex-col">
								<ActiveExercise
									exerciseId={currentExercise.id}
									exerciseName={currentExercise.name}
									onSave={handleSaveSet}
								/>

								<button
									onClick={() => { triggerHaptic('light'); setIsCatalogOpen(true); }}
									className="w-full mt-6 flex items-center justify-center gap-2 py-4 border border-dashed border-slate-600 rounded-2xl text-slate-400 hover:text-white hover:border-slate-500 hover:bg-slate-800/50 transition-colors"
								>
									<Plus size={20} />
									Сменить упражнение
								</button>
							</div>
						) : (
							<div className="flex-1 flex flex-col items-center justify-center text-center">
								<div className="w-16 h-16 bg-slate-800 rounded-full flex items-center justify-center mb-4 text-slate-500">
									<Dumbbell size={32} />
								</div>
								<h3 className="text-xl font-bold text-slate-300 mb-2">Начни первый подход</h3>
								<p className="text-slate-500 text-sm mb-6 max-w-62.5">
									Выбери упражнение из справочника, чтобы начать записывать результаты.
								</p>
								<button
									onClick={() => { triggerHaptic('medium'); setIsCatalogOpen(true); }}
									className="flex items-center gap-2 bg-blue-600 hover:bg-blue-500 text-white font-bold py-3 px-8 rounded-full shadow-lg shadow-blue-600/20 active:scale-95 transition-all"
								>
									<Plus size={20} />
									ВЫБРАТЬ УПРАЖНЕНИЕ
								</button>
							</div>
						)}

						<div className="mt-8 pt-4 border-t border-slate-800">
							<button
								onClick={handleFinishWorkout}
								disabled={isLoading}
								className="w-full py-4 rounded-2xl font-bold text-red-400 bg-red-500/10 border border-red-500/20 active:bg-red-500/20 transition-colors"
							>
								ЗАВЕРШИТЬ ТРЕНИРОВКУ
							</button>
						</div>
					</div>
				)}
			</div>

			{isCatalogOpen && (
				<ExerciseCatalog
					onSelect={handleExerciseSelect}
					onClose={() => setIsCatalogOpen(false)}
				/>
			)}
		</div>
	);
};