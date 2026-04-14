import React, { useState, useEffect } from 'react';
import { GymCounter } from './GymCounter';
import { apiClient } from '../api/client';
import { useTelegram } from '../hooks/useTelegram';
import { Loader2, Info } from 'lucide-react';

interface PlateRequirement {
	weight: number;
	count: number;
}

interface HintResult {
	last_weight: number;
	last_reps: number;
	plates_required: PlateRequirement[];
}

interface Props {
	exerciseId: number;
	exerciseName: string;
	onSave: (weight: number, reps: number) => void;
}

export const ActiveExercise: React.FC<Props> = ({ exerciseId, exerciseName, onSave }) => {
	const { triggerHaptic } = useTelegram();
	const [hint, setHint] = useState<HintResult | null>(null);
	const [weight, setWeight] = useState(0);
	const [reps, setReps] = useState(0);
	const [loading, setLoading] = useState(true);

	// Загружаем подсказку при смене упражнения
	useEffect(() => {
		const fetchHint = async () => {
			try {
				setLoading(true);
				const { data } = await apiClient.get<HintResult>(`/workouts/hints?exercise_id=${exerciseId}`);
				if (data.last_weight) {
					setHint(data);
					setWeight(data.last_weight);
					setReps(data.last_reps);
				}
			} catch (e) {
				console.error("Не удалось загрузить подсказку", e);
			} finally {
				setLoading(false);
			}
		};

		fetchHint();
	}, [exerciseId]);

	const handleSave = () => {
		triggerHaptic('heavy');
		onSave(weight, reps);
	};

	if (loading) return (
		<div className="flex justify-center p-12">
			<Loader2 className="animate-spin text-blue-500" size={32} />
		</div>
	);

	return (
		<div className="flex flex-col gap-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
			{/* Шапка упражнения */}
			<div>
				<h2 className="text-2xl font-bold text-white leading-tight">{exerciseName}</h2>
				{hint && (
					<div className="flex items-center gap-2 mt-2 text-blue-400 bg-blue-500/10 w-fit px-3 py-1 rounded-full border border-blue-500/20">
						<Info size={14} />
						<span className="text-xs font-medium uppercase tracking-wider">
							Прошлый раз: {hint.last_weight}кг × {hint.last_reps}
						</span>
					</div>
				)}
			</div>

			{/* Визуализация блинов (если есть подсказка) */}
			{hint && hint.plates_required.length > 0 && (
				<div className="bg-slate-800/50 border border-slate-700 rounded-2xl p-4">
					<span className="text-[10px] text-slate-500 font-bold uppercase tracking-widest mb-3 block">Навеска на сторону:</span>
					<div className="flex items-end gap-1 h-12">
						{/* Гриф */}
						<div className="w-2 h-4 bg-slate-600 rounded-sm" />
						{/* Блины */}
						{hint.plates_required.map((p, i) => (
							<div key={i} className="flex gap-1 items-end">
								{[...Array(p.count)].map((_, idx) => (
									<div
										key={idx}
										className="w-3 bg-blue-500 rounded-sm border border-blue-400/50 shadow-[0_0_10px_rgba(59,130,246,0.3)]"
										style={{ height: `${Math.max(20, p.weight * 1.5)}px` }}
									/>
								))}
							</div>
						))}
					</div>
				</div>
			)}

			{/* Счетчики */}
			<div className="space-y-4">
				<GymCounter
					label="ВЕС (КГ)"
					value={weight}
					step={0.5}
					onChange={setWeight}
					unit="кг"
				/>
				<GymCounter
					label="ПОВТОРЕНИЯ"
					value={reps}
					step={1}
					onChange={setReps}
					unit=""
				/>
			</div>

			<button
				onClick={handleSave}
				className="w-full bg-blue-600 hover:bg-blue-500 text-white font-bold py-4 rounded-2xl shadow-lg shadow-blue-600/20 transition-all active:scale-[0.98] mt-2"
			>
				ЗАПИСАТЬ ПОДХОД
			</button>
		</div>
	);
};