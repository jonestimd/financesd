import React from 'react';
import * as scrollHooks from '../lib/component/scroll/scrollHooks';
import * as selectionHooks from '../lib/component/scroll/selectionHooks';

type ScrollHook = Pick<ReturnType<typeof scrollHooks['useScroll']>, 'startRow' | 'rowHeight' | 'headerHeight'>;

export function mockScrollHook<T extends HTMLElement>(overrides: Partial<ScrollHook> = {}) {
    const scroll = {
        startRow: 0,
        rowHeight: 0,
        headerHeight: 0,
        listRef: {current: null} as unknown as React.MutableRefObject<T>,
        endRow: jest.fn<number, [number]>(),
        onScroll: jest.fn<void, [React.UIEvent<HTMLElement>]>(),
        ...overrides,
    };
    jest.spyOn(scrollHooks, 'useScroll').mockReturnValue(scroll);
    return scroll;
}

export function mockSelectionHook<T extends HTMLElement>(row = 0, column = 0, overrides: Partial<ScrollHook> = {}) {
    const selection = {
        row,
        column,
        setCell: jest.fn(),
        onKeyDown: jest.fn(),
        onMouseDown: jest.fn(),
        scroll: mockScrollHook<T>(overrides),
    };
    jest.spyOn(selectionHooks, 'useSelection').mockReturnValue(selection);
    return selection;
}

const effects: string[] = [];
const disposers: (() => void)[] = [];

export function mockUseEffect() {
    jest.spyOn(React, 'useEffect').mockImplementation((effect) => {
        if (!effects.includes(effect.toString())) {
            const disposer = effect();
            if (disposer) disposers.push(disposer);
            effects.push(effect.toString());
        }
    });
}

afterEach(() => {
    disposers.forEach((disposer) => disposer());
    disposers.splice(0, disposers.length);
    effects.splice(0, effects.length);
});
