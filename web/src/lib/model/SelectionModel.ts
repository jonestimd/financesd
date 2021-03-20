import {action, makeObservable, observable} from "mobx";

const defaultRowHeight = 24;

function getHeight(container: HTMLElement, selector: string) {
    return container.querySelector(selector)?.getBoundingClientRect().height;
}

function getPageSize(container: HTMLElement, rowSelector: string) {
    const {height} = container.getBoundingClientRect();
    const rowHeight = getHeight(container, rowSelector) ?? defaultRowHeight;
    return Math.floor(height / rowHeight);
}

function isTableCell(elem: EventTarget): elem is HTMLTableCellElement {
    return elem instanceof Element && elem.tagName.toLocaleLowerCase() === 'td';
}

interface IOptions {
    rows: number;
    columns: number;
    rowSelector?: string;
    headerSelector?: string;
}

export interface ICell {
    row: number;
    column: number;
}

const defaultOptions = {
    rowSelector: 'tbody tr.MuiTableRow-root',
};

function clamp(value: number, max: number) {
    if (value < 0) return 0;
    return (value >= max) ? max-1 : value;
}

function wrap(value: number, max: number) {
    if (value < 0) return max - 1;
    return (value >= max) ? 0 : value;
}

export default class SelectionModel {
    @observable cell: ICell = {row: 0, column: 0};
    @observable editCell?: ICell;
    @observable rowOffset = 0;
    @observable rows: number;
    readonly columns: number;
    readonly rowSelector: string;
    readonly headerSelector?: string;
    private container: HTMLElement | null = null;

    constructor(options: IOptions) {
        this.rows = options.rows;
        this.columns = options.columns;
        this.rowSelector = options.rowSelector ?? defaultOptions.rowSelector;
        this.headerSelector = options.headerSelector;
        makeObservable(this);
    }

    containerRef = (container: HTMLElement | null) => {
        this.container = container;
        this.container?.focus();
    };

    @action
    setRows(rows: number) {
        this.rows = rows;
    }

    @action
    setCell(row: number, column: number) {
        this.cell = {row, column};
    }

    @action
    setEditCell(cell?: ICell) {
        if (cell) this.cell = {...cell};
        this.editCell = cell;
    }

    private move(nRows: number, nCols: number) {
        this.cell.row = clamp(this.cell.row + nRows, this.rows);
        this.cell.column = wrap(this.cell.column + nCols, this.columns);
        this.ensureVisible();
    }

    private ensureVisible() {
        const container = this.container;
        if (container) {
            const rowElement = container.querySelectorAll(this.rowSelector)[this.cell.row - this.rowOffset];
            rowElement?.scrollIntoView();
        }
    }

    @action
    onKeyDown = (event: React.KeyboardEvent<HTMLElement>) => {
        const {currentTarget, key, ctrlKey} = event;
        if (key !== 'Escape') event.stopPropagation();
        switch (key) {
            case 'ArrowRight': this.move(0, 1); break;
            case 'ArrowLeft': this.move(0, -1); break;
            case 'Tab':
                if (event.shiftKey) {
                    if (this.cell.row > 0 || this.cell.column > 0) {
                        event.preventDefault();
                        this.move(this.cell.column === 0 ? -1 : 0, -1);
                    }
                }
                else if (this.cell.row < this.rows-1 || this.cell.column < this.columns-1) {
                    event.preventDefault();
                    this.move(this.cell.column === this.columns - 1 ? 1 : 0, 1);
                }
                break;
            case 'ArrowUp':
                event.preventDefault();
                this.move(-1, 0);
                break;
            case 'ArrowDown':
                event.preventDefault();
                this.move(1, 0);
                break;
            case 'PageUp':
                event.preventDefault();
                this.move(-getPageSize(currentTarget, this.rowSelector), 0);
                break;
            case 'PageDown':
                event.preventDefault();
                this.move(getPageSize(currentTarget, this.rowSelector), 0);
                break;
            case 'Home':
                event.preventDefault();
                if (ctrlKey) {
                    currentTarget.scrollTo({top: 0});
                    this.cell.row = 0;
                }
                else this.cell.column = 0;
                break;
            case 'End':
                event.preventDefault();
                if (ctrlKey && this.rows > 0) {
                    this.cell.row = this.rows - 1;
                    this.ensureVisible();
                }
                else this.cell.column = this.columns - 1;
                break;
        }
    };

    @action
    onMouseDown = (event: React.MouseEvent<HTMLElement>) => {
        if (isTableCell(event.target)) {
            const tr = event.target.parentElement as HTMLTableRowElement;
            const cellIndex = Array.from(tr.querySelectorAll('td'))
                .slice(0, event.target.cellIndex)
                .reduce((count, cell) => count + (cell.colSpan || 1), 0);
            this.cell.row = tr.sectionRowIndex + this.rowOffset;
            this.cell.column = cellIndex;
        }
        else {
            const container = event.currentTarget;
            const index = Array.from(container.querySelectorAll(this.rowSelector)).findIndex((r) => r.getBoundingClientRect().bottom >= event.clientY);
            if (index >= 0) this.cell.row = index + this.rowOffset;
        }
    };

    @action
    stopEditing = () => {
        this.setEditCell(undefined);
        this.container?.focus();
    };
}
