import { Node } from "../index.js";
export default function swap(components?: Record<string, string | ((props: Record<string, any>, ...children: any[]) => any)>): (doc: Node) => Node;
