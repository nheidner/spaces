import {
  Address,
  Location,
  Space,
  Uuid,
  UserUid,
  User,
  Sorting,
  Thread,
} from "../types";
import { fetchApi } from "./fetch_api";
import { parseQuery } from "./parse_query";

const radius = 500;

export const getSpacesByLocation = async (loc: Location) => {
  const locationParamValue = `${loc.longitude},${loc.latitude}`;
  const queryStr = parseQuery({ location: locationParamValue, radius });
  const url = `/spaces${queryStr}`;

  return fetchApi<Space[]>(url);
};

export const getSpaceById = async (spaceId: Uuid) => {
  const url = `/spaces/${encodeURIComponent(spaceId)}`;

  return fetchApi<Space>(url);
};

export const getAddress = async (loc: Location) => {
  const locationParamValue = `${loc.longitude},${loc.latitude}`;
  const queryStr = parseQuery({ location: locationParamValue });
  const url = `/address${queryStr}`;

  return fetchApi<Address>(url);
};

export const getToplevelThreads = async (spaceId: Uuid) => {
  const queryStr = parseQuery({ sort: "recent", offset: 0, count: 10 });
  const url = `/spaces/${encodeURIComponent(
    spaceId
  )}/toplevel-threads${queryStr}`;

  return fetchApi<Thread[]>(url);
};

export const getThreadWithMessages = async (
  spaceId: Uuid,
  threadId: Uuid,
  sorting: Sorting,
  count: number,
  offset: number
) => {
  const queryStr = parseQuery({
    sort: sorting,
    messages_offset: offset,
    messages_count: count,
  });

  const url = `/spaces/${encodeURIComponent(
    spaceId
  )}/threads/${encodeURIComponent(threadId)}${queryStr}`;

  return fetchApi<Thread>(url);
};

export const getUser = async (userId: UserUid) => {
  const url = `/users/${encodeURIComponent(userId)}`;

  return fetchApi<User>(url);
};

type SpaceParams = {
  name: string;
  themeColorHexaCode: string;
  radius: number;
  location: Location;
};

export const createSpace = async (spaceParams: SpaceParams) => {
  const url = "/spaces";

  return fetchApi<{ spaceId: Uuid }>(url, {
    method: "POST",
    body: JSON.stringify(spaceParams),
  });
};

export const createToplevelThread = async (spaceId: Uuid, content: string) => {
  const url = `/spaces/${encodeURIComponent(spaceId)}/toplevel-threads`;

  return fetchApi<{ threadId: Uuid; firstMessageId: Uuid }>(url, {
    method: "POST",
    body: JSON.stringify({ content, type: "text" }),
  });
};

export const createThread = async (
  spaceId: Uuid,
  threadId: Uuid,
  messageId: Uuid
) => {
  const url = `/spaces/${encodeURIComponent(
    spaceId
  )}/threads/${encodeURIComponent(threadId)}/messages/${encodeURIComponent(
    messageId
  )}/threads`;

  return fetchApi<{ threadId: Uuid }>(url, { method: "POST" });
};

export const createMessage = async (
  spaceId: Uuid,
  threadId: Uuid,
  content: string
) => {
  const url = `/spaces/${spaceId}/threads/${threadId}/messages`;

  return fetchApi<{ messageId: Uuid }>(url, {
    method: "POST",
    body: JSON.stringify({ content, type: "text" }),
  });
};
