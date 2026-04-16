"use client";

import React, { useState, useTransition } from 'react';
import { useRouter } from 'next/navigation';
import { ChevronRight, Loader2, Check } from 'lucide-react';
import { toast } from 'sonner';

import { Button } from '@/components/ui/button';
import { StepWorkspaceName } from './_components/step-workspace-name';
import { StepInviteMembers } from './_components/step-invite-members';

import { completeOnboardAction } from '@/app/(auth)/actions';

const TOTAL_STEPS = 2;

export default function CreateWorkplacePage() {
  const router = useRouter();
  const [step, setStep] = useState(1);
  const [isPending, startTransition] = useTransition();

  // Step 1 state
  const [workspaceName, setWorkspaceName] = useState('');
  const [workspaceIcon, setWorkspaceIcon] = useState('');
  const [companySize, setCompanySize] = useState('');

  // Step 2 state
  const [emailInput, setEmailInput] = useState('');

  const handleNext = () => {
    if (step === 1) {
      if (!workspaceName.trim()) {
        toast.error('Please enter a workspace name');
        return;
      }
      if (!companySize) {
        toast.error('Please select your company size');
        return;
      }

      startTransition(async () => {
        try {
          // Normalize size for backend enum
          const sizeMapping: Record<string, string> = {
            '2 - 5': '2-5',
            '6 - 10': '6-10',
            '11 - 20': '11-20',
            '21 - 50': '21-50',
            '51 - 100': '51-100',
            '101 - 250': '101-250',
            '250 - more': '250-more'
          };

          const response = await fetch('/api/v1/workplaces', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              name: workspaceName,
              size: sizeMapping[companySize] || companySize,
              icon_url: workspaceIcon,
            }),
          });

          if (!response.ok) {
            throw new Error('Failed to create workplace');
          }

          toast.success('Workspace created successfully!');
          setStep(2);
        } catch (error) {
          toast.error('Something went wrong. Please try again.');
          console.error(error);
        }
      });
    } else {
      handleFinish();
    }
  };

  const handleFinish = () => {
    startTransition(async () => {
      try {
        const result = await completeOnboardAction();
        
        if (result.error) {
          throw new Error(result.error);
        }

        toast.success('Onboarding complete!');
        router.push('/');
      } catch (error: any) {
        toast.error(error.message || 'Failed to finalize onboarding');
        console.error(error);
      }
    });
  };

  return (
    <div className="flex flex-col w-full h-full font-inter">
      {/* Header & Progress */}
      <div className="mb-10">
        <div className="flex items-center justify-between mb-8">
          <span className="text-[12px] font-bold text-slate-400 uppercase tracking-[0.2em]">
            Step {step} of {TOTAL_STEPS}
          </span>
          <div className="flex gap-2">
            {[1, 2].map((i) => (
              <div
                key={i}
                className="h-1.5 rounded-full transition-all duration-500 ease-out"
                style={{
                  width: i === step ? '40px' : '12px',
                  background: i <= step ? '#2563eb' : '#e2e8f0',
                }}
              />
            ))}
          </div>
        </div>

        {step === 1 ? (
          <div className="space-y-2">
            <h1 className="text-[36px] font-bold text-slate-900 tracking-tight leading-tight">Create Your Workspace</h1>
            <p className="text-[15px] text-slate-500 font-medium">Your workspace contains your projects, teams and AI models.</p>
          </div>
        ) : (
          <div className="space-y-2">
            <h1 className="text-[36px] font-bold text-slate-900 tracking-tight leading-tight">Who else is on your workspace?</h1>
          </div>
        )}
      </div>

      <div className="flex-1 flex flex-col justify-between">
        <div className="mb-8">
          {step === 1 ? (
            <StepWorkspaceName
              workspaceName={workspaceName}
              setWorkspaceName={setWorkspaceName}
              workspaceIcon={workspaceIcon}
              setWorkspaceIcon={setWorkspaceIcon}
              companySize={companySize}
              setCompanySize={setCompanySize}
            />
          ) : (
            <StepInviteMembers
              emailInput={emailInput}
              setEmailInput={setEmailInput}
              onSkip={handleFinish}
            />
          )}
        </div>

        {/* Actions */}
        <div className="mt-auto space-y-4">
          <Button
            id="workspace-action-btn"
            onClick={handleNext}
            disabled={isPending}
            className="w-full h-14 bg-slate-400 hover:bg-slate-500 text-white rounded-xl font-bold text-[16px] shadow-sm transition-all active:scale-[0.98] disabled:opacity-50"
          >
            {isPending ? (
              <Loader2 className="w-5 h-5 animate-spin" />
            ) : step === 1 ? (
              'Create New Workspace'
            ) : (
              'Next'
            )}
          </Button>
          
          {step === 1 && (
            <button
              onClick={() => router.back()}
              className="w-full text-[15px] font-medium text-slate-400 hover:text-slate-600 transition-colors"
            >
              Go back
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
