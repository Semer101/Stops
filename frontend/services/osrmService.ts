import type { WalkingRoute, WalkingStep } from "@/types/directions";

const osrmBaseUrl = process.env.NEXT_PUBLIC_OSRM_BASE_URL ?? "https://router.project-osrm.org";

interface OsrmRouteResponse {
  code: string;
  routes?: OsrmRoute[];
}

interface OsrmRoute {
  distance: number;
  duration: number;
  geometry: {
    coordinates: [number, number][];
  };
  legs: OsrmLeg[];
}

interface OsrmLeg {
  steps: OsrmStep[];
}

interface OsrmStep {
  distance: number;
  duration: number;
  name: string;
  maneuver: {
    type: string;
    modifier?: string;
  };
}

export async function getWalkingDirections(
  origin: [number, number],
  destination: [number, number],
): Promise<WalkingRoute> {
  const [originLat, originLng] = origin;
  const [destinationLat, destinationLng] = destination;
  const coordinates = `${originLng},${originLat};${destinationLng},${destinationLat}`;
  const url = `${osrmBaseUrl}/route/v1/foot/${coordinates}?overview=full&geometries=geojson&steps=true`;

  const response = await fetch(url);
  if (!response.ok) {
    throw new Error("Walking directions are unavailable.");
  }

  const body = (await response.json()) as OsrmRouteResponse;
  const route = body.routes?.[0];
  if (body.code !== "Ok" || !route) {
    throw new Error("Walking directions are unavailable.");
  }

  return {
    distanceMeters: route.distance,
    durationSeconds: route.duration,
    coordinates: route.geometry.coordinates.map(([lng, lat]) => [lat, lng]),
    steps: route.legs.flatMap((leg) => leg.steps.map(formatStep)),
    source: "osrm",
  };
}

export function buildFallbackWalkingRoute(
  origin: [number, number],
  destination: [number, number],
): WalkingRoute {
  return {
    distanceMeters: 0,
    durationSeconds: 0,
    coordinates: [origin, destination],
    steps: [
      {
        instruction: "Walk toward the selected stop using the map line as a visual guide.",
        distanceMeters: 0,
        durationSeconds: 0,
      },
    ],
    source: "fallback",
  };
}

function formatStep(step: OsrmStep): WalkingStep {
  const action = step.maneuver.modifier
    ? `${step.maneuver.type} ${step.maneuver.modifier}`
    : step.maneuver.type;
  const roadName = step.name ? ` on ${step.name}` : "";

  return {
    instruction: `${capitalize(action.replaceAll("_", " "))}${roadName}.`,
    distanceMeters: step.distance,
    durationSeconds: step.duration,
  };
}

function capitalize(value: string): string {
  return value.length === 0 ? value : `${value[0].toUpperCase()}${value.slice(1)}`;
}
