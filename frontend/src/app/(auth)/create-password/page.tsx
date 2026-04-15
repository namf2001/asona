"use client";

import React, { useState, useTransition, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm, useWatch } from 'react-hook-form';
import { z } from 'zod';
import { toast } from 'sonner';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Form, FormControl, FormField, FormItem, FormMessage } from '@/components/ui/form';
import { User, Lock, EyeOff, Eye, CheckCircle2, Loader2 } from 'lucide-react';

import { registerStep3Action } from '../actions';
import { passwordSchema } from '../schema';

function CreatePasswordContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [isPending, startTransition] = useTransition();
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const email = searchParams.get("email");
  const otp = searchParams.get("otp");

  const form = useForm<z.infer<typeof passwordSchema>>({
    resolver: zodResolver(passwordSchema),
    defaultValues: {
      name: "",
      password: "",
      confirmPassword: "",
    },
    mode: "onChange"
  });

  const passwordValue = useWatch({
    control: form.control,
    name: "password",
    defaultValue: ""
  });
  const hasMinLength = passwordValue.length >= 6;
  const hasSpecialChar = /[^a-zA-Z0-9]/.test(passwordValue);

  function onSubmit(values: z.infer<typeof passwordSchema>) {
    if (!email || !otp) {
      toast.error("Vui lòng quay lại bước nhập email hoặc OTP hợp lệ!");
      return;
    }

    startTransition(async () => {
      const result = await registerStep3Action({
        email,
        otp,
        name: values.name,
        password: values
      });

      if (result?.error) {
        toast.error(result.error);
      } else {
        toast.success("Đăng ký thành công!");
        router.push("/onboard");
      }
    });
  }

  return (
    <div className="flex flex-col w-full h-full font-inter">
      <div className="mb-8">
        <h1 className="text-[32px] font-semibold text-slate-900 tracking-tight">Create Password</h1>
        <p className="text-sm text-slate-500 mt-2">Please set your account name and a strong password.</p>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem className="space-y-2">
                <Label htmlFor="name" className="text-xs font-semibold text-slate-800">Name <span className="text-red-500">*</span></Label>
                <div className="relative">
                  <User className="absolute left-3.5 top-3.5 h-4 w-4 text-slate-400" />
                  <FormControl>
                    <Input id="name" type="text" placeholder="e.g John" className="pl-10 h-12 border-slate-200 rounded-xl bg-white shadow-sm text-sm focus:ring-emerald-500/20 focus:border-emerald-500 transition-all" {...field} />
                  </FormControl>
                </div>
                <FormMessage className="text-[10px]" />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="password"
            render={({ field }) => (
              <FormItem className="space-y-2">
                <Label htmlFor="create-password" className="text-xs font-semibold text-slate-800">Create Password <span className="text-red-500">*</span></Label>
                <div className="relative">
                  <Lock className="absolute left-3.5 top-3.5 h-4 w-4 text-slate-400" />
                  <FormControl>
                    <Input id="create-password" type={showPassword ? "text" : "password"} placeholder="Create a password" className="pl-10 pr-10 h-12 border-slate-200 rounded-xl bg-white shadow-sm text-sm focus:ring-emerald-500/20 focus:border-emerald-500 transition-all" {...field} />
                  </FormControl>
                  <button 
                    type="button" 
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3.5 top-3.5 text-slate-300 hover:text-slate-500 transition-colors"
                  >
                    {showPassword ? <Eye className="h-4 w-4" /> : <EyeOff className="h-4 w-4" />}
                  </button>
                </div>
                <FormMessage className="text-[10px]" />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="confirmPassword"
            render={({ field }) => (
              <FormItem className="space-y-2">
                <Label htmlFor="confirm-password" className="text-xs font-semibold text-slate-800">Confirm Password <span className="text-red-500">*</span></Label>
                <div className="relative">
                  <Lock className="absolute left-3.5 top-3.5 h-4 w-4 text-slate-400" />
                  <FormControl>
                    <Input id="confirm-password" type={showConfirmPassword ? "text" : "password"} placeholder="Confirm your password" className="pl-10 pr-10 h-12 border-slate-200 rounded-xl bg-white shadow-sm text-sm focus:ring-emerald-500/20 focus:border-emerald-500 transition-all" {...field} />
                  </FormControl>
                  <button 
                    type="button" 
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    className="absolute right-3.5 top-3.5 text-slate-300 hover:text-slate-500 transition-colors"
                  >
                    {showConfirmPassword ? <Eye className="h-4 w-4" /> : <EyeOff className="h-4 w-4" />}
                  </button>
                </div>
                <FormMessage className="text-[10px]" />
              </FormItem>
            )}
          />

          <div className="space-y-2 pt-2">
            <div className="flex items-center gap-2">
              <CheckCircle2 className={`w-3.5 h-3.5 transition-colors ${hasMinLength ? 'text-emerald-500 fill-emerald-100' : 'text-slate-400 fill-slate-200'}`} />
              <span className={`text-xs font-medium transition-colors ${hasMinLength ? 'text-emerald-600' : 'text-slate-500'}`}>At least 6 characters</span>
            </div>
            <div className="flex items-center gap-2">
              <CheckCircle2 className={`w-3.5 h-3.5 transition-colors ${hasSpecialChar ? 'text-emerald-500 fill-emerald-100' : 'text-slate-400 fill-slate-200'}`} />
              <span className={`text-xs font-medium transition-colors ${hasSpecialChar ? 'text-emerald-600' : 'text-slate-500'}`}>At least 1 special character</span>
            </div>
          </div>

          <Button disabled={isPending} className="w-full h-12 bg-[#9CA3AF] hover:bg-[#868e96] text-white rounded-xl font-semibold mt-8 shadow-md transition-all active:scale-[0.98]" type="submit">
            {isPending ? <Loader2 className="w-4 h-4 mr-2 animate-spin" /> : null}
            Continue
          </Button>
        </form>
      </Form>
    </div>
  );
}

export default function CreatePasswordPage() {
  return (
    <Suspense fallback={<div className="flex w-full justify-center"><Loader2 className="animate-spin text-slate-400" /></div>}>
      <CreatePasswordContent />
    </Suspense>
  );
}
