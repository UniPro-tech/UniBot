import { Collection } from "discord.js";
import { Command } from "./commands/types/Command";
import timeUtils from "@/lib/timeUtils";
import logUtils from "@/lib/logUtils";

declare module "discord.js" {
  interface Client {
    commands: Collection<string, Command>;
    fs: typeof import("fs");
    config: typeof import("@/config");
    function: {
      timeUtils: typeof timeUtils;
      logUtils: typeof logUtils;
    };
  }
}
