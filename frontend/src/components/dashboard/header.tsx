"use client"

import * as React from "react"
import { Search, Bell, ChevronDown } from "lucide-react"

import { Input } from "@/components/ui/input"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { cn } from "@/lib/utils"
import { SidebarTrigger } from "@/components/ui/sidebar"

const tabs = ["Overview", "Project", "Team", "Members"]

export function DashboardHeader() {
  const [activeTab, setActiveTab] = React.useState("Project")

  return (
    <header className="flex items-center justify-between px-8 py-4 bg-transparent">
      <SidebarTrigger className="-ml-1" />
      {/* Left: Tab Navigation */}
      <div className="flex bg-white/50 backdrop-blur-sm p-1 rounded-full border border-gray-200">
        {tabs.map((tab) => {
          const isActive = activeTab === tab
          return (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={cn(
                "px-6 py-2 rounded-full text-sm font-medium transition-all duration-300",
                isActive 
                  ? "bg-[#111827] text-white shadow-md scale-100" 
                  : "text-gray-500 hover:text-gray-900"
              )}
            >
              {tab}
            </button>
          )
        })}
      </div>

      {/* Right: Actions */}
      <div className="flex items-center gap-4">
        {/* Notifications */}
        <button className="relative p-2 rounded-full bg-white border border-gray-200 hover:bg-gray-50 transition-colors">
          <Bell className="size-5 text-gray-600" />
          <span className="absolute top-2 right-2 size-2 bg-red-500 rounded-full border-2 border-white" />
        </button>

        {/* Search */}
        <div className="relative w-64 group">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-gray-400 group-focus-within:text-blue-500 transition-colors" />
          <Input 
            placeholder="Search for..." 
            className="pl-10 pr-4 h-10 rounded-full bg-white border-gray-200 focus-visible:ring-blue-500 transition-all"
          />
        </div>

        {/* Org Switcher */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button className="flex items-center gap-3 px-4 py-2 bg-white border border-gray-200 rounded-full hover:bg-gray-50 transition-all font-medium text-sm">
              <div className="size-6 bg-orange-100 rounded flex items-center justify-center text-orange-600 font-bold text-xs">
                B
              </div>
              <span>Betasoft</span>
              <ChevronDown className="size-4 text-gray-400" />
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-56">
            <DropdownMenuItem>Betasoft Org</DropdownMenuItem>
            <DropdownMenuItem>Personal Space</DropdownMenuItem>
            <DropdownMenuItem className="border-t">Add Organization</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  )
}
