/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable @typescript-eslint/ban-types */
import {Component} from 'react';

declare module 'enzyme' {
    interface ShallowWrapper<P = {}, S = {}, C = Component> {
        rerender(props: P): void;
    }
}

declare global {
    namespace jest {
        interface Matchers<R, T> {
            toHaveProps(props: Record<string, unknown>): R;
        }
    }
}
