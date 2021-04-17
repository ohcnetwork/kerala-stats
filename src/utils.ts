import fuseJs from "https://deno.land/x/fuse@v6.4.1/dist/fuse.esm.min.js";
import {
  DATE_REGEX,
  DEFAULT_HEADERS,
  District,
  GEO_LSG,
  INDEX_PAGE,
} from "./constants.ts";
import { ERRORS } from "./errors.ts";

let sessid = "";

export const getHTMLString = async (url: string, referer = INDEX_PAGE) => {
  const headers = DEFAULT_HEADERS;
  headers["Referer"] = referer;
  if (sessid) {
    headers["Cookie"] = sessid;
  }
  const res = await fetch(url, {
    method: "GET",
    headers,
  });
  sessid = res.headers.get("Set-Cookie")?.split(";")[0] || "";
  return await res.text();
};

export const lastUpdatedToDate = (lastUpdated: string) => {
  const date = lastUpdated.split(" ")[0];
  if (!DATE_REGEX.test(date)) {
    throw new Error(`invalid date: ${date} last_updated: ${lastUpdated}`);
  }
  return date;
};

export const summarizeDistrictInfo = <T>(histories: Record<string, number>[]) =>
  histories.reduce((p, c) => {
    Object.keys(c).forEach((k) => {
      p[k] = ((p[k]) || 0) + (c[k]);
    });
    return p;
  }, {} as Record<string, number>) as unknown as T;

export const DistrictFuse = new fuseJs(Object.values(District), {
  includeScore: true,
  threshold: 0.8,
});

export const LSGFuse = Object.entries(GEO_LSG)
  .reduce((a, [key, lsgs]) => {
    a[key as District] = new fuseJs(lsgs, {
      includeScore: true,
      threshold: 0.8,
    });
    return a;
  }, {} as Record<District, fuseJs>);

export const sendWebhookMessage = async (msg: string) => {
  try {
    const url = Deno.env.get("WEBHOOK_URL");
    if (!url) {
      throw ERRORS.WEBHOOK_NO_URL;
    }
    await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ content: msg }),
    });
  } catch (error) {
    console.error(ERRORS.WEBHOOK_MESSAGE);
    console.error(error);
  }
};

export const readJson = async <T>(filePath: string) => {
  try {
    const jsonString = await Deno.readTextFile(filePath);
    return JSON.parse(jsonString) as T;
  } catch (error) {
    error.message = `${filePath}: ${error.message}`;
    throw error;
  }
};

export const writeJson = async (filePath: string, object: any) => {
  try {
    const jsonString = JSON.stringify(object);
    await Deno.writeTextFile(filePath, jsonString);
  } catch (error) {
    error.message = `${filePath}: ${error.message}`;
    throw error;
  }
};

export const handler = async (
  fn: (...args: any[]) => void,
  ...args: any[]
) => {
  try {
    return await fn(...args);
  } catch (error) {
    sendWebhookMessage(error.toString());
  }
};
