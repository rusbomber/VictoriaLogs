import { TimeParams } from "../types";
import dayjs from "dayjs";
import { LOGS_BARS_VIEW, LOGS_GROUP_BY } from "../constants/logs";
import { LogHits, Logs } from "../api/types";
import { OTHER_HITS_LABEL } from "../components/Chart/BarHitsChart/hooks/useBarHitsOptions";

export const getStreamPairs = (value: string): string[] => {
  const pairs = /^{.+}$/.test(value) ? value.slice(1, -1).split(",") : [value];
  return pairs.filter(Boolean);
};

export const getHitsTimeParams = (period: TimeParams) => {
  const start = dayjs(period.start * 1000);
  const end = dayjs(period.end * 1000);
  const totalMs = end.diff(start, "milliseconds");
  const step = Math.ceil(totalMs / LOGS_BARS_VIEW) || 1;
  return { start, end, step };
};

export const convertToFieldFilter = (value: string, field = LOGS_GROUP_BY) => {
  const isKeyValue = /(.+)?=(".+")/.test(value);

  if (isKeyValue) {
    return value.replace(/=/, ": ");
  }

  // Escape double quotes in the field value
  return `${field}: ${JSON.stringify(value)}`;
};

export const calculateTotalHits = (hits: LogHits[]): number => {
  return hits.reduce((acc, item) => acc + (item.total || 0), 0);
};

export const sortLogHits = <T extends { label?: string }>(key: keyof T) => (a: T, b: T): number => {
  if (a.label === OTHER_HITS_LABEL) return 1;
  if (b.label === OTHER_HITS_LABEL) return -1;

  const aValue = a[key] as unknown as number;
  const bValue = b[key] as unknown as number;

  return bValue - aValue;
};

export const isSameLog = (a: Logs, b: Logs): boolean => {
  if (a._time !== b._time) return false;
  const keys = new Set([...Object.keys(a), ...Object.keys(b)]);
  return [...keys].every(k => a[k] === b[k]);
};

export const removeExactLog = (logs: Logs[], target: Logs): Logs[] => {
  return logs.filter(log => !isSameLog(log, target));
};
