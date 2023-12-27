declare const _default: any;
export default _default;
declare global {
    interface Document {
        startViewTransition?: (callback: () => Promise<void> | void) => {
            finished: Promise<void>;
            updateCallbackDone: Promise<void>;
            ready: Promise<void>;
        };
    }
}
