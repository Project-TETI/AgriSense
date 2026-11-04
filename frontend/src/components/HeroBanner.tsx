'use client';

import React, { useState } from 'react';
import { Mascot } from './Mascot';

export const HeroBanner: React.FC = () => {
  const [isExaminingHovered, setIsExaminingHovered] = useState<boolean>(false);

  return (
    <section 
      className="relative overflow-hidden rounded-[28px] bg-gradient-to-br from-[#E2F5EE] via-[#D8EFE7] to-[#CDE8DD] border border-white/80 p-7 md:p-9 shadow-glass" 
      data-purpose="hero-banner"
    >
      {/* Decorative background glow */}
      <div className="absolute left-1/3 bottom-0 w-64 h-64 bg-emerald-300/20 rounded-full blur-2xl pointer-events-none" />

      <div className="relative z-10 grid grid-cols-1 md:grid-cols-12 items-center gap-8">
        {/* Hero Text & CTA Buttons (8 cols) */}
        <div className="md:col-span-8 space-y-4">
          <div className="space-y-2">
            <h2 className="text-2xl sm:text-3xl lg:text-4xl font-extrabold text-forest-900 tracking-tight leading-tight">
              Tanaman Sehat, Panen Lebih Berlimpah
            </h2>
            <p className="text-xs sm:text-sm text-forest-800/80 font-medium max-w-xl leading-relaxed">
              Identifikasi penyakit daun hanya dalam hitungan detik. Dapatkan diagnosa komprehensif dan panduan tindakan praktis untuk hasil panen terbaik.
            </p>
          </div>

          {/* Action Buttons Group */}
          <div className="flex flex-wrap items-center gap-3 pt-1">
            <a
              href="#periksa"
              onPointerEnter={(e) => {
                if (e.pointerType !== 'touch') {
                  setIsExaminingHovered(true);
                }
              }}
              onPointerLeave={() => setIsExaminingHovered(false)}
              onTouchEnd={() => setIsExaminingHovered(false)}
              className="inline-flex items-center gap-2.5 px-6 py-3.5 rounded-2xl text-white font-semibold text-sm sm:text-base shadow-lg transition-all transform hover:-translate-y-0.5 active:translate-y-0 hover:brightness-105"
              style={{
                background: 'linear-gradient(135deg, rgb(46, 94, 90) 0%, rgb(89, 169, 163) 100%)',
                boxShadow: 'rgba(46, 94, 90, 0.3) 0px 4px 18px',
              }}
            >
              <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" strokeWidth="2.2" viewBox="0 0 24 24">
                <path d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" strokeLinecap="round" strokeLinejoin="round" />
                <path d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
              <span>Mulai Periksa</span>
            </a>

            <a
              href="#cara-kerja"
              className="inline-flex items-center gap-2.5 px-5 py-3.5 rounded-2xl backdrop-blur-md text-[#184F48] font-semibold text-sm sm:text-base shadow-sm hover:shadow-md transition-all duration-200 active:scale-[0.98] hover:-translate-y-0.5 group"
              style={{
                background: 'linear-gradient(135deg, rgba(240, 253, 248, 0.95) 0%, rgba(209, 245, 233, 0.85) 100%)',
                border: '1.5px solid rgba(43, 122, 112, 0.35)',
                boxShadow: 'rgba(24, 79, 72, 0.12) 0px 4px 14px -2px',
              }}
            >
              <div className="w-6 h-6 rounded-lg bg-emerald-100 flex items-center justify-center text-forest-700 group-hover:scale-105 transition-transform">
                <svg className="w-4 h-4 text-forest-700" fill="none" stroke="currentColor" strokeWidth="2.2" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                </svg>
              </div>
              <span className="tracking-wide">Pelajari cara kerja</span>
            </a>
          </div>
        </div>

        {/* Mascot Character Illustration (4 cols) */}
        <div className="md:col-span-4 flex flex-col items-center justify-center">
          <Mascot isSmiling={isExaminingHovered} />
        </div>
      </div>
    </section>
  );
};
