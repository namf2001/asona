import React from 'react';
import { Shield, Building2, Users, Check } from 'lucide-react';

export interface RoleOption {
  id: string;
  label: string;
  desc: string;
  icon: React.ElementType;
}

export const roles: RoleOption[] = [
  { id: 'owner', label: 'Owner', desc: 'Full control of workspace and billing', icon: Shield },
  { id: 'admin', label: 'Administrator', desc: 'Manage members and settings', icon: Users },
  { id: 'member', label: 'Member', desc: 'Collaborate on projects and tasks', icon: Building2 },
];

export function RoleCard({
  role,
  selected,
  onSelect,
}: {
  role: RoleOption;
  selected: boolean;
  onSelect: () => void;
}) {
  const Icon = role.icon;
  return (
    <button
      id={`role-${role.id}`}
      type="button"
      onClick={onSelect}
      className="w-full flex items-center gap-4 rounded-xl border-2 p-4 text-left transition-all duration-150 focus:outline-none"
      style={{
        borderColor: selected ? '#10b981' : '#e5e7eb',
        background: selected ? '#f0fdf4' : '#f9fafb',
      }}
    >
      <div
        className="w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0 transition-colors"
        style={{
          background: selected ? '#10b981' : '#e5e7eb',
        }}
      >
        <Icon className="w-5 h-5" style={{ color: selected ? 'white' : '#6b7280' }} />
      </div>
      <div className="flex-1 min-w-0">
        <p className="text-sm font-semibold text-slate-800">{role.label}</p>
        <p className="text-xs text-slate-500 mt-0.5">{role.desc}</p>
      </div>
      <div
        className="w-5 h-5 rounded-full border-2 flex items-center justify-center flex-shrink-0 transition-all"
        style={{
          borderColor: selected ? '#10b981' : '#d1d5db',
          background: selected ? '#10b981' : 'transparent',
        }}
      >
        {selected && <Check className="w-3 h-3 text-white" />}
      </div>
    </button>
  );
}

interface StepRoleSelectionProps {
  selectedRole: string;
  setSelectedRole: (role: string) => void;
}

export function StepRoleSelection({ selectedRole, setSelectedRole }: StepRoleSelectionProps) {
  return (
    <div className="space-y-3 flex-1">
      {roles.map((role) => (
        <RoleCard
          key={role.id}
          role={role}
          selected={selectedRole === role.id}
          onSelect={() => setSelectedRole(role.id)}
        />
      ))}
    </div>
  );
}
