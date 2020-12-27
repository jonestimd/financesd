import React from 'react';

interface IScrollOptions {
    defaultRowHeight?: number;
    prototypeSelector?: string;
    headerSelector?: string;
    defaultHeaderHeight?: number;
}

const getHeight = (list: HTMLElement, itemSelector: string, defaultHeight: number) => {
    return list?.querySelector(itemSelector)?.getBoundingClientRect().height ?? defaultHeight;
};

const nextEven = (num: number) => num % 2 === 0 ? num : num + 1;

export function useScroll<T extends HTMLElement>(options: IScrollOptions) {
    const {defaultRowHeight = 24, prototypeSelector = '*', headerSelector, defaultHeaderHeight} = options;
    const [startRow, setStartRow] = React.useState(0);
    const listRef = React.useRef<T>(null);
    const rowHeight = getHeight(listRef.current, prototypeSelector, defaultRowHeight);
    const headerHeight = headerSelector ? getHeight(listRef.current, headerSelector, defaultHeaderHeight) : 0;
    return {
        startRow, listRef, rowHeight, headerHeight,
        onScroll: React.useCallback(({currentTarget}: React.UIEvent<HTMLElement>) => {
            const {clientHeight, scrollTop} = currentTarget;
            const overscan = Math.ceil(clientHeight / rowHeight);
            setStartRow(nextEven(Math.max(0, Math.floor(scrollTop / rowHeight) - overscan)));
        }, [rowHeight]),
    };
}
