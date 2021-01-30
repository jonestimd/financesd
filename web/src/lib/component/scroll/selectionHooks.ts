import React from 'react';

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
    return row;
}

function isTableCell(elem: EventTarget): elem is HTMLTableCellElement {
    return elem instanceof Element && elem.tagName.toLocaleLowerCase() === 'td';
}

export interface ISelectionOptions {
    /** initial selection */
    initialRow?: number;
    /** total number of rows */
    rows: number;
    /** total number of columns */
    columns?: number;
    /** index of 1st row element in DOM */
    rowOffset?: number;
    /** CSS selector for row elements */
    rowSelector: string;
    /** CSS selector for header */
    headerSelector?: string;
}

export function useSelection({initialRow = 0, rows, columns = 1, rowOffset = 0, rowSelector, headerSelector}: ISelectionOptions) {
    const [row, setRow] = React.useState(initialRow);
    const [column, setColumn] = React.useState(0);
    return {
        row, column,
        onKeyDown(event: React.KeyboardEvent<HTMLElement>) {
            const {currentTarget, key, ctrlKey} = event;
            const pageSize = getPageSize(currentTarget, rowSelector);
            event.stopPropagation();
            // TODO always make new selection visible
            switch (key) {
                case 'ArrowRight': setColumn((c) => c === columns - 1 ? 0 : c + 1); break;
                case 'ArrowLeft': setColumn((c) => c === 0 ? columns - 1 : c - 1); break;
                case 'ArrowUp':
                    event.preventDefault();
                    setRow((r) => ensureVisible(currentTarget, Math.max(0, r - 1), rowSelector, headerSelector));
                    break;
                case 'ArrowDown':
                    event.preventDefault();
                    setRow((r) => ensureVisible(currentTarget, Math.min(r + 1, rows - 1), rowSelector, headerSelector));
                    break;
                case 'PageUp':
                    event.preventDefault();
                    setRow((r) => ensureVisible(currentTarget, Math.max(0, r - pageSize), rowSelector, headerSelector));
                    break;
                case 'PageDown':
                    event.preventDefault();
                    if (rows > 0) setRow((r) => ensureVisible(currentTarget, Math.min(r + pageSize, rows - 1), rowSelector, headerSelector));
                    break;
                case 'Home':
                    event.preventDefault();
                    if (ctrlKey) {
                        currentTarget.scrollTo({top: 0});
                        setRow(0);
                    }
                    else setColumn(0);
                    break;
                case 'End':
                    event.preventDefault();
                    if (ctrlKey && rows > 0) {
                        ensureVisible(currentTarget, rows - 1, rowSelector, headerSelector);
                        setRow(rows - 1);
                    }
                    else setColumn(columns - 1);
                    break;
            }
        },
        onMouseDown(event: React.MouseEvent<HTMLElement>) {
            if (isTableCell(event.target)) {
                const tr = event.target.parentElement as HTMLTableRowElement;
                const cellIndex = Array.from(tr.querySelectorAll('td'))
                    .slice(0, event.target.cellIndex)
                    .reduce((count, cell) => count + (cell.colSpan || 1), 0);
                setRow(tr.sectionRowIndex + rowOffset);
                setColumn(cellIndex);
            }
            else {
                const container = event.currentTarget;
                const index = Array.from(container.querySelectorAll(rowSelector)).findIndex((r) => r.getBoundingClientRect().bottom >= event.clientY);
                if (index >= 0) setRow(index + rowOffset);
            }
        },
    };
}
