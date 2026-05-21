import Link from "next/link";
import { Button } from "@/components/ui/button";

export const metadata = {
  title: "ReadyUp — Dashboard",
};

export default function DashboardPage() {
  return (
    <main className="flex min-h-screen items-center justify-center bg-black px-5 text-white">
      <section className="w-full max-w-xl text-center">
        <h1 className="text-4xl font-semibold tracking-normal">
          Добро пожаловать в ReadyUp
        </h1>
        <p className="mt-4 text-[#86868b]">
          Авторизация прошла успешно. Следующий шаг — собрать здесь личный
          кабинет игрока.
        </p>
        <Button asChild className="mt-8">
          <Link href="/">На главную</Link>
        </Button>
      </section>
    </main>
  );
}
