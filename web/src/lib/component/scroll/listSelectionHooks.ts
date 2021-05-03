import React from 'react';
import {IScrollOptions, useScroll} from './scrollHooks';

const defaultRowHeight = 24;

function getHeight(container: HTMLElement, selector: string) {
    return container.querySelector(selector)?.getBoundingClientRect().height;
}

function getPageSize(container: HTMLElement, rowSelector: string) {
    const {height} = container.getBoundingClientRect();
    const rowHeight = getHeight(container, rowSelector) ?? defaultRowHeight;
    return Math.floor(height / rowHeight);
}

function ensureVisible(container: HTMLElement, row: number, rowSelector: string, headerSelector?: string) {
    const {clientHeight, scrollTop} = container;
    const headerHeight = headerSelector ? getHeight(container, headerSelector) ?? defaultRowHeight : 0;
    const rowHeight = getHeight(container, rowSelector) ?? defaultRowHeight;
    const rowTop = row * rowHeight;
    if (rowTop < scrollTop) container.scrollTo({top: rowTop});
    else if (rowTop + rowHeight > scrollTop + clientHeight) container.scrollTo({top: rowTop + headerHeight - clientHeight + rowHeight});
}

export interface ICell {
    row: number;
    column: number;
}

export interface ISelectionOptions extends IScrollOptions {
    /** total number of rows */
    rows: number;
    /** CSS selector for row elements */
    rowSelector?: string;
    /** CSS selector for header */
    headerSelector?: string;
    /** Select new column when row changes */
    getColumn?: (cell: ICell) => number;
}

const defaultOptions = {
    rowSelector: 'tbody tr.MuiTableRow-root',
    headerSelector: 'thead',
    getColumn: (cell: ICell) => cell.column,
};

function clamp(value: number, max: number) {
    if (value < 0) return 0;
    return (value >= max) ? max-1 : value;
}

function isContained(element: HTMLElement, parentSelector: string) {
    let parent: HTMLElement | null = element;
    while (parent) {
        if (parent.matches(parentSelector)) return true;
        parent = parent.parentElement;
    }
    return false;
}

export function useSelection<T extends HTMLElement>({rows, ...options}: ISelectionOptions) {
    const scroll = useScroll<T>({defaultRowHeight, ...options});
    const {rowSelector, headerSelector, getColumn} = {...defaultOptions, ...options};
    const [cell, setCell] = React.useState<ICell>({row: scroll.startRow, column: 0});
    const move = (target: HTMLElement, nRows: number) => setCell(({row: r, column: c}) => {
        const row = clamp(r + nRows, rows);
        const column = getColumn({row, column: c});
        ensureVisible(target, row, rowSelector, headerSelector);
        return {row, column};
    });
    return {
        ...cell,
        scroll,
        setCell(row: number, column: number) {
            setCell({row, column});
        },
        onKeyDown(event: React.KeyboardEvent<HTMLElement>) {
            const {currentTarget, key, ctrlKey} = event;
            const pageSize = getPageSize(currentTarget, rowSelector);
            switch (key) {
                case 'ArrowUp':
                    event.stopPropagation();
                    event.preventDefault();
                    move(currentTarget, -1);
                    break;
                case 'ArrowDown':
                    event.stopPropagation();
                    event.preventDefault();
                    move(currentTarget, 1);
                    break;
                case 'PageUp':
                    event.stopPropagation();
                    event.preventDefault();
                    move(currentTarget, -pageSize);
                    break;
                case 'PageDown':
                    event.stopPropagation();
                    event.preventDefault();
                    move(currentTarget, pageSize);
                    break;
                case 'Home':
                    if (ctrlKey) {
                        event.stopPropagation();
                        event.preventDefault();
                        currentTarget.scrollTo({top: 0});
                        setCell({row: 0, column: getColumn({...cell, row: 0})});
                    }
                    break;
                case 'End':
                    if (ctrlKey && rows > 0) {
                        event.stopPropagation();
                        event.preventDefault();
                        ensureVisible(currentTarget, rows - 1, rowSelector, headerSelector);
                        setCell({row: rows - 1, column: getColumn({...cell, row: rows - 1})});
                    }
                    break;
            }
        },
        onMouseDown(event: React.MouseEvent<HTMLElement>) {
            if (event.button === 0 && isContained(event.target as HTMLElement, rowSelector)) {
                const container = event.currentTarget;
                const index = Array.from(container.querySelectorAll(rowSelector)).findIndex((r) => r.getBoundingClientRect().bottom >= event.clientY);
                if (index >= 0 && cell.row !== index + scroll.startRow) {
                    const row = index + scroll.startRow;
                    setCell({row, column: getColumn({row, column: cell.column})});
                }
            }
        },
    };
}
