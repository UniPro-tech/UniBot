import {
  AudioPlayer,
  createAudioPlayer,
  createAudioResource,
  getVoiceConnection,
  VoiceConnection,
  VoiceConnectionReadyState,
  StreamType,
  NoSubscriberBehavior,
} from "@discordjs/voice";
import { Readable } from "stream";
import { RPC, Generate, Query } from "voicevox.js";
import { ALStorage, loggingSystem } from "..";

interface TTSQueueItem {
  text: string;
  styleId: number;
  priority?: number; // 0: 高優先度（ボイスチャンネル通知など）, 1: 通常（デフォルト）
  // 事前取得したVoiceVoxクエリ
  query?: any;
  // 事前取得中/完了のクエリPromise
  queryPromise?: Promise<any>;
  // 事前生成した音声データ（Buffer | Uint8Array 想定）
  audio?: any;
  // 事前生成中/完了の音声Promise
  audioPromise?: Promise<any>;
}

export class TTSQueue {
  private static instances: Map<string, TTSQueue> = new Map();
  private static voiceVoxInitialized = false;
  private queue: TTSQueueItem[] = [];
  private isProcessing = false;
  private guildId: string;
  private player?: AudioPlayer;

  private constructor(guildId: string) {
    this.guildId = guildId;
    this.initializeVoiceVox();
  }

  /**
   * VoiceVoxの初期化（非同期）
   */
  private initializeVoiceVox(): void {
    if (!TTSQueue.voiceVoxInitialized && process.env.VOICEVOX_API_URL) {
      TTSQueue.voiceVoxInitialized = true;
      // 非同期で初期化し、エラーはログのみ
      this.connectVoiceVox().catch((error) => {
        const ctx = ALStorage.getStore();
        const logger = loggingSystem.getLogger({ ...ctx, function: "TTSQueue.initializeVoiceVox" });
        logger.warn("Failed to initialize VoiceVox connection:", error as any);
        TTSQueue.voiceVoxInitialized = false; // 失敗したらリセット
      });
    }
  }

  /**
   * VoiceVox APIに接続
   */
  private async connectVoiceVox(): Promise<void> {
    if (!RPC.rpc && process.env.VOICEVOX_API_URL) {
      const headers = { Authorization: `ApiKey ${process.env.VOICEVOX_API_KEY}` };
      await RPC.connect(process.env.VOICEVOX_API_URL, headers);
    }
  }

  /**
   * ギルドごとのTTSQueueインスタンスを取得
   */
  public static getInstance(guildId: string): TTSQueue {
    if (!this.instances.has(guildId)) {
      this.instances.set(guildId, new TTSQueue(guildId));
    }
    return this.instances.get(guildId)!;
  }

  /**
   * キューにTTSアイテムを追加（クエリ事前取得）
   */
  public enqueue(text: string, styleId: number = 0, priority: number = 1): void {
    const ctx = ALStorage.getStore();
    const logger = loggingSystem.getLogger({ ...ctx, function: "TTSQueue.enqueue" });

    if (!text.trim()) return;

    const item: TTSQueueItem = { text, styleId, priority };

    // 優先度に基づいてキューに挿入
    if (priority === 0) {
      // 高優先度の場合は先頭に挿入（既に処理中でない場合）
      this.queue.unshift(item);
    } else {
      // 通常優先度の場合は末尾に追加
      this.queue.push(item);
    }

    logger.info(
      { text },
      `TTS item enqueued for guild ${this.guildId}. Queue length: ${this.queue.length}`
    );

    // クエリ→音声までを事前取得（非同期で先に走らせる）
    this.preloadItem(item).catch((error) => {
      logger.warn(`Failed to preload item for guild ${this.guildId}:`, error as any);
    });

    // 処理が停止していれば開始
    if (!this.isProcessing) {
      this.processQueue();
    }
  }

  /**
   * クエリ→音声の事前取得（非同期処理）
   * 再生順は変えず、重い処理だけ先に走らせてキャッシュする
   */
  private async preloadItem(item: TTSQueueItem): Promise<void> {
    const ctx = ALStorage.getStore();
    const logger = loggingSystem.getLogger({ ...ctx, function: "TTSQueue.preloadItem" });
    try {
      // VoiceVox接続確認
      await this.ensureVoiceVoxConnection();

      // クエリの事前取得（既に走っていれば使う）
      if (!item.queryPromise) {
        item.queryPromise = Query.getTalkQuery(item.text, item.styleId)
          .then((q) => {
            item.query = q;
            logger.debug(`Query preloaded for guild ${this.guildId}, styleId: ${item.styleId}`);
            return q;
          })
          .catch((err) => {
            // ここでの失敗は致命的ではない（再生直前に再取得する）
            logger.warn(
              { extra_context: { styleId: item.styleId } },
              "Query preload failed",
              err as any
            );
            // エラーは表に伝播しない（再生側でフォールバック）
            return undefined as any;
          });
      }

      // 音声の事前生成（クエリが解決してから）
      if (!item.audioPromise) {
        item.audioPromise = (async () => {
          const q = item.query ?? (await item.queryPromise);
          if (!q) {
            // クエリがない場合はここでは諦める（再生側で取得）
            return undefined as any;
          }
          const audio = await Generate.generate(item.styleId, q);
          item.audio = audio;
          logger.debug(
            `Audio pre-generated for guild ${this.guildId}, text length: ${item.text.length}`
          );
          return audio;
        })().catch((err) => {
          // ここでの失敗も致命的ではない（再生直前に再生成）
          logger.warn(
            { extra_context: { styleId: item.styleId } },
            "Audio preload failed",
            err as any
          );
          return undefined as any;
        });
      }
    } catch (error) {
      // 事前取得フェーズの例外は握りつぶしてOK（再生時にフォールバック）
      logger.warn("Preload pipeline failed (will fallback at play time)", error as any);
    }
  }

  /**
   * skipコマンドで現在の音声をスキップ
   */
  public skip(): boolean {
    const connection = getVoiceConnection(this.guildId);
    if (!connection) return false;

    const player = (connection.state as VoiceConnectionReadyState).subscription
      ?.player as AudioPlayer;
    if (player && player.state.status === "playing") {
      player.stop(true);
      return true;
    }
    return false;
  }

  /**
   * キューをクリア
   */
  public clear(): void {
    this.queue = [];
    const ctx = ALStorage.getStore();
    const logger = loggingSystem.getLogger({ ...ctx, function: "TTSQueue.clear" });
    logger.info(`TTS queue cleared for guild ${this.guildId}`);
  }

  /**
   * キューの長さを取得
   */
  public getQueueLength(): number {
    return this.queue.length;
  }

  /**
   * キューを順次処理
   */
  private async processQueue(): Promise<void> {
    if (this.isProcessing || this.queue.length === 0) return;

    this.isProcessing = true;
    const ctx = ALStorage.getStore();
    const logger = loggingSystem.getLogger({ ...ctx, function: "TTSQueue.processQueue" });

    try {
      while (this.queue.length > 0) {
        const connection = getVoiceConnection(this.guildId);
        if (!connection || connection.state.status === "destroyed") {
          logger.warn(`Voice connection not found or destroyed for guild ${this.guildId}`);
          break;
        }
        const item = this.queue.shift()!;
        logger.debug(
          `Processing TTS item for guild ${this.guildId}. Queue length: ${this.queue.length}`
        );
        await this.synthesizeAndPlay(item, connection);
      }
    } catch (error) {
      logger.error(
        {
          extra_context: {
            guildId: this.guildId,
            queueLength: this.queue.length,
            isProcessing: this.isProcessing,
          },
          stack_trace: (error as Error).stack,
        },
        "Error processing TTS queue:",
        error as any
      );
    } finally {
      this.isProcessing = false;
    }
  }

  /**
   * 音声合成と再生
   */
  private async synthesizeAndPlay(item: TTSQueueItem, connection: VoiceConnection): Promise<void> {
    const ctx = ALStorage.getStore();
    const logger = loggingSystem.getLogger({ ...ctx, function: "TTSQueue.synthesizeAndPlay" });

    try {
      if (!process.env.VOICEVOX_API_URL) {
        logger.warn("VOICEVOX_API_URL is not set.");
        return;
      }

      // VoiceVox接続の再確認（必要に応じて再接続）
      await this.ensureVoiceVoxConnection();

      // できるだけ事前生成済みの音声を使う（なければPromiseを待機、さらにダメならその場生成）
      let audio = item.audio;
      let path: "preloaded-audio" | "awaited-audioPromise" | "generated-on-demand" =
        "preloaded-audio";

      if (!audio && item.audioPromise) {
        path = "awaited-audioPromise";
        try {
          audio = await item.audioPromise;
        } catch {
          audio = undefined as any;
        }
      }

      if (!audio) {
        // 音声がまだなければクエリ（preload済みまたはその場生成）を使って生成
        path = "generated-on-demand";
        let query = item.query;
        if (!query && item.queryPromise) {
          try {
            query = await item.queryPromise;
          } catch {
            query = undefined as any;
          }
        }
        if (!query) {
          logger.debug(`Query not preloaded, generating on-demand for guild ${this.guildId}`);
          query = await Query.getTalkQuery(item.text, item.styleId);
          item.query = query;
        }
        audio = await Generate.generate(item.styleId, query);
        // 再利用できるようにキャッシュ
        item.audio = audio;
      }

      const startRequestTs = Date.now();
      logger.debug(
        {
          extra_context: {
            path,
            textLength: item.text.length,
            guildId: this.guildId,
          },
        },
        `Using audio via '${path}' for guild ${this.guildId}, text length: ${item.text.length}`
      );

      const audioStream = Readable.from(audio);
      const resource = createAudioResource(audioStream, {
        inputType: StreamType.Arbitrary,
        inlineVolume: false,
      });

      // プレイヤーの取得または作成（キャッシュされたものを使用）
      if (!this.player || this.player.state.status === "idle") {
        // 既存のプレイヤーがある場合はクリーンアップ
        if (this.player) {
          this.player.removeAllListeners();
        }
        this.player = createAudioPlayer({
          behaviors: {
            noSubscriber: NoSubscriberBehavior.Pause,
            maxMissedFrames: 5,
          },
        });
        connection.subscribe(this.player);
      }

      // 再生開始／完了の詳細ログを取る
      const playStartTs = Date.now();
      await this.playAudio(this.player, resource, {
        meta: { path, textLength: item.text.length, textSnippet: item.text.slice(0, 120) },
      });
      const playEndTs = Date.now();
      logger.info(
        {
          extra_context: {
            guildId: this.guildId,
            path,
            textLength: item.text.length,
            timeToStartMs: playStartTs - startRequestTs,
            playDurationMs: playEndTs - playStartTs,
          },
        },
        `Playback finished for guild ${this.guildId} (path=${path})`
      );
      // 再生後は音声キャッシュを解放してメモリ節約（クエリは残してOK）
      item.audio = undefined;
    } catch (error) {
      logger.error(
        {
          extra_context: {
            guildId: this.guildId,
            text: item.text,
            styleId: item.styleId,
            hasQuery: !!item.query,
            playerState: this.player?.state.status,
          },
          stack_trace: (error as Error).stack,
        },
        "Error in synthesizeAndPlay:",
        error as any
      );
      throw error;
    }
  }

  /**
   * VoiceVox接続を確認し、必要に応じて再接続
   */
  private async ensureVoiceVoxConnection(): Promise<void> {
    const ctx = ALStorage.getStore();
    const logger = loggingSystem.getLogger({
      ...ctx,
      function: "TTSQueue.ensureVoiceVoxConnection",
    });

    try {
      if (!RPC.rpc && process.env.VOICEVOX_API_URL) {
        logger.debug(`Reconnecting to VoiceVox for guild ${this.guildId}`);
        await this.connectVoiceVox();
        logger.debug(`VoiceVox reconnected for guild ${this.guildId}`);
      }
    } catch (error) {
      logger.error(
        {
          extra_context: {
            guildId: this.guildId,
            voicevoxUrl: process.env.VOICEVOX_API_URL,
          },
          stack_trace: (error as Error).stack,
        },
        "Failed to ensure VoiceVox connection:",
        error as any
      );
      throw error;
    }
  }

  /**
   * 音声リソースを再生し、完了まで待機
   */
  private async playAudio(
    player: AudioPlayer,
    resource: any,
    options?: { meta?: { path?: string; textLength?: number; textSnippet?: string } }
  ): Promise<void> {
    const meta = options?.meta;
    const loggerCtx = ALStorage.getStore();
    const logger = loggingSystem.getLogger({ ...loggerCtx, function: "TTSQueue.playAudio" });

    return new Promise<void>((resolve, reject) => {
      // 再生開始待ちタイムアウト（startTimeout）と、再生開始後は完了まで待つ挙動に分離
      // startTimeout: 再生が開始されるまでに何らかの問題で待ち続けないようにする（15秒）
      const START_TIMEOUT_MS = 15000;

      let startTimeout: NodeJS.Timeout | undefined;
      let started = false; // 再生が開始されたかどうか

      const cleanup = () => {
        if (startTimeout) clearTimeout(startTimeout);
        player.removeListener("stateChange", onStateChange);
        player.removeListener("error", onError);
      };

      const onStateChange = (oldState: any, newState: any) => {
        // 再生が開始（playing または buffer 挙動）した瞬間
        if (!started && newState.status === "playing") {
          started = true;
          // startTimeout を解除してから、完了（idle）を待つ。
          if (startTimeout) {
            clearTimeout(startTimeout);
            startTimeout = undefined;
          }
          // ここでは resolve しない（完了まで待つ）
          return;
        }

        // 再生完了
        if (started && newState.status === "idle") {
          cleanup();
          resolve();
        }
      };

      const onError = (error: Error) => {
        cleanup();
        reject(error);
      };

      try {
        player.on("stateChange", onStateChange);
        player.on("error", onError);

        // 再生開始
        player.play(resource);

        // 再生が 'playing' になるまでの待ち時間を制限（15秒）
        startTimeout = setTimeout(() => {
          // 再生が始まらなかった -> 強制クリーンアップしてエラーにする
          cleanup();
          logger.error(
            { extra_context: { meta } },
            "Audio playback did not start within 15 seconds"
          );
          reject(new Error("Audio playback did not start within 15 seconds"));
        }, START_TIMEOUT_MS);
      } catch (error) {
        const ctx = ALStorage.getStore();
        const logger = loggingSystem.getLogger({ ...ctx, function: "TTSQueue.playAudio" });
        logger.error(
          {
            extra_context: {
              playerState: player.state.status,
              resourcePlayable: resource?.playable,
            },
            stack_trace: (error as Error).stack,
          },
          "Error starting audio playback:",
          error as any
        );
        cleanup();
        reject(error);
      }
    });
  }

  /**
   * インスタンスを削除（ボイスチャンネル切断時など）
   */
  public static removeInstance(guildId: string): void {
    const instance = this.instances.get(guildId);
    if (instance) {
      instance.clear();
      // プレイヤーをクリーンアップ
      if (instance.player) {
        instance.player.stop(true);
        instance.player = undefined;
      }
      this.instances.delete(guildId);
    }
  }

  /**
   * グローバルVoiceVox初期化（アプリ起動時に呼び出し）
   */
  public static async initializeGlobal(): Promise<void> {
    const ctx = ALStorage.getStore();
    const logger = loggingSystem.getLogger({ ...ctx, function: "TTSQueue.initializeGlobal" });
    if (!this.voiceVoxInitialized && process.env.VOICEVOX_API_URL) {
      try {
        const headers = { Authorization: `ApiKey ${process.env.VOICEVOX_API_KEY}` };
        await RPC.connect(process.env.VOICEVOX_API_URL, headers);
        this.voiceVoxInitialized = true;
        logger.info("VoiceVox initialized successfully");

        // よく使われるメッセージのプリロード（同期実行で確実に準備）
        logger.info("Preloading common messages...");
        await this.preloadCommonMessages();
        logger.info("Common messages preloaded successfully");
      } catch (error) {
        logger.error("Failed to initialize VoiceVox:" + error);
      }
    }
  }

  /**
   * よく使われるメッセージをプリロードして初回読み上げを高速化
   */
  private static async preloadCommonMessages(): Promise<void> {
    const ctx = ALStorage.getStore();
    const logger = loggingSystem.getLogger({ ...ctx, function: "TTSQueue.preloadCommonMessages" });
    const commonMessages = [
      "に接続しました。",
      "が参加しました。",
      "から退出しました。",
      "に切り替えました。",
    ];

    for (const message of commonMessages) {
      try {
        // バックグラウンドでクエリと音声を事前生成（キャッシュ効果を期待）
        const query = await Query.getTalkQuery(message, 0);
        // 初回音声合成も事前に実行してサーバーキャッシュをウォームアップ
        await Generate.generate(0, query);
        logger.info(`Preloaded: "${message}"`);
      } catch {
        // エラーは無視
      }
    }
  }

  /**
   * VoiceVox初期化を待機してから接続メッセージを再生（joinコマンド用）
   */
  public static async enqueueConnectionMessage(
    guildId: string,
    channelName: string
  ): Promise<void> {
    const ctx = ALStorage.getStore();
    const logger = loggingSystem.getLogger({
      ...ctx,
      function: "TTSQueue.enqueueConnectionMessage",
    });
    // VoiceVoxの初期化を待機（最大3秒）
    const startTime = Date.now();
    while (!this.voiceVoxInitialized && Date.now() - startTime < 3000) {
      await new Promise((resolve) => setTimeout(resolve, 100));
    }

    const text = `${channelName}に接続しました。`;
    const ttsQueue = TTSQueue.getInstance(guildId);

    // 高優先度でキューに追加
    ttsQueue.enqueue(text, 0, 0);

    logger.info(`Connection message queued for guild ${guildId}: "${text}"`);
  }
}
