'use client';

import React from 'react';
import Image from 'next/image';
import { DiseaseEducationItem } from '../types';

const diseaseList: DiseaseEducationItem[] = [
  {
    id: 'bercak',
    title: 'Bercak Daun',
    description: 'Bercak konsentris kuning/cokelat. Segera terapkan fungisida nabati & pangkas daun bawah.',
    imageUrl: '/images/bercak-daun.jpg',
    actionText: 'Pelajari Solusi',
  },
  {
    id: 'embun',
    title: 'Embun Tepung',
    description: 'Lapisan bedak putih halus. Jaga kelembapan sirkulasi udara dan semprot sulfur larut air.',
    imageUrl: '/images/embun-tepung.jpg',
    actionText: 'Pelajari Solusi',
  },
  {
    id: 'hawar',
    title: 'Hawar Daun',
    description: 'Daun mengering cokelat dari tepi helai. Optimalkan drainase parit dan berikan pupuk kalium.',
    imageUrl: '/images/hawar-daun.jpg',
    actionText: 'Pelajari Solusi',
  },
];

export const DiseaseEducation: React.FC = () => {
  return (
    <section 
      className="p-6 rounded-[28px] shadow-glass flex flex-col justify-between relative overflow-hidden" 
      data-purpose="disease-education"
      style={{
        background: 'linear-gradient(145deg, rgba(235, 252, 245, 0.95) 0%, rgba(215, 245, 235, 0.85) 50%, rgba(225, 248, 240, 0.9) 100%)',
        border: '1.5px solid rgba(43, 122, 112, 0.35)',
        boxShadow: 'rgba(24, 79, 72, 0.12) 0px 12px 36px -6px',
      }}
    >
      {/* Header */}
      <div className="flex items-center justify-between pb-3.5">
        <div className="flex items-center gap-2.5">
          <div className="w-8 h-8 rounded-xl bg-forest-700/10 flex items-center justify-center text-forest-700 shadow-sm">
            <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2.2" viewBox="0 0 24 24">
              <path d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" strokeLinecap="round" strokeLinejoin="round" />
            </svg>
          </div>
          <div>
            <h3 className="text-base font-bold text-forest-900 leading-tight">Edukasi Penyakit Tanaman</h3>
            <p className="text-[11px] font-medium text-forest-700/75">Panduan identifikasi gejala &amp; solusi terpadu</p>
          </div>
        </div>
        <a className="inline-flex items-center gap-1 text-xs font-bold text-forest-700 hover:text-forest-900 transition-colors group" href="#semua-penyakit">
          Lihat semua
          <svg className="w-3.5 h-3.5 group-hover:translate-x-0.5 transition-transform" fill="none" stroke="currentColor" strokeWidth="2.2" viewBox="0 0 24 24">
            <path d="M9 5l7 7-7 7" strokeLinecap="round" strokeLinejoin="round" />
          </svg>
        </a>
      </div>

      {/* Cards Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-3.5 pt-1">
        {diseaseList.map((disease) => (
          <article 
            key={disease.id}
            className="flex flex-col rounded-2xl p-3.5 transition-all duration-300 hover:shadow-xl hover:-translate-y-0.5 group relative overflow-hidden" 
            style={{
              background: 'linear-gradient(145deg, rgba(255, 255, 255, 0.95) 0%, rgba(225, 248, 240, 0.9) 55%, rgba(195, 240, 226, 0.85) 100%)',
              border: '1.5px solid rgba(43, 122, 112, 0.3)',
              boxShadow: 'rgba(24, 79, 72, 0.1) 0px 6px 20px -3px',
            }}
          >
            <div className="relative w-full h-24 rounded-xl overflow-hidden mb-2.5 shadow-sm ring-1 ring-teal-800/15">
              <Image 
                alt={disease.title} 
                className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" 
                src={disease.imageUrl} 
                width={300}
                height={160}
              />
            </div>
            <h4 className="text-xs font-extrabold text-forest-900 group-hover:text-forest-700 mb-1 flex items-center justify-between">
              <span>{disease.title}</span>
            </h4>
            <p className="text-[11px] text-forest-800/80 leading-relaxed flex-1 mb-3 font-medium">
              {disease.description}
            </p>
            <a 
              className="inline-flex items-center justify-center w-full px-3.5 py-2 rounded-xl text-xs font-bold text-white shadow-sm transition-all duration-200 text-center hover:shadow-md hover:brightness-105 active:scale-[0.98]" 
              href={`#detail-${disease.id}`} 
              style={{
                background: 'linear-gradient(135deg, rgb(32, 99, 91) 0%, rgb(56, 167, 131) 100%)',
                boxShadow: 'rgba(32, 99, 91, 0.25) 0px 3px 10px',
              }}
            >
              {disease.actionText}
            </a>
          </article>
        ))}
      </div>
    </section>
  );
};
