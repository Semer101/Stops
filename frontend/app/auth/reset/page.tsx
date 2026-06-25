"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";

export default function ResetPasswordPage() {
  const router = useRouter();
  const [emailOrPhone, setEmailOrPhone] = useState("");
  const [step, setStep] = useState(1); // 1 = Request, 2 = Reset
  const [token, setToken] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [simulatedToken, setSimulatedToken] = useState("");

  const handleRequestToken = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!emailOrPhone.trim()) {
      setError("Email or Phone number is required.");
      return;
    }

    setLoading(true);
    setError("");
    setSuccess("");

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}/api/auth/forgot-password`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email_or_phone: emailOrPhone }),
      });

      const resData = await response.json();
      if (response.ok && resData.success) {
        setSuccess("Reset token generated successfully!");
        setSimulatedToken(resData.token); // Store simulated token returned directly by API
        setToken(resData.token); // Pre-fill token for easy testing
        setStep(2);
      } else {
        setError(resData.message || "Failed to generate reset token.");
      }
    } catch {
      setError("Failed to connect to the server.");
    } finally {
      setLoading(false);
    }
  };

  const handleResetPassword = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token.trim()) {
      setError("Reset token is required.");
      return;
    }
    if (newPassword.length < 8) {
      setError("Password must be at least 8 characters long.");
      return;
    }
    if (newPassword !== confirmPassword) {
      setError("Passwords do not match.");
      return;
    }

    setLoading(true);
    setError("");
    setSuccess("");

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}/api/auth/reset-password`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          token: token,
          new_password: newPassword,
        }),
      });

      const resData = await response.json();
      if (response.ok && resData.success) {
        setSuccess("Password has been reset successfully! Redirecting to login...");
        setError("");
        setTimeout(() => {
          router.push("/auth");
        }, 2500);
      } else {
        setError(resData.message || "Failed to reset password.");
      }
    } catch {
      setError("Failed to connect to the server.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className="min-h-screen bg-slate-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-3xl shadow-xl border border-emerald-800/10 p-8 max-w-md w-full">
        <div className="text-center mb-6">
          <h2 className="text-2xl font-black text-emerald-800 tracking-tight">STOPS Addis</h2>
          <p className="text-sm text-slate-500 font-medium mt-1">Reset your password account</p>
        </div>

        {error && (
          <div className="bg-red-50 text-red-700 border border-red-200 rounded-xl p-3 text-xs font-semibold mb-4">
            ⚠️ {error}
          </div>
        )}

        {success && (
          <div className="bg-emerald-50 text-emerald-700 border border-emerald-200 rounded-xl p-3 text-xs font-semibold mb-4">
            ✅ {success}
          </div>
        )}

        {step === 1 ? (
          <form onSubmit={handleRequestToken} className="space-y-4">
            <div>
              <label className="text-xs font-bold text-slate-700" htmlFor="email-phone">Email or Phone Number</label>
              <input
                id="email-phone"
                type="text"
                required
                placeholder="e.g. user@example.com or +251..."
                value={emailOrPhone}
                onChange={(e) => setEmailOrPhone(e.target.value)}
                className="w-full text-xs rounded-xl border border-slate-300 px-3.5 py-3 outline-none mt-1 focus:border-emerald-700 focus:ring-1 focus:ring-emerald-700"
              />
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full bg-emerald-800 hover:bg-emerald-950 text-white font-bold text-xs py-3 rounded-xl shadow transition disabled:opacity-50"
            >
              {loading ? "Generating Token..." : "Generate Reset Token"}
            </button>
          </form>
        ) : (
          <form onSubmit={handleResetPassword} className="space-y-4">
            {simulatedToken && (
              <div className="bg-amber-500/10 border border-amber-500/20 text-amber-800 rounded-xl p-3.5 text-xs font-medium leading-relaxed">
                📢 <strong>Simulation Note:</strong> In production, a reset link is sent. For local testing, your reset token has been captured and filled below:<br />
                <code className="bg-white/80 border border-amber-500/30 px-1.5 py-0.5 rounded font-mono text-[10px] block mt-1 break-all select-all">
                  {simulatedToken}
                </code>
              </div>
            )}

            <div>
              <label className="text-xs font-bold text-slate-700" htmlFor="reset-token">Reset Token</label>
              <input
                id="reset-token"
                type="text"
                required
                placeholder="Enter token"
                value={token}
                onChange={(e) => setToken(e.target.value)}
                className="w-full text-xs rounded-xl border border-slate-300 px-3.5 py-3 outline-none mt-1 focus:border-emerald-700 font-mono"
              />
            </div>

            <div>
              <label className="text-xs font-bold text-slate-700" htmlFor="new-password">New Password</label>
              <input
                id="new-password"
                type="password"
                required
                placeholder="Minimum 8 characters"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                className="w-full text-xs rounded-xl border border-slate-300 px-3.5 py-3 outline-none mt-1 focus:border-emerald-700"
              />
            </div>

            <div>
              <label className="text-xs font-bold text-slate-700" htmlFor="confirm-password">Confirm Password</label>
              <input
                id="confirm-password"
                type="password"
                required
                placeholder="Repeat password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                className="w-full text-xs rounded-xl border border-slate-300 px-3.5 py-3 outline-none mt-1 focus:border-emerald-700"
              />
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full bg-emerald-800 hover:bg-emerald-950 text-white font-bold text-xs py-3 rounded-xl shadow transition disabled:opacity-50"
            >
              {loading ? "Resetting Password..." : "Change Password"}
            </button>
          </form>
        )}

        <div className="text-center mt-6 border-t border-slate-100 pt-4 flex justify-between text-[11px] text-slate-500 font-semibold">
          <Link href="/auth" className="hover:underline">
            Back to Sign In
          </Link>
          <Link href="/" className="hover:underline">
            Back to Map
          </Link>
        </div>
      </div>
    </main>
  );
}
