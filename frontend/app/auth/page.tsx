import Link from "next/link";

import { AuthPanel } from "./AuthPanel";

export default function AuthPage() {
  return (
    <main className="min-h-screen bg-mist text-ink">
      <section className="mx-auto flex min-h-screen w-full max-w-6xl flex-col px-4 py-5 sm:px-6 lg:px-8">
        <header className="flex items-center justify-between border-b border-ink/10 pb-5">
          <Link className="text-2xl font-bold tracking-normal" href="/">
            STOPS
          </Link>
          <Link
            className="rounded-md border border-road/20 bg-white px-3 py-2 text-sm font-medium text-road shadow-sm transition hover:border-road hover:bg-road hover:text-white"
            href="/"
          >
            Map
          </Link>
        </header>

        <div className="grid flex-1 items-center gap-8 py-8 lg:grid-cols-[minmax(0,1fr)_420px]">
          <section className="max-w-2xl">
            <p className="text-sm font-semibold uppercase tracking-wide text-road">Secure access</p>
            <h1 className="mt-3 text-4xl font-bold leading-tight sm:text-5xl">
              Plan trips, save places, and contribute verified transit data.
            </h1>
            <div className="mt-8 grid gap-3 sm:grid-cols-3">
              {["JWT sessions", "bcrypt passwords", "Role-based access"].map((item) => (
                <div key={item} className="rounded-md border border-road/15 bg-white p-4 shadow-sm">
                  <p className="text-sm font-semibold text-road">{item}</p>
                </div>
              ))}
            </div>
          </section>

          <AuthPanel />
        </div>
      </section>
    </main>
  );
}
