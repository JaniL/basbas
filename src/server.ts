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

const referralKey = "t7c5hxB";
const restaurantId = "589";
const baseUrl = "https://v2.tableonline.fi/instabook/bookings";
const sourceUrl =
  `${baseUrl}/availabilities/${referralKey}/${restaurantId}/` +
  "{persons}?locale=fi&date={date}";
const bookingUrl = `${baseUrl}/${referralKey}/`;

const MAX_PERSONS = 8;
const MIN_PERSONS = 1;
const CACHE_TTL_MS = 60_000;
const DAY_COUNT = 33;

const availabilityCache = new TtlCache<string>(CACHE_TTL_MS);

const buildSourceUrl = (date: Date, persons: number): string =>
  sourceUrl
    .replace("{persons}", String(persons))
    .replace("{date}", formatReverseIsoDate(date));

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

const fallbackMarkup = (date: Date): string =>
  `<section><p class="alert">Emme voineet hakea vapaita aikoja ${formatDate(date)}. <a href="${bookingUrl}" target="_blank" rel="noreferrer">Varaa suoraan TableOnline-palvelusta</a>.</p></section>`;

const renderArticle = (day: DayAvailability): string => {
  const details = day.markup ?? fallbackMarkup(day.date);
  return `<article>\n<h1>${formatDate(day.date)}</h1>\n${details}\n</article>`;
};

const fetchDayAvailability = async (
  date: Date,
  persons: number
): Promise<DayAvailability> => {
  const source = buildSourceUrl(date, persons);
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
  persons: number
): Promise<DayAvailability[]> => {
  const settled = await Promise.allSettled(
    dates.map((date) => fetchDayAvailability(date, persons))
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

const writePage = async (res: ServerResponse, persons: number): Promise<void> => {
  res.writeHead(200, { "Content-Type": "text/html; charset=utf-8" });
  res.write("<!DOCTYPE html>");
  res.write(head);

  const dayAvailabilities = await combineDays(getOpenDates(), persons);
  const errorCount = dayAvailabilities.filter((day) => day.error).length;
  if (errorCount > 0) {
    console.warn(`[availability] ${errorCount} days failed for persons=${persons}`);
  }

  res.end(
    dayAvailabilities.map(renderArticle).join("\n") +
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

  if (isGet && uri === "/") {
    void writePage(res, persons).catch((error: Error) => {
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

const head = `<html>
<head>
<!-- Server started: ${startedTime}-->
<script async src="https://www.googletagmanager.com/gtag/js?id=UA-4154602-10"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());
  gtag('config', 'UA-4154602-10');
</script>
<title>Basbas</title>
<meta name="apple-mobile-web-app-capable" content="yes"/>
<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">
<link rel="apple-touch-icon" href="menu.png"/>
<link id="page_favicon" href="/favicon.ico" rel="icon" type="image/x-icon"/>
<link rel="stylesheet" type="text/css" href="public/styles.css"/>
</head>
<body>
<article>
<header>
<h1>Basbas p&ouml;yt&auml;varaukset</h1>
<p class="subtitle"><a href="http://basbas.fi/bistro#reservation" target="_blank" rel="noreferrer">Varaa t&auml;&auml;lt&auml;</a></p>
</header>
</article>`;
