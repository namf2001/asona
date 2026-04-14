import React from 'react';
import Image from 'next/image';

import { AsonaLogo, AsonaLoginGroupItems, AsonaBg } from '@/../public/images';

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen relative overflow-hidden flex bg-[#0f172a] font-inter">
      {/* Background Image Optimization (per Next.js best practices) */}
      <div className="absolute inset-0 z-0">
        <Image 
          src={AsonaBg} 
          alt="Asona Background" 
          fill 
          className="object-cover" 
          priority 
          sizes="100vw"
          placeholder="blur"
        />
        {/* Subtle overlay to soften the image and improve form readability if needed */}
        <div className="absolute inset-0 bg-slate-900/10 backdrop-blur-[2px]"></div>
      </div>

      {/* Main Container */}
      <div className="relative z-10 w-full max-w-[1600px] mx-auto flex flex-col lg:flex-row min-h-screen">
        
        {/* Left Side: Illustration & Logo */}
        <div className="w-full lg:w-[55%] flex flex-col relative p-8 lg:p-16">
          <div className="flex items-center gap-[18px] z-10">
            {/* Logo Wrapper */}
            <div className="relative w-16 h-16 lg:w-20 lg:h-20 flex items-center justify-center">
              {/* Specialized Background Layer (prevents image resizing issues) */}
              <div 
                className="absolute inset-0 rounded-[28px] shadow-[0_12px_30px_-10px_rgba(0,0,0,0.2)] border border-white/50"
                style={{ background: 'linear-gradient(40deg, #ffffff 0%, #f3f4f6 100%)' }}
              />
              
              {/* Logo Image Layer */}
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

            <h2 className="text-[42px] text-white font-semibold flex items-baseline">
              Asona<span className="font-bold">.ai</span>
            </h2>
          </div>

          <div className="flex-1 w-full relative min-h-[400px] hidden lg:flex items-center justify-center">
            {/* The Login Group Items Illustration */}
            <div className="relative w-full h-full max-w-[800px] max-h-[600px] flex items-center justify-center">
              <Image 
                src={AsonaLoginGroupItems} 
                alt="Asona Illustration" 
                fill 
                className="object-contain drop-shadow-2xl" 
                priority 
              />
            </div>
          </div>
        </div>

        {/* Right Side: Form Container */}
        <div className="w-full lg:w-[45%] flex items-center justify-center px-4 py-8 lg:p-12 z-10">
          <div className="bg-white w-full max-w-[480px] rounded-[32px] p-8 sm:p-10 shadow-[0_20px_60px_-15px_rgba(0,0,0,0.15)] border border-white/50 backdrop-blur-sm self-center">
            {children}
          </div>
        </div>

      </div>
    </div>
  );
}