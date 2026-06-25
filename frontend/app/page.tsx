"use client";

import { useEffect, useState } from "react";
import dynamic from "next/dynamic";
import Link from "next/link";
import type { TransitStop } from "@/types/transitStop";
import type { PublicUser } from "@/types/auth";
import type { WalkingRoute } from "@/types/directions";
import { getNearbyStops } from "@/services/transitStopService";
import { buildFallbackWalkingRoute, getWalkingDirections } from "@/services/osrmService";
import Chatbot from "@/components/Chatbot";

interface Route {
  id: string;
  routeName: string;
  startPoint: string;
  destinationPoint: string;
  fare: number;
  city: string;
  estimatedTime: number;
  status: string;
}

interface Prediction {
  type: string;
  predictedValue: string;
  confidence: number;
  generatedAt: string;
}

// Dynamically import Map component to prevent server-side rendering errors (window is not defined)
const Map = dynamic(() => import("../components/Map"), { ssr: false });

export default function HomePage() {
  const [token, setToken] = useState<string | null>(null);
  const [user, setUser] = useState<PublicUser | null>(null);
  const [userLocation, setUserLocation] = useState<[number, number] | null>([9.0108, 38.7612]); // default Meskel Square
  const [mapCenter, setMapCenter] = useState<[number, number]>([9.0108, 38.7612]);
  const [manualLocation, setManualLocation] = useState("9.0108, 38.7612");
  const [locationError, setLocationError] = useState("");
  
  const [stops, setStops] = useState<TransitStop[]>([]);
  const [selectedStop, setSelectedStop] = useState<TransitStop | null>(null);
  const [filterType, setFilterType] = useState<"all" | "bus_stop" | "taxi_stand">("all");
  const [walkingRoute, setWalkingRoute] = useState<WalkingRoute | null>(null);
  const [walkingStatus, setWalkingStatus] = useState("");
  
  // Search & Routing state
  const [searchStart, setSearchStart] = useState("");
  const [searchDest, setSearchDest] = useState("");
  const [routes, setRoutes] = useState<Route[]>([]);
  const [selectedRoute, setSelectedRoute] = useState<Route | null>(null);
  const [routeFares, setRouteFares] = useState<number | null>(null);
  const [predictions, setPredictions] = useState<Prediction[]>([]);
  const [predictionLoading, setPredictionLoading] = useState(false);

  // Crowdsourcing Reporting Modal state
  const [isReportOpen, setIsReportOpen] = useState(false);
  const [reportName, setReportName] = useState("");
  const [reportType, setReportType] = useState<"bus_stop" | "taxi_stand">("bus_stop");
  const [reportLat, setReportLat] = useState<number>(9.0108);
  const [reportLng, setReportLng] = useState<number>(38.7612);
  const [reportArea, setReportArea] = useState("");
  const [reportError, setReportError] = useState("");
  const [reportSuccess, setReportSuccess] = useState("");

  // Chatbot Drawer state
  const [isChatOpen, setIsChatOpen] = useState(false);

  // Load Auth Session
  useEffect(() => {
    const savedToken = localStorage.getItem("stops.auth.token");
    const savedUser = localStorage.getItem("stops.auth.user");
    if (savedToken && savedUser) {
      setToken(savedToken);
      setUser(JSON.parse(savedUser));
    }

    // Attempt Geolocation
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const loc: [number, number] = [position.coords.latitude, position.coords.longitude];
          setUserLocation(loc);
          setMapCenter(loc);
        },
        () => {
          console.log("Using default Addis Ababa location coordinates.");
        }
      );
    }
  }, []);

  // Fetch Nearby Stops
  const fetchStops = async (lat: number, lng: number) => {
    try {
      const res = await getNearbyStops({
        latitude: lat,
        longitude: lng,
        radiusMeters: 1500,
        limit: 20,
      });
      if (res.success) {
        setStops(res.data);
      }
    } catch (err) {
      console.error("Failed to load nearby stops", err);
    }
  };

  useEffect(() => {
    fetchStops(mapCenter[0], mapCenter[1]);
  }, [mapCenter]);

  useEffect(() => {
    if (!userLocation || !selectedStop) {
      setWalkingRoute(null);
      setWalkingStatus("");
      return;
    }

    let isMounted = true;
    const destination: [number, number] = [selectedStop.latitude, selectedStop.longitude];

    setWalkingStatus("Loading walking directions...");
    getWalkingDirections(userLocation, destination)
      .then((route) => {
        if (!isMounted) return;
        setWalkingRoute(route);
        setWalkingStatus("");
      })
      .catch(() => {
        if (!isMounted) return;
        setWalkingRoute(buildFallbackWalkingRoute(userLocation, destination));
        setWalkingStatus("OSRM walking directions are unavailable. Showing a map guide line.");
      });

    return () => {
      isMounted = false;
    };
  }, [userLocation, selectedStop]);

  // Handle Log Out
  const handleLogout = () => {
    localStorage.removeItem("stops.auth.token");
    localStorage.removeItem("stops.auth.user");
    setToken(null);
    setUser(null);
  };

  const handleManualLocationSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    const parsed = parseManualCoordinates(manualLocation);
    if (!parsed) {
      setLocationError("Enter coordinates as latitude, longitude.");
      return;
    }

    setLocationError("");
    setSelectedStop(null);
    setUserLocation(parsed);
    setMapCenter(parsed);
  };

  // Handle Map Click (open report stop form)
  const handleMapClick = (lat: number, lng: number) => {
    setReportLat(lat);
    setReportLng(lng);
    setReportSuccess("");
    setReportError("");
    setReportName("");
    setReportArea("");
    setIsReportOpen(true);
  };

  // Submit Crowdsourced Stop Report
  const handleReportSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!reportName.trim()) {
      setReportError("Stop Name is required.");
      return;
    }

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}/api/stops/report`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`,
        },
        body: JSON.stringify({
          name: reportName,
          type: reportType,
          latitude: reportLat,
          longitude: reportLng,
          area: reportArea,
        }),
      });

      const resData = await response.json();
      if (response.ok && resData.success) {
        setReportSuccess("Stop reported successfully! It is pending community verification.");
        setReportError("");
        fetchStops(mapCenter[0], mapCenter[1]);
        setTimeout(() => setIsReportOpen(false), 2000);
      } else {
        setReportError(resData.message || "Failed to submit report.");
      }
    } catch {
      setReportError("Failed to connect to the server.");
    }
  };

  // Submit Vote on reported Stop
  const handleVote = async (stopId: string, value: number) => {
    if (!token) {
      alert("Please login to vote and contribute data.");
      return;
    }

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}/api/stops/vote`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`,
        },
        body: JSON.stringify({
          report_id: stopId, // the repo looks up stop ID
          vote_value: value,
        }),
      });

      const resData = await response.json();
      if (response.ok && resData.success) {
        fetchStops(mapCenter[0], mapCenter[1]);
      } else {
        alert(resData.message || "Failed to submit vote.");
      }
    } catch {
      alert("Error submitting vote.");
    }
  };

  // Handle Route Search
  const handleRouteSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}/api/routes?start=${encodeURIComponent(
          searchStart
        )}&dest=${encodeURIComponent(searchDest)}`
      );
      const resData = await response.json();
      if (response.ok && resData.success) {
        setRoutes(resData.data || []);
        setSelectedRoute(null);
        setPredictions([]);
      }
    } catch (err) {
      console.error("Failed to query routes", err);
    }
  };

  // Select Route & load Fares/Predictions
  const handleSelectRoute = async (route: Route) => {
    setSelectedRoute(route);
    setPredictionLoading(true);

    try {
      // 1. Fetch Fare
      const fareResponse = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}/api/fare?route_id=${route.id}`
      );
      const fareData = await fareResponse.json();
      if (fareResponse.ok && fareData.success && fareData.data) {
        setRouteFares(fareData.data.fare);
      } else {
        setRouteFares(null);
      }

      // 2. Fetch Heuristic Predictions
      const predResponse = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}/api/predictions?route_id=${route.id}`
      );
      const predData = await predResponse.json();
      if (predResponse.ok && predData.success) {
        setPredictions(predData.data || []);
      }
    } catch (err) {
      console.error("Error loading route details", err);
    } finally {
      setPredictionLoading(false);
    }
  };

  // Filtered stops to display in sidebar list
  const filteredStops = stops.filter((stop) => {
    if (filterType === "all") return true;
    return stop.type === filterType;
  });

  return (
    <main className="min-h-screen bg-slate-100 text-slate-900 flex flex-col font-sans relative overflow-hidden">
      {/* Header Panel */}
      <header className="absolute top-4 left-4 right-4 z-10 bg-white/95 backdrop-blur-md rounded-2xl shadow-lg border border-emerald-800/10 p-4 flex flex-col md:flex-row items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="h-10 w-10 bg-emerald-800 rounded-xl flex items-center justify-center text-white font-black text-xl tracking-wider shadow">
            S
          </div>
          <div>
            <h1 className="text-xl font-black text-emerald-800 tracking-tight flex items-center gap-2">
              STOPS <span className="text-xs bg-amber-500 text-emerald-950 font-bold px-2 py-0.5 rounded-full uppercase tracking-wide">Addis Ababa</span>
            </h1>
            <p className="text-xs text-slate-500 font-medium">Smart Transit Optimization & Planning</p>
          </div>
        </div>

        {/* Filters */}
        <div className="flex items-center gap-1.5 bg-slate-100 p-1.5 rounded-xl border border-slate-200">
          <button
            onClick={() => setFilterType("all")}
            className={`px-3 py-1.5 text-xs font-semibold rounded-lg transition-all ${
              filterType === "all" ? "bg-white text-emerald-800 shadow-sm" : "text-slate-600 hover:text-slate-900"
            }`}
          >
            All Stops
          </button>
          <button
            onClick={() => setFilterType("bus_stop")}
            className={`px-3 py-1.5 text-xs font-semibold rounded-lg transition-all ${
              filterType === "bus_stop" ? "bg-white text-emerald-800 shadow-sm" : "text-slate-600 hover:text-slate-900"
            }`}
          >
            Buses
          </button>
          <button
            onClick={() => setFilterType("taxi_stand")}
            className={`px-3 py-1.5 text-xs font-semibold rounded-lg transition-all ${
              filterType === "taxi_stand" ? "bg-white text-emerald-800 shadow-sm" : "text-slate-600 hover:text-slate-900"
            }`}
          >
            Taxis
          </button>
        </div>

        {/* Quick Actions & Auth */}
        <div className="flex items-center gap-2">
          <button
            onClick={() => setIsChatOpen(true)}
            className="rounded-xl border border-emerald-800/20 bg-emerald-50 px-4 py-2 text-xs font-bold text-emerald-800 shadow-sm transition hover:bg-emerald-800 hover:text-white"
            type="button"
          >
            🤖 AI Assistant
          </button>

          {token ? (
            <div className="flex items-center gap-3">
              <span className="text-xs font-semibold text-slate-600 bg-slate-100 px-3 py-1.5 rounded-xl border border-slate-200">
                👤 {user?.name}
              </span>
              <button
                onClick={handleLogout}
                className="rounded-xl bg-red-50 text-red-600 border border-red-200 px-3.5 py-2 text-xs font-bold hover:bg-red-600 hover:text-white transition shadow-sm"
              >
                Logout
              </button>
            </div>
          ) : (
            <Link
              className="rounded-xl bg-emerald-800 px-4 py-2 text-xs font-bold text-white shadow hover:bg-emerald-950 transition"
              href="/auth"
            >
              Sign In
            </Link>
          )}
        </div>
      </header>

      {/* Main Map & Side Drawer Split */}
      <div className="flex-1 w-full h-[100vh] relative z-0">
        {/* Fullscreen Map */}
        <Map
          stops={stops}
          onMapClick={handleMapClick}
          onCenterChange={(lat, lng) => setMapCenter([lat, lng])}
          userLocation={userLocation}
          selectedStop={selectedStop}
          walkingRoute={walkingRoute}
        />

        {/* Floating Side Drawer (collapsible or scrollable container) */}
        <aside className="absolute top-28 left-4 bottom-4 w-full max-w-[380px] z-10 bg-white/95 backdrop-blur-md rounded-2xl shadow-xl border border-emerald-800/10 p-5 flex flex-col gap-4 overflow-y-auto max-h-[calc(100vh-140px)]">
          {/* Instructions banner */}
          <div className="bg-amber-500/10 border border-amber-500/20 rounded-xl p-3 text-xs text-amber-800 font-medium">
            💡 Click anywhere on the map to suggest/report a new transit stop or taxi stand.
          </div>

          <section className="bg-white border border-slate-200 rounded-xl p-4">
            <h2 className="text-sm font-black text-emerald-800 uppercase tracking-wider mb-3">Manual Location</h2>
            <form onSubmit={handleManualLocationSubmit} className="space-y-2">
              <input
                type="text"
                value={manualLocation}
                onChange={(event) => setManualLocation(event.target.value)}
                placeholder="Latitude, longitude"
                className="w-full text-xs rounded-lg border border-slate-300 px-3 py-2.5 outline-none transition focus:border-emerald-700 focus:ring-1 focus:ring-emerald-700 bg-white"
              />
              <button
                type="submit"
                className="w-full bg-slate-900 hover:bg-emerald-950 text-white font-bold text-xs py-2 rounded-lg transition shadow"
              >
                Search Nearby Stops
              </button>
            </form>
            {locationError && <p className="mt-2 text-xs font-semibold text-red-600">{locationError}</p>}
          </section>

          {/* Route Planner Module */}
          <section className="bg-slate-50 border border-slate-200 rounded-xl p-4">
            <h2 className="text-sm font-black text-emerald-800 uppercase tracking-wider mb-3">📍 Addis Route Planner</h2>
            <form onSubmit={handleRouteSearch} className="space-y-3">
              <div>
                <input
                  type="text"
                  placeholder="From (e.g. Bole)"
                  value={searchStart}
                  onChange={(e) => setSearchStart(e.target.value)}
                  className="w-full text-xs rounded-lg border border-slate-300 px-3 py-2.5 outline-none transition focus:border-emerald-700 focus:ring-1 focus:ring-emerald-700 bg-white"
                />
              </div>
              <div>
                <input
                  type="text"
                  placeholder="To (e.g. Piazza)"
                  value={searchDest}
                  onChange={(e) => setSearchDest(e.target.value)}
                  className="w-full text-xs rounded-lg border border-slate-300 px-3 py-2.5 outline-none transition focus:border-emerald-700 focus:ring-1 focus:ring-emerald-700 bg-white"
                />
              </div>
              <button
                type="submit"
                className="w-full bg-emerald-800 hover:bg-emerald-950 text-white font-bold text-xs py-2 rounded-lg transition shadow"
              >
                Find Best Route
              </button>
            </form>

            {/* Routes List */}
            {routes.length > 0 && (
              <div className="mt-4 space-y-2">
                <p className="text-[10px] font-bold text-slate-400 uppercase tracking-wide">Available Routes</p>
                {routes.map((route) => (
                  <button
                    key={route.id}
                    onClick={() => handleSelectRoute(route)}
                    className={`w-full text-left p-3 rounded-lg border transition-all ${
                      selectedRoute?.id === route.id
                        ? "bg-emerald-800 text-white border-emerald-950 shadow-md"
                        : "bg-white text-slate-800 border-slate-200 hover:border-slate-300"
                    }`}
                  >
                    <div className="flex items-center justify-between">
                      <h4 className="font-bold text-xs">{route.routeName}</h4>
                      <span className="text-[10px] bg-emerald-700/10 text-emerald-800 font-semibold px-2 py-0.5 rounded-full border border-emerald-800/10">
                        {route.estimatedTime}m
                      </span>
                    </div>
                  </button>
                ))}
              </div>
            )}
          </section>

          {/* Route Details, Fares & Predictions */}
          {selectedRoute && (
            <section className="bg-emerald-900 text-white border border-emerald-950 rounded-xl p-4 space-y-3 shadow-md animate-fadeIn">
              <h3 className="font-black text-xs uppercase tracking-wider text-amber-400 border-b border-white/10 pb-2">
                📊 Travel Predictions & Fares
              </h3>
              
              <div className="grid grid-cols-2 gap-2 text-xs">
                <div className="bg-white/5 rounded-lg p-2.5 border border-white/5">
                  <span className="text-[10px] text-white/50 block font-medium">Estimated Time</span>
                  <span className="font-black text-sm text-amber-300">{selectedRoute.estimatedTime} min</span>
                </div>
                <div className="bg-white/5 rounded-lg p-2.5 border border-white/5">
                  <span className="text-[10px] text-white/50 block font-medium">Base Fare</span>
                  <span className="font-black text-sm text-emerald-300">
                    {routeFares ? `${routeFares} ETB` : "Unavailable"}
                  </span>
                </div>
              </div>

              {predictionLoading ? (
                <div className="text-center text-xs py-3 text-white/50">Calculating predictions...</div>
              ) : (
                predictions.length > 0 && (
                  <div className="space-y-2 pt-2">
                    <p className="text-[10px] text-white/60 font-bold uppercase tracking-wider">Transit Indicators</p>
                    <div className="space-y-1.5">
                      {predictions.map((p, idx) => {
                        let colorClass = "text-emerald-400";
                        if (p.predictedValue.includes("High") || p.predictedValue.includes("Busy")) {
                          colorClass = "text-amber-400";
                        }
                        return (
                          <div key={idx} className="flex justify-between items-center text-xs bg-white/5 px-3 py-2 rounded-lg border border-white/5">
                            <span className="text-white/70 font-semibold">{p.type}</span>
                            <span className={`font-black ${colorClass}`}>{p.predictedValue}</span>
                          </div>
                        );
                      })}
                    </div>
                  </div>
                )
              )}
            </section>
          )}

          {selectedStop && (
            <section className="bg-white border border-slate-200 rounded-xl p-4 space-y-3">
              <div>
                <h2 className="text-sm font-black text-emerald-800 uppercase tracking-wider">Walking Directions</h2>
                <p className="mt-1 text-xs font-semibold text-slate-600">{selectedStop.name}</p>
              </div>

              {walkingStatus && (
                <p className="rounded-lg bg-amber-50 border border-amber-200 px-3 py-2 text-xs font-semibold text-amber-800">
                  {walkingStatus}
                </p>
              )}

              {walkingRoute && (
                <>
                  <div className="grid grid-cols-2 gap-2 text-xs">
                    <div className="bg-slate-50 rounded-lg p-2.5 border border-slate-200">
                      <span className="text-[10px] text-slate-400 block font-bold uppercase">Distance</span>
                      <span className="font-black text-emerald-800">
                        {walkingRoute.distanceMeters > 0 ? `${Math.round(walkingRoute.distanceMeters)}m` : "Map guide"}
                      </span>
                    </div>
                    <div className="bg-slate-50 rounded-lg p-2.5 border border-slate-200">
                      <span className="text-[10px] text-slate-400 block font-bold uppercase">Walk time</span>
                      <span className="font-black text-emerald-800">
                        {walkingRoute.durationSeconds > 0 ? `${Math.max(1, Math.round(walkingRoute.durationSeconds / 60))} min` : `${selectedStop.walkingTimeMinutes} min`}
                      </span>
                    </div>
                  </div>

                  <ol className="space-y-2 text-xs">
                    {walkingRoute.steps.slice(0, 5).map((step, index) => (
                      <li key={`${step.instruction}-${index}`} className="flex gap-2 rounded-lg bg-slate-50 border border-slate-200 p-2">
                        <span className="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-emerald-800 text-[10px] font-black text-white">
                          {index + 1}
                        </span>
                        <span className="font-medium text-slate-700">{step.instruction}</span>
                      </li>
                    ))}
                  </ol>
                </>
              )}
            </section>
          )}

          {/* Nearby Stops List */}
          <section className="flex-1 flex flex-col min-h-[220px]">
            <h2 className="text-sm font-black text-emerald-800 uppercase tracking-wider mb-3">🚏 Nearby Transit Stops</h2>
            
            {filteredStops.length === 0 ? (
              <div className="flex-1 border-2 border-dashed border-slate-200 rounded-xl p-6 flex flex-col items-center justify-center text-center">
                <span className="text-3xl mb-2">🗺️</span>
                <p className="text-xs text-slate-500 font-semibold leading-relaxed">
                  This area hasn&apos;t been mapped yet. Help your community by adding the first transit stop.
                </p>
              </div>
            ) : (
              <div className="space-y-2 flex-1 max-h-[300px] overflow-y-auto">
                {filteredStops.map((stop) => (
                  <div
                    key={stop.id}
                    onClick={() => setSelectedStop(stop)}
                    className={`p-3 rounded-xl border transition-all cursor-pointer ${
                      selectedStop?.id === stop.id
                        ? "bg-slate-100 border-emerald-800/30 ring-1 ring-emerald-850/10 shadow-sm"
                        : "bg-white border-slate-200 hover:border-slate-300"
                    }`}
                  >
                    <div className="flex justify-between items-start">
                      <div>
                        <h4 className="font-bold text-xs text-slate-800">{stop.name}</h4>
                        <p className="text-[10px] text-slate-500 mt-0.5 capitalize">{stop.type.replace("_", " ")}</p>
                        {stop.distanceMeters && (
                          <p className="text-[10px] text-emerald-600 font-semibold mt-1">
                            🏃 {Math.round(stop.distanceMeters)}m (~{Math.round(stop.distanceMeters / 80)}m walk)
                          </p>
                        )}
                      </div>
                      
                      {/* Voting and Verifying Actions */}
                      <div className="flex items-center gap-1.5 bg-slate-100 px-2 py-1 rounded-lg border border-slate-200">
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            handleVote(stop.id, 1);
                          }}
                          className="text-xs hover:text-emerald-700 p-0.5"
                          title="Upvote Stop"
                        >
                          👍
                        </button>
                        <span className="text-xs font-black text-slate-700">{stop.votes}</span>
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            handleVote(stop.id, -1);
                          }}
                          className="text-xs hover:text-red-700 p-0.5"
                          title="Downvote Stop"
                        >
                          👎
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </section>
        </aside>
      </div>

      {/* Report Stop Modal Dialog */}
      {isReportOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4">
          <div className="bg-white rounded-2xl shadow-2xl border border-slate-200 max-w-md w-full p-6 animate-scaleIn">
            <h3 className="text-lg font-black text-emerald-800 mb-2">📢 Report Missing Transit Stop</h3>
            <p className="text-xs text-slate-500 mb-4">Suggest adding this location to the Addis Ababa transit database.</p>

            {reportError && (
              <div className="bg-red-50 text-red-700 border border-red-200 rounded-xl p-3 text-xs font-semibold mb-4">
                ⚠️ {reportError}
              </div>
            )}

            {reportSuccess && (
              <div className="bg-emerald-50 text-emerald-700 border border-emerald-200 rounded-xl p-3 text-xs font-semibold mb-4">
                ✅ {reportSuccess}
              </div>
            )}

            {!token ? (
              <div className="text-center py-4 space-y-3">
                <p className="text-xs text-slate-600 font-medium">You must be logged in to contribute transit data.</p>
                <Link
                  href="/auth"
                  className="inline-block bg-emerald-800 hover:bg-emerald-950 text-white font-bold text-xs px-6 py-2.5 rounded-xl shadow"
                >
                  Go to Sign In
                </Link>
                <button
                  type="button"
                  onClick={() => setIsReportOpen(false)}
                  className="block mx-auto text-xs text-slate-500 underline font-medium mt-2"
                >
                  Cancel
                </button>
              </div>
            ) : (
              <form onSubmit={handleReportSubmit} className="space-y-4">
                <div>
                  <label className="text-xs font-bold text-slate-700" htmlFor="stop-name">Stop Name</label>
                  <input
                    id="stop-name"
                    type="text"
                    required
                    placeholder="e.g. Churchill Road Bus Stop"
                    value={reportName}
                    onChange={(e) => setReportName(e.target.value)}
                    className="w-full text-xs rounded-xl border border-slate-300 px-3.5 py-2.5 outline-none mt-1 focus:border-emerald-700 focus:ring-1 focus:ring-emerald-700"
                  />
                </div>

                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="text-xs font-bold text-slate-700" htmlFor="stop-type">Type</label>
                    <select
                      id="stop-type"
                      value={reportType}
                      onChange={(e) => setReportType(e.target.value as "bus_stop" | "taxi_stand")}
                      className="w-full text-xs rounded-xl border border-slate-300 px-3.5 py-2.5 outline-none mt-1 focus:border-emerald-700 bg-white"
                    >
                      <option value="bus_stop">Bus Stop</option>
                      <option value="taxi_stand">Taxi Stand</option>
                    </select>
                  </div>
                  <div>
                    <label className="text-xs font-bold text-slate-700" htmlFor="stop-area">Area/Neighborhood</label>
                    <input
                      id="stop-area"
                      type="text"
                      placeholder="e.g. Piazza"
                      value={reportArea}
                      onChange={(e) => setReportArea(e.target.value)}
                      className="w-full text-xs rounded-xl border border-slate-300 px-3.5 py-2.5 outline-none mt-1 focus:border-emerald-700"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-3 text-xs bg-slate-50 border border-slate-200 p-3 rounded-xl">
                  <div>
                    <span className="text-slate-400 block font-medium">Latitude</span>
                    <span className="font-bold text-slate-700">{reportLat.toFixed(6)}</span>
                  </div>
                  <div>
                    <span className="text-slate-400 block font-medium">Longitude</span>
                    <span className="font-bold text-slate-700">{reportLng.toFixed(6)}</span>
                  </div>
                </div>

                <div className="flex gap-2 justify-end pt-2">
                  <button
                    type="button"
                    onClick={() => setIsReportOpen(false)}
                    className="border border-slate-200 bg-white hover:bg-slate-50 text-slate-700 font-bold text-xs px-4 py-2.5 rounded-xl transition"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    className="bg-emerald-800 hover:bg-emerald-950 text-white font-bold text-xs px-5 py-2.5 rounded-xl shadow transition"
                  >
                    Submit Report
                  </button>
                </div>
              </form>
            )}
          </div>
        </div>
      )}

      {/* Floating Chatbot Assistant Drawer */}
      <Chatbot
        isOpen={isChatOpen}
        onClose={() => setIsChatOpen(false)}
        token={token}
      />
    </main>
  );
}

function parseManualCoordinates(value: string): [number, number] | null {
  const [latRaw, lngRaw] = value.split(",").map((part) => part.trim());
  if (!latRaw || !lngRaw) {
    return null;
  }

  const latitude = Number(latRaw);
  const longitude = Number(lngRaw);
  if (!Number.isFinite(latitude) || !Number.isFinite(longitude)) {
    return null;
  }
  if (latitude < -90 || latitude > 90 || longitude < -180 || longitude > 180) {
    return null;
  }

  return [latitude, longitude];
}
