import { Location, Space, Uuid } from "../types";
import { fetchApi } from "./fetch_api";
import { parseQuery } from "./parse_query";

const radius = 500;

export const getSpacesByLocation = async (loc: Location) => {
  const locationParamValue = `${loc.longitude},${loc.latitude}`;
  const queryStr = parseQuery({ location: locationParamValue, radius });

  const spaces = await fetchApi<Space[]>(`/spaces${queryStr}`);
  return spaces;
};

export const getSpaceById = async (spaceId: Uuid) => {
  return fetchApi<Space>(`/spaces/${encodeURIComponent(spaceId)}`);
};
