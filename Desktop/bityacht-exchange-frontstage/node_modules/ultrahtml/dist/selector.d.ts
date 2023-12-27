import type { Node } from './index.js';
export declare function specificity(selector: string): number;
export declare function matches(node: Node, selector: string): boolean;
export declare function querySelector(node: Node, selector: string): Node;
export declare function querySelectorAll(node: Node, selector: string): Node[];
export default querySelectorAll;
