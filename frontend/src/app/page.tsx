'use client';

import React from 'react';
import { Sidebar } from '../components/Sidebar';
import { Header } from '../components/Header';
import { Greeting } from '../components/Greeting';
import { HeroBanner } from '../components/HeroBanner';
import { SummaryCards } from '../components/SummaryCards';
import { ActivityChart } from '../components/ActivityChart';
import { DiseaseEducation } from '../components/DiseaseEducation';
import { InspectionHistory } from '../components/InspectionHistory';

export default function DashboardPage() {
  return (
    <>
      {/* Navigation Sidebar */}
      <Sidebar />

      {/* Main Content Area */}
      <main className="flex-1 space-y-6 min-w-0" data-purpose="dashboard-main-view">
        {/* Top Header */}
        <Header />

        {/* User Greeting Bar */}
        <Greeting />

        {/* Hero Banner with Animated Mascot */}
        <HeroBanner />

        {/* Summary Statistics Metric Row */}
        <SummaryCards />

        {/* Middle Section: Bar Chart & Disease Education */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-7">
          <div className="lg:col-span-5">
            <ActivityChart />
          </div>
          <div className="lg:col-span-7">
            <DiseaseEducation />
          </div>
        </div>

        {/* Recent Inspection History */}
        <InspectionHistory />
      </main>
    </>
  );
}
