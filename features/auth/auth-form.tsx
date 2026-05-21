"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useState, type ReactNode } from "react";
import { useForm, type Resolver } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { submitAuth, type AuthRequest } from "@/lib/auth-api";
import { cn } from "@/lib/utils";
import { useAuthStore } from "@/store/auth-store";

type Message = {
  type: "error" | "success";
  text: string;
} | null;

const loginSchema = z.object({
  email: z.string().email("Введите корректный email"),
  password: z.string().min(1, "Введите пароль"),
  username: z.string().optional(),
});

const registerSchema = loginSchema.extend({
  username: z.string().min(2, "Имя пользователя должно быть длиннее"),
  password: z.string().min(6, "Пароль должен быть не короче 6 символов"),
});

type AuthFormValues = {
  email: string;
  password: string;
  username?: string;
};

export function AuthForm() {
  const router = useRouter();
  const { mode, toggleMode, setSession } = useAuthStore();
  const [message, setMessage] = useState<Message>(null);

  const schema = useMemo(
    () => (mode === "login" ? loginSchema : registerSchema),
    [mode],
  );

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<AuthFormValues>({
    resolver: zodResolver(schema) as Resolver<AuthFormValues>,
    defaultValues: {
      email: "",
      password: "",
      username: "",
    },
  });

  useEffect(() => {
    setMessage(null);
    reset(undefined, { keepValues: true });
  }, [mode, reset]);

  const mutation = useMutation({
    mutationFn: (payload: AuthRequest) => submitAuth(mode, payload),
    onSuccess: (data) => {
      if (mode === "login") {
        setSession(data.token, data.token_type);
        localStorage.setItem("token", data.token);
        localStorage.setItem("token_type", data.token_type);
        localStorage.setItem("expires_in", String(data.expires_in));
        setMessage({
          type: "success",
          text: "Вход выполнен успешно! Перенаправление...",
        });
        window.setTimeout(() => router.push("/dashboard"), 1000);
        return;
      }

      setMessage({
        type: "success",
        text: "Регистрация успешна! Теперь вы можете войти.",
      });
      window.setTimeout(() => toggleMode(), 1500);
    },
    onError: (error) => {
      setMessage({
        type: "error",
        text:
          error instanceof Error
            ? error.message
            : "Что-то пошло не так. Попробуйте снова.",
      });
    },
  });

  const onSubmit = (values: AuthFormValues) => {
    setMessage(null);
    const payload: AuthRequest = {
      email: values.email.trim(),
      password: values.password,
    };

    if (mode === "register") {
      payload.username = values.username?.trim();
    }

    mutation.mutate(payload);
  };

  const isLoginMode = mode === "login";

  return (
    <div className="w-full max-w-[380px] animate-in fade-in slide-in-from-bottom-2 duration-500">
      <h1 className="mb-2 text-center text-[32px] font-semibold tracking-normal text-white">
        {isLoginMode ? "Вход в ReadyUp" : "Создание ReadyUp ID"}
      </h1>
      <p className="mb-8 text-center text-[15px] text-[#86868b]">
        {isLoginMode
          ? "Используй свой аккаунт для поиска сквада"
          : "Один аккаунт для всех игровых дисциплин"}
      </p>

      <form className="space-y-3" onSubmit={handleSubmit(onSubmit)}>
        {!isLoginMode && (
          <FieldError message={errors.username?.message}>
            <Input
              autoComplete="username"
              placeholder="Имя пользователя"
              type="text"
              {...register("username")}
            />
          </FieldError>
        )}

        <FieldError message={errors.email?.message}>
          <Input
            autoComplete="email"
            placeholder="Email"
            type="email"
            {...register("email")}
          />
        </FieldError>

        <FieldError message={errors.password?.message}>
          <Input
            autoComplete={isLoginMode ? "current-password" : "new-password"}
            placeholder="Пароль"
            type="password"
            {...register("password")}
          />
        </FieldError>

        <Button
          className="mt-3 w-full rounded-xl"
          disabled={mutation.isPending}
          size="lg"
          type="submit"
        >
          {mutation.isPending
            ? isLoginMode
              ? "Вход..."
              : "Создание..."
            : isLoginMode
              ? "Продолжить"
              : "Создать аккаунт"}
        </Button>
      </form>

      <div
        className={cn(
          "mt-4 min-h-5 text-center text-[13px]",
          message?.type === "error" && "text-[#ff453a]",
          message?.type === "success" && "text-[#30d158]",
        )}
      >
        {message?.text}
      </div>

      <div className="mt-6 text-center text-sm text-[#86868b]">
        <span>{isLoginMode ? "Еще нет аккаунта?" : "Уже есть аккаунт?"}</span>{" "}
        <button
          className="text-primary hover:underline"
          onClick={toggleMode}
          type="button"
        >
          {isLoginMode ? "Создать сейчас" : "Войти"}
        </button>
      </div>
    </div>
  );
}

function FieldError({
  children,
  message,
}: {
  children: ReactNode;
  message?: string;
}) {
  return (
    <div>
      {children}
      <div className="mt-1 min-h-4 px-1 text-left text-xs text-[#ff453a]">
        {message}
      </div>
    </div>
  );
}
