export const upperCaseFirst = (string) =>
  string.charAt(0).toUpperCase() + string.slice(1);

export const prettyName = (string) =>
  string
    .split("_")
    .map((s) => upperCaseFirst(s))
    .join(" ");
