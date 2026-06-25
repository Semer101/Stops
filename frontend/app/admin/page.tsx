"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import type { PublicUser } from "@/types/auth";

interface PendingReport {
  reportId: string;
  stopId: string;
  stopName: string;
  stopType: string;
  latitude: number;
  longitude: number;
  area: string | null;
  netVotes: number;
  createdBy: string;
  creatorName: string;
  createdAt: string;
}

interface AuditLog {
  id: string;
  adminUserId: string;
  adminName: string;
  action: string;
  entityType: string;
  entityId: string | null;
  metadata: string;
  createdAt: string;
}

export default function AdminDashboardPage() {
  const [token, setToken] = useState<string | null>(null);
  const [user, setUser] = useState<PublicUser | null>(null);
  const [authorized, setAuthorized] = useState<boolean | null>(null);

  const [activeTab, setActiveTab] = useState<"moderation" | "logs" | "analytics">("moderation");
  const [pendingReports, setPendingReports] = useState<PendingReport[]>([]);
  const [auditLogs, setAuditLogs] = useState<AuditLog[]>([]);
  
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  // Check Authorization
  useEffect(() => {
    const savedToken = localStorage.getItem("stops.auth.token");
    const savedUser = localStorage.getItem("stops.auth.user");

    if (!savedToken || !savedUser) {
      setAuthorized(false);
      return;
    }

    const parsedUser = JSON.parse(savedUser) as PublicUser;
    setToken(savedToken);
    setUser(parsedUser);

    if (parsedUser.role !== "admin") {
      setAuthorized(false);
    } else {
      setAuthorized(true);
    }
  }, []);

  // Fetch Data based on active tab
  const fetchModerationQueue = async (authToken: string) => {
    setLoading(true);
    setError("");
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}/api/admin/reports`, {
        headers: {
          "Authorization": `Bearer ${authToken}`,
        },
      });
      const resData = await response.json();
      if (response.ok && resData.success) {
        setPendingReports(resData.data || []);
      } else {
        setError(resData.message || "Failed to load moderation queue.");
      }
    } catch {
      setError("Failed to connect to the server.");
    } finally {
      setLoading(false);
    }
  };

  const fetchAuditLogs = async (authToken: string) => {
    setLoading(true);
    setError("");
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}/api/admin/audit-logs`, {
        headers: {
          "Authorization": `Bearer ${authToken}`,
        },
      });
      const resData = await response.json();
      if (response.ok && resData.success) {
        setAuditLogs(resData.data || []);
      } else {
        setError(resData.message || "Failed to load audit logs.");
      }
    } catch {
      setError("Failed to connect to the server.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (authorized && token) {
      if (activeTab === "moderation") {
        fetchModerationQueue(token);
      } else if (activeTab === "logs") {
        fetchAuditLogs(token);
      }
    }
  }, [authorized, activeTab, token]);

  const handleVerify = async (reportId: string) => {
    if (!token) return;
    setSuccess("");
    setError("");
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}/api/admin/reports/${reportId}/verify`, {
        method: "POST",
        headers: {
          "Authorization": `Bearer ${token}`,
        },
      });
      const resData = await response.json();
      if (response.ok && resData.success) {
        setSuccess("Report verified and stop activated successfully.");
        fetchModerationQueue(token);
      } else {
        setError(resData.message || "Failed to verify stop.");
      }
    } catch {
      setError("Failed to verify stop.");
    }
  };

  const handleReject = async (reportId: string) => {
    if (!token) return;
    setSuccess("");
    setError("");
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}/api/admin/reports/${reportId}/reject`, {
        method: "POST",
        headers: {
          "Authorization": `Bearer ${token}`,
        },
      });
      const resData = await response.json();
      if (response.ok && resData.success) {
        setSuccess("Report rejected and stop deleted successfully.");
        fetchModerationQueue(token);
      } else {
        setError(resData.message || "Failed to reject stop.");
      }
    } catch {
      setError("Failed to reject stop.");
    }
  };

  if (authorized === false) {
    return (
      <main className="min-h-screen bg-slate-50 flex items-center justify-center p-4">
        <div className="bg-white rounded-3xl shadow-xl border border-red-200 p-8 max-w-md w-full text-center">
          <span className="text-4xl">🚫</span>
          <h2 className="text-xl font-black text-red-700 mt-3">Access Denied</h2>
          <p className="text-xs text-slate-500 font-semibold mt-2 leading-relaxed">
            You must be logged in as an administrator to access the STOPS Moderation Dashboard.
          </p>
          <div className="mt-6 flex flex-col gap-2">
            <Link
              href="/auth"
              className="w-full bg-emerald-800 hover:bg-emerald-950 text-white font-bold text-xs py-3 rounded-xl transition shadow"
            >
              Sign In with Admin Account
            </Link>
            <Link
              href="/"
              className="w-full bg-slate-100 hover:bg-slate-200 text-slate-700 font-bold text-xs py-3 rounded-xl transition border border-slate-200"
            >
              Back to Map
            </Link>
          </div>
        </div>
      </main>
    );
  }

  if (authorized === null) {
    return (
      <div className="min-h-screen bg-slate-50 flex items-center justify-center">
        <div className="h-8 w-8 border-4 border-emerald-800 border-t-transparent rounded-full animate-spin" />
      </div>
    );
  }

  return (
    <main className="min-h-screen bg-slate-50 text-slate-900 flex flex-col font-sans">
      {/* Top Header */}
      <header className="bg-emerald-900 text-white p-5 shadow-md flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="h-10 w-10 bg-amber-500 rounded-xl flex items-center justify-center text-emerald-950 font-black text-xl shadow">
            A
          </div>
          <div>
            <h1 className="text-xl font-black tracking-tight">STOPS Moderation</h1>
            <p className="text-xs text-emerald-300 font-semibold">Admin Panel • Hello, {user?.name}</p>
          </div>
        </div>
        <Link
          href="/"
          className="rounded-xl border border-white/20 bg-white/10 px-4 py-2 text-xs font-bold text-white shadow hover:bg-white hover:text-emerald-950 transition"
        >
          🗺️ Return to Map
        </Link>
      </header>

      {/* Main Contents Grid */}
      <div className="max-w-7xl w-full mx-auto p-4 md:p-6 flex-1 flex flex-col md:flex-row gap-6">
        {/* Navigation Tabs */}
        <aside className="w-full md:w-64 flex flex-row md:flex-col gap-2">
          <button
            onClick={() => setActiveTab("moderation")}
            className={`flex-1 md:flex-initial text-left px-4 py-3 rounded-xl font-bold text-xs transition border ${
              activeTab === "moderation"
                ? "bg-emerald-800 text-white border-emerald-900 shadow"
                : "bg-white text-slate-700 border-slate-200 hover:border-slate-300"
            }`}
          >
            🚏 Moderation Queue ({pendingReports.length})
          </button>
          <button
            onClick={() => setActiveTab("logs")}
            className={`flex-1 md:flex-initial text-left px-4 py-3 rounded-xl font-bold text-xs transition border ${
              activeTab === "logs"
                ? "bg-emerald-800 text-white border-emerald-900 shadow"
                : "bg-white text-slate-700 border-slate-200 hover:border-slate-300"
            }`}
          >
            📋 Admin Audit Logs
          </button>
          <button
            onClick={() => setActiveTab("analytics")}
            className={`flex-1 md:flex-initial text-left px-4 py-3 rounded-xl font-bold text-xs transition border ${
              activeTab === "analytics"
                ? "bg-emerald-800 text-white border-emerald-900 shadow"
                : "bg-white text-slate-700 border-slate-200 hover:border-slate-300"
            }`}
          >
            📈 System Analytics
          </button>
        </aside>

        {/* Dashboard Panels */}
        <section className="flex-1 bg-white rounded-3xl shadow-sm border border-slate-200 p-5 md:p-6">
          {success && (
            <div className="bg-emerald-50 text-emerald-700 border border-emerald-200 rounded-xl p-3.5 text-xs font-semibold mb-4 animate-fadeIn">
              ✅ {success}
            </div>
          )}

          {error && (
            <div className="bg-red-50 text-red-700 border border-red-200 rounded-xl p-3.5 text-xs font-semibold mb-4 animate-fadeIn">
              ⚠️ {error}
            </div>
          )}

          {loading ? (
            <div className="flex items-center justify-center py-12">
              <div className="h-8 w-8 border-4 border-emerald-800 border-t-transparent rounded-full animate-spin" />
            </div>
          ) : activeTab === "moderation" ? (
            <div>
              <h2 className="text-lg font-black text-emerald-800 mb-4 uppercase tracking-wider">🚏 Crowdsourced Moderation Queue</h2>
              {pendingReports.length === 0 ? (
                <div className="text-center py-12 border-2 border-dashed border-slate-200 rounded-2xl">
                  <span className="text-3xl">🎉</span>
                  <p className="text-xs text-slate-500 font-bold mt-2">Moderation queue is clean. No pending stop submissions.</p>
                </div>
              ) : (
                <div className="overflow-x-auto">
                  <table className="w-full text-left text-xs border-collapse">
                    <thead>
                      <tr className="border-b border-slate-200 text-slate-400 font-bold uppercase tracking-wider">
                        <th className="pb-3">Stop Info</th>
                        <th className="pb-3">Coordinates</th>
                        <th className="pb-3">Area</th>
                        <th className="pb-3">Votes</th>
                        <th className="pb-3">Created By</th>
                        <th className="pb-3 text-right">Actions</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-100">
                      {pendingReports.map((report) => (
                        <tr key={report.reportId} className="hover:bg-slate-50/50">
                          <td className="py-4.5 font-bold text-slate-800">
                            {report.stopName}
                            <span className="block text-[10px] text-slate-400 capitalize font-medium">{report.stopType.replace("_", " ")}</span>
                          </td>
                          <td className="py-4.5 font-mono text-[10px] text-slate-500">
                            {report.latitude.toFixed(5)}, {report.longitude.toFixed(5)}
                          </td>
                          <td className="py-4.5 text-slate-600 font-medium">{report.area || "N/A"}</td>
                          <td className="py-4.5 font-black text-slate-700">{report.netVotes}</td>
                          <td className="py-4.5 text-slate-600">
                            {report.creatorName}
                            <span className="block text-[9px] text-slate-400 font-mono">{new Date(report.createdAt).toLocaleDateString()}</span>
                          </td>
                          <td className="py-4.5 text-right space-x-1">
                            <button
                              onClick={() => handleVerify(report.reportId)}
                              className="bg-emerald-800 hover:bg-emerald-950 text-white font-bold text-[10px] px-3.5 py-2 rounded-lg shadow-sm transition"
                            >
                              Approve
                            </button>
                            <button
                              onClick={() => handleReject(report.reportId)}
                              className="bg-red-50 hover:bg-red-600 text-red-600 hover:text-white border border-red-200 font-bold text-[10px] px-3.5 py-2 rounded-lg transition"
                            >
                              Reject
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          ) : activeTab === "logs" ? (
            <div>
              <h2 className="text-lg font-black text-emerald-800 mb-4 uppercase tracking-wider">📋 Administrative Audit Logs</h2>
              {auditLogs.length === 0 ? (
                <div className="text-center py-12">
                  <p className="text-xs text-slate-500 font-bold">No admin actions have been logged yet.</p>
                </div>
              ) : (
                <div className="overflow-x-auto">
                  <table className="w-full text-left text-xs border-collapse">
                    <thead>
                      <tr className="border-b border-slate-200 text-slate-400 font-bold uppercase tracking-wider">
                        <th className="pb-3">Admin</th>
                        <th className="pb-3">Action</th>
                        <th className="pb-3">Target Entity</th>
                        <th className="pb-3">Metadata</th>
                        <th className="pb-3 text-right">Timestamp</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-100">
                      {auditLogs.map((log) => (
                        <tr key={log.id} className="hover:bg-slate-50/50">
                          <td className="py-4 font-bold text-slate-800">{log.adminName}</td>
                          <td className="py-4">
                            <span className="bg-emerald-700/10 text-emerald-800 font-black px-2.5 py-1 rounded-full text-[10px]">
                              {log.action}
                            </span>
                          </td>
                          <td className="py-4 text-slate-500">
                            {log.entityType}
                            {log.entityId && <span className="block text-[10px] font-mono mt-0.5">{log.entityId}</span>}
                          </td>
                          <td className="py-4 font-mono text-[10px] text-slate-500 max-w-[200px] truncate" title={log.metadata}>
                            {log.metadata}
                          </td>
                          <td className="py-4 text-right text-slate-400 font-medium">
                            {new Date(log.createdAt).toLocaleString()}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          ) : (
            <div>
              <h2 className="text-lg font-black text-emerald-800 mb-4 uppercase tracking-wider">📈 System Analytics</h2>
              
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                <div className="bg-slate-50 border border-slate-200 rounded-2xl p-5 shadow-sm">
                  <span className="text-[10px] text-slate-400 font-bold uppercase tracking-wider">Total Queue count</span>
                  <span className="block font-black text-3xl text-emerald-800 mt-1">{pendingReports.length}</span>
                  <p className="text-[10px] text-slate-500 mt-2 font-medium">Stops waiting for moderation.</p>
                </div>
                <div className="bg-slate-50 border border-slate-200 rounded-2xl p-5 shadow-sm">
                  <span className="text-[10px] text-slate-400 font-bold uppercase tracking-wider">Total Audit logs</span>
                  <span className="block font-black text-3xl text-emerald-800 mt-1">{auditLogs.length}</span>
                  <p className="text-[10px] text-slate-500 mt-2 font-medium">Actions logged in current instance.</p>
                </div>
                <div className="bg-slate-50 border border-slate-200 rounded-2xl p-5 shadow-sm">
                  <span className="text-[10px] text-slate-400 font-bold uppercase tracking-wider">Database Status</span>
                  <span className="block font-black text-xl text-emerald-800 mt-2 flex items-center gap-1.5">
                    <span className="h-3.5 w-3.5 bg-emerald-500 rounded-full inline-block animate-pulse" />
                    Connected
                  </span>
                  <p className="text-[10px] text-slate-500 mt-2 font-medium">PostgreSQL/PostGIS Port 5433</p>
                </div>
              </div>
            </div>
          )}
        </section>
      </div>
    </main>
  );
}
