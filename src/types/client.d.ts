import { Collection } from "discord.js";
import { Command } from "@/commands/types/Command";
import timeUtils from "@/lib/timeUtils";
import logUtils from "@/lib/dataUtils";
import { StringSelectMenuDefineType } from "@/selectMenus/types/SelectMenuDefineType";

declare module "discord.js" {
  interface Client {
    fs: typeof import("fs");
    config: typeof import("@/config");
    function: {
      timeUtils: typeof timeUtils;
      logUtils: typeof logUtils;
    };
    interactionExecutorsCollections: {
      chatInputCommands: Collection<string, Command>;
      stringSelectMenus: Collection<string, StringSelectMenuDefineType>;
    };
  }
}
