import React from 'react';

const getHeight = (list: HTMLElement, itemSelector: string, defaultHeight: number) => {
    return list?.querySelector(itemSelector)?.clientHeight ?? defaultHeight;
};

export function useScroll(defaultRowHeight = 0, prototypeSelector = '*') {
    const [startRow, setStartRow] = React.useState(0);
    const listRef = React.useRef<HTMLDivElement>(null);
    const rowHeight = getHeight(listRef.current, prototypeSelector, defaultRowHeight);
    return {
        startRow, listRef, rowHeight,
        onScroll: React.useCallback(({currentTarget}: React.UIEvent<HTMLElement>) => {
            const {clientHeight, scrollTop} = currentTarget;
            const overscan = Math.ceil(clientHeight / rowHeight);
            setStartRow(Math.max(0, Math.floor(scrollTop / rowHeight) - overscan));
        }, [rowHeight])
    };
}
