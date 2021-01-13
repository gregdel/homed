export const upperCaseFirst = (string) =>
  string.charAt(0).toUpperCase() + string.slice(1);

export const prettyName = (string) => upperCaseFirst(string).replace("_", " ");
