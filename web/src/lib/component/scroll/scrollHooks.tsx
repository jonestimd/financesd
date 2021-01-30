import React from 'react';

export interface IScrollOptions {
    overscan?: number;
    defaultRowHeight?: number;
    prototypeSelector?: string;
    headerSelector?: string;
    defaultHeaderHeight?: number;
}

const getHeight = (list: HTMLElement | null, itemSelector: string, defaultHeight: number) => {
    return list?.querySelector(itemSelector)?.getBoundingClientRect().height ?? defaultHeight;
};

const nextEven = (num: number) => num % 2 === 0 ? num : num + 1;

export const defaultOverscan = 0.25;

export function useScroll<T extends HTMLElement>(options: IScrollOptions) {
    const {overscan = defaultOverscan, defaultRowHeight = 24, prototypeSelector = '*', headerSelector, defaultHeaderHeight} = options;
    const [startRow, setStartRow] = React.useState(0);
    const listRef = React.useRef<T | null>(null);
    const rowHeight = getHeight(listRef.current, prototypeSelector, defaultRowHeight);
    const headerHeight = headerSelector ? getHeight(listRef.current, headerSelector, defaultHeaderHeight ?? defaultRowHeight) : 0;
    return {
        startRow, listRef, rowHeight, headerHeight,
        endRow: (scrollHeight: number) => startRow + Math.ceil(scrollHeight / rowHeight * (1 + 2 * overscan)),
        onScroll: React.useCallback(({currentTarget}: React.UIEvent<HTMLElement>) => {
            const {clientHeight, scrollTop} = currentTarget;
            const visibleRows = clientHeight / rowHeight;
            const leadingRows = Math.floor(visibleRows * overscan);
            setStartRow(nextEven(Math.max(0, Math.floor(scrollTop / rowHeight) - leadingRows)));
        }, [overscan, rowHeight]),
    };
}
