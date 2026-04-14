import React, { useState, useEffect, useRef } from 'react';
import { Minus, Plus } from 'lucide-react';
import { useTelegram } from '../hooks/useTelegram';

interface CounterProps {
  label: string;
  value: number;
  step?: number;
  unit?: string;
  minValue?: number;
  maxValue?: number;
  onChange: (value: number) => void;
}

export const GymCounter: React.FC<CounterProps> = ({
  label,
  value,
  step = 1,
  unit = '',
  minValue = 0,
  maxValue = 999,
  onChange,
}) => {
  const { triggerHaptic } = useTelegram();
  const [isEditing, setIsEditing] = useState(false);
  const [inputValue, setInputValue] = useState(String(value));
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    setInputValue(String(value));
  }, [value]);

  useEffect(() => {
    if (isEditing && inputRef.current) {
      inputRef.current.focus();
      inputRef.current.select();
    }
  }, [isEditing]);

  const handleDecrement = () => {
    const newValue = Math.max(minValue, value - step);
    triggerHaptic('light');
    onChange(newValue);
  };

  const handleIncrement = () => {
    const newValue = Math.min(maxValue, value + step);
    triggerHaptic('medium');
    onChange(newValue);
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value.replace(/[^0-9.]/g, '');
    setInputValue(val);
  };

  const handleInputBlur = () => {
    setIsEditing(false);
    let parsed = parseFloat(inputValue);
    
    if (isNaN(parsed)) parsed = minValue;
    if (parsed < minValue) parsed = minValue;
    if (parsed > maxValue) parsed = maxValue;
    
    onChange(parsed);
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      inputRef.current?.blur();
    }
  };

  return (
    <div className="flex flex-col items-center w-full bg-slate-800 rounded-2xl p-4 shadow-sm border border-slate-700/50">
      <span className="text-slate-400 text-xs font-bold uppercase tracking-wider mb-3">{label}</span>
      
      <div className="flex items-center justify-between w-full">
        <button
          onClick={handleDecrement}
          className="w-14 h-14 flex items-center justify-center bg-slate-700 hover:bg-slate-600 active:bg-slate-500 rounded-full text-white transition-colors touch-manipulation shrink-0"
        >
          <Minus size={24} />
        </button>

        <div 
          className="flex-1 flex justify-center items-center cursor-pointer min-w-25"
          onClick={() => !isEditing && setIsEditing(true)}
        >
          {isEditing ? (
            <div className="flex items-baseline">
              <input
                ref={inputRef}
                type="text"
                inputMode="decimal"
                value={inputValue}
                onChange={handleInputChange}
                onBlur={handleInputBlur}
                onKeyDown={handleKeyDown}
                className="w-20 bg-transparent text-center text-4xl font-bold text-white outline-none caret-blue-500"
              />
              {unit && <span className="text-slate-400 ml-1 text-xl">{unit}</span>}
            </div>
          ) : (
            <div className="flex items-baseline">
              <span className="text-4xl font-bold text-white tracking-tight">
                {value}
              </span>
              {unit && <span className="text-slate-400 ml-1 text-xl font-medium">{unit}</span>}
            </div>
          )}
        </div>

        <button
          onClick={handleIncrement}
          className="w-14 h-14 flex items-center justify-center bg-blue-600 hover:bg-blue-500 active:bg-blue-700 rounded-full text-white transition-colors touch-manipulation shadow-lg shadow-blue-600/20 shrink-0"
        >
          <Plus size={24} />
        </button>
      </div>
    </div>
  );
};