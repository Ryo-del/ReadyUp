import Link from "next/link";
import { Button } from "@/components/ui/button";

const games = [
  {
    title: "CS 2",
    description:
      "Максимальная концентрация, тактический раскид и выверенные тайминги.",
    online: "14,200",
  },
  {
    title: "DOTA 2",
    description:
      "Сложная макро-игра, драфты, координация линий и битва за Рошана.",
    online: "9,800",
  },
  {
    title: "Fortnite",
    description:
      "Динамичная королевская битва, безумное строительство и яркие ивенты.",
    online: "4,500",
  },
  {
    title: "PUBG",
    description:
      "Реалистичный баттл-рояль. Лутинг, тактическое позиционирование и борьба за топ-1.",
    online: "2,100",
  },
];

export default function HomePage() {
  return (
    <main className="min-h-screen bg-black text-[#f5f5f7]">
      <header className="fixed top-0 z-50 h-12 w-full border-b border-white/5 bg-black/80 backdrop-blur-xl">
        <nav className="mx-auto flex h-full max-w-5xl items-center justify-between px-5">
          <Link
            href="/"
            className="text-[19px] font-bold tracking-normal text-white"
          >
            ReadyUp
          </Link>

          <div className="hidden items-center gap-6 text-xs text-[#86868b] sm:flex">
            <Link className="transition-colors hover:text-white" href="#player">
              Поиск Тимейта
            </Link>
            <Link className="transition-colors hover:text-white" href="#teams">
              Поиск Сквадов
            </Link>
            <Link
              className="transition-colors hover:text-white"
              href="#tournaments"
            >
              Турниры
            </Link>
          </div>

          <Button asChild variant="secondary" size="sm">
            <Link href="/reg">Войти</Link>
          </Button>
        </nav>
      </header>

      <section className="flex min-h-screen flex-col items-center justify-center bg-[radial-gradient(circle_at_50%_40%,rgba(0,255,102,0.07)_0%,rgba(0,113,227,0.05)_30%,#000_70%)] px-5 pb-16 pt-28 text-center">
        <div className="mb-4 text-xs font-semibold uppercase tracking-[0.1em] text-[#00ff66]">
          Твой новый уровень кооперации
        </div>
        <h1 className="max-w-4xl bg-gradient-to-b from-white from-30% to-[#86868b] bg-clip-text text-5xl font-bold leading-[1.05] tracking-normal text-transparent md:text-7xl">
          Твоя идеальная команда уже ждет тебя.
        </h1>
        <p className="mt-5 max-w-xl text-xl leading-relaxed text-[#86868b]">
          Забудь про одиночные матчи с незнакомцами. Находи единомышленников
          для ярких побед и комфортной игры в удобное время.
        </p>
        <div className="mt-9">
          <Button asChild size="lg">
            <Link href="#games">Найти команду</Link>
          </Button>
        </div>
      </section>

      <section id="games" className="mx-auto max-w-5xl px-5 py-16 pb-28">
        <h2 className="mb-10 text-center text-3xl font-bold tracking-normal">
          Популярные игры
        </h2>
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
          {games.map((game) => (
            <article
              key={game.title}
              className="flex h-[220px] cursor-pointer flex-col justify-between overflow-hidden rounded-[20px] border border-white/[0.03] bg-[#161617]/70 p-6 transition duration-300 hover:scale-[1.02] hover:border-[#00ff66]/30 hover:bg-[#161617]/90"
            >
              <div>
                <h3 className="mb-1.5 text-[22px] font-semibold tracking-normal">
                  {game.title}
                </h3>
                <p className="text-[13px] leading-relaxed text-[#86868b]">
                  {game.description}
                </p>
              </div>
              <div className="flex items-center gap-1.5 text-xs font-medium text-[#00ff66]">
                <span className="size-1.5 rounded-full bg-[#00ff66] shadow-[0_0_8px_#00ff66]" />
                {game.online} в поиске
              </div>
            </article>
          ))}
        </div>
      </section>

      <footer className="border-t border-white/5 bg-black px-5 py-10 text-xs text-[#86868b]">
        <div className="mx-auto flex max-w-5xl flex-col items-center gap-4 md:flex-row md:justify-between">
          <div>Copyright &copy; 2026 ReadyUp Gaming Inc.</div>
          <div className="flex gap-4">
            <Link className="transition-colors hover:text-white" href="#">
              Политика конфиденциальности
            </Link>
            <Link className="transition-colors hover:text-white" href="#">
              Правила сообщества
            </Link>
          </div>
        </div>
      </footer>
    </main>
  );
}
