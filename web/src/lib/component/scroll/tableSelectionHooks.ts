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

function isTableCell(elem: EventTarget): elem is HTMLTableCellElement {
    return elem instanceof Element && elem.tagName.toLocaleLowerCase() === 'td';
}

export interface ICell {
    row: number;
    column: number;
}

export interface ISelectionOptions extends IScrollOptions {
    /** total number of rows */
    rows: number;
    /** total number of columns */
    columns?: number;
    /** CSS selector for row elements */
    rowSelector?: string;
    /** CSS selector for header */
    headerSelector?: string;
}

const defaultOptions = {
    rowSelector: 'tbody tr.MuiTableRow-root',
    headerSelector: 'thead',
};

function clamp(value: number, max: number) {
    if (value < 0) return 0;
    return (value >= max) ? max-1 : value;
}

function wrap(value: number, max: number) {
    if (value < 0) return max - 1;
    return (value >= max) ? 0 : value;
}

function isContained(element: HTMLElement, parentSelector: string) {
    let parent: HTMLElement | null = element;
    while (parent) {
        if (parent.matches(parentSelector)) return true;
        parent = parent.parentElement;
    }
    return false;
}

export function useSelection<T extends HTMLElement>({rows, columns = 1, ...options}: ISelectionOptions) {
    const scroll = useScroll<T>({defaultRowHeight, ...options});
    const {rowSelector, headerSelector} = {...defaultOptions, ...options};
    const [cell, setCell] = React.useState<ICell>({row: scroll.startRow, column: 0});
    const move = (target: HTMLElement, nRows: number, nCols: number) => setCell(({row: r, column: c}) => {
        const cell = {row: clamp(r + nRows, rows), column: wrap(c + nCols, columns)};
        ensureVisible(target, cell.row, rowSelector, headerSelector);
        return cell;
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
            if (key !== 'Escape') event.stopPropagation();
            switch (key) {
                case 'ArrowRight': move(currentTarget, 0, 1); break;
                case 'ArrowLeft': move(currentTarget, 0, -1); break;
                case 'Tab':
                    if (event.shiftKey) {
                        if (cell.row > 0 || cell.column > 0) {
                            event.preventDefault();
                            move(currentTarget, cell.column === 0 ? -1 : 0, -1);
                        }
                    }
                    else if (cell.row < rows-1 || cell.column < columns-1) {
                        event.preventDefault();
                        move(currentTarget, cell.column === columns - 1 ? 1 : 0, 1);
                    }
                    break;
                case 'ArrowUp':
                    event.preventDefault();
                    move(currentTarget, -1, 0);
                    break;
                case 'ArrowDown':
                    event.preventDefault();
                    move(currentTarget, 1, 0);
                    break;
                case 'PageUp':
                    event.preventDefault();
                    move(currentTarget, -pageSize, 0);
                    break;
                case 'PageDown':
                    event.preventDefault();
                    move(currentTarget, pageSize, 0);
                    break;
                case 'Home':
                    event.preventDefault();
                    if (ctrlKey) {
                        currentTarget.scrollTo({top: 0});
                        setCell({...cell, row: 0});
                    }
                    else setCell({...cell, column: 0});
                    break;
                case 'End':
                    event.preventDefault();
                    if (ctrlKey && rows > 0) {
                        ensureVisible(currentTarget, rows - 1, rowSelector, headerSelector);
                        setCell({...cell, row: rows - 1});
                    }
                    else setCell({...cell, column: columns - 1});
                    break;
            }
        },
        onMouseDown(event: React.MouseEvent<HTMLElement>) {
            if (event.button === 0 && isContained(event.target as HTMLElement, rowSelector)) {
                if (isTableCell(event.target)) {
                    const tr = event.target.parentElement as HTMLTableRowElement;
                    const cellIndex = Array.from(tr.querySelectorAll('td'))
                        .slice(0, event.target.cellIndex)
                        .reduce((count, cell) => count + (cell.colSpan || 1), 0);
                    setCell({row: tr.sectionRowIndex + scroll.startRow, column: cellIndex});
                }
                else {
                    const container = event.currentTarget;
                    const index = Array.from(container.querySelectorAll(rowSelector)).findIndex((r) => r.getBoundingClientRect().bottom >= event.clientY);
                    if (index >= 0) setCell({...cell, row: index + scroll.startRow});
                }
            }
        },
    };
}
