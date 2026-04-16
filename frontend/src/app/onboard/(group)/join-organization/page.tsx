"use client";

import React, { useState, useTransition } from 'react';
import { useRouter } from 'next/navigation';
import { Hash, Loader2, ArrowRight } from 'lucide-react';
import { toast } from 'sonner';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

import { completeOnboardAction } from '@/app/(auth)/actions';

// JoinOrganizationPage is the onboarding screen where users join an existing workspace via invite code.
export default function JoinOrganizationPage() {
  const router = useRouter();
  const [inviteCode, setInviteCode] = useState('');
  const [isPending, startTransition] = useTransition();

  function handleJoin() {
    if (!inviteCode.trim()) {
      toast.error('Please enter an invite code');
      return;
    }
    startTransition(async () => {
      try {
        // TODO: call API to join organization (if needed)
        // For now, we just complete the onboarding
        const result = await completeOnboardAction();
        
        if (result.error) {
          throw new Error(result.error);
        }

        toast.success('Joined organization successfully!');
        router.push('/');
      } catch (error: any) {
        toast.error(error.message || 'Failed to join organization');
        console.error(error);
      }
    });
  }

  return (
    <div className="flex flex-col w-full h-full font-inter">
      {/* Header */}
      <div className="mb-8">
        <div className="mb-5">
          <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-widest">
            Join Organization
          </span>
        </div>
        <h1 className="text-[28px] font-semibold text-slate-900 tracking-tight">Join Organization</h1>
        <p className="text-sm text-slate-500 mt-1.5">
          Enter the invite code shared by your team to join their workspace.
        </p>
      </div>

      {/* Invite Code Input */}
      <div className="flex-1 space-y-6">
        <div className="space-y-2">
          <Label htmlFor="invite-code" className="text-xs font-semibold text-slate-800">
            Invite Code <span className="text-red-500">*</span>
          </Label>
          <div className="relative">
            <Hash className="absolute left-3.5 top-3.5 h-4 w-4 text-slate-400" />
            <Input
              id="invite-code"
              placeholder="e.g. TEAM-XXXX-XXXX"
              value={inviteCode}
              onChange={(e) => setInviteCode(e.target.value.toUpperCase())}
              onKeyDown={(e) => e.key === 'Enter' && handleJoin()}
              className="h-12 pl-10 border-slate-200 rounded-xl bg-white shadow-sm text-sm font-mono tracking-widest focus:ring-blue-500/20 focus:border-blue-500 transition-all"
            />
          </div>
          <p className="text-[11px] text-slate-400">
            Ask your workspace admin for the invite code or link.
          </p>
        </div>

        {/* Info card */}
        <div className="rounded-xl border border-blue-100 bg-blue-50/60 p-4">
          <p className="text-xs font-semibold text-blue-700 mb-1">How it works</p>
          <ul className="text-xs text-blue-600 space-y-1 leading-relaxed">
            <li>1. Get an invite code from your workspace admin</li>
            <li>2. Enter the code above and click Join</li>
            <li>3. You&apos;ll be added to the workspace instantly</li>
          </ul>
        </div>
      </div>

      {/* Actions */}
      <div className="mt-8 flex flex-col gap-3">
        <Button
          id="join-org-btn"
          onClick={handleJoin}
          disabled={isPending || !inviteCode.trim()}
          className="w-full h-12 bg-blue-500 hover:bg-blue-600 text-white rounded-xl font-semibold shadow-md transition-all active:scale-[0.98] flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isPending ? <Loader2 className="w-4 h-4 animate-spin" /> : <ArrowRight className="w-4 h-4" />}
          {isPending ? 'Joining…' : 'Join Organization'}
        </Button>

        <button
          id="back-join-btn"
          type="button"
          onClick={() => router.back()}
          className="text-sm text-slate-400 hover:text-slate-600 font-medium transition-colors"
        >
          ← Go back
        </button>
      </div>
    </div>
  );
}
