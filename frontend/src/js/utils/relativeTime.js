const rtf = new Intl.RelativeTimeFormat("en", { numeric: "auto" });

const DIVISIONS = [
  { amount: 60, name: "seconds" },
  { amount: 60, name: "minutes" },
  { amount: 24, name: "hours" },
  { amount: 7, name: "days" },
  { amount: 4.34524, name: "weeks" },
  { amount: 12, name: "months" },
  { amount: Infinity, name: "years" },
];

export function relativeTime(date) {
  if (!date) return "";

  const d = typeof date === "string" ? new Date(date) : date;
  let duration = (d - new Date()) / 1000;

  for (const division of DIVISIONS) {
    if (Math.abs(duration) < division.amount) {
      return rtf.format(Math.round(duration), division.name);
    }
    duration /= division.amount;
  }

  return "";
}

export default relativeTime;
