import React from 'react';
import { Upload, User, Users } from 'lucide-react';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';

interface StepWorkspaceNameProps {
  workspaceName: string;
  setWorkspaceName: (name: string) => void;
  workspaceIcon: string;
  setWorkspaceIcon: (icon: string) => void;
  companySize: string;
  setCompanySize: (size: string) => void;
}

const companySizes = ['2 - 5', '6 - 10', '11 - 20', '21 - 50', '51 - 100', '101 - 250', '250 - more'];

export function StepWorkspaceName({
  workspaceName,
  setWorkspaceName,
  workspaceIcon,
  setWorkspaceIcon,
  companySize,
  setCompanySize,
}: StepWorkspaceNameProps) {
  return (
    <div className="space-y-8 flex-1">
      {/* Icon Selection */}
      <div className="flex items-center gap-4">
        <div className="w-16 h-16 rounded-2xl bg-blue-600 flex items-center justify-center shadow-lg">
          <Users className="w-8 h-8 text-white" />
        </div>
        <div className="space-y-1.5">
          <p className="text-[13px] font-medium text-slate-400">Set an icon</p>
          <button
            type="button"
            className="flex items-center gap-2 px-4 py-2 bg-white border border-slate-200 rounded-lg text-sm font-semibold text-slate-700 hover:bg-slate-50 transition-all shadow-sm"
          >
            <Upload className="w-4 h-4" />
            Upload image
          </button>
        </div>
      </div>

      {/* Name Input */}
      <div className="space-y-2.5">
        <Label htmlFor="workspace-name" className="text-sm font-bold text-slate-900">
          Company Name <span className="text-rose-500">*</span>
        </Label>
        <div className="relative group">
          <div className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-blue-500 transition-colors">
            <User className="w-5 h-5" />
          </div>
          <Input
            id="workspace-name"
            placeholder="Enter name"
            value={workspaceName}
            onChange={(e) => setWorkspaceName(e.target.value)}
            className="h-14 pl-12 border-slate-200 rounded-xl bg-white shadow-sm text-[16px] focus-visible:ring-blue-500/10 focus-visible:border-blue-500 transition-all placeholder:text-slate-300"
          />
        </div>
      </div>

      {/* Company Size */}
      <div className="space-y-3.5">
        <Label className="text-sm font-bold text-slate-900">
          Company Size
        </Label>
        <div className="flex flex-wrap gap-2">
          {companySizes.map((size) => (
            <button
              key={size}
              type="button"
              onClick={() => setCompanySize(size)}
              className={cn(
                "px-5 py-2.5 rounded-full text-[13px] font-semibold transition-all border shadow-sm",
                companySize === size 
                  ? "bg-blue-600 border-blue-600 text-white shadow-blue-200" 
                  : "bg-white border-slate-200 text-slate-600 hover:border-slate-300 hover:bg-slate-50"
              )}
            >
              {size}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}
