"use client";

import React, { useTransition } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { toast } from 'sonner';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Checkbox } from '@/components/ui/checkbox';
import { Form, FormControl, FormField, FormItem, FormMessage } from '@/components/ui/form';
import { Mail, Loader2 } from 'lucide-react';

import { emailSchema } from '../schema';
import { registerStep1Action } from '../actions';

export default function RegisterPage() {
  const router = useRouter();
  const [isPending, startTransition] = useTransition();

  const form = useForm<z.infer<typeof emailSchema>>({
    resolver: zodResolver(emailSchema),
    defaultValues: {
      email: "",
      terms: false,
    },
  });

  function onSubmit(values: z.infer<typeof emailSchema>) {
    startTransition(async () => {
      const result = await registerStep1Action(values.email);
      if (result.error) {
        toast.error(result.error);
      } else {
        toast.success("Verification code sent!");
        router.push(`/verify?email=${encodeURIComponent(values.email)}`);
      }
    });
  }

  return (
    <div className="flex flex-col w-full h-full font-inter">
      <div className="mb-8">
        <h1 className="text-[32px] font-semibold text-slate-900 tracking-tight">Sign Up Free</h1>
        <p className="text-sm text-slate-500 mt-2 leading-relaxed">
          Join us today to unlock exclusive features that will boost your team&apos;s productivity!
        </p>
      </div>

      <div className="mb-6">
        <Button 
          variant="outline" 
          className="w-full h-12 bg-white border-slate-200 text-slate-700 hover:bg-slate-50 rounded-xl font-medium shadow-sm transition-all" 
          type="button"
          onClick={async () => {
            try {
              const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000'}/api/v1/auth/google`, {
                credentials: 'include'
              });
              const body = await response.json();
              if (body.data?.url) {
                window.location.href = body.data.url;
              } else {
                toast.error("Failed to get Google login URL");
              }
            } catch (err) {
              console.error("Google register error:", err);
              toast.error("Network error during Google login");
            }
          }}
        >
          <svg className="w-5 h-5 mr-3" viewBox="0 0 24 24">
            <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4" />
            <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" />
            <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05" />
            <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" />
          </svg>
          Continue with Google
        </Button>
      </div>

      <div className="relative flex py-2 items-center mb-6">
        <div className="flex-grow border-t border-slate-100"></div>
        <span className="flex-shrink-0 mx-4 text-[10px] text-slate-400 font-medium uppercase tracking-wider">Or continue with</span>
        <div className="flex-grow border-t border-slate-100"></div>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
          <FormField
            control={form.control}
            name="email"
            render={({ field }) => (
              <FormItem className="space-y-2">
                <Label htmlFor="email" className="text-xs font-semibold text-slate-800">Email <span className="text-red-500">*</span></Label>
                <div className="relative">
                  <Mail className="absolute left-3.5 top-3.5 h-4 w-4 text-slate-400" />
                  <FormControl>
                    <Input id="email" type="email" placeholder="e.g somebody@example.com" className="pl-10 h-12 border-slate-200 rounded-xl bg-white shadow-sm text-sm focus:ring-emerald-500/20 focus:border-emerald-500 transition-all" {...field} />
                  </FormControl>
                </div>
                <FormMessage className="text-[10px]" />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="terms"
            render={({ field }) => (
              <FormItem className="flex flex-col space-y-1 pt-2">
                <div className="flex items-start space-x-3">
                  <FormControl>
                    <Checkbox id="terms" className="rounded-md border-slate-300 w-4 h-4 mt-0.5 data-[state=checked]:bg-emerald-500 data-[state=checked]:border-emerald-500" checked={field.value} onCheckedChange={field.onChange} />
                  </FormControl>
                  <Label htmlFor="terms" className="text-xs text-slate-500 font-normal leading-relaxed cursor-pointer">
                    Agree to the <Link href="/terms" className="text-emerald-500 hover:underline">terms of use</Link> and our <Link href="/privacy" className="text-emerald-500 hover:underline">privacy policy</Link>, which describes how we collect, use, store, and share your data.
                  </Label>
                </div>
                <FormMessage className="text-[10px]" />
              </FormItem>
            )}
          />

          <Button disabled={isPending} className="w-full h-12 bg-[#9CA3AF] hover:bg-[#868e96] text-white rounded-xl font-semibold mt-6 shadow-md transition-all active:scale-[0.98]" type="submit">
            {isPending ? <Loader2 className="w-4 h-4 mr-2 animate-spin" /> : null}
            Sign Up
          </Button>
        </form>
      </Form>

      <div className="mt-8 text-center text-xs text-slate-500">
        Already have an account? <Link href="/login" className="text-emerald-500 hover:text-emerald-600 font-semibold ml-1 transition-colors">Log in</Link>
      </div>
    </div>
  );
}
