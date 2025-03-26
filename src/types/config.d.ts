declare module "@/config" {
  export const adminRoleId: string;
  export const color: {
    s: number;
    e: number;
  };
  export const token: string;
  export const dev: {
    testGuild: string;
  };
  export const logch: {
    ready: string;
    command: string;
    error: string;
    guildCreate: string;
    guildDelete: string;
  };
}
