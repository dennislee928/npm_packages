import { Node } from "../index.js";
export interface ScopeOptions {
    hash?: string;
    attribute?: string;
}
export default function scope(opts?: ScopeOptions): (doc: Node) => Promise<Node>;
