import type { Ref, VNode } from 'vue';
export declare const RouteProvider: import("vue").DefineComponent<{
    vnode: {
        type: () => VNode;
        required: true;
    };
    route: {
        type: () => RouteLocationNormalizedLoaded;
        required: true;
    };
    vnodeRef: () => Ref<any>;
    renderKey: StringConstructor;
    trackRootNodes: BooleanConstructor;
}, () => VNode<import("vue").RendererNode, import("vue").RendererElement, {
    [key: string]: any;
}>, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
    vnode: {
        type: () => VNode;
        required: true;
    };
    route: {
        type: () => RouteLocationNormalizedLoaded;
        required: true;
    };
    vnodeRef: () => Ref<any>;
    renderKey: StringConstructor;
    trackRootNodes: BooleanConstructor;
}>>, {
    trackRootNodes: boolean;
}, {}>;
