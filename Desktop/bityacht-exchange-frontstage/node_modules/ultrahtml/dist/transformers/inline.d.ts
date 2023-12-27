import { Node } from "../index.js";
export interface InlineOptions {
    /** Emit `style` attributes as objects rather than strings. */
    useObjectSyntax?: boolean;
}
export default function inline(opts?: InlineOptions): (doc: Node) => Node;
