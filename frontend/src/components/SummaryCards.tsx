'use client';

import React from 'react';
import { SummaryMetric } from '../types';

const defaultMetrics: SummaryMetric[] = [
  {
    id: 'total',
    title: 'TOTAL PEMERIKSAAN',
    value: '128',
    changeText: '12% dari minggu lalu',
    changeType: 'up',
    icon: 'search',
  },
  {
    id: 'sehat',
    title: 'HASIL SEHAT',
    value: '96',
    changeText: '18% dari minggu lalu',
    changeType: 'up',
    icon: 'shield',
  },
  {
    id: 'perhatian',
    title: 'PERLU PERHATIAN',
    value: '32',
    changeText: '6% dari minggu lalu',
    changeType: 'down',
    icon: 'alert',
  },
];

export const SummaryCards: React.FC = () => {
  return (
    <section className="grid grid-cols-1 md:grid-cols-3 gap-5" data-purpose="metrics-summary">
      {defaultMetrics.map((metric) => (
        <div
          key={metric.id}
          className="relative overflow-hidden p-5 rounded-3xl backdrop-blur-md flex items-center gap-4 transition-all duration-300 hover:shadow-lg"
          style={{
            background: 'linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(50, 145, 135, 0.16) 55%, rgba(50, 145, 135, 0.3) 100%)',
            border: '1.5px solid rgba(50, 145, 135, 0.4)',
            boxShadow: 'rgba(50, 145, 135, 0.15) 0px 8px 24px -4px',
          }}
        >
          {/* Metric Icon Container */}
          <div
            className="w-14 h-14 rounded-2xl flex items-center justify-center shrink-0 shadow-sm"
            style={{
              background: 'linear-gradient(135deg, rgba(255, 255, 255, 0.9) 0%, rgba(50, 145, 135, 0.22) 100%)',
              border: '1px solid rgba(50, 145, 135, 0.35)',
              color: '#20635B',
            }}
          >
            {metric.icon === 'search' && (
              <svg className="w-7 h-7" fill="none" stroke="currentColor" strokeWidth="2.2" viewBox="0 0 24 24">
                <path d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" strokeLinecap="round" strokeLinejoin="round" />
                <path d="M10 10l.01 0" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
            )}
            {metric.icon === 'shield' && (
              <svg className="w-7 h-7" fill="none" stroke="currentColor" strokeWidth="2.2" viewBox="0 0 24 24">
                <path d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
            )}
            {metric.icon === 'alert' && (
              <svg className="w-7 h-7" fill="none" stroke="currentColor" strokeWidth="2.2" viewBox="0 0 24 24">
                <path d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
            )}
          </div>

          {/* Metric Details */}
          <div className="space-y-0.5">
            <p className="text-xs font-bold tracking-wide uppercase" style={{ color: '#184F48' }}>
              {metric.title}
            </p>
            <div className="flex items-baseline gap-2.5">
              <span className="text-3xl font-extrabold tracking-tight" style={{ color: '#143d38' }}>
                {metric.value}
              </span>
              <span className="inline-flex items-center text-xs font-bold text-emerald-700">
                {metric.changeType === 'up' ? (
                  <svg className="w-3.5 h-3.5 mr-0.5" fill="none" stroke="currentColor" strokeWidth="2.5" viewBox="0 0 24 24">
                    <path d="M5 10l7-7m0 0l7 7m-7-7v18" strokeLinecap="round" strokeLinejoin="round" />
                  </svg>
                ) : (
                  <svg className="w-3.5 h-3.5 mr-0.5 text-emerald-600" fill="none" stroke="currentColor" strokeWidth="2.5" viewBox="0 0 24 24">
                    <path d="M19 14l-7 7m0 0l-7-7m7 7V3" strokeLinecap="round" strokeLinejoin="round" />
                  </svg>
                )}
                {metric.changeText.split(' ')[0]}{' '}
                <span className="font-medium ml-1" style={{ color: '#20635B' }}>
                  {metric.changeText.substring(metric.changeText.indexOf(' ') + 1)}
                </span>
              </span>
            </div>
          </div>
        </div>
      ))}
    </section>
  );
};
