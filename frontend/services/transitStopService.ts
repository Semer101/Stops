import { getJson } from "@/services/apiClient";
import type { NearbyStopsParams, NearbyStopsResponse } from "@/types/transitStop";

export function getNearbyStops(params: NearbyStopsParams): Promise<NearbyStopsResponse> {
  return getJson<NearbyStopsResponse>("/api/stops/nearby", {
    lat: params.latitude,
    lng: params.longitude,
    radius: params.radiusMeters,
    limit: params.limit,
  });
}
