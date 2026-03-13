import { createServer, ServerResponse } from "node:http";
import { parse } from "node:url";
import {
  addDay,
  formatDate,
  formatReverseIsoDate,
  isOpen,
  toMidnight
} from "./dateUtil";
import { fetchText } from "./tablesClient";
import { extractAvailabilityMarkup } from "./availabilityParser";
import { TtlCache } from "./cache";
import { DayAvailability } from "./types";
import { serveStatic } from "./static";

const port = Number(process.env.PORT ?? 5000);
const startMsg =
  "\u001b[33mServer started in \u001b[36mhttp://localhost:" +
  port +
  ", \u001b[33mtook \u001b[39m";
const startedTime = new Date().toString();

const baseUrl = "https://v2.tableonline.fi/instabook/bookings";

type Venue = {
  name: string;
  referralKey: string;
  restaurantId: string;
};

const venues: Record<"bistro" | "kulma", Venue> = {
  bistro: {
    name: "BasBas",
    referralKey: "t7c5hxB",
    restaurantId: "589"
  },
  kulma: {
    name: "BasBas Kulma",
    referralKey: "kxSrGpi",
    restaurantId: "1100"
  }
};

const MAX_PERSONS = 8;
const MIN_PERSONS = 1;
const CACHE_TTL_MS = 60_000;
const DAY_COUNT = 33;

const availabilityCache = new TtlCache<string>(CACHE_TTL_MS);

const buildSourceUrl = (date: Date, persons: number, venue: Venue): string =>
  `${baseUrl}/availabilities/${venue.referralKey}/${venue.restaurantId}/${persons}?locale=fi&date=${formatReverseIsoDate(date)}`;

const buildBookingUrl = (venue: Venue): string => `${baseUrl}/${venue.referralKey}/`;

const clampPersons = (value: number): number => {
  if (Number.isNaN(value)) {
    return 2;
  }
  if (value < MIN_PERSONS) {
    return MIN_PERSONS;
  }
  if (value > MAX_PERSONS) {
    return MAX_PERSONS;
  }
  return value;
};

const notFound = (res: ServerResponse): void => {
  res.writeHead(404);
  res.end();
};

const fallbackMarkup = (date: Date, venue: Venue): string =>
  `<section><p class="alert">Emme voineet hakea vapaita aikoja ${formatDate(date)}. <a href="${buildBookingUrl(venue)}" target="_blank" rel="noreferrer">Varaa suoraan TableOnline-palvelusta</a>.</p></section>`;

const renderArticle = (day: DayAvailability, venue: Venue): string => {
  const details = day.markup ?? fallbackMarkup(day.date, venue);
  return `<article>\n<h1>${formatDate(day.date)}</h1>\n${details}\n</article>`;
};

const fetchDayAvailability = async (
  date: Date,
  persons: number,
  venue: Venue
): Promise<DayAvailability> => {
  const source = buildSourceUrl(date, persons, venue);
  const cached = availabilityCache.get(source);
  if (cached) {
    return { date, markup: cached };
  }

  const response = await fetchText(source);
  if (!response.ok) {
    return {
      date,
      markup: null,
      error: response.error ?? "request_failed"
    };
  }

  const markup = extractAvailabilityMarkup(response.body);
  if (!markup) {
    return { date, markup: null, error: "parse_failed" };
  }

  availabilityCache.set(source, markup);
  return { date, markup };
};

const getOpenDates = (): Date[] =>
  [...Array(DAY_COUNT).keys()]
    .map((delta) => addDay(toMidnight(new Date()), delta))
    .filter(isOpen);

const combineDays = async (
  dates: Date[],
  persons: number,
  venue: Venue
): Promise<DayAvailability[]> => {
  const settled = await Promise.allSettled(
    dates.map((date) => fetchDayAvailability(date, persons, venue))
  );

  return settled.map((result, index) => {
    if (result.status === "fulfilled") {
      return result.value;
    }

    return {
      date: dates[index],
      markup: null,
      error: "unexpected_error"
    };
  });
};

const writePage = async (
  res: ServerResponse,
  persons: number,
  venue: Venue,
  showKulma: boolean
): Promise<void> => {
  res.writeHead(200, { "Content-Type": "text/html; charset=utf-8" });
  res.write("<!DOCTYPE html>");
  res.write(head(venue, showKulma));

  const dayAvailabilities = await combineDays(getOpenDates(), persons, venue);
  const errorCount = dayAvailabilities.filter((day) => day.error).length;
  if (errorCount > 0) {
    console.warn(
      `[availability] ${errorCount} days failed for persons=${persons} venue=${venue.referralKey}`
    );
  }

  res.end(
    dayAvailabilities.map((day) => renderArticle(day, venue)).join("\n") +
      '<article><a href="https://github.com/eeroan/basbas" target="_blank" rel="noreferrer">L&auml;hdekoodi</a></article></body></html>'
  );
};

console.time(startMsg);
createServer((req, res) => {
  const parsedUrl = parse(req.url ?? "", true);
  const uri = parsedUrl.pathname ?? "/";
  const isGet = req.method === "GET";
  const personsRaw = parsedUrl.query?.persons;
  const persons = clampPersons(typeof personsRaw === "string" ? Number(personsRaw) : 2);
  const showKulma = parsedUrl.query?.kulma !== undefined;
  const venue = showKulma ? venues.kulma : venues.bistro;

  if (isGet && uri === "/") {
    void writePage(res, persons, venue, showKulma).catch((error: Error) => {
      console.error(`[server] failed to render page: ${error.message}`);
      res.writeHead(502, { "Content-Type": "text/plain; charset=utf-8" });
      res.end("Failed to render page");
    });
    return;
  }

  if (
    isGet &&
    (uri.startsWith("/public/") || uri === "/favicon.ico" || uri === "/menu.png")
  ) {
    const remappedUri = uri.startsWith("/public/") ? uri : `/public${uri}`;
    serveStatic(remappedUri, req, res);
    return;
  }

  notFound(res);
}).listen(port, () => console.timeEnd(startMsg));

const head = (venue: Venue, showKulma: boolean): string => `<html>
<head>
<!-- Server started: ${startedTime}-->
<title>BasBas</title>
<meta name="apple-mobile-web-app-capable" content="yes"/>
<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">
<link rel="apple-touch-icon" href="menu.png"/>
<link id="page_favicon" href="/favicon.ico" rel="icon" type="image/x-icon"/>
<link rel="stylesheet" type="text/css" href="public/styles.css"/>
</head>
<body>
<article>
<header>
<h1>${venue.name} p&ouml;yt&auml;&shy;varaukset</h1>
<p class="subtitle"><a href="http://basbas.fi/bistro#reservation" target="_blank" rel="noreferrer">Varaa t&auml;&auml;lt&auml;</a></p>
${
  showKulma
    ? '<p class="subtitle"><a href="/">BasBas Bistro p&ouml;yt&auml;&shy;varaukset</a></p>'
    : '<p class="subtitle"><a href="/?kulma=1">BasBas Kulma p&ouml;yt&auml;&shy;varaukset</a></p>'
}
</header>
</article>`;
