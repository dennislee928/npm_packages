import { D as Driver, S as StorageValue, a as Storage } from './types-e2b6a05c.js';
export { c as StorageMeta, T as TransactionOptions, U as Unwatch, b as WatchCallback, W as WatchEvent } from './types-e2b6a05c.js';

interface CreateStorageOptions {
    driver?: Driver;
}
declare function createStorage<T extends StorageValue>(options?: CreateStorageOptions): Storage<T>;
type Snapshot<T = string> = Record<string, T>;
declare function snapshot(storage: Storage, base: string): Promise<Snapshot<string>>;
declare function restoreSnapshot(driver: Storage, snapshot: Snapshot<StorageValue>, base?: string): Promise<void>;

declare function prefixStorage<T extends StorageValue>(storage: Storage<T>, base: string): Storage<T>;
declare function normalizeKey(key?: string): string;
declare function joinKeys(...keys: string[]): string;
declare function normalizeBaseKey(base?: string): string;

type DriverFactory<T> = (opts: T) => Driver;
declare function defineDriver<T = any>(factory: DriverFactory<T>): DriverFactory<T>;

declare const builtinDrivers: {
    azureAppConfiguration: string;
    azureCosmos: string;
    azureKeyVault: string;
    azureStorageBlob: string;
    azureStorageTable: string;
    cloudflareKVBinding: string;
    cloudflareKVHTTP: string;
    cloudflareR2Binding: string;
    fs: string;
    github: string;
    http: string;
    localStorage: string;
    lruCache: string;
    memory: string;
    mongodb: string;
    overlay: string;
    planetscale: string;
    redis: string;
    sessionStorage: string;
    vercelKV: string;
    /** @deprecated */
    "cloudflare-kv-binding": string;
    /** @deprecated */
    "cloudflare-kv-http": string;
};
type BuiltinDriverName = keyof typeof builtinDrivers;

export { BuiltinDriverName, CreateStorageOptions, Driver, Snapshot, Storage, StorageValue, builtinDrivers, createStorage, defineDriver, joinKeys, normalizeBaseKey, normalizeKey, prefixStorage, restoreSnapshot, snapshot };
