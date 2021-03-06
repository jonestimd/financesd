import {createDiv} from 'test/htmlUtils';
import {testAction} from 'test/mobxUtils';
import SelectionModel from './SelectionModel';

describe('SelectionModel', () => {
    const rows = 1245;
    const columns = 4;
    const rowSelector = '.row';
    const rowHeight = 10;
    const options = {rows, columns, rowSelector};
    const container = createDiv();

    function setupModel(rowIndex: number, headerSelector?: string) {
        const model = new SelectionModel({...options, headerSelector});
        model.cell.row = rowIndex;
        model.containerRef(container);
        return model;
    }

    describe('constructor', () => {
        it('initializes defaults', () => {
            const model = new SelectionModel({rows: 19, columns: 3, headerSelector: 'thead'});

            expect(model).toEqual(expect.objectContaining({
                cell: {row: 0, column: 0},
                rows: 19,
                columns: 3,
                container: null,
                editCell: undefined,
                rowSelector: 'tbody tr.MuiTableRow-root',
                headerSelector: 'thead',
                rowOffset: 0,
            }));
        });
    });
    describe('containerRef', () => {
        it('focuses the container', () => {
            const container = createDiv();
            jest.spyOn(container, 'focus');
            const model = new SelectionModel({rows: 19, columns: 3, headerSelector: 'thead'});

            model.containerRef(container);

            expect(container.focus).toBeCalled();
            expect(model['container']).toBe(container);
        });
    });
    describe('setCell', () => {
        it('sets the selected cell', () => {
            const model = new SelectionModel({rows: 19, columns: 3});

            testAction(() => model.cell, () => model.setCell(5, 1));

            expect(model.cell).toEqual({row: 5, column: 1});
        });
    });
    describe('setEditCell', () => {
        it('sets the selected cell and edit cell', () => {
            const model = new SelectionModel({rows: 19, columns: 3});

            testAction(() => model.editCell, () => model.setEditCell({row: 5, column: 1}));

            expect(model.cell).toEqual({row: 5, column: 1});
            expect(model.editCell).toEqual({row: 5, column: 1});
        });
        it('clears the edit cell', () => {
            const model = new SelectionModel({rows: 19, columns: 3});
            model.setEditCell({row: 5, column: 1});

            testAction(() => model.editCell, () => model.setEditCell());

            expect(model.cell).toEqual({row: 5, column: 1});
            expect(model.editCell).toBeUndefined();
        });
    });
    describe('stopEditing', () => {
        it('clears edit cell', () => {
            const model = new SelectionModel({rows: 19, columns: 3});
            model.setEditCell({row: 5, column: 1});

            testAction(() => model.editCell, () => model.stopEditing());

            expect(model.editCell).toBeUndefined();
        });
        it('returns focus to container', () => {
            const model = new SelectionModel({rows: 19, columns: 3});
            model.container = createDiv();
            jest.spyOn(model.container, 'focus');

            model.stopEditing();

            expect(model.container.focus).toBeCalledTimes(1);
        });
    });
    describe('onKeyDown', () => {
        const headerSelector = '.header';
        const row = createDiv();
        const header = createDiv();
        const children: Record<string, Element> = {
            [rowSelector]: row,
            [headerSelector]: header,
        };
        const containerHeight = 200;
        const headerHeight = 15;
        const pageSize = containerHeight / rowHeight;

        const mockKeyEvent = (key: string, ctrlKey: boolean, shiftKey = false) => ({
            currentTarget: container, key, ctrlKey, shiftKey,
            stopPropagation: jest.fn(),
            preventDefault: jest.fn(),
        } as unknown as React.KeyboardEvent<HTMLElement>);

        beforeEach(() => {
            jest.spyOn(container, 'getBoundingClientRect').mockReturnValue({height: containerHeight} as DOMRect);
            jest.spyOn(container, 'clientHeight', 'get').mockReturnValue(containerHeight);
            jest.spyOn(container, 'querySelector').mockImplementation((selector: keyof typeof children) => children[selector]);
            jest.spyOn(row, 'getBoundingClientRect').mockReturnValue({height: rowHeight} as DOMRect);
            jest.spyOn(header, 'getBoundingClientRect').mockReturnValue({height: headerHeight} as DOMRect);
        });
        it('stops event propagation', () => {
            const model = setupModel(0);
            const event = mockKeyEvent('', false);

            model.onKeyDown(event);

            expect(event.stopPropagation).toBeCalledTimes(1);
            expect(event.preventDefault).not.toBeCalled();
            expect(container.getBoundingClientRect).not.toBeCalled();
            expect(container.querySelector).not.toBeCalled();
            expect(row.getBoundingClientRect).not.toBeCalled();
        });
        it('scrolls up to selected row', () => {
            const initialRow = pageSize * 2;
            jest.spyOn(container, 'scrollTop', 'get').mockReturnValue((initialRow - pageSize / 2) * rowHeight);

            setupModel(initialRow, headerSelector).onKeyDown(mockKeyEvent('PageUp', false));

            expect(container.scrollTo).toBeCalledWith({top: (initialRow - pageSize) * rowHeight});
        });
        it('scrolls down to selected row', () => {
            const initialRow = pageSize * 2;
            jest.spyOn(container, 'scrollTop', 'get').mockReturnValue((initialRow - pageSize / 2) * rowHeight);

            setupModel(initialRow, headerSelector).onKeyDown(mockKeyEvent('PageDown', false));

            const finalRow = initialRow + pageSize;
            expect(container.scrollTo).toBeCalledWith({top: finalRow * rowHeight + headerHeight - containerHeight + rowHeight});
        });
        const rowTests = [
            {name: 'selects previous row on up arrow', key: 'ArrowUp', ctrl: false, initialRow: 9, expectedRow: 8},
            {name: 'remains at top on up arrow', key: 'ArrowUp', ctrl: false, initialRow: 0, expectedRow: 0},
            {name: 'selects next row on down arrow', key: 'ArrowDown', ctrl: false, initialRow: 9, expectedRow: 10},
            {name: 'remains at bottom on down arrow', key: 'ArrowDown', ctrl: false, initialRow: rows - 1, expectedRow: rows - 1},
            {name: 'selects previous page on page up', key: 'PageUp', ctrl: false, initialRow: pageSize * 2, expectedRow: pageSize},
            {name: 'selects first row on page up', key: 'PageUp', ctrl: false, initialRow: pageSize * .5, expectedRow: 0},
            {name: 'selects next page on page down', key: 'PageDown', ctrl: false, initialRow: pageSize * 2, expectedRow: pageSize * 3},
            {name: 'selects last row on page down', key: 'PageDown', ctrl: false, initialRow: rows - pageSize * .5, expectedRow: rows - 1},
            {name: 'selects first row on ctrl+home', key: 'Home', ctrl: true, initialRow: rows - 1, expectedRow: 0},
            {name: 'remains on current row on home', key: 'Home', ctrl: false, initialRow: 9, expectedRow: 9},
            {name: 'selects last row on ctrl+end', key: 'End', ctrl: true, initialRow: 0, expectedRow: rows - 1},
            {name: 'remains on current row on end', key: 'End', ctrl: false, initialRow: 9, expectedRow: 9},
        ];
        rowTests.forEach(({name, initialRow, key, ctrl, expectedRow}) =>
            it(name, () => {
                const event = mockKeyEvent(key, ctrl);
                const model = setupModel(initialRow);

                model.onKeyDown(event);

                expect(model.cell.row).toEqual(expectedRow);
                expect(event.preventDefault).toBeCalled();
            })
        );
        const columnTests = [
            {name: 'selects previous column on left arrow', key: 'ArrowLeft', ctrl: false,
                initial: {column: 2}, expected: {column: 1}},
            {name: 'wraps to last column on left arrow', key: 'ArrowLeft', ctrl: false,
                initial: {column: 0}, expected: {column: columns - 1}},
            {name: 'selects next column on right arrow', key: 'ArrowRight', ctrl: false,
                initial: {column: 0}, expected: {column: 1}},
            {name: 'wraps to first column on right arrow', key: 'ArrowRight', ctrl: false,
                initial: {column: columns - 1}, expected: {column: 0}},
            {name: 'selects first column on home', key: 'Home', ctrl: false,
                initial: {column: columns - 1}, expected: {column: 0}},
            {name: 'remains at current column on ctrl+home', key: 'Home', ctrl: true,
                initial: {column: 1}, expected: {column: 1}},
            {name: 'selects last column on end', key: 'End', ctrl: false,
                initial: {column: 0}, expected: {column: columns - 1}},
            {name: 'remains at current column on ctrl+end', key: 'End', ctrl: true,
                initial: {column: 1}, expected: {column: 1}},
            {name: 'selects next column on tab', key: 'Tab', ctrl: false,
                initial: {column: 0}, expected: {column: 1}},
            {name: 'wraps to next row on tab', key: 'Tab', ctrl: false,
                initial: {column: columns - 1, row: 0}, expected: {column: 0, row: 1}},
            {name: 'remains on last cell row on tab', key: 'Tab', ctrl: false,
                initial: {column: columns - 1, row: rows - 1}, expected: {column: columns - 1, row: rows - 1}},
            {name: 'selects previous column on shift tab', key: 'Tab', ctrl: false, shift: true,
                initial: {column: 2}, expected: {column: 1}},
            {name: 'wraps to previous row on shift tab', key: 'Tab', ctrl: false, shift: true,
                initial: {column: 0, row: 1}, expected: {column: columns - 1, row: 0}},
            {name: 'remains on first cell on shift tab', key: 'Tab', ctrl: false, shift: true,
                initial: {column: 0, row: 0}, expected: {column: 0, row: 0}},
        ];
        columnTests.forEach(({name, initial, key, ctrl, shift, expected}) =>
            it(name, () => {
                const model = setupModel(initial.row ?? 0);
                model.cell.column = initial.column;

                model.onKeyDown(mockKeyEvent(key, ctrl, shift));

                expect(model.cell.column).toEqual(expected.column);
                if (expected.row) expect(model.cell.row).toEqual(expected.row);
            })
        );
    });
    describe('onMouseDown', () => {
        const mockMouseEvent = (target: HTMLElement | null, clientY = -1) => ({
            target,
            currentTarget: container,
            clientY,
            stopPropagation: jest.fn(),
            preventDefault: jest.fn(),
        } as unknown as React.MouseEvent<HTMLElement>);

        describe('table', () => {
            const cells = new Array(columns).fill(null).map(() => document.createElement('td'));
            const tr = document.createElement('tr');
            const rowIndex = Math.floor(Math.random() * 100);

            beforeAll(() => {
                cells.forEach((cell, index) => {
                    Object.defineProperty(cell, 'cellIndex', {value: index});
                    tr.appendChild(cell);
                });
                Object.defineProperty(tr, 'sectionRowIndex', {value: rowIndex});
            });
            it('selects row and column', () => {
                const column = columns - 2;
                cells.forEach((cell) => cell.colSpan = 1);
                const model = setupModel(0);

                model.onMouseDown(mockMouseEvent(cells[column]));

                expect(model.cell.row).toEqual(rowIndex);
                expect(model.cell.column).toEqual(column);
            });
            it('uses colSpan to calculate column index', () => {
                Object.defineProperty(cells[1], 'colSpan', {value: 2});
                const model = setupModel(0);

                model.onMouseDown(mockMouseEvent(cells[2]));

                expect(model.cell.row).toEqual(rowIndex);
                expect(model.cell.column).toEqual(3);
            });
        });
        describe('non-table', () => {
            const rows = new Array(50).fill(null).map(() => container.appendChild(document.createElement('div')));
            const boundingRect = (index: number) => ({bottom: (index + 1) * rowHeight} as DOMRect);

            beforeAll(() => rows.forEach((row) => row.setAttribute('class', rowSelector.substring(1))));
            beforeEach(() => {
                jest.spyOn(container, 'querySelectorAll').mockReturnValue(rows as unknown as NodeListOf<Element>);
                rows.forEach((row, index) => jest.spyOn(row, 'getBoundingClientRect').mockReturnValue(boundingRect(index)));
            });
            it('selects row at click', () => {
                const model = setupModel(0);

                model.onMouseDown(mockMouseEvent(null, rowHeight * 4.1));

                expect(model.cell.row).toEqual(4);
            });
            it('ignores click past end row', () => {
                const initialRow = 5;
                const model = setupModel(initialRow);

                model.onMouseDown(mockMouseEvent(null, rowHeight * (rows.length + 1)));

                expect(model.cell.row).toEqual(initialRow);
            });
        });
    });
});
