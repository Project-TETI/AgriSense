import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'AgriSense - Dasbor AI Kesehatan Tanaman',
  description: 'Deteksi dini kesehatan tanaman padi, cabai & tomat Anda dengan akurasi AI terkini.',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="id">
      <body className="flex min-h-screen overflow-x-hidden antialiased select-none p-4 md:p-6 lg:p-7 gap-6 bg-gradient-to-br from-[#E8F5F2] via-[#DCEDE8] to-[#E5F5F0]">
        {children}
      </body>
    </html>
  );
}
