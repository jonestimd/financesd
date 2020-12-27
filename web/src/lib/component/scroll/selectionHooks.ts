import React from 'react';

function getHeight(container: HTMLElement, selector: string) {
    return container.querySelector(selector).getBoundingClientRect().height;
}

function isPreviousRowVisible(container: HTMLElement, rowSelector: string, headerSelector?: string) {
    const headerHeight = headerSelector ? getHeight(container, headerSelector) : 0;
    const rect = container.querySelector(`${rowSelector}.selected`).getBoundingClientRect();
    return rect.top - rect.height >= container.getBoundingClientRect().top + headerHeight;
}

function isNextRowVisible(container: HTMLElement, rowSelector: string) {
    const {top, height} = container.getBoundingClientRect();
    const rect = container.querySelector(`${rowSelector}.selected`).getBoundingClientRect();
    return rect.top + rect.height * 2 <= top + height;
}

function getPageSize(container: HTMLElement, rowSelector: string) {
    const {height} = container.getBoundingClientRect();
    const rowHeight = getHeight(container, `${rowSelector}.selected`); // TODO why selected row?
    return Math.floor(height / rowHeight);
}

function scrollTo(container: HTMLElement, row: number, rowOffset: number, rowSelector: string, headerSelector?: string) {
    const headerHeight = headerSelector ? getHeight(container, headerSelector) : 0;
    const tr = container.querySelector<HTMLTableRowElement>(`${rowSelector}:nth-child(${row + 1 - rowOffset})`);
    if (tr) container.scrollTo({top: tr.offsetTop - headerHeight});
    else {
        const rowHeight = getHeight(container, `${rowSelector}.selected`);
        container.scrollTo({top: rowHeight * row});
    }
    return row;
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
            switch (key) {
                case 'ArrowRight': setColumn(c => c === columns - 1 ? 0 : c + 1); break;
                case 'ArrowLeft': setColumn(c => c === 0 ? columns - 1 : c - 1); break;
                case 'ArrowUp':
                    if (isPreviousRowVisible(currentTarget, rowSelector, headerSelector)) event.preventDefault();
                    setRow(r => Math.max(0, r - 1));
                    break;
                case 'ArrowDown':
                    if (isNextRowVisible(currentTarget, rowSelector)) event.preventDefault();
                    setRow(r => Math.min(r + 1, rows - 1));
                    break;
                case 'PageUp':
                    event.preventDefault();
                    setRow(r => scrollTo(currentTarget, Math.max(0, r - pageSize), rowOffset, rowSelector, headerSelector));
                    break;
                case 'PageDown':
                    event.preventDefault();
                    if (rows > 0) setRow(r => scrollTo(currentTarget, Math.min(r + pageSize, rows - 1), rowOffset, rowSelector, headerSelector));
                    break;
                case 'Home':
                    event.preventDefault();
                    if (ctrlKey) {
                        scrollTo(currentTarget, 0, 0, rowSelector, headerSelector);
                        setRow(0);
                    }
                    else setColumn(0);
                    break;
                case 'End':
                    event.preventDefault();
                    if (ctrlKey && rows > 0) {
                        scrollTo(currentTarget, rows - 1, 0, rowSelector, headerSelector);
                        setRow(rows - 1);
                    }
                    else setColumn(columns - 1);
                    break;
            }
        },
        onMouseDown(event: React.MouseEvent<HTMLElement>) {
            if (event.target instanceof HTMLTableCellElement) {
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
