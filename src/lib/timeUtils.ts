/**
 * Converts a timestamp to JST (Japan Standard Time) timestamp.
 * @param {number} timestamp - The timestamp to convert.
 * @returns {Date} - The converted JST timestamp.
 */
const timeToJST = (timestamp: number | string | Date): Date => {
  const dt = new Date(
    typeof timestamp === "number" ? timestamp : new Date(timestamp).getTime()
  );
  const utc = new Date(dt.getTime() - dt.getTimezoneOffset() * 60 * 1000); // Convert to UTC
  const tzOffset = 540 * 60 * 1000; // JST is UTC+9
  return new Date(utc.getTime() + tzOffset); // Convert to JST
};

/**
 * Converts a JST (Japan Standard Time) timestamp to a formatted date string or an object with individual date components.
 * @param {number} time - The JST timestamp to convert.
 * @param {boolean} [format=false] - Determines whether to return a formatted date string or an object with individual date components. Default is false.
 * @returns {string|Object} - The formatted date string or an object with individual date components.
 */
export const timeToJSTstamp = (
  time: number | string | Date,
  format = false
): Object => {
  const dt = new Date(time);
  const year = dt.getFullYear();
  const month = dt.getMonth() + 1;
  const date = dt.getDate();
  const hour = dt.getHours();
  const minute = dt.getMinutes();
  const second = dt.getSeconds();

  let return_str;
  if (format == true) {
    return_str = `${year}/${month}/${date} ${hour}:${minute}:${second}`;
  } else {
    return_str = {
      year: year,
      month: month,
      date: date,
      hour: hour,
      minute: minute,
      second: second,
    };
  }
  return return_str;
};

export default {
  timeToJST,
  timeToJSTstamp,
};
