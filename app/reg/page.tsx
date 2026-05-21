import Link from "next/link";
import { AuthForm } from "@/features/auth/auth-form";

export const metadata = {
  title: "ReadyUp — Вход в аккаунт",
};

export default function AuthPage() {
  return (
    <main className="flex min-h-screen flex-col bg-[radial-gradient(circle_at_50%_50%,#1d1d1f_0%,#000_80%)] text-[#f5f5f7]">
      <header className="absolute top-0 flex h-12 w-full items-center px-5">
        <Link
          href="/"
          className="text-[19px] font-bold tracking-normal text-white"
        >
          ReadyUp
        </Link>
      </header>

      <section className="flex flex-1 items-center justify-center px-5 py-20">
        <AuthForm />
      </section>

      <footer className="border-t border-white/5 p-6 text-center text-xs text-[#86868b]">
        Copyright &copy; 2026 ReadyUp Inc. Все права защищены.
      </footer>
    </main>
  );
}
