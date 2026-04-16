"use client";

import React, { useState, useTransition } from 'react';
import { useRouter } from 'next/navigation';
import { Search, Check, X, Loader2, MessageCircle, AtSign } from 'lucide-react';
import { toast } from 'sonner';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

import { completeOnboardAction } from '@/app/(auth)/actions';

// MockFriend is a mock data type for a friend suggestion.
interface MockFriend {
  id: string;
  name: string;
  username: string;
  avatar: string;
  mutual: number;
}

const MOCK_FRIENDS: MockFriend[] = [
  { id: '1', name: 'Alice Johnson', username: 'alice.j', avatar: 'AJ', mutual: 5 },
  { id: '2', name: 'Bob Smith', username: 'bob.smith', avatar: 'BS', mutual: 3 },
  { id: '3', name: 'Carol White', username: 'carol.w', avatar: 'CW', mutual: 8 },
  { id: '4', name: 'Dan Lee', username: 'dan.lee', avatar: 'DL', mutual: 2 },
  { id: '5', name: 'Eva Chen', username: 'eva.chen', avatar: 'EC', mutual: 12 },
];

const AVATAR_COLORS = [
  'bg-rose-400',
  'bg-sky-400',
  'bg-violet-400',
  'bg-amber-400',
  'bg-teal-400',
];

// ChatWithFriendPage is the multi-step onboarding flow for starting a chat with friends.
export default function ChatWithFriendPage() {
  const router = useRouter();
  const [step, setStep] = useState(1);
  const [isPending, startTransition] = useTransition();

  // Step 1: Search
  const [query, setQuery] = useState('');

  // Step 2: Select friends
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());

  // Step 3: Preview selected friend
  const [previewId, setPreviewId] = useState<string>('');

  const TOTAL_STEPS = 3;

  const filteredFriends = MOCK_FRIENDS.filter(
    (f) =>
      f.name.toLowerCase().includes(query.toLowerCase()) ||
      f.username.toLowerCase().includes(query.toLowerCase())
  );

  function toggleSelect(id: string) {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  }

  function handleNext() {
    if (step === 2 && selectedIds.size === 0) {
      toast.error('Please select at least one friend');
      return;
    }
    if (step === 2) {
      // Pick first selected for preview
      setPreviewId([...selectedIds][0]);
    }
    if (step < TOTAL_STEPS) setStep((s) => s + 1);
  }

  function handleStartChat() {
    startTransition(async () => {
      try {
        // TODO: call API to initiate chat
        const result = await completeOnboardAction();
        
        if (result.error) {
          throw new Error(result.error);
        }

        toast.success('Chat started successfully!');
        router.push('/');
      } catch (error: any) {
        toast.error(error.message || 'Failed to start chat');
        console.error(error);
      }
    });
  }

  const previewFriend = MOCK_FRIENDS.find((f) => f.id === previewId);
  const previewColorIdx = MOCK_FRIENDS.findIndex((f) => f.id === previewId);

  return (
    <div className="flex flex-col w-full h-full font-inter">
      {/* Header & Progress */}
      <div className="mb-7">
        <div className="flex items-center justify-between mb-5">
          <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-widest">
            Step {step} of {TOTAL_STEPS}
          </span>
          <div className="flex gap-1.5">
            {Array.from({ length: TOTAL_STEPS }).map((_, i) => (
              <div
                key={i}
                className="h-1.5 rounded-full transition-all duration-300"
                style={{
                  width: i < step ? '28px' : '10px',
                  background: i < step ? '#8b5cf6' : '#e2e8f0',
                }}
              />
            ))}
          </div>
        </div>

        {step === 1 && (
          <>
            <h1 className="text-[28px] font-semibold text-slate-900 tracking-tight">Chat With Friend</h1>
            <p className="text-sm text-slate-500 mt-1.5">Search for friends to start a conversation.</p>
          </>
        )}
        {step === 2 && (
          <>
            <h1 className="text-[28px] font-semibold text-slate-900 tracking-tight">Chat With Friend</h1>
            <p className="text-sm text-slate-500 mt-1.5">Select people you want to chat with.</p>
          </>
        )}
        {step === 3 && previewFriend && (
          <>
            <h1 className="text-[28px] font-semibold text-slate-900 tracking-tight">Chat With Friend</h1>
            <p className="text-sm text-slate-500 mt-1.5">Ready to start chatting!</p>
          </>
        )}
      </div>

      {/* Step 1: Search */}
      {step === 1 && (
        <div className="flex-1 space-y-4">
          <div className="space-y-2">
            <Label htmlFor="search-friend" className="text-xs font-semibold text-slate-800">
              Search by name or username
            </Label>
            <div className="relative">
              <Search className="absolute left-3.5 top-3.5 h-4 w-4 text-slate-400" />
              <Input
                id="search-friend"
                placeholder="e.g. Alice Johnson or @alice.j"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                className="h-12 pl-10 border-slate-200 rounded-xl bg-white shadow-sm text-sm focus:ring-violet-500/20 focus:border-violet-500 transition-all"
              />
            </div>
          </div>

          <div className="space-y-2 max-h-[260px] overflow-y-auto pr-1">
            {filteredFriends.map((friend, idx) => (
              <FriendRow
                key={friend.id}
                friend={friend}
                colorClass={AVATAR_COLORS[idx % AVATAR_COLORS.length]}
                selected={false}
                showCheckbox={false}
                onSelect={() => {
                  setSelectedIds(new Set([friend.id]));
                  setStep(2);
                }}
              />
            ))}
            {filteredFriends.length === 0 && (
              <div className="text-center py-8">
                <Search className="w-8 h-8 text-slate-200 mx-auto mb-2" />
                <p className="text-sm text-slate-400">No results for &quot;{query}&quot;</p>
              </div>
            )}
          </div>
        </div>
      )}

      {/* Step 2: Select friends */}
      {step === 2 && (
        <div className="flex-1 space-y-4">
          {selectedIds.size > 0 && (
            <div className="flex gap-2 flex-wrap pb-2 border-b border-slate-100">
              {[...selectedIds].map((id) => {
                const f = MOCK_FRIENDS.find((fr) => fr.id === id)!;
                const ci = MOCK_FRIENDS.findIndex((fr) => fr.id === id);
                return (
                  <div key={id} className="flex items-center gap-1.5 bg-violet-50 border border-violet-200 rounded-full px-3 py-1">
                    <span className={`w-5 h-5 rounded-full ${AVATAR_COLORS[ci % AVATAR_COLORS.length]} flex items-center justify-center text-[9px] font-bold text-white`}>
                      {f.avatar}
                    </span>
                    <span className="text-xs font-medium text-violet-700">{f.name.split(' ')[0]}</span>
                    <button type="button" onClick={() => toggleSelect(id)}>
                      <X className="w-3 h-3 text-violet-400 hover:text-violet-600" />
                    </button>
                  </div>
                );
              })}
            </div>
          )}

          <div className="space-y-2 max-h-[240px] overflow-y-auto pr-1">
            {MOCK_FRIENDS.map((friend, idx) => (
              <FriendRow
                key={friend.id}
                friend={friend}
                colorClass={AVATAR_COLORS[idx % AVATAR_COLORS.length]}
                selected={selectedIds.has(friend.id)}
                showCheckbox={true}
                onSelect={() => toggleSelect(friend.id)}
              />
            ))}
          </div>
        </div>
      )}

      {/* Step 3: Profile preview */}
      {step === 3 && previewFriend && (
        <div className="flex-1 flex flex-col items-center justify-center space-y-4">
          <div
            className={`w-20 h-20 rounded-full ${AVATAR_COLORS[previewColorIdx % AVATAR_COLORS.length]} flex items-center justify-center shadow-lg`}
          >
            <span className="text-2xl font-bold text-white">{previewFriend.avatar}</span>
          </div>
          <div className="text-center">
            <p className="text-xl font-semibold text-slate-900">{previewFriend.name}</p>
            <p className="text-sm text-slate-400 flex items-center justify-center gap-1 mt-1">
              <AtSign className="w-3.5 h-3.5" />
              {previewFriend.username}
            </p>
          </div>
          <div className="flex gap-4 text-center">
            <div className="px-4 py-2 bg-violet-50 rounded-xl">
              <p className="text-lg font-bold text-violet-600">{previewFriend.mutual}</p>
              <p className="text-xs text-slate-500">Mutual connections</p>
            </div>
          </div>
          {selectedIds.size > 1 && (
            <p className="text-xs text-slate-400 text-center">
              + {selectedIds.size - 1} more friend{selectedIds.size > 2 ? 's' : ''} selected
            </p>
          )}
        </div>
      )}

      {/* Actions */}
      <div className="mt-8 flex flex-col gap-3">
        {step < TOTAL_STEPS ? (
          <Button
            id="chat-next-btn"
            onClick={step === 1 ? () => setStep(2) : handleNext}
            className="w-full h-12 bg-violet-500 hover:bg-violet-600 text-white rounded-xl font-semibold shadow-md transition-all active:scale-[0.98]"
          >
            Continue
          </Button>
        ) : (
          <Button
            id="start-chat-btn"
            onClick={handleStartChat}
            disabled={isPending}
            className="w-full h-12 bg-violet-500 hover:bg-violet-600 text-white rounded-xl font-semibold shadow-md transition-all active:scale-[0.98] flex items-center justify-center gap-2"
          >
            {isPending ? <Loader2 className="w-4 h-4 animate-spin" /> : <MessageCircle className="w-4 h-4" />}
            {isPending ? 'Starting chat…' : 'Start Chatting'}
          </Button>
        )}

        {step > 1 ? (
          <button
            id="chat-back-btn"
            type="button"
            onClick={() => setStep((s) => s - 1)}
            className="text-sm text-slate-400 hover:text-slate-600 font-medium transition-colors"
          >
            ← Back
          </button>
        ) : (
          <button
            id="back-from-chat-btn"
            type="button"
            onClick={() => router.back()}
            className="text-sm text-slate-400 hover:text-slate-600 font-medium transition-colors"
          >
            ← Go back
          </button>
        )}
      </div>
    </div>
  );
}

// FriendRow renders a single friend list item.
function FriendRow({
  friend,
  colorClass,
  selected,
  showCheckbox,
  onSelect,
}: {
  friend: MockFriend;
  colorClass: string;
  selected: boolean;
  showCheckbox: boolean;
  onSelect: () => void;
}) {
  return (
    <button
      id={`friend-${friend.id}`}
      type="button"
      onClick={onSelect}
      className="w-full flex items-center gap-3 rounded-xl border px-4 py-3 text-left transition-all duration-150 focus:outline-none"
      style={{
        borderColor: selected ? '#8b5cf6' : '#f1f5f9',
        background: selected ? '#f5f3ff' : '#f8fafc',
      }}
    >
      <div className={`w-9 h-9 rounded-full ${colorClass} flex items-center justify-center flex-shrink-0`}>
        <span className="text-xs font-bold text-white">{friend.avatar}</span>
      </div>
      <div className="flex-1 min-w-0">
        <p className="text-sm font-semibold text-slate-800 truncate">{friend.name}</p>
        <p className="text-xs text-slate-400 truncate">@{friend.username} · {friend.mutual} mutual</p>
      </div>
      {showCheckbox && (
        <div
          className="w-5 h-5 rounded-full border-2 flex items-center justify-center flex-shrink-0 transition-all"
          style={{
            borderColor: selected ? '#8b5cf6' : '#d1d5db',
            background: selected ? '#8b5cf6' : 'transparent',
          }}
        >
          {selected && <Check className="w-3 h-3 text-white" />}
        </div>
      )}
    </button>
  );
}
