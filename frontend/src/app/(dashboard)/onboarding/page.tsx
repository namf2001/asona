"use client"

import * as React from "react"
import { useRouter } from "next/navigation"
import {
  FolderPlus,
  Users,
  MessageSquare,
  ArrowRight,
  Sparkles,
  LayoutDashboard,
  Zap,
} from "lucide-react"
import { cn } from "@/lib/utils"

/* ─── Action Card Data ─── */

interface ActionCard {
  id: string
  icon: React.ElementType
  title: string
  description: string
  gradient: string
  iconColor: string
  href: string
  badge?: string
}

const actions: ActionCard[] = [
  {
    id: "create-project",
    icon: FolderPlus,
    title: "Create a Project",
    description:
      "Start a new project from scratch with customizable boards, tasks and deadlines.",
    gradient: "from-emerald-500/10 via-emerald-500/5 to-transparent",
    iconColor: "text-emerald-600 bg-emerald-100",
    href: "/",
    badge: "Recommended",
  },
  {
    id: "join-team",
    icon: Users,
    title: "Join a Team",
    description:
      "Accept an invite or search for your organization to start collaborating instantly.",
    gradient: "from-blue-500/10 via-blue-500/5 to-transparent",
    iconColor: "text-blue-600 bg-blue-100",
    href: "/",
  },
  {
    id: "explore-chat",
    icon: MessageSquare,
    title: "Start a Conversation",
    description:
      "Chat in real time with teammates, share files and keep discussions organized.",
    gradient: "from-violet-500/10 via-violet-500/5 to-transparent",
    iconColor: "text-violet-600 bg-violet-100",
    href: "/",
  },
]

/* ─── Quick-tip items ─── */

const tips = [
  {
    icon: LayoutDashboard,
    label: "Navigate projects from the sidebar",
  },
  {
    icon: Zap,
    label: "Use keyboard shortcuts for speed",
  },
  {
    icon: Sparkles,
    label: "Invite your team to unlock full potential",
  },
]

/* ─── Page Component ─── */

export default function OnboardingPage() {
  const router = useRouter()
  const [hoveredCard, setHoveredCard] = React.useState<string | null>(null)

  return (
    <div className="flex flex-col items-center justify-center min-h-[calc(100vh-120px)] py-8">
      {/* ─ Hero Section ─ */}
      <div className="w-full max-w-3xl text-center mb-12 animate-in fade-in slide-in-from-bottom-4 duration-700">
        {/* Greeting badge */}
        <div className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-emerald-50 border border-emerald-200 text-emerald-700 text-sm font-medium mb-6">
          <Sparkles className="size-4" />
          Welcome aboard!
        </div>

        <h1 className="text-4xl font-semibold tracking-tight text-gray-900 mb-3">
          What would you like to do first?
        </h1>
        <p className="text-lg text-gray-500 max-w-xl mx-auto leading-relaxed">
          Pick an action below to get started. You can always come back to these
          later from the sidebar.
        </p>
      </div>

      {/* ─ Action Cards Grid ─ */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-5 w-full max-w-4xl mb-14 animate-in fade-in slide-in-from-bottom-6 duration-700 delay-150">
        {actions.map((action) => {
          const Icon = action.icon
          const isHovered = hoveredCard === action.id

          return (
            <button
              key={action.id}
              onClick={() => router.push(action.href)}
              onMouseEnter={() => setHoveredCard(action.id)}
              onMouseLeave={() => setHoveredCard(null)}
              className={cn(
                "group relative flex flex-col items-start p-6 rounded-2xl",
                "bg-white border border-gray-100 text-left",
                "transition-all duration-300 ease-out",
                "hover:shadow-[0_20px_40px_-12px_rgba(0,0,0,0.08)]",
                "hover:border-gray-200 hover:-translate-y-1",
                "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:ring-offset-2"
              )}
            >
              {/* Gradient overlay */}
              <div
                className={cn(
                  "absolute inset-0 rounded-2xl bg-gradient-to-br opacity-0 transition-opacity duration-300",
                  action.gradient,
                  isHovered && "opacity-100"
                )}
              />

              {/* Badge */}
              {action.badge && (
                <span className="absolute top-4 right-4 px-2.5 py-0.5 text-[11px] font-semibold tracking-wide uppercase rounded-full bg-emerald-500 text-white">
                  {action.badge}
                </span>
              )}

              {/* Content */}
              <div className="relative z-10 flex flex-col h-full">
                {/* Icon */}
                <div
                  className={cn(
                    "flex items-center justify-center size-12 rounded-xl mb-5 transition-transform duration-300",
                    action.iconColor,
                    isHovered && "scale-110"
                  )}
                >
                  <Icon className="size-6" />
                </div>

                {/* Text */}
                <h3 className="text-lg font-semibold text-gray-900 mb-2">
                  {action.title}
                </h3>
                <p className="text-sm text-gray-500 leading-relaxed flex-1">
                  {action.description}
                </p>

                {/* Action indicator */}
                <div
                  className={cn(
                    "flex items-center gap-1.5 mt-5 text-sm font-medium transition-all duration-300",
                    isHovered
                      ? "text-gray-900 translate-x-0.5"
                      : "text-gray-400"
                  )}
                >
                  Get started
                  <ArrowRight
                    className={cn(
                      "size-4 transition-transform duration-300",
                      isHovered && "translate-x-1"
                    )}
                  />
                </div>
              </div>
            </button>
          )
        })}
      </div>

      {/* ─ Quick Tips ─ */}
      <div className="w-full max-w-2xl animate-in fade-in slide-in-from-bottom-8 duration-700 delay-300">
        <div className="flex items-center gap-2 mb-4">
          <div className="h-px flex-1 bg-gradient-to-r from-transparent via-gray-200 to-transparent" />
          <span className="text-xs font-medium text-gray-400 uppercase tracking-widest">
            Quick tips
          </span>
          <div className="h-px flex-1 bg-gradient-to-r from-transparent via-gray-200 to-transparent" />
        </div>

        <div className="flex flex-wrap justify-center gap-3">
          {tips.map((tip, i) => {
            const TipIcon = tip.icon
            return (
              <div
                key={i}
                className="flex items-center gap-2.5 px-4 py-2.5 rounded-full bg-gray-50 border border-gray-100 text-sm text-gray-600 transition-colors hover:bg-gray-100"
              >
                <TipIcon className="size-4 text-gray-400 shrink-0" />
                {tip.label}
              </div>
            )
          })}
        </div>
      </div>

      {/* ─ Skip link ─ */}
      <button
        onClick={() => router.push("/")}
        className="mt-10 text-sm text-gray-400 hover:text-gray-600 transition-colors underline underline-offset-4 decoration-gray-300 hover:decoration-gray-400 animate-in fade-in duration-1000 delay-500"
      >
        Skip, take me to the dashboard
      </button>
    </div>
  )
}
