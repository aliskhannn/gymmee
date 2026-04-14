import { useState, useEffect } from 'react';
import { apiClient } from '../api/client';
import { useTelegram } from '../hooks/useTelegram';
import { Save, Scale, Ruler } from 'lucide-react';

export const ProfilePage = () => {
	const { triggerHaptic, showPopup } = useTelegram();
	const [profile, setProfile] = useState<any>(null);

	useEffect(() => {
		apiClient.get('/me').then(res => setProfile(res.data));
	}, []);

	const handleSave = async () => {
		try {
			await apiClient.post('/me', profile);
			triggerHaptic('medium');
			showPopup('Профиль сохранен!');
		} catch (e) {
			showPopup('Ошибка сохранения');
		}
	};

	if (!profile) return null;

	return (
		<div className="p-4 pb-24 space-y-8 animate-in fade-in duration-500">
			<header className="mt-4">
				<h1 className="text-3xl font-bold text-white tracking-tight">Профиль</h1>
			</header>

			<section className="space-y-4">
				<h2 className="text-sm font-bold text-slate-500 uppercase tracking-widest">Твои данные</h2>
				<div className="grid grid-cols-2 gap-4">
					<div className="bg-slate-800/50 p-4 rounded-2xl border border-slate-700/50">
						<div className="flex items-center gap-2 text-slate-400 mb-2 text-xs">
							<Scale size={14} /> Вес (кг)
						</div>
						<input
							type="number"
							value={profile.weight || ''}
							onChange={e => setProfile({ ...profile, weight: parseFloat(e.target.value) })}
							className="bg-transparent text-2xl font-bold text-white w-full outline-none"
						/>
					</div>
					<div className="bg-slate-800/50 p-4 rounded-2xl border border-slate-700/50">
						<div className="flex items-center gap-2 text-slate-400 mb-2 text-xs">
							<Ruler size={14} /> Рост (см)
						</div>
						<input
							type="number"
							value={profile.height || ''}
							onChange={e => setProfile({ ...profile, height: parseFloat(e.target.value) })}
							className="bg-transparent text-2xl font-bold text-white w-full outline-none"
						/>
					</div>
				</div>
			</section>

			<section className="space-y-4">
				<h2 className="text-sm font-bold text-slate-500 uppercase tracking-widest">Оборудование зала</h2>
				<div className="bg-slate-800/50 p-4 rounded-2xl border border-slate-700/50 space-y-4">
					<div>
						<label className="text-xs text-slate-400 block mb-2">Вес пустого грифа (кг)</label>
						<input
							type="number"
							value={profile.barbell_weight}
							onChange={e => setProfile({ ...profile, barbell_weight: parseFloat(e.target.value) })}
							className="bg-slate-900/50 w-full p-3 rounded-xl border border-slate-700 text-white font-bold"
						/>
					</div>
					<div>
						<label className="text-xs text-slate-400 block mb-2">Доступные блины (через запятую)</label>
						<input
							type="text"
							value={profile.available_plates}
							onChange={e => setProfile({ ...profile, available_plates: e.target.value })}
							className="bg-slate-900/50 w-full p-3 rounded-xl border border-slate-700 text-white font-bold text-sm"
							placeholder="25, 20, 15, 10, 5, 2.5"
						/>
					</div>
				</div>
			</section>

			<button
				onClick={handleSave}
				className="w-full bg-blue-600 py-4 rounded-2xl font-bold text-white flex items-center justify-center gap-2 shadow-lg shadow-blue-500/20 active:scale-95 transition-all"
			>
				<Save size={20} /> СОХРАНИТЬ
			</button>
		</div>
	);
};