export interface WalkingStep {
  instruction: string;
  distanceMeters: number;
  durationSeconds: number;
}

export interface WalkingRoute {
  distanceMeters: number;
  durationSeconds: number;
  coordinates: [number, number][];
  steps: WalkingStep[];
  source: "osrm" | "fallback";
}
