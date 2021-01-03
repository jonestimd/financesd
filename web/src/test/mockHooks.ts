import * as scrollHooks from '../lib/component/scroll/scrollHooks';
import * as selectionHooks from '../lib/component/scroll/selectionHooks';

type ScrollHook = Pick<ReturnType<typeof scrollHooks['useScroll']>, 'startRow' | 'rowHeight' | 'headerHeight'>;

export function mockScrollHook(overrides: Partial<ScrollHook> = {}) {
    const scroll = {
        startRow: 0,
        rowHeight: 0,
        headerHeight: 0,
        listRef: {current: null} as React.MutableRefObject<HTMLElement>,
        endRow: jest.fn<number, [number]>(),
        onScroll: jest.fn<void, [React.UIEvent<HTMLElement>]>(),
        ...overrides,
    };
    jest.spyOn(scrollHooks, 'useScroll').mockReturnValue(scroll);
    return scroll;
}

export function mockSelectionHook() {
    const selection = {
        row: 0,
        column: 0,
        onKeyDown: jest.fn(),
        onMouseDown: jest.fn(),
    };
    jest.spyOn(selectionHooks, 'useSelection').mockReturnValue(selection);
    return selection;
}

export function mockHooks(scrollOverrides: Partial<ScrollHook> = {}) {
    const scroll = mockScrollHook(scrollOverrides);
    const selection = mockSelectionHook();
    return {scroll, selection};
}
