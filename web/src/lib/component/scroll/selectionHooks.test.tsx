import React from 'react';
import {shallow} from 'enzyme';
import {useSelection, ISelectionOptions} from './selectionHooks';

const hook: Partial<ReturnType<typeof useSelection>> = {};

function setup(options: ISelectionOptions) {
    const TestComponent: React.FC = () => {
        Object.assign(hook, useSelection(options));
        return null;
    };
    const component = shallow(<TestComponent />);
    return component;
}
const document = new Document();

function createDiv() {
    const div = document.createElement('div');
    div.scrollTo = jest.fn();
    return div;
}

describe('selectionHooks', () => {
    const rows = 1245;
    const columns = 4;
    const headerSelector = '.header';
    const rowSelector = '.row';
    const options = {rows, columns, rowSelector};
    const rowHeight = 10;

    describe('useSelection', () => {
        it('defaults to first row and column', () => {
            setup(options);

            expect(hook).toEqual(expect.objectContaining({row: 0, column: 0}));
        });
        it('defaults to initial row', () => {
            setup({...options, initialRow: 9});

            expect(hook).toEqual(expect.objectContaining({row: 9, column: 0}));
        });
        describe('onKeyDown', () => {
            const container = createDiv();
            const row = createDiv();
            const header = createDiv();
            const children: Record<string, Element> = {
                [rowSelector]: row,
                [headerSelector]: header,
            };
            const containerHeight = 200;
            const headerHeight = 15;
            const pageSize = containerHeight / rowHeight;

            const mockKeyEvent = (key: string, ctrlKey: boolean) => ({
                currentTarget: container, key, ctrlKey,
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
                setup(options);
                const event = mockKeyEvent('', false);

                hook.onKeyDown!(event);

                expect(event.stopPropagation).toBeCalledTimes(1);
                expect(event.preventDefault).not.toBeCalled();
                expect(container.getBoundingClientRect).toBeCalled();
                expect(container.querySelector).toBeCalledWith(rowSelector);
                expect(row.getBoundingClientRect).toBeCalled();
            });
            it('scrolls up to selected row', () => {
                const initialRow = pageSize * 2;
                jest.spyOn(container, 'scrollTop', 'get').mockReturnValue((initialRow - pageSize / 2) * rowHeight);
                setup({...options, headerSelector, initialRow});
                const event = mockKeyEvent('PageUp', false);

                hook.onKeyDown!(event);

                expect(container.scrollTo).toBeCalledWith({top: (initialRow - pageSize) * rowHeight});
            });
            it('scrolls down to selected row', () => {
                const initialRow = pageSize * 2;
                jest.spyOn(container, 'scrollTop', 'get').mockReturnValue((initialRow - pageSize / 2) * rowHeight);
                setup({...options, headerSelector, initialRow});
                const event = mockKeyEvent('PageDown', false);

                hook.onKeyDown!(event);

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
                    setup({...options, initialRow});
                    const event = mockKeyEvent(key, ctrl);

                    hook.onKeyDown!(event);

                    expect(hook.row).toEqual(expectedRow);
                    expect(event.preventDefault).toBeCalled();
                })
            );

            const columnTests = [
                {name: 'selects previous column on left arrow', key: 'ArrowLeft', ctrl: false, initialColumn: 2, expectedColumn: 1},
                {name: 'wraps to last column on left arrow', key: 'ArrowLeft', ctrl: false, initialColumn: 0, expectedColumn: columns - 1},
                {name: 'selects next column on right arrow', key: 'ArrowRight', ctrl: false, initialColumn: 0, expectedColumn: 1},
                {name: 'wraps to first column on right arrow', key: 'ArrowRight', ctrl: false, initialColumn: columns - 1, expectedColumn: 0},
                {name: 'selects first column on home', key: 'Home', ctrl: false, initialColumn: columns - 1, expectedColumn: 0},
                {name: 'remains at current column on ctrl+home', key: 'Home', ctrl: true, initialColumn: 1, expectedColumn: 1},
                {name: 'selects last column on end', key: 'End', ctrl: false, initialColumn: 0, expectedColumn: columns - 1},
                {name: 'remains at current column on ctrl+end', key: 'End', ctrl: true, initialColumn: 1, expectedColumn: 1},
            ];
            columnTests.forEach(({name, initialColumn, key, ctrl, expectedColumn}) =>
                it(name, () => {
                    setup(options);
                    for (let index = 0; index < initialColumn; index++) hook.onKeyDown!(mockKeyEvent('ArrowRight', false));
                    const event = mockKeyEvent(key, ctrl);

                    hook.onKeyDown!(event);

                    expect(hook.column).toEqual(expectedColumn);
                })
            );
        });
        describe('onMouseDown', () => {
            const container = document.createElement('div');
            const rows = new Array(50).fill(null).map(() => container.appendChild(document.createElement('div')));

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
                    setup(options);

                    hook.onMouseDown!(mockMouseEvent(cells[column]));

                    expect(hook.row).toEqual(rowIndex);
                    expect(hook.column).toEqual(column);
                });
                it('uses colSpan to calculate column index', () => {
                    Object.defineProperty(cells[1], 'colSpan', {value: 2});
                    setup(options);

                    hook.onMouseDown!(mockMouseEvent(cells[2]));

                    expect(hook.row).toEqual(rowIndex);
                    expect(hook.column).toEqual(3);
                });
            });
            describe('non-table', () => {
                const boundingRect = (index: number) => ({bottom: (index + 1) * rowHeight} as DOMRect);

                beforeAll(() => rows.forEach((row) => row.setAttribute('class', rowSelector.substring(1))));
                beforeEach(() => {
                    rows.forEach((row, index) => jest.spyOn(row, 'getBoundingClientRect').mockReturnValue(boundingRect(index)));
                });
                it('selects row at click', () => {
                    setup(options);

                    hook.onMouseDown!(mockMouseEvent(null, rowHeight * 4.1));

                    expect(hook.row).toEqual(4);
                });
                it('ignores click past end row', () => {
                    const initialRow = 5;
                    setup({...options, initialRow});

                    hook.onMouseDown!(mockMouseEvent(null, rowHeight * (rows.length + 1)));

                    expect(hook.row).toEqual(initialRow);
                });
            });
        });
    });
});
