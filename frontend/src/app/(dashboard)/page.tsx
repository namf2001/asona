"use client"

import Image from "next/image"
import { Plus } from "lucide-react"
import { Button } from "@/components/ui/button"

export default function DashboardPage() {
  return (
    <div className="flex flex-col items-center justify-center min-h-[calc(100vh-160px)]">
      <div className="w-full max-w-4xl p-12 bg-white rounded-[48px] border border-gray-100 shadow-[0_32px_64px_-16px_rgba(0,0,0,0.05)] flex flex-col items-center text-center animate-in fade-in slide-in-from-bottom-8 duration-700">
        
        {/* Illustration Container */}
        <div className="relative w-full aspect-[16/10] mb-8 overflow-hidden rounded-[32px]">
          <Image
            src="/images/new_project_illustration.png"
            alt="New Project Illustration"
            fill
            className="object-cover"
            priority
          />
        </div>

        {/* Text Content */}
        <h1 className="text-3xl font-bold text-gray-900 mb-4 tracking-tight">
          New Project
        </h1>
        <p className="text-gray-500 max-w-lg mb-8 text-lg leading-relaxed">
          Project management to ensure all relevant tasks are performed correctly, on time and to achieve the end goal.
        </p>

        {/* CTA Button */}
        <Button 
          size="lg" 
          className="h-14 px-8 rounded-full bg-[#3b82f6] hover:bg-blue-600 text-white font-semibold text-lg shadow-[0_8px_20px_rgba(59,130,246,0.3)] transition-all hover:scale-105 active:scale-95"
        >
          <Plus className="mr-2 size-6 stroke-[3px]" />
          Create project
        </Button>
      </div>
    </div>
  )
}
