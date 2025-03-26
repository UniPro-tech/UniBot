declare module "@/lib/logUtils" {
  export interface LogUtils {
    write: (post_data: Object, api_name: string) => Promise<Object>;
    read: (api_name: string) => Promise<Object>;
  }
}
