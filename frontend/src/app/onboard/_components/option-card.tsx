"use client";

import React from 'react';
import Image, { StaticImageData } from 'next/image';
import { cn } from "@/lib/utils";
import {
  Card,
  CardHeader,
  CardContent,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";

export interface OnboardOption {
  id: string;
  label: string;
  sublabel: string;
  image: StaticImageData | string;
  href: string;
}

interface OptionCardProps {
  option: OnboardOption;
  onClick: () => void;
  className?: string;
}

export function OptionCard({ option, onClick, className }: OptionCardProps) {
  return (
    <button
      id={`onboard-option-${option.id}`}
      type="button"
      onClick={onClick}
      className={cn(
        "group outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 rounded-3xl transition-all duration-300 active:scale-95",
        className
      )}
    >
      <Card className="relative overflow-hidden flex flex-col items-center justify-center text-center p-0 w-[180px] h-[220px] sm:w-[280px] sm:h-[296px] border-white/10 bg-white hover:bg-emerald-500/10 hover:border-emerald-500/50 transition-all duration-300 shadow-[0_8px_32px_0_rgba(0,0,0,0.1)] hover:shadow-[0_8px_32px_0_rgba(16,185,129,0.15)] rounded-[32px]">
        
        {/* Illustration Section - Using div instead of CardHeader to avoid grid constraints */}
        <div className="flex flex-col items-center justify-center flex-1 w-full pt-8">
          <div className="relative size-44 transition-transform duration-500 ease-out group-hover:scale-110 group-hover:-translate-y-1">
            <Image
              src={option.image}
              alt={option.label}
              width={176}
              height={176}
              className="object-contain filter drop-shadow-[0_10px_10px_rgba(0,0,0,0.2)]"
              priority
            />
          </div>
          <div className="absolute bottom-6 opacity-0 group-hover:opacity-100 transition-all duration-300 translate-y-2 group-hover:translate-y-0"/>
        </div>

        {/* Labels Section */}
        <CardContent className="p-6 pt-0 flex flex-col gap-0.5">
          <CardTitle className="text-sm font-bold group-hover:text-emerald-400 transition-colors uppercase tracking-wider leading-none">
            {option.label}
          </CardTitle>
          <CardDescription className="text-sm font-bold group-hover:text-emerald-400 transition-colors uppercase tracking-wider leading-none">
            {option.sublabel}
          </CardDescription>
        </CardContent>
      </Card>
    </button>
  );
}
