import { createReadStream, stat } from "node:fs";
import { extname, resolve } from "node:path";
import { IncomingMessage, ServerResponse } from "node:http";

const MIME_BY_EXT: Record<string, string> = {
  ".css": "text/css; charset=utf-8",
  ".ico": "image/x-icon",
  ".png": "image/png"
};

const publicRoot = resolve(process.cwd(), "public");

const publicPathForUri = (uri: string): string | null => {
  const sanitized = uri.startsWith("/") ? uri.slice(1) : uri;
  const candidate = resolve(publicRoot, sanitized.replace(/^public\//, ""));
  if (!candidate.startsWith(publicRoot)) {
    return null;
  }
  return candidate;
};

export const serveStatic = (uri: string, _req: IncomingMessage, res: ServerResponse): void => {
  const filePath = publicPathForUri(uri);
  if (!filePath) {
    res.writeHead(404);
    res.end();
    return;
  }

  stat(filePath, (statsError, fileStats) => {
    if (statsError || !fileStats.isFile()) {
      res.writeHead(404);
      res.end();
      return;
    }

    const contentType = MIME_BY_EXT[extname(filePath)] ?? "application/octet-stream";
    res.writeHead(200, { "Content-Type": contentType });
    const stream = createReadStream(filePath);
    stream.pipe(res);
    stream.on("error", () => {
      if (!res.headersSent) {
        res.writeHead(404);
      }
      res.end();
    });
  });
};
