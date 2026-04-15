import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/dashboard/sidebar/app-sidebar"
import { DashboardHeader } from "@/components/dashboard/header"

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <SidebarProvider defaultOpen={false}>
        <AppSidebar />
        <SidebarInset className="flex flex-col bg-transparent">
          <DashboardHeader />
          <main className="flex-1 p-8 pt-4 overflow-auto">
            {children}
          </main>
        </SidebarInset>
    </SidebarProvider>
  );
}