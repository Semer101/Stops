"use client";

import { useEffect, useRef } from "react";
import L from "leaflet";
import type { TransitStop } from "@/types/transitStop";
import type { WalkingRoute } from "@/types/directions";

interface MapProps {
  stops: TransitStop[];
  onMapClick?: (lat: number, lng: number) => void;
  onCenterChange?: (lat: number, lng: number) => void;
  userLocation: [number, number] | null;
  selectedStop: TransitStop | null;
  walkingRoute: WalkingRoute | null;
}

export default function Map({
  stops,
  onMapClick,
  onCenterChange,
  userLocation,
  selectedStop,
  walkingRoute,
}: MapProps) {
  const mapContainerRef = useRef<HTMLDivElement>(null);
  const mapRef = useRef<L.Map | null>(null);
  const markersRef = useRef<{ [key: string]: L.Marker }>({});
  const userMarkerRef = useRef<L.Marker | null>(null);
  const routingPolylineRef = useRef<L.Polyline | null>(null);

  // Initialize Map
  useEffect(() => {
    if (!mapContainerRef.current || mapRef.current) return;

    // Default center to Addis Ababa (Meskel Square)
    const initialCenter: [number, number] = userLocation || [9.0108, 38.7612];

    const map = L.map(mapContainerRef.current, {
      zoomControl: false,
    }).setView(initialCenter, 14);

    L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
      attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
    }).addTo(map);

    L.control.zoom({
      position: "bottomright",
    }).addTo(map);

    // Map events
    map.on("dragend", () => {
      const center = map.getCenter();
      if (onCenterChange) {
        onCenterChange(center.lat, center.lng);
      }
    });

    map.on("zoomend", () => {
      const center = map.getCenter();
      if (onCenterChange) {
        onCenterChange(center.lat, center.lng);
      }
    });

    map.on("click", (e: L.LeafletMouseEvent) => {
      if (onMapClick) {
        onMapClick(e.latlng.lat, e.latlng.lng);
      }
    });

    mapRef.current = map;

    return () => {
      if (mapRef.current) {
        mapRef.current.remove();
        mapRef.current = null;
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Update user location marker
  useEffect(() => {
    if (!mapRef.current) return;

    if (userLocation) {
      const icon = L.divIcon({
        className: "custom-user-marker",
        html: `
          <div class="relative flex h-5 w-5">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-sky-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-5 w-5 bg-sky-500 border-2 border-white shadow-md"></span>
          </div>
        `,
        iconSize: [20, 20],
        iconAnchor: [10, 10],
      });

      if (userMarkerRef.current) {
        userMarkerRef.current.setLatLng(userLocation);
      } else {
        userMarkerRef.current = L.marker(userLocation, { icon }).addTo(mapRef.current);
      }
    } else {
      if (userMarkerRef.current) {
        userMarkerRef.current.remove();
        userMarkerRef.current = null;
      }
    }
  }, [userLocation]);

  // Update stop markers
  useEffect(() => {
    if (!mapRef.current) return;

    const map = mapRef.current;
    const currentMarkers = markersRef.current;
    const newMarkers: { [key: string]: L.Marker } = {};

    stops.forEach((stop) => {
      const isSelected = selectedStop && selectedStop.id === stop.id;
      const markerColorClass = stop.type === "bus_stop" ? "bg-emerald-600" : "bg-amber-500";
      const borderClass = isSelected ? "border-4 border-red-500 scale-125" : "border-2 border-white";
      const label = stop.type === "bus_stop" ? "B" : "T";

      const icon = L.divIcon({
        className: "custom-stop-marker",
        html: `
          <div class="flex items-center justify-center h-8 w-8 rounded-full ${markerColorClass} ${borderClass} text-white font-bold text-xs shadow-lg transition-all duration-300 transform hover:scale-110">
            ${label}
          </div>
        `,
        iconSize: [32, 32],
        iconAnchor: [16, 16],
      });

      const latlng: [number, number] = [stop.latitude, stop.longitude];

      if (currentMarkers[stop.id]) {
        // Update existing marker icon and position
        currentMarkers[stop.id].setLatLng(latlng);
        currentMarkers[stop.id].setIcon(icon);
        newMarkers[stop.id] = currentMarkers[stop.id];
        delete currentMarkers[stop.id];
      } else {
        // Create new marker
        const marker = L.marker(latlng, { icon }).addTo(map);
        marker.bindPopup(`
          <div class="p-1 font-sans">
            <h3 class="font-bold text-sm text-slate-800">${stop.name}</h3>
            <p class="text-xs text-slate-500 mt-1 capitalize">${stop.type.replace("_", " ")}</p>
            <p class="text-xs text-emerald-600 font-semibold mt-1">Status: ${stop.availabilityStatus}</p>
          </div>
        `);
        newMarkers[stop.id] = marker;
      }
    });

    // Remove obsolete markers
    Object.keys(currentMarkers).forEach((id) => {
      currentMarkers[id].remove();
    });

    markersRef.current = newMarkers;
  }, [stops, selectedStop]);

  // Pan to selected stop
  useEffect(() => {
    if (!mapRef.current || !selectedStop) return;
    mapRef.current.setView([selectedStop.latitude, selectedStop.longitude], 16, {
      animate: true,
      duration: 1.0,
    });
  }, [selectedStop]);

  // Render walking path polyline between user and selected stop
  useEffect(() => {
    if (!mapRef.current) return;

    // Clear previous polyline
    if (routingPolylineRef.current) {
      routingPolylineRef.current.remove();
      routingPolylineRef.current = null;
    }

    if (userLocation && selectedStop) {
      const latlngs: [number, number][] =
        walkingRoute?.coordinates.length
          ? walkingRoute.coordinates
          : [userLocation, [selectedStop.latitude, selectedStop.longitude]];

      routingPolylineRef.current = L.polyline(latlngs, {
        color: "#0f766e", // emerald-700
        weight: 4,
        dashArray: walkingRoute?.source === "osrm" ? undefined : "5, 10",
        lineCap: "round",
        lineJoin: "round"
      }).addTo(mapRef.current);

      // Fit bounds to show both user and stop
      const bounds = L.latLngBounds(latlngs);
      mapRef.current.fitBounds(bounds, { padding: [50, 50] });
    }
  }, [userLocation, selectedStop, walkingRoute]);

  return (
    <div className="h-full w-full relative z-0 min-h-[400px]">
      <div ref={mapContainerRef} className="h-full w-full min-h-[400px]" style={{ height: '100%', minHeight: '400px' }} />
    </div>
  );
}
