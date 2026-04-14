"use client";

import React, { useState, useTransition, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { InputOTP, InputOTPGroup, InputOTPSlot } from '@/components/ui/input-otp';
import { Loader2 } from 'lucide-react';
import { toast } from 'sonner';
import { registerStep2Action } from '../actions';

function VerifyContent() {
  const [value, setValue] = useState("");
  const router = useRouter();
  const searchParams = useSearchParams();
  const [isPending, startTransition] = useTransition();

  const email = searchParams.get("email") || "";

  const handleVerify = () => {
    if (!email) {
      toast.error("Email is missing, please go back");
      return;
    }

    startTransition(async () => {
      const result = await registerStep2Action(email, value);
      
      if (result.error) {
        toast.error(result.error);
        return;
      }

      toast.success("Email verified!");
      router.push(`/create-password?email=${encodeURIComponent(email)}&otp=${encodeURIComponent(value)}`);
    });
  };

  return (
    <div className="flex flex-col items-center w-full h-full font-inter text-center justify-center">
      <div className="mb-8 w-full">
        <h1 className="text-[32px] font-semibold text-slate-900 tracking-tight">Check Your Email</h1>
        <p className="text-sm text-slate-500 mt-3 leading-relaxed">
          Enter the verification code sent to:<br/>
          <span className="font-semibold text-slate-800">{email}</span>
        </p>
      </div>

      <div className="py-6 flex justify-center w-full">
        <InputOTP maxLength={6} value={value} onChange={setValue} disabled={isPending}>
          <InputOTPGroup className="gap-2 sm:gap-3 flex w-full justify-center">
            <InputOTPSlot index={0} className="w-10 h-14 sm:w-12 sm:h-16 rounded-xl border-slate-200 text-xl font-semibold bg-white shadow-sm ring-offset-background focus:ring-emerald-500/20 focus:border-emerald-500 transition-all" />
            <InputOTPSlot index={1} className="w-10 h-14 sm:w-12 sm:h-16 rounded-xl border-slate-200 text-xl font-semibold bg-white shadow-sm ring-offset-background focus:ring-emerald-500/20 focus:border-emerald-500 transition-all" />
            <InputOTPSlot index={2} className="w-10 h-14 sm:w-12 sm:h-16 rounded-xl border-slate-200 text-xl font-semibold bg-white shadow-sm ring-offset-background focus:ring-emerald-500/20 focus:border-emerald-500 transition-all" />
            <InputOTPSlot index={3} className="w-10 h-14 sm:w-12 sm:h-16 rounded-xl border-slate-200 text-xl font-semibold bg-white shadow-sm ring-offset-background focus:ring-emerald-500/20 focus:border-emerald-500 transition-all" />
            <InputOTPSlot index={4} className="w-10 h-14 sm:w-12 sm:h-16 rounded-xl border-slate-200 text-xl font-semibold bg-white shadow-sm ring-offset-background focus:ring-emerald-500/20 focus:border-emerald-500 transition-all" />
            <InputOTPSlot index={5} className="w-10 h-14 sm:w-12 sm:h-16 rounded-xl border-slate-200 text-xl font-semibold bg-white shadow-sm ring-offset-background focus:ring-emerald-500/20 focus:border-emerald-500 transition-all" />
          </InputOTPGroup>
        </InputOTP>
      </div>

      <Button 
        className="w-full h-12 bg-[#9CA3AF] hover:bg-[#868e96] text-white rounded-xl font-semibold mt-6 shadow-md transition-all active:scale-[0.98]" 
        onClick={handleVerify}
        disabled={value.length !== 6 || isPending}
      >
        {isPending ? <Loader2 className="w-4 h-4 mr-2 animate-spin" /> : null}
        Verify
      </Button>

      <div className="mt-8 text-center text-xs text-slate-500 font-medium tracking-tight">
        Didn&apos;t receive the code? <button type="button" className="text-emerald-500 hover:text-emerald-600 font-semibold ml-1 bg-transparent border-none p-0 cursor-pointer transition-colors">Resend code (13)</button>
      </div>
    </div>
  );
}

export default function VerifyPage() {
  return (
    <Suspense fallback={<div className="flex w-full justify-center"><Loader2 className="animate-spin text-slate-400" /></div>}>
      <VerifyContent />
    </Suspense>
  );
}
