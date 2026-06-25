export type TransitStopType = "bus_stop" | "taxi_stand";

export interface TransitStop {
  id: string;
  name: string;
  type: TransitStopType;
  latitude: number;
  longitude: number;
  area: string | null;
  verified: boolean;
  votes: number;
  availabilityStatus: string;
  distanceMeters: number;
  walkingTimeMinutes: number;
}

export interface NearbyStopsResponse {
  success: boolean;
  message: string;
  data: TransitStop[];
}

export interface NearbyStopsParams {
  latitude: number;
  longitude: number;
  radiusMeters?: number;
  limit?: number;
}
