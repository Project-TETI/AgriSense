'use client';

import React from 'react';

export const Header: React.FC = () => {
  return (
    <header className="flex flex-col md:flex-row md:items-center justify-between gap-4 py-1" data-purpose="top-navigation-header">
      {/* Weather & Field Condition Pill */}
      <div 
        className="relative flex-1 max-w-xl flex flex-wrap sm:flex-nowrap items-center justify-between gap-3 px-4 py-2.5 rounded-2xl shadow-sm transition-all" 
        style={{
          background: 'linear-gradient(135deg, rgba(226, 245, 238, 0.9) 0%, rgba(188, 232, 218, 0.5) 100%)',
          border: '1px solid rgba(49, 158, 145, 0.3)',
          backdropFilter: 'blur(12px)',
          boxShadow: 'rgba(24, 79, 72, 0.06) 0px 4px 16px',
        }}
      >
        <div className="flex items-center gap-2.5 min-w-0">
          <div className="w-8 h-8 rounded-xl bg-forest-700/10 flex items-center justify-center text-forest-700 shrink-0">
            <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <circle cx="12" cy="12" r="4" />
              <path d="M12 2v2m0 16v2M4.93 4.93l1.41 1.41m11.32 11.32l1.41 1.41M2 12h2m16 0h2M6.34 17.66l-1.41 1.41M19.07 4.93l-1.41 1.41" strokeLinecap="round" strokeLinejoin="round" />
            </svg>
          </div>
          <div className="min-w-0">
            <div className="flex items-center gap-1.5">
              <span className="text-[10px] font-bold uppercase tracking-wider text-forest-700/70">
                Kondisi Lahan Hari Ini
              </span>
            </div>
            <p className="text-xs font-bold text-forest-900 leading-tight truncate flex items-center gap-1.5 mt-0.5">
              <span>29°C Cerah Berawan</span>
              <span className="text-forest-700/40">•</span>
              <span className="text-forest-800/90 font-medium flex items-center gap-1">
                <svg className="w-3 h-3 text-forest-600 inline" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                  <path d="M12 2.69l5.66 5.66a8 8 0 1 1-11.31 0z" strokeLinecap="round" strokeLinejoin="round" />
                </svg>
                Kelembapan 74%
              </span>
            </p>
          </div>
        </div>
      </div>

      {/* User Profile & Notification Badges */}
      <div className="flex items-center gap-3 self-end md:self-auto">
        {/* Notification Bell Pill */}
        <button 
          aria-label="Notifikasi" 
          className="relative p-2.5 rounded-full bg-white/80 border border-white text-forest-700 hover:text-forest-900 shadow-sm hover:shadow transition-all" 
          type="button"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
            <path d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" strokeLinecap="round" strokeLinejoin="round" />
          </svg>
          {/* DIHAPUS: lingkaran hijau kecil di ikon lonceng */}
        </button>

        {/* User Profile Pill with Online Indicator */}
        <div className="flex items-center gap-2.5 pl-1.5 pr-3.5 py-1 bg-white/80 border border-white rounded-full shadow-sm cursor-pointer hover:shadow hover:bg-white transition-all">
          <div className="w-8 h-8 rounded-full overflow-hidden bg-mint-200 border border-forest-600/20 shrink-0 flex items-center justify-center">
            <svg className="w-full h-full text-forest-800" fill="currentColor" viewBox="0 0 36 36">
              <circle cx="18" cy="18" fill="#A8DEC9" r="18" />
              <circle cx="18" cy="14" fill="#1F4A44" r="6" />
              <path d="M7 32c1.5-6 6.5-9 11-9s9.5 3 11 9" fill="#1F4A44" />
              <ellipse cx="18" cy="10" fill="#EAB308" rx="9" ry="2.2" />
            </svg>
          </div>
          <div className="text-left">
            <p className="text-xs font-bold text-forest-900 leading-none">Setyanto</p>
            {/* DIHAPUS: tulisan "Online" dan titik hijaunya */}
          </div>
          <svg className="w-3.5 h-3.5 text-forest-700/60 ml-0.5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
            <path d="M19 9l-7 7-7-7" strokeLinecap="round" strokeLinejoin="round" />
          </svg>
        </div>
      </div>
    </header>
  );
};