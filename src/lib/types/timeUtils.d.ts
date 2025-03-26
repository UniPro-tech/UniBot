declare module "@/lib/timeUtils" {
  export const timeToJST: (
    time: string | number | Date,
    format: boolean
  ) => string | Object;
}
