import React from 'react';

function isPreviousRowVisible(container: HTMLElement) {
    const headerHeight = container.querySelector('thead').clientHeight;
    const rect = container.querySelector('tr.selected').getBoundingClientRect();
    return rect.top - rect.height >= container.getBoundingClientRect().top + headerHeight;
}

function isNextRowVisible(container: HTMLElement) {
    const {top, height} = container.getBoundingClientRect();
    const rect = container.querySelector('tr.selected').getBoundingClientRect();
    return rect.top + rect.height * 2 <= top + height;
}

function getPageSize(container: HTMLElement) {
    const {height} = container.getBoundingClientRect();
    const rowHeight = container.querySelector('tr.selected').getBoundingClientRect().height;
    return Math.floor(height / rowHeight);
}

function scrollTo(container: HTMLElement, row: number) {
    const headerHeight = container.querySelector('thead').clientHeight;
    const tr = container.querySelector(`tbody tr:nth-child(${row + 1})`) as HTMLTableRowElement;
    container.scrollTo({top: tr.offsetTop - headerHeight});
    return row;
}

export function useSelection(initialRow: number, rows: number, columns: number) {
    const [row, setRow] = React.useState(initialRow);
    const [column, setColumn] = React.useState(0);
    return {
        row, column,
        onKeyDown(event: React.KeyboardEvent<HTMLElement>) {
            const {currentTarget, key, ctrlKey} = event;
            const pageSize = getPageSize(currentTarget);
            switch (key) {
                case 'ArrowRight': setColumn(c => c === columns - 1 ? 0 : c + 1); break;
                case 'ArrowLeft': setColumn(c => c === 0 ? columns - 1 : c - 1); break;
                case 'ArrowUp':
                    if (isPreviousRowVisible(currentTarget)) event.preventDefault();
                    setRow(r => Math.max(0, r - 1));
                    break;
                case 'ArrowDown':
                    if (isNextRowVisible(currentTarget)) event.preventDefault();
                    setRow(r => Math.min(r + 1, rows - 1));
                    break;
                case 'PageUp':
                    event.preventDefault();
                    setRow(r => scrollTo(currentTarget, Math.max(0, r - pageSize)));
                    break;
                case 'PageDown':
                    event.preventDefault();
                    if (rows > 0) setRow(r => scrollTo(currentTarget, Math.min(r + pageSize, rows - 1)));
                    break;
                case 'Home':
                    event.preventDefault();
                    if (ctrlKey) {
                        scrollTo(currentTarget, 0);
                        setRow(0);
                    }
                    else setColumn(0);
                    break;
                case 'End':
                    event.preventDefault();
                    if (ctrlKey && rows > 0) {
                        scrollTo(currentTarget, rows - 1);
                        setRow(rows - 1);
                    }
                    else setColumn(columns - 1);
                    break;
            }
        },
        onMouseDown(event: React.MouseEvent<HTMLElement>) {
            if (event.target instanceof HTMLTableCellElement) {
                const tr = event.target.parentElement as HTMLTableRowElement;
                setRow(tr.sectionRowIndex);
                setColumn(event.target.cellIndex);
            }
        },
    };
}
