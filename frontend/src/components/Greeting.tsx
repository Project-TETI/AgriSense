'use client';

import React from 'react';

export const Greeting: React.FC = () => {
  return (
    <div className="flex flex-col md:flex-row md:items-center justify-between gap-3 pt-1" data-purpose="page-greeting">
      <div>
        <h1 className="text-2xl md:text-3xl font-extrabold tracking-tight text-forest-900">
          Selamat sore, Setyanto
        </h1>
        <p className="text-xs md:text-sm text-forest-700/80 font-medium mt-0.5">
          Deteksi dini kesehatan tanaman padi, cabai &amp; tomat Anda dengan akurasi AI terkini.
        </p>
      </div>
      {/* Friendly encouraging mascot speech badge */}
      <div className="inline-flex items-center gap-2 self-start md:self-auto px-3.5 py-1.5 rounded-full bg-white/70 border border-white/90 shadow-sm text-xs font-semibold text-forest-800">
        <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
        <span>Ayo periksa tanaman hari ini</span>
      </div>
    </div>
  );
};
