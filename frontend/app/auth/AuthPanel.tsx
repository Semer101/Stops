"use client";

import { FormEvent, useMemo, useState } from "react";
import Link from "next/link";

import { login, register, saveAuthSession } from "@/services/authService";
import type { AuthResponse, LoginPayload, RegisterPayload } from "@/types/auth";

type AuthMode = "login" | "register";

export function AuthPanel() {
  const [mode, setMode] = useState<AuthMode>("login");
  const [name, setName] = useState("");
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const title = mode === "login" ? "Welcome back" : "Create account";
  const submitLabel = mode === "login" ? "Log in" : "Register";
  const alternateLabel = mode === "login" ? "Create account" : "Use existing account";

  const contactPayload = useMemo(() => {
    const value = identifier.trim();
    if (value.includes("@")) {
      return { email: value };
    }

    return { phone: value };
  }, [identifier]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setMessage("");
    setIsSubmitting(true);

    try {
      const response = mode === "login" ? await submitLogin() : await submitRegister();
      saveAuthSession(response);
      setMessage(`${response.data.user.name} is signed in.`);
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Authentication failed.");
    } finally {
      setIsSubmitting(false);
    }
  }

  function submitLogin(): Promise<AuthResponse> {
    const payload: LoginPayload = {
      ...contactPayload,
      password,
    };

    return login(payload);
  }

  function submitRegister(): Promise<AuthResponse> {
    const payload: RegisterPayload = {
      name,
      ...contactPayload,
      password,
    };

    return register(payload);
  }

  return (
    <section className="rounded-md border border-road/15 bg-white p-5 shadow-sm">
      <div className="grid grid-cols-2 rounded-md bg-mist p-1">
        {(["login", "register"] as AuthMode[]).map((item) => (
          <button
            className={`rounded px-3 py-2 text-sm font-semibold transition ${
              mode === item ? "bg-white text-road shadow-sm" : "text-ink/65 hover:text-road"
            }`}
            key={item}
            onClick={() => {
              setMode(item);
              setMessage("");
            }}
            type="button"
          >
            {item === "login" ? "Login" : "Register"}
          </button>
        ))}
      </div>

      <h2 className="mt-5 text-2xl font-bold">{title}</h2>

      <form className="mt-5 grid gap-4" onSubmit={handleSubmit}>
        {mode === "register" ? (
          <label className="grid gap-2 text-sm font-semibold text-ink" htmlFor="name">
            Full name
            <input
              className="rounded-md border border-ink/15 px-3 py-3 text-base font-normal outline-none transition focus:border-road focus:ring-2 focus:ring-road/20"
              id="name"
              minLength={2}
              onChange={(event) => setName(event.target.value)}
              required
              type="text"
              value={name}
            />
          </label>
        ) : null}

        <label className="grid gap-2 text-sm font-semibold text-ink" htmlFor="identifier">
          Email or phone
          <input
            className="rounded-md border border-ink/15 px-3 py-3 text-base font-normal outline-none transition focus:border-road focus:ring-2 focus:ring-road/20"
            id="identifier"
            onChange={(event) => setIdentifier(event.target.value)}
            required
            type="text"
            value={identifier}
          />
        </label>

        <label className="grid gap-2 text-sm font-semibold text-ink" htmlFor="password">
          Password
          <input
            className="rounded-md border border-ink/15 px-3 py-3 text-base font-normal outline-none transition focus:border-road focus:ring-2 focus:ring-road/20"
            id="password"
            minLength={mode === "register" ? 8 : undefined}
            onChange={(event) => setPassword(event.target.value)}
            required
            type="password"
            value={password}
          />
        </label>

        <button
          className="rounded-md bg-road px-4 py-3 text-sm font-bold text-white shadow-sm transition hover:bg-ink disabled:cursor-not-allowed disabled:bg-ink/35"
          disabled={isSubmitting}
          type="submit"
        >
          {isSubmitting ? "Working..." : submitLabel}
        </button>
      </form>

      <div className="mt-4 flex flex-col gap-2">
        <Link href="/auth/reset" className="text-xs font-semibold text-road hover:underline self-start">
          Forgot password?
        </Link>
        <button
          className="text-sm font-semibold text-road hover:text-ink self-start text-left"
          onClick={() => {
            setMode(mode === "login" ? "register" : "login");
            setMessage("");
          }}
          type="button"
        >
          {alternateLabel}
        </button>
      </div>

      {message ? (
        <p className="mt-4 rounded-md border border-signal/40 bg-signal/10 px-3 py-2 text-sm font-medium text-ink">
          {message}
        </p>
      ) : null}
    </section>
  );
}
