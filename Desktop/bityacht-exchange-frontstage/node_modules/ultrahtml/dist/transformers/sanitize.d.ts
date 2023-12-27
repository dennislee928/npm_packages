import { Node } from "../index.js";
export interface SanitizeOptions {
    /** An Array of strings indicating elements that the sanitizer should not remove. All elements not in the array will be dropped. */
    allowElements?: string[];
    /** An Array of strings indicating elements that the sanitizer should remove, but keeping their child elements. */
    blockElements?: string[];
    /** An Array of strings indicating elements (including nested elements) that the sanitizer should remove. */
    dropElements?: string[];
    /** An Object where each key is the attribute name and the value is an Array of allowed tag names. Matching attributes will not be removed. All attributes that are not in the array will be dropped. */
    allowAttributes?: Record<string, string[]>;
    /** An Object where each key is the attribute name and the value is an Array of dropped tag names. Matching attributes will be removed. */
    dropAttributes?: Record<string, string[]>;
    /** A Boolean value set to false (default) to remove components and their children. If set to true, components will be subject to built-in and custom configuration checks (and will be retained or dropped based on those checks). */
    allowComponents?: boolean;
    /** A Boolean value set to false (default) to remove custom elements and their children. If set to true, custom elements will be subject to built-in and custom configuration checks (and will be retained or dropped based on those checks). */
    allowCustomElements?: boolean;
    /** A Boolean value set to false (default) to remove HTML comments. Set to true in order to keep comments. */
    allowComments?: boolean;
}
export default function sanitize(opts?: SanitizeOptions): (doc: Node) => Node;
