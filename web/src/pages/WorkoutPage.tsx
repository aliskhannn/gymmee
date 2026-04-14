import { useState } from 'react';
import { apiClient } from '../api/client';
import { useTelegram } from '../hooks/useTelegram';
import { Play } from 'lucide-react';

interface WorkoutSession {
  id: number;
  user_id: number;
  started_at: string;
}

export const WorkoutPage = () => {
  const { user, triggerHaptic } = useTelegram();
  const [session, setSession] = useState<WorkoutSession | null>(null);
  const [isLoading, setIsLoading] = useState(false);

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

  return (
    <div className="p-4 flex flex-col h-full">
      <header className="mb-8 mt-4">
        <h1 className="text-3xl font-bold text-white tracking-tight">Тренировка</h1>
        <p className="text-slate-400 mt-1">Готов поработать, {user?.first_name}?</p>
      </header>

      <div className="flex-1 flex flex-col items-center justify-center pb-20">
        {!session ? (
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
        ) : (
          <div className="w-full text-center space-y-4">
            <div className="inline-flex items-center justify-center px-4 py-1.5 rounded-full bg-green-500/20 border border-green-500/30 text-green-400 text-sm font-medium">
              Тренировка активна
            </div>
            <p className="text-slate-500 text-sm">
              Сессия: #{session.id} <br />
              Начата: {new Date(session.started_at).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})}
            </p>
          </div>
        )}
      </div>
    </div>
  );
};