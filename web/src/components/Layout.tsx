import { Outlet, NavLink } from 'react-router-dom';
import { Dumbbell, CalendarDays, CheckSquare, User } from 'lucide-react';
import { useTelegram } from '../hooks/useTelegram';

export const Layout = () => {
	const { triggerHaptic } = useTelegram();

	const navItems = [
		{ path: '/', icon: Dumbbell, label: 'Треня' },
		{ path: '/history', icon: CalendarDays, label: 'История' },
		{ path: '/habits', icon: CheckSquare, label: 'Привычки' },
		{ path: '/profile', icon: User, label: 'Профиль' },
	];

	return (
		<div className="flex flex-col h-screen bg-slate-900 text-slate-100 overflow-hidden">
			<main className="flex-1 overflow-y-auto pb-20">
				<Outlet />
			</main>

			<nav className="fixed bottom-0 w-full bg-slate-900/90 backdrop-blur-md border-t border-slate-800 pb-safe">
				<ul className="flex justify-around items-center h-16 px-2">
					{navItems.map((item) => (
						<li key={item.path} className="w-full">
							<NavLink
								to={item.path}
								onClick={() => triggerHaptic('light')}
								className={({ isActive }) =>
									`flex flex-col items-center justify-center w-full h-full space-y-1 transition-colors ${isActive ? 'text-blue-500' : 'text-slate-500 hover:text-slate-400'
									}`
								}
							>
								<item.icon size={24} />
								<span className="text-[10px] font-medium">{item.label}</span>
							</NavLink>
						</li>
					))}
				</ul>
			</nav>
		</div>
	);
};