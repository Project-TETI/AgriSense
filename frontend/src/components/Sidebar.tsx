'use client';

import React, { useState } from 'react';

interface SidebarProps {
  activeTab?: string;
  onNavigate?: (tab: string) => void;
}

export const Sidebar: React.FC<SidebarProps> = ({ activeTab = 'beranda', onNavigate }) => {
  const [collapsed, setCollapsed] = useState(false);

  const menuItems = [
    {
      id: 'beranda',
      label: 'Beranda',
      icon: (
        <svg className="w-5 h-5 filter drop-shadow-sm" fill="none" stroke="currentColor" strokeWidth="2.2" viewBox="0 0 24 24">
          <path d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" strokeLinecap="round" strokeLinejoin="round" />
        </svg>
      ),
    },
    {
      id: 'periksa',
      label: 'Mulai Periksa',
      icon: (
        <svg className="w-5 h-5 text-forest-700/70 group-hover:text-forest-700 transition-colors" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
          <path d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" strokeLinecap="round" strokeLinejoin="round" />
          <path d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" strokeLinecap="round" strokeLinejoin="round" />
        </svg>
      ),
    },
    {
      id: 'riwayat',
      label: 'Riwayat',
      icon: (
        <svg className="w-5 h-5 text-forest-700/70 group-hover:text-forest-700 transition-colors" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
          <path d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" strokeLinecap="round" strokeLinejoin="round" />
        </svg>
      ),
    },
    {
      id: 'informasi',
      label: 'Informasi Penyakit',
      icon: (
        <svg className="w-5 h-5 text-forest-700/70 group-hover:text-forest-700 transition-colors" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
          <path d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" strokeLinecap="round" strokeLinejoin="round" />
        </svg>
      ),
    },
    {
      id: 'pengaturan',
      label: 'Pengaturan',
      icon: (
        <svg className="w-5 h-5 text-forest-700/70 group-hover:text-forest-700 transition-colors" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
          <path d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" strokeLinecap="round" strokeLinejoin="round" />
          <path d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" strokeLinecap="round" strokeLinejoin="round" />
        </svg>
      ),
    },
  ];

  return (
    <aside 
      // PERBAIKAN 1: transition-all -> transition-[width]
      className={`sticky top-4 md:top-6 lg:top-7 self-start shrink-0 flex flex-col justify-between p-5 rounded-[28px] bg-gradient-to-b from-white/95 via-emerald-50/80 to-mint-100/70 backdrop-blur-md border border-emerald-800/20 shadow-lg shadow-forest-900/15 ring-2 ring-white/80 transition-[width] duration-300 ease-in-out z-30 min-h-[calc(100vh-3.5rem)] ${collapsed ? 'w-[92px]' : 'w-64'}`} 
      data-purpose="main-navigation"
    >
      {/* Outer floating collapse pill indicator like reference */}
      <button 
        type="button"
        onClick={() => setCollapsed((prev) => !prev)}
        className="absolute -right-3 top-1/2 -translate-y-1/2 w-6 h-6 rounded-full bg-white shadow-md border border-white/80 flex items-center justify-center text-forest-700 hover:scale-105 cursor-pointer z-40 transition-transform"
        aria-label="Toggle Sidebar"
      >
        <svg className={`w-3.5 h-3.5 text-forest-700 transition-transform duration-300 ${collapsed ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" strokeWidth="2.5" viewBox="0 0 24 24">
          <path d="M15 19l-7-7 7-7" strokeLinecap="round" strokeLinejoin="round" />
        </svg>
      </button>

      {/* Top Branding & Navigation Links */}
      <div className="space-y-6">
        {/* ================= Logo Branding (DIUBAH) ================= */}
        {/* PERBAIKAN 2: gap-4 dihapus di parent, jarak dipindah ke div teks (ml-4) */}
        <div className="flex items-center px-2 pt-1">
          {/* DIUBAH: warna gradient & bayangan ikon */}
          <div className="flex items-center justify-center w-10 h-10 shrink-0 rounded-2xl bg-gradient-to-tr from-[#1F7A63] to-[#35A583] text-white shadow-md shadow-emerald-900/25">
            {/* DIUBAH: daun solid -> daun outline (stroke, tanpa fill) */}
            <svg
              className="w-[22px] h-[22px]"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="1.8"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M11 20A7 7 0 0 1 9.8 6.1C15.5 5 17 4.48 19 2c1 2 2 4.18 2 8 0 5.5-4.78 10-10 10Z" />
              <path d="M2 21c0-3 1.85-5.36 5.08-6C9.5 14.52 12 13 13 12" />
            </svg>
          </div>
          <div className={`overflow-hidden whitespace-nowrap transition-[max-width,opacity,margin] duration-300 ease-in-out ${collapsed ? 'max-w-0 opacity-0 ml-0' : 'max-w-[170px] opacity-100 ml-4'}`}>
            {/* DIUBAH: ukuran teks text-xl -> text-2xl, warna "Agri" dan "Sense" dipisah */}
            <span className="text-2xl font-bold tracking-tight leading-none text-[#0F3D33] flex items-center gap-1">
              Agri<span className="text-[#1F6B5C]">Sense</span>
            </span>
            {/* DIUBAH: ukuran 11px -> 13px, font-normal, warna lebih pudar */}
            <p className="text-[13px] font-normal text-[#7D9C95] leading-none mt-1.5">
              AI Plant Diagnostic
            </p>
          </div>
        </div>
        {/* ============== Akhir Logo Branding (DIUBAH) ============== */}

        {/* Main Navigation Menu */}
        <nav aria-label="Menu Utama" className="space-y-1.5 pt-1">
          {menuItems.map((item) => {
            const isActive = activeTab === item.id;
            return (
              <a
                key={item.id}
                title={item.label}
                href={`#${item.id}`}
                onClick={(e) => {
                  e.preventDefault();
                  if (onNavigate) onNavigate(item.id);
                }}
                className={`relative flex items-center justify-between px-4 py-3 text-sm rounded-2xl transition-all duration-200 group active:scale-[0.98] ${
                  isActive
                    ? 'font-semibold text-white shadow-lg'
                    : 'font-medium text-forest-800/80 hover:bg-white/70 hover:text-forest-900'
                }`}
                style={
                  isActive
                    ? {
                        background: 'linear-gradient(135deg, rgb(46, 94, 90) 0%, rgb(89, 169, 163) 100%)',
                        boxShadow: 'rgba(46, 94, 90, 0.3) 0px 4px 18px',
                      }
                    : undefined
                }
              >
                {/* PERBAIKAN 3: gap-3.5 dihapus, jarak dipindah ke span teks (ml-3.5) */}
                <div className="flex items-center">
                  <span className={`shrink-0 ${isActive ? 'text-white' : ''}`}>{item.icon}</span>
                  <span className={`tracking-wide overflow-hidden whitespace-nowrap transition-[max-width,opacity,margin] duration-300 ease-in-out ${collapsed ? 'max-w-0 opacity-0 ml-0' : 'max-w-[160px] opacity-100 ml-3.5'}`}>
                    {item.label}
                  </span>
                </div>
              </a>
            );
          })}
        </nav>
      </div>

      {/* Bottom Sidebar Section */}
      <div className="pt-4 border-t border-forest-800/10 space-y-3">
        {/* PERBAIKAN 4: gap-3 dihapus, jarak dipindah ke span teks (ml-3) */}
        <a 
          title="Help & Resources"
          className="flex items-center px-3.5 py-2.5 rounded-2xl text-xs font-semibold text-forest-700 hover:bg-white/80 transition-all" 
          href="#bantuan"
        >
          <div className="w-7 h-7 shrink-0 rounded-xl bg-forest-700/10 flex items-center justify-center text-forest-700">
            <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2.2" viewBox="0 0 24 24">
              <path d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" strokeLinecap="round" strokeLinejoin="round" />
            </svg>
          </div>
          <span className={`overflow-hidden whitespace-nowrap transition-[max-width,opacity,margin] duration-300 ease-in-out ${collapsed ? 'max-w-0 opacity-0 ml-0' : 'max-w-[160px] opacity-100 ml-3'}`}>
            Help &amp; Resources
          </span>
        </a>
        {/* PERBAIKAN 5: border selalu ada, hanya warnanya yang berubah (tidak snap) */}
        <div className={`rounded-2xl bg-white/50 text-center overflow-hidden whitespace-nowrap border transition-all duration-300 ease-in-out ${
          collapsed ? 'max-h-0 opacity-0 p-0 border-transparent' : 'max-h-16 opacity-100 p-3 border-white/60'
        }`}>
          <p className="text-[11px] font-medium text-forest-700/80">AgriSense</p>
        </div>
      </div>
    </aside>
  );
};