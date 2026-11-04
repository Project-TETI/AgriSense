'use client';

import React from 'react';
import { ActivityDay } from '../types';

const weekData: ActivityDay[] = [
  { day: 'Sen', count: 12, heightPx: 48 },
  { day: 'Sel', count: 18, heightPx: 72 },
  { day: 'Rab', count: 24, heightPx: 96 },
  { day: 'Kam', count: 20, heightPx: 80 },
  { day: 'Jum', count: 28, heightPx: 112, isPeak: true },
  { day: 'Sab', count: 16, heightPx: 64 },
  { day: 'Min', count: 10, heightPx: 40 },
];

export const ActivityChart: React.FC = () => {
  return (
    <section 
      className="p-6 rounded-[28px] shadow-glass flex flex-col justify-between" 
      data-purpose="chart-activity"
      style={{
        background: 'linear-gradient(135deg, rgba(255, 255, 255, 0.85) 0%, rgba(216, 239, 231, 0.7) 60%, rgba(200, 235, 222, 0.5) 100%)',
        border: '1.5px solid rgba(43, 122, 112, 0.28)',
        boxShadow: 'rgba(24, 79, 72, 0.08) 0px 8px 30px -4px',
      }}
    >
      {/* Chart Header */}
      <div className="flex items-center justify-between pb-4">
        <div className="flex items-center gap-2.5">
          <svg className="w-5 h-5 text-forest-700" fill="none" stroke="currentColor" strokeWidth="2.2" viewBox="0 0 24 24">
            <path d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" strokeLinecap="round" strokeLinejoin="round" />
          </svg>
          <h3 className="text-base font-bold text-forest-900">Aktivitas Pemeriksaan</h3>
        </div>
        <div className="flex items-center gap-1.5 text-xs text-forest-700/80 font-medium">
          <span className="w-2.5 h-2.5 rounded-full bg-forest-600" />
          <span>Jumlah pemeriksaan</span>
        </div>
      </div>

      {/* Bar Chart Container */}
      <div className="relative w-full h-56 pt-2">
        {/* Horizontal Grid Lines with Values */}
        <div className="absolute inset-0 flex flex-col justify-between pointer-events-none text-[11px] text-forest-700/50 pb-7">
          <div className="flex items-center gap-2"><span className="w-6 text-right">40</span><div className="flex-1 border-b border-forest-600/10" /></div>
          <div className="flex items-center gap-2"><span className="w-6 text-right">30</span><div className="flex-1 border-b border-forest-600/10" /></div>
          <div className="flex items-center gap-2"><span className="w-6 text-right">20</span><div className="flex-1 border-b border-forest-600/10" /></div>
          <div className="flex items-center gap-2"><span className="w-6 text-right">10</span><div className="flex-1 border-b border-forest-600/10" /></div>
          <div className="flex items-center gap-2"><span className="w-6 text-right">0</span><div className="flex-1 border-b border-forest-600/20" /></div>
        </div>

        {/* Vertical Columns Bars */}
        <div className="relative h-full flex items-end justify-between pl-8 pr-2 pb-7 z-10">
          {weekData.map((item) => (
            <div key={item.day} className="flex flex-col items-center gap-1.5 group cursor-pointer">
              <span className={`text-[11px] transition-opacity ${item.isPeak ? 'font-extrabold text-forest-900' : 'font-bold text-forest-700 opacity-0 group-hover:opacity-100'}`}>
                {item.count}
              </span>
              <div 
                className={`w-8 sm:w-9 rounded-t-lg group-hover:brightness-110 transition-all shadow-sm ${
                  item.isPeak 
                    ? 'bg-gradient-to-t from-forest-800 to-mint-500 shadow-md ring-2 ring-forest-700/20' 
                    : 'bg-gradient-to-t from-forest-700 to-mint-400'
                }`}
                style={{ height: `${item.heightPx}px` }}
              />
              <span className={`text-xs ${item.isPeak ? 'font-bold text-forest-900' : 'font-semibold text-forest-800'}`}>
                {item.day}
              </span>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};
