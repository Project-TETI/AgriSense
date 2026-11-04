'use client';

import React, { useState } from 'react';
import Image from 'next/image';
import { InspectionRecord } from '../types';

const defaultRecords: InspectionRecord[] = [
  {
    id: 'tomat-cherry',
    cropName: 'Tomat Cherry',
    conditionDetail: 'Bercak Daun (Alternaria)',
    dateStr: '12 Okt 2025, 14:32 WIB',
    statusBadge: {
      label: 'Perlu Perhatian',
      variant: 'warning',
    },
    accuracyPercent: 94,
    actionButton: {
      label: 'Fungisida Nabati',
      variant: 'action',
    },
    imageUrl: '/images/bercak-daun.jpg',
  },
  {
    id: 'cabai-rawit',
    cropName: 'Cabai Rawit',
    conditionDetail: 'Daun Hijau Segar (Bebas Hama)',
    dateStr: '11 Okt 2025, 09:17 WIB',
    statusBadge: {
      label: 'Sehat Prima',
      variant: 'success',
    },
    accuracyPercent: 98,
    actionButton: {
      label: 'Kondisi Optimal',
      variant: 'optimal',
    },
    imageUrl: '/images/embun-tepung.jpg',
  },
  {
    id: 'padi-ciherang',
    cropName: 'Padi Ciherang',
    conditionDetail: 'Hawar Daun Bakteri (Xanthomonas)',
    dateStr: '10 Okt 2025, 16:45 WIB',
    statusBadge: {
      label: 'Perlu Perhatian',
      variant: 'warning',
    },
    accuracyPercent: 89,
    actionButton: {
      label: 'Drainase & Kalium',
      variant: 'action',
    },
    imageUrl: '/images/hawar-daun.jpg',
  },
];

export const InspectionHistory: React.FC = () => {
  const [activeFilter, setActiveFilter] = useState<'all' | 'warning' | 'healthy'>('all');

  const filteredRecords = defaultRecords.filter((record) => {
    if (activeFilter === 'warning') return record.statusBadge.variant === 'warning';
    if (activeFilter === 'healthy') return record.statusBadge.variant === 'success';
    return true;
  });

  return (
    <section 
      className="p-6 rounded-[28px] shadow-glass space-y-5 relative overflow-hidden" 
      data-purpose="recent-inspections"
      style={{
        background: 'linear-gradient(135deg, rgba(226, 248, 241, 0.94) 0%, rgba(206, 240, 230, 0.82) 50%, rgba(195, 233, 222, 0.88) 100%)',
        border: '1.5px solid rgba(43, 122, 112, 0.32)',
        boxShadow: 'rgba(24, 79, 72, 0.12) 0px 12px 36px -6px',
      }}
    >
      {/* Header with title, description, filter pills, and navigation */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 pb-1 border-b border-forest-800/10">
        <div className="flex items-start sm:items-center gap-3.5">
          <div className="w-10 h-10 rounded-2xl bg-gradient-to-br from-[#20635B] to-[#38A783] text-white flex items-center justify-center shrink-0 shadow-md shadow-forest-900/15">
            <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="2.2" viewBox="0 0 24 24">
              <path d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" strokeLinecap="round" strokeLinejoin="round" />
            </svg>
          </div>
          <div>
            <h3 className="text-base sm:lg font-extrabold text-forest-900 tracking-tight leading-tight">
              Riwayat Pemeriksaan Terbaru
            </h3>
            <p className="text-xs text-forest-700/80 font-medium mt-0.5">
              Hasil deteksi AI dan rekomendasi tindakan cepat untuk lahan Anda
            </p>
          </div>
        </div>

        {/* Quick Filter Pills & Link */}
        <div className="flex flex-wrap items-center gap-2">
          <div className="inline-flex p-1 rounded-2xl bg-white/70 border border-white/90 shadow-sm">
            <button 
              type="button" 
              onClick={() => setActiveFilter('all')}
              className={`px-3 py-1 text-xs rounded-xl shadow-sm transition-all ${
                activeFilter === 'all' 
                  ? 'bg-forest-700 text-white font-bold' 
                  : 'font-semibold text-forest-800/80 hover:text-forest-900 hover:bg-white/60'
              }`}
            >
              Semua (3)
            </button>
            <button 
              type="button" 
              onClick={() => setActiveFilter('warning')}
              className={`px-3 py-1 text-xs rounded-xl transition-all ${
                activeFilter === 'warning' 
                  ? 'bg-forest-700 text-white font-bold shadow-sm' 
                  : 'font-semibold text-forest-800/80 hover:text-forest-900 hover:bg-white/60'
              }`}
            >
              Perlu Tindakan (2)
            </button>
            <button 
              type="button" 
              onClick={() => setActiveFilter('healthy')}
              className={`px-3 py-1 text-xs rounded-xl transition-all ${
                activeFilter === 'healthy' 
                  ? 'bg-forest-700 text-white font-bold shadow-sm' 
                  : 'font-semibold text-forest-800/80 hover:text-forest-900 hover:bg-white/60'
              }`}
            >
              Sehat (1)
            </button>
          </div>
          <a 
            href="#semua-riwayat" 
            className="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-2xl text-xs font-bold text-forest-700 hover:text-forest-900 bg-white/60 hover:bg-white border border-white/80 shadow-sm transition-all group"
          >
            <span>Lihat Semua Riwayat</span>
            <svg className="w-3.5 h-3.5 group-hover:translate-x-0.5 transition-transform" fill="none" stroke="currentColor" strokeWidth="2.2" viewBox="0 0 24 24">
              <path d="M9 5l7 7-7 7" strokeLinecap="round" strokeLinejoin="round" />
            </svg>
          </a>
        </div>
      </div>

      {/* Card List */}
      <div className="grid grid-cols-1 gap-3.5">
        {filteredRecords.map((item) => (
          <article 
            key={item.id}
            className="p-4 rounded-2xl transition-all duration-300 hover:shadow-lg hover:-translate-y-0.5 group grid grid-cols-1 gap-4 md:grid-cols-[260px_160px_minmax(140px,1fr)_220px] md:items-center"
            style={{
              background: 'linear-gradient(135deg, rgba(255, 255, 255, 0.92) 0%, rgba(246, 253, 250, 0.88) 55%, rgba(230, 248, 240, 0.8) 100%)',
              border: '1.5px solid rgba(43, 122, 112, 0.28)',
              boxShadow: 'rgba(24, 79, 72, 0.08) 0px 4px 18px',
            }}
          >
            {/* Left: Thumbnail & Plant Info */}
            <div className="flex items-center gap-3.5">
              <div className="relative w-14 h-14 rounded-2xl overflow-hidden shrink-0 ring-2 ring-forest-700/20 shadow-sm group-hover:scale-105 transition-transform duration-300">
                <Image 
                  src={item.imageUrl} 
                  alt={item.cropName} 
                  width={56}
                  height={56}
                  className="w-full h-full object-cover" 
                />
              </div>
              <div>
                <div className="flex items-center gap-2">
                  <h4 className="text-sm font-extrabold text-forest-900 group-hover:text-forest-700 transition-colors">
                    {item.cropName}
                  </h4>
                </div>
                <p className={`text-xs font-bold mt-0.5 flex items-center gap-1 ${
                  item.statusBadge.variant === 'warning' ? 'text-amber-800' : 'text-emerald-700'
                }`}>
                  <span>{item.conditionDetail}</span>
                </p>
                <p className="text-[11px] text-forest-700/70 font-medium mt-0.5 flex items-center gap-1.5">
                  <svg className="w-3.5 h-3.5 text-forest-700/50" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                    <path d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" strokeLinecap="round" strokeLinejoin="round" />
                  </svg>
                  {item.dateStr}
                </p>
              </div>
            </div>

            {/* Status */}
            <div>
              <span className="text-[10px] font-bold uppercase tracking-wider text-forest-700/60 block mb-1">
                Kondisi Diagnosa
              </span>
              {item.statusBadge.variant === 'warning' ? (
                <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold bg-amber-100 text-amber-800 border border-amber-200/80 shadow-sm">
                  <span className="w-2 h-2 rounded-full bg-amber-500 animate-pulse" />
                  {item.statusBadge.label}
                </span>
              ) : (
                <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold bg-emerald-100 text-emerald-800 border border-emerald-200/80 shadow-sm">
                  <svg className="w-3.5 h-3.5 text-emerald-600" fill="currentColor" viewBox="0 0 20 20">
                    <path clipRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" fillRule="evenodd" />
                  </svg>
                  {item.statusBadge.label}
                </span>
              )}
            </div>

            {/* Confidence Bar */}
            <div className="w-full">
              <div className="flex items-center justify-between text-xs font-bold mb-1">
                <span className="text-forest-700/80 font-medium text-[11px]">Akurasi AI</span>
                <span className="text-forest-900 font-extrabold">{item.accuracyPercent}%</span>
              </div>
              <div className="w-full h-2 rounded-full bg-forest-700/10 overflow-hidden p-0.5 ring-1 ring-forest-700/15">
                <div 
                  className={`h-full rounded-full ${
                    item.statusBadge.variant === 'warning' 
                      ? 'bg-gradient-to-r from-emerald-500 to-amber-500' 
                      : 'bg-gradient-to-r from-mint-500 to-emerald-600'
                  }`} 
                  style={{ width: `${item.accuracyPercent}%` }} 
                />
              </div>
            </div>

            {/* Right: Action CTA & Detail Icon */}
            <div className="flex items-center justify-between md:justify-end gap-3 pt-2 md:pt-0 border-t md:border-t-0 border-forest-800/10">
              {item.actionButton.variant === 'action' ? (
                <a 
                  href={`#action-${item.id}`} 
                  className="inline-flex items-center justify-center gap-1.5 min-w-[150px] px-3.5 py-2 rounded-xl text-xs font-bold text-white shadow-sm transition-all duration-200 hover:shadow-md hover:brightness-105 active:scale-[0.98]" 
                  style={{ background: 'linear-gradient(135deg, rgb(32, 99, 91) 0%, rgb(56, 167, 131) 100%)' }}
                >
                  <span>{item.actionButton.label}</span>
                  <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" strokeWidth="2.2" viewBox="0 0 24 24">
                    <path d="M14 5l7 7m0 0l-7 7m7-7H3" strokeLinecap="round" strokeLinejoin="round" />
                  </svg>
                </a>
              ) : (
                <span className="inline-flex items-center justify-center gap-1.5 min-w-[150px] px-3 py-1.5 rounded-xl text-xs font-bold text-forest-800 bg-white/70 border border-forest-600/20">
                  <svg className="w-3.5 h-3.5 text-emerald-600" fill="none" stroke="currentColor" strokeWidth="2.2" viewBox="0 0 24 24">
                    <path d="M5 13l4 4L19 7" strokeLinecap="round" strokeLinejoin="round" />
                  </svg>
                  {item.actionButton.label}
                </span>
              )}
              <button 
                type="button" 
                aria-label={`Lihat rincian riwayat ${item.cropName}`} 
                className="p-2 rounded-xl bg-white/80 hover:bg-white text-forest-700 hover:text-forest-900 border border-forest-600/15 shadow-sm transition-all"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2.2" viewBox="0 0 24 24">
                  <path d="M9 5l7 7-7 7" strokeLinecap="round" strokeLinejoin="round" />
                </svg>
              </button>
            </div>
          </article>
        ))}
      </div>
    </section>
  );
};