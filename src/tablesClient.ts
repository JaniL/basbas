import { get as httpGet } from "node:http";
import { get as httpsGet } from "node:https";
import { URL } from "node:url";
import { FetchResult } from "./types";

const REQUEST_TIMEOUT_MS = 8_000;

const USER_AGENT =
  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36";

export const fetchText = async (uri: string): Promise<FetchResult> =>
  new Promise((resolve) => {
    const parsed = new URL(uri);
    const options = {
      hostname: parsed.hostname,
      path: `${parsed.pathname}${parsed.search}`,
      headers: {
        "user-agent": USER_AGENT,
        Accept: "*/*"
      }
    };

    const request = (uri.startsWith("https") ? httpsGet : httpGet)(
      options,
      (response) => {
        const chunks: string[] = [];
        response.setEncoding("utf8");
        response.on("data", (chunk) => chunks.push(chunk));
        response.on("end", () => {
          const body = chunks.join("");
          const status = response.statusCode ?? 0;
          resolve({
            ok: status >= 200 && status < 300,
            status,
            body,
            error: status >= 200 && status < 300 ? undefined : `status_${status}`
          });
        });
      }
    );

    request.setTimeout(REQUEST_TIMEOUT_MS, () => {
      request.destroy(new Error("timeout"));
    });

    request.on("error", (error: Error) => {
      resolve({
        ok: false,
        status: 0,
        body: "",
        error: error.message
      });
    });
  });
