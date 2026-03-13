const weekdays = [
  "sunnuntai",
  "maanantai",
  "tiistai",
  "keskiviikko",
  "torstai",
  "perjantai",
  "lauantai"
].map((value) => value.toUpperCase());

const pad = (value: number): string => (value < 10 ? "0" : "") + value;

export const formatDate = (date: Date): string =>
  `${weekdays[date.getDay()]} ${date.getDate()}.${date.getMonth() + 1}`;

export const formatReverseIsoDate = (date: Date): string =>
  `${pad(date.getDate())}-${pad(date.getMonth() + 1)}-${date.getFullYear()}`;

export const addDay = (date: Date, days: number): Date => {
  const newDate = new Date(date.getTime());
  newDate.setDate(newDate.getDate() + days);
  return newDate;
};

export const isOpen = (date: Date): boolean => ![0, 1, 6].includes(date.getDay());

export const toMidnight = (date: Date): Date => {
  const newDate = new Date(date.getTime());
  newDate.setHours(0, 0, 0, 0);
  return newDate;
};
