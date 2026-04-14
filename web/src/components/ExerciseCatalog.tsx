import React, { useState, useEffect, useMemo } from 'react';
import { Search, ChevronDown, ChevronRight, Dumbbell, X, Loader2 } from 'lucide-react';
import { apiClient } from '../api/client';
import { useTelegram } from '../hooks/useTelegram';

export interface Exercise {
	id: number;
	name: string;
	muscle_group: string;
}

interface Props {
	onSelect: (exercise: Exercise) => void;
	onClose: () => void;
}

export const ExerciseCatalog: React.FC<Props> = ({ onSelect, onClose }) => {
	const { triggerHaptic } = useTelegram();
	const [exercises, setExercises] = useState<Exercise[]>([]);
	const [searchQuery, setSearchQuery] = useState('');
	const [loading, setLoading] = useState(true);

	const [expandedGroups, setExpandedGroups] = useState<Record<string, boolean>>({});

	useEffect(() => {
		const fetchExercises = async () => {
			try {
				const { data } = await apiClient.get<Exercise[]>('/exercises');
				setExercises(data);

				const initialExpanded: Record<string, boolean> = {};
				const groups = Array.from(new Set(data.map(e => e.muscle_group)));
				groups.slice(0, 2).forEach(g => initialExpanded[g] = true);
				setExpandedGroups(initialExpanded);
			} catch (error) {
				console.error('Ошибка загрузки упражнений:', error);
			} finally {
				setLoading(false);
			}
		};
		fetchExercises();
	}, []);

	const groupedExercises = useMemo(() => {
		const filtered = exercises.filter(ex =>
			ex.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
			ex.muscle_group.toLowerCase().includes(searchQuery.toLowerCase())
		);

		return filtered.reduce((acc, ex) => {
			if (!acc[ex.muscle_group]) {
				acc[ex.muscle_group] = [];
			}
			acc[ex.muscle_group].push(ex);
			return acc;
		}, {} as Record<string, Exercise[]>);
	}, [exercises, searchQuery]);

	const toggleGroup = (group: string) => {
		triggerHaptic('light');
		setExpandedGroups(prev => ({ ...prev, [group]: !prev[group] }));
	};

	const handleSelect = (exercise: Exercise) => {
		triggerHaptic('medium');
		onSelect(exercise);
	};

	return (
		<div className="fixed inset-0 z-50 flex flex-col bg-slate-900 animate-in slide-in-from-bottom-full duration-300">
			<div className="flex items-center justify-between p-4 border-b border-slate-800 bg-slate-900/95 backdrop-blur z-10">
				<h2 className="text-xl font-bold text-white flex items-center gap-2">
					<Dumbbell className="text-blue-500" size={24} />
					Упражнения
				</h2>
				<button
					onClick={() => { triggerHaptic('light'); onClose(); }}
					className="p-2 bg-slate-800 rounded-full text-slate-400 hover:text-white transition-colors"
				>
					<X size={20} />
				</button>
			</div>

			<div className="p-4 bg-slate-900">
				<div className="relative">
					<Search className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500" size={18} />
					<input
						type="text"
						placeholder="Поиск упражнения..."
						value={searchQuery}
						onChange={(e) => setSearchQuery(e.target.value)}
						className="w-full bg-slate-800 text-white rounded-xl py-3 pl-10 pr-4 outline-none focus:ring-2 focus:ring-blue-500/50 transition-all placeholder:text-slate-500"
					/>
				</div>
			</div>

			<div className="flex-1 overflow-y-auto px-4 pb-24">
				{loading ? (
					<div className="flex justify-center mt-10">
						<Loader2 className="animate-spin text-blue-500" size={32} />
					</div>
				) : Object.keys(groupedExercises).length === 0 ? (
					<div className="text-center text-slate-500 mt-10">
						Ничего не найдено 🤷‍♂️
					</div>
				) : (
					Object.entries(groupedExercises).map(([group, groupExercises]) => {
						const isExpanded = searchQuery ? true : expandedGroups[group];

						return (
							<div key={group} className="mb-4 bg-slate-800/50 rounded-2xl overflow-hidden border border-slate-800/80">
								<button
									onClick={() => toggleGroup(group)}
									className="w-full flex items-center justify-between p-4 bg-slate-800 hover:bg-slate-700/80 transition-colors"
								>
									<span className="font-bold text-slate-200">{group}</span>
									<div className="flex items-center gap-2">
										<span className="text-xs text-slate-500 bg-slate-900 px-2 py-1 rounded-md">
											{groupExercises.length}
										</span>
										{isExpanded ? (
											<ChevronDown size={18} className="text-slate-400" />
										) : (
											<ChevronRight size={18} className="text-slate-400" />
										)}
									</div>
								</button>

								{isExpanded && (
									<div className="divide-y divide-slate-700/30">
										{groupExercises.map(exercise => (
											<button
												key={exercise.id}
												onClick={() => handleSelect(exercise)}
												className="w-full flex items-center justify-between p-4 text-left hover:bg-blue-500/10 active:bg-blue-500/20 transition-colors group"
											>
												<span className="text-slate-300 group-hover:text-white transition-colors">
													{exercise.name}
												</span>
												<div className="w-6 h-6 rounded-full bg-slate-700 flex items-center justify-center group-hover:bg-blue-500 group-hover:text-white text-slate-400 transition-colors">
													<ChevronRight size={14} />
												</div>
											</button>
										))}
									</div>
								)}
							</div>
						);
					})
				)}
			</div>
		</div>
	);
};