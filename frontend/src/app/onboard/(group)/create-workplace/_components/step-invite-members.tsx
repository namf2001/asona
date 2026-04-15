import React, { useState } from 'react';
import { Link, Check } from 'lucide-react';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { cn } from '@/lib/utils';

interface StepInviteMembersProps {
  emailInput: string;
  setEmailInput: (email: string) => void;
  onSkip: () => void;
}

export function StepInviteMembers({
  emailInput,
  setEmailInput,
  onSkip,
}: StepInviteMembersProps) {
  const [copied, setCopied] = useState(false);

  const handleCopyLink = () => {
    navigator.clipboard.writeText('https://asona.io/join/temp-token');
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="space-y-6 flex-1">
      <div className="space-y-2">
        <p className="text-[13px] font-medium text-slate-400">Add member, coworker by email</p>
        <div className="relative group">
          <Textarea
            id="invite-emails"
            placeholder="e.g somebody@example.com"
            value={emailInput}
            onChange={(e) => setEmailInput(e.target.value)}
            className="min-h-[220px] p-5 border-2 border-slate-200 rounded-xl bg-white shadow-sm text-[16px] focus-visible:ring-blue-500/10 focus-visible:border-blue-500 transition-all placeholder:text-slate-300 resize-none"
          />
        </div>
      </div>

      <div className="flex items-center gap-3">
        <button
          type="button"
          onClick={handleCopyLink}
          className={cn(
            "flex-1 h-14 rounded-xl font-bold flex items-center justify-center gap-2 transition-all border-2",
            copied 
              ? "bg-white border-blue-500 text-blue-600 shadow-sm" 
              : "bg-white border-slate-200 text-slate-800 hover:border-slate-300 shadow-sm"
          )}
        >
          {copied ? <Check className="w-5 h-5" /> : <Link className="w-5 h-5 rotate-45" />}
          {copied ? 'Copied!' : 'Copy Invite Link'}
        </button>
      </div>

      <div className="text-center">
        <button
          onClick={onSkip}
          className="text-[15px] font-medium text-blue-500 hover:text-blue-600 transition-colors"
        >
          Skip this step
        </button>
      </div>
    </div>
  );
}
