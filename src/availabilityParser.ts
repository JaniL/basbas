const MARKUP_PATTERNS = [
  /\$\("#availabilities"\)\.html\('([\s\S]*?)'\);/,
  /\$\('#availabilities'\)\.html\('([\s\S]*?)'\);/
];

const decodeMarkup = (value: string): string =>
  value
    .replace(/\\n/g, "")
    .replace(/\\\//g, "/")
    .replace(/\\'/g, "'")
    .replace(/\\"/g, '"')
    .replace(/\\\\/g, "\\");

export const extractAvailabilityMarkup = (body: string): string | null => {
  for (const pattern of MARKUP_PATTERNS) {
    const match = body.match(pattern);
    if (match?.[1]) {
      return decodeMarkup(match[1]);
    }
  }
  return null;
};
