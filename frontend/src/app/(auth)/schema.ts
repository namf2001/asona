import { z } from "zod";

export const loginSchema = z.object({
  email: z.string().min(1, { message: "Email is required" }).email({ message: "Invalid email address" }),
  password: z.string().min(1, { message: "Password is required" }),
  remember: z.boolean().default(false).optional(),
});

export const emailSchema = z.object({
  email: z.string().min(1, { message: "Email is required" }).email({ message: "Invalid email address" }),
  terms: z.boolean().refine((val) => val === true, {
    message: "You must agree to the terms",
  }),
});

export const passwordSchema = z.object({
  name: z.string().min(2, { message: "Name must be at least 2 characters" }),
  password: z
    .string()
    .min(6, { message: "Must be at least 6 characters" })
    .regex(/[^a-zA-Z0-9]/, { message: "Must contain 1 special character" }),
  confirmPassword: z.string().min(1, { message: "Please confirm your password" }),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords do not match",
  path: ["confirmPassword"],
});
