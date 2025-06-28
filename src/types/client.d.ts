import { Collection } from "discord.js";
import { ChatInputCommand } from "@/executors/types/ChatInputCommand";
import timeUtils from "@/lib/timeUtils";
import logUtils from "@/lib/dataUtils";
import { StringSelectMenu } from "@/executors/types/StringSelectMenu";
import { Button } from "@/executors/types/Button";

declare module "discord.js" {
  interface Client {
    fs: typeof import("fs");
    config: typeof import("@/config");
    functions: {
      timeUtils: typeof timeUtils;
      logUtils: typeof logUtils;
    };
    interactionExecutorsCollections: {
      chatInputCommands: Collection<string, ChatInputCommand>;
      stringSelectMenus: Collection<string, StringSelectMenu>;
      messageContextMenuCommands: Collection<string, ChatInputCommand>;
      buttons: Collection<string, Button>;
    };
  }
}
