"use client";

import React from 'react';
import { useRouter } from 'next/navigation';
import Image from 'next/image';
import {
  Card,
} from "@/components/ui/card";

import {
  AsonaUserPick1,
  AsonaUserPick2,
  AsonaUserPick3,
  AsonaBg,
} from '@/../public/images';

import { OptionCard, OnboardOption } from './_components/option-card';

const options: OnboardOption[] = [
  {
    id: 'create-workspace',
    label: 'Create',
    sublabel: 'Workspace',
    image: AsonaUserPick1,
    href: '/onboard/create-workplace',
  },
  {
    id: 'join-organization',
    label: 'Join',
    sublabel: 'Organization',
    image: AsonaUserPick2,
    href: '/onboard/join-organization',
  },
  {
    id: 'chat-with-friend',
    label: 'Chat',
    sublabel: 'With Friend',
    image: AsonaUserPick3,
    href: '/onboard/chat-with-friend',
  },
];

export default function OnboardPage() {
  const router = useRouter();

  return (
    <div className="min-h-screen relative flex items-center justify-center font-inter">
      {/* Background Image Layer */}
      <div className="absolute inset-0 z-0">
        <Image 
          src={AsonaBg} 
          alt="Asona Background" 
          fill 
          className="object-cover" 
          priority 
          quality={100}
        />
      </div>

      {/* shadcn Card Container */}
      <Card className="relative z-10 flex flex-col items-center max-w-6xl mx-auto px-12 py-16 bg-white/10 rounded-[48px] backdrop-blur-3xl border-white/10 shadow-2xl">
        {/* Title Section */}
        <div className="text-center mb-16 space-y-4 w-full px-4"> 
          <h1 className="font-bold text-white tracking-tight text-[32px] md:text-[42px] drop-shadow-2xl text-balance leading-tight">
            Already have a Workspace?
          </h1>
        </div>

        {/* Options grid */}
        <div className="flex flex-wrap justify-center gap-6 md:gap-10 w-full">
          {options.map((opt) => (
            <OptionCard
              key={opt.id}
              option={opt}
              onClick={() => router.push(opt.href)}
            />
          ))}
        </div>
      </Card>
    </div>
  );
}