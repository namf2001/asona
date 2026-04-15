import React from 'react';
import Image from 'next/image';

import { AsonaLogo, AsonaBg, AsonaCreateYourWorkspace } from '@/../public/images';

// OnBoardLayout wraps all multi-step onboarding screens (Create Workspace, Join, etc.) with the shared split-panel layout.
export default function OnBoardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen relative overflow-hidden flex bg-[#0f172a] font-inter">
      {/* Absolute Background Layer */}
      <div className="absolute inset-0 z-0">
        <Image
          src={AsonaCreateYourWorkspace}
          alt="Asona Background"
          fill
          className="object-cover"
          priority
          quality={100}
        />
      </div>

      {/* Main Split Layout Container */}
      <div className="relative z-10 w-full max-w-[1600px] mx-auto flex flex-col lg:flex-row min-h-screen">
        {/* Left Side: Branding & Illustration */}
        <div className="w-full lg:w-[55%] flex flex-col relative p-8 lg:p-16">
          {/* Logo & Platform Name */}
          <div className="flex items-center gap-[18px] z-10">
            <div className="relative w-16 h-16 lg:w-20 lg:h-20 flex items-center justify-center">
              <div
                className="absolute inset-0 rounded-[28px] shadow-[0_12px_30px_-10px_rgba(0,0,0,0.2)] border border-white/50"
                style={{ background: 'linear-gradient(40deg, #ffffff 0%, #f3f4f6 100%)' }}
              />
              <div className="relative z-10">
                <Image
                  src={AsonaLogo}
                  alt="Asona Logo"
                  width={44}
                  height={44}
                  className="object-contain lg:w-[52px] lg:h-[52px]"
                  priority
                />
              </div>
            </div>
            <h2 className="text-[42px] text-white font-semibold flex items-baseline tracking-tight">
              Asona<span className="font-bold">.ai</span>
            </h2>
          </div>

          {/* Foreground Illustration - Restored and Fixed */}
          <div className="flex-1 w-full relative min-h-[400px]"/>
        </div>

        {/* Right Side: Step Form Container */}
        <div className="w-full lg:w-[45%] flex items-center justify-center px-4 py-8 lg:p-12 z-10">
          <div className="bg-white w-full max-w-[480px] rounded-[32px] p-8 sm:p-10 shadow-[0_25px_60px_-15px_rgba(0,0,0,0.25)] border border-white/60 self-center">
            {children}
          </div>
        </div>
      </div>      
    </div>
  );
}